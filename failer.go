package failer

import (
	"encoding/binary"
	"fmt"
)

const (
	ErrCodeDecodeFail uint32 = 1
)

type Encoder interface {
	Encode() ([]byte)
}

type Decoder interface {
	Decode([]byte) Fail
}

var emptyString = ""

type Fail interface {
	Code() uint32
	Message() string
	Error() string
	Tag() string
	Stack() string

	Encoder
	Decoder



}

type fail struct {
	code    uint32
	tag     string
	message string
	stack   string
}

func (f *fail) Code() uint32 {
	return f.code
}

func (f *fail) Tag() string {
	return f.tag
}

func (f *fail) Stack() string {
	return f.stack
}

func (f *fail) Decode(data []byte) Fail {
	if len(data) < 12 {
		return WithTag("decode,fail", "expect headers", ErrCodeDecodeFail)
	}
	f.code = binary.BigEndian.Uint32(data[:4])
	msgLen := int(binary.BigEndian.Uint16(data[4:6]))
	stackLen := int(binary.BigEndian.Uint32(data[8:12]))
	tagLen := int(binary.BigEndian.Uint16(data[6:8]))
	if msgLen+stackLen+tagLen+12 != len(data) {
		return WithTag("decode,fail", "invalid body length", ErrCodeDecodeFail)
	}

	f.message = string(data[12:12+msgLen])
	f.tag = string(data[12+msgLen:12+msgLen+tagLen])
	f.stack = string(data[12+msgLen+tagLen:12+msgLen+tagLen+stackLen])
	return nil
}

func (f *fail) Encode() ([]byte) {
	result := make([]byte, 12)
	binary.BigEndian.PutUint32(result[:4], f.code)
	binary.BigEndian.PutUint16(result[4:6], uint16(len(f.message)))
	binary.BigEndian.PutUint16(result[6:8], uint16(len(f.tag)))
	binary.BigEndian.PutUint32(result[8:12], uint32(len(f.stack)))

	if len(f.message) > 0 {
		result = append(result, []byte(f.message)...)
	}

	if len(f.tag) > 0 {
		result = append(result, []byte(f.tag)...)
	}

	if len(f.stack) > 0 {
		result = append(result, []byte(f.stack)...)
	}

	return result
}

func (f *fail) Message() string {
	return f.message
}

func (f *fail) Error() string {
	result := fmt.Sprintf("%d: %s", f.code, f.message)
	if f.tag != emptyString {
		result = "[" + f.tag + "] " + result
	}
	return result
}

func New(message string, code uint32) Fail {
	return &fail{
		code:    code,
		message: message,
	}
}

func Decode(data []byte) (Fail, Fail) {
	f := &fail{}
	if err := f.Decode(data); err != nil {
		return nil, err
	}

	return f, nil
}

func WithTag(tag, message string, code uint32) Fail {
	return &fail{
		code:    code,
		message: message,
		tag:     tag,
	}
}
