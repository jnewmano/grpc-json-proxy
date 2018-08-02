package codec

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(JSON{
		Marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			OrigName:     true,
		},
	})
}

type JSON struct {
	jsonpb.Marshaler
	jsonpb.Unmarshaler
}

func (_ JSON) Name() string {
	return "json"
}

func (j JSON) Marshal(v interface{}) (out []byte, err error) {
	if pm, ok := v.(proto.Message); ok {
		b := new(bytes.Buffer)
		err := j.Marshaler.Marshal(b, pm)
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}
	return json.Marshal(v)
}

func (j JSON) Unmarshal(data []byte, v interface{}) (err error) {
	if pm, ok := v.(proto.Message); ok {
		b := bytes.NewBuffer(data)
		return j.Unmarshaler.Unmarshal(b, pm)
	}
	return json.Unmarshal(data, v)
}
