/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package alert

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/redhat-cip/skydive/config"
	"github.com/redhat-cip/skydive/logging"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 1024 * 1024
)

type Server struct {
	AlertManager *AlertManager
	Router       *mux.Router
	wsServer     *WSServer
	Host         string
	wg           sync.WaitGroup
	listening    atomic.Value
}

type WSClient struct {
	conn   *websocket.Conn
	read   chan []byte
	send   chan []byte
	server *WSServer
}

type WSServer struct {
	AlertManager *AlertManager
	clients      map[*WSClient]bool
	broadcast    chan string
	quit         chan bool
	register     chan *WSClient
	unregister   chan *WSClient
	pongWait     time.Duration
	pingPeriod   time.Duration
}

/* Called by alert.EvalNodes() */
func (c *WSClient) OnAlert(amsg *AlertMessage) {
	c.send <- amsg.Marshal()
}

func (c *WSClient) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.server.pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.server.pongWait))
		return nil
	})

	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.read <- m
	}
}

func (c *WSClient) writePump(wg *sync.WaitGroup, quit chan struct{}) {
	ticker := time.NewTicker(c.server.pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				wg.Done()
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				logging.GetLogger().Warningf("Error while writing to the websocket: %s", err.Error())
				wg.Done()
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				wg.Done()
				return
			}
		case <-quit:
			wg.Done()
			return
		}
	}
}

func (c *WSClient) write(mt int, message []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, message)
}

func (s *WSServer) ListenAndServe() {
	for {
		select {
		case <-s.quit:
			return
		case c := <-s.register:
			s.clients[c] = true
			s.AlertManager.AddEventListener(c)
		case c := <-s.unregister:
			if _, ok := s.clients[c]; ok {
				s.AlertManager.DelEventListener(c)
				delete(s.clients, c)
			}
		case m := <-s.broadcast:
			s.broadcastMessage(m)
		}
	}
}

func (s *WSServer) broadcastMessage(m string) {
	for c := range s.clients {
		select {
		case c.send <- []byte(m):
		default:
			delete(s.clients, c)
		}
	}
}

func (s *Server) serveMessages(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &WSClient{
		read:   make(chan []byte, maxMessageSize),
		send:   make(chan []byte, maxMessageSize),
		conn:   conn,
		server: s.wsServer,
	}
	logging.GetLogger().Infof("New WebSocket Connection from %s : URI path %s", conn.RemoteAddr().String(), r.URL.Path)

	s.wsServer.register <- c

	var wg sync.WaitGroup
	wg.Add(2)

	quit := make(chan struct{})

	go c.writePump(&wg, quit)

	c.readPump()

	quit <- struct{}{}
	quit <- struct{}{}

	close(c.read)
	close(c.send)

	wg.Wait()
}

func (s *Server) ListenAndServe() {
	s.wg.Add(1)
	defer s.wg.Done()

	s.listening.Store(true)
	s.wsServer.ListenAndServe()
}

func (s *Server) Stop() {
	s.wsServer.quit <- true
	if s.listening.Load() == true {
		s.wg.Wait()
	}
}

func NewServer(a *AlertManager, router *mux.Router, pongWait time.Duration) *Server {
	s := &Server{
		AlertManager: a,
		Router:       router,
		wsServer: &WSServer{
			AlertManager: a,
			broadcast:    make(chan string, 500),
			quit:         make(chan bool, 1),
			register:     make(chan *WSClient),
			unregister:   make(chan *WSClient),
			clients:      make(map[*WSClient]bool),
			pongWait:     pongWait,
			pingPeriod:   (pongWait * 8) / 10,
		},
	}

	s.Router.HandleFunc("/ws/alert", s.serveMessages)

	return s
}

func NewServerFromConfig(a *AlertManager, router *mux.Router) (*Server, error) {
	w := config.GetConfig().GetInt("ws_pong_timeout")

	return NewServer(a, router, time.Duration(w)*time.Second), nil
}