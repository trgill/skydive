// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package topology

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2220f231DecodeGithubComSkydiveProjectSkydiveTopology(in *jlexer.Lexer, out *InterfaceMetric) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Collisions":
			out.Collisions = int64(in.Int64())
		case "Multicast":
			out.Multicast = int64(in.Int64())
		case "RxBytes":
			out.RxBytes = int64(in.Int64())
		case "RxCompressed":
			out.RxCompressed = int64(in.Int64())
		case "RxCrcErrors":
			out.RxCrcErrors = int64(in.Int64())
		case "RxDropped":
			out.RxDropped = int64(in.Int64())
		case "RxErrors":
			out.RxErrors = int64(in.Int64())
		case "RxFifoErrors":
			out.RxFifoErrors = int64(in.Int64())
		case "RxFrameErrors":
			out.RxFrameErrors = int64(in.Int64())
		case "RxLengthErrors":
			out.RxLengthErrors = int64(in.Int64())
		case "RxMissedErrors":
			out.RxMissedErrors = int64(in.Int64())
		case "RxOverErrors":
			out.RxOverErrors = int64(in.Int64())
		case "RxPackets":
			out.RxPackets = int64(in.Int64())
		case "TxAbortedErrors":
			out.TxAbortedErrors = int64(in.Int64())
		case "TxBytes":
			out.TxBytes = int64(in.Int64())
		case "TxCarrierErrors":
			out.TxCarrierErrors = int64(in.Int64())
		case "TxCompressed":
			out.TxCompressed = int64(in.Int64())
		case "TxDropped":
			out.TxDropped = int64(in.Int64())
		case "TxErrors":
			out.TxErrors = int64(in.Int64())
		case "TxFifoErrors":
			out.TxFifoErrors = int64(in.Int64())
		case "TxHeartbeatErrors":
			out.TxHeartbeatErrors = int64(in.Int64())
		case "TxPackets":
			out.TxPackets = int64(in.Int64())
		case "TxWindowErrors":
			out.TxWindowErrors = int64(in.Int64())
		case "Start":
			out.Start = int64(in.Int64())
		case "Last":
			out.Last = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2220f231EncodeGithubComSkydiveProjectSkydiveTopology(out *jwriter.Writer, in InterfaceMetric) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Collisions != 0 {
		const prefix string = ",\"Collisions\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Collisions))
	}
	if in.Multicast != 0 {
		const prefix string = ",\"Multicast\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Multicast))
	}
	if in.RxBytes != 0 {
		const prefix string = ",\"RxBytes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxBytes))
	}
	if in.RxCompressed != 0 {
		const prefix string = ",\"RxCompressed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxCompressed))
	}
	if in.RxCrcErrors != 0 {
		const prefix string = ",\"RxCrcErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxCrcErrors))
	}
	if in.RxDropped != 0 {
		const prefix string = ",\"RxDropped\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxDropped))
	}
	if in.RxErrors != 0 {
		const prefix string = ",\"RxErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxErrors))
	}
	if in.RxFifoErrors != 0 {
		const prefix string = ",\"RxFifoErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxFifoErrors))
	}
	if in.RxFrameErrors != 0 {
		const prefix string = ",\"RxFrameErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxFrameErrors))
	}
	if in.RxLengthErrors != 0 {
		const prefix string = ",\"RxLengthErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxLengthErrors))
	}
	if in.RxMissedErrors != 0 {
		const prefix string = ",\"RxMissedErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxMissedErrors))
	}
	if in.RxOverErrors != 0 {
		const prefix string = ",\"RxOverErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxOverErrors))
	}
	if in.RxPackets != 0 {
		const prefix string = ",\"RxPackets\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RxPackets))
	}
	if in.TxAbortedErrors != 0 {
		const prefix string = ",\"TxAbortedErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxAbortedErrors))
	}
	if in.TxBytes != 0 {
		const prefix string = ",\"TxBytes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxBytes))
	}
	if in.TxCarrierErrors != 0 {
		const prefix string = ",\"TxCarrierErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxCarrierErrors))
	}
	if in.TxCompressed != 0 {
		const prefix string = ",\"TxCompressed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxCompressed))
	}
	if in.TxDropped != 0 {
		const prefix string = ",\"TxDropped\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxDropped))
	}
	if in.TxErrors != 0 {
		const prefix string = ",\"TxErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxErrors))
	}
	if in.TxFifoErrors != 0 {
		const prefix string = ",\"TxFifoErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxFifoErrors))
	}
	if in.TxHeartbeatErrors != 0 {
		const prefix string = ",\"TxHeartbeatErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxHeartbeatErrors))
	}
	if in.TxPackets != 0 {
		const prefix string = ",\"TxPackets\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxPackets))
	}
	if in.TxWindowErrors != 0 {
		const prefix string = ",\"TxWindowErrors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TxWindowErrors))
	}
	if in.Start != 0 {
		const prefix string = ",\"Start\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Start))
	}
	if in.Last != 0 {
		const prefix string = ",\"Last\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Last))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InterfaceMetric) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2220f231EncodeGithubComSkydiveProjectSkydiveTopology(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InterfaceMetric) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2220f231EncodeGithubComSkydiveProjectSkydiveTopology(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InterfaceMetric) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2220f231DecodeGithubComSkydiveProjectSkydiveTopology(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InterfaceMetric) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2220f231DecodeGithubComSkydiveProjectSkydiveTopology(l, v)
}
