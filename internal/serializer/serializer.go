package serializer

import (
	"bytes"
	"encoding/gob"
	"time"
)

type MessageHeader struct {
	Time uint64
	ID   uint32
	Data []byte
}

func EncodeMessage(msgid uint32, inmsg any) ([]byte, error) {
	msg := MessageHeader{Time: uint64(time.Now().Unix()), ID: msgid}

	inmsgBuffer := &bytes.Buffer{}
	err := gob.NewEncoder(inmsgBuffer).Encode(inmsg)
	if err != nil {
		return []byte{}, err
	}
	msg.Data = inmsgBuffer.Bytes()

	headerBuffer := &bytes.Buffer{}
	err = gob.NewEncoder(headerBuffer).Encode(msg)
	if err != nil {
		return []byte{}, err
	}

	return append(headerBuffer.Bytes(), inmsgBuffer.Bytes()...), nil
}

func DecodeMessage(data []byte, msg any) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(msg); err != nil {
		return err
	}
	return nil
}
