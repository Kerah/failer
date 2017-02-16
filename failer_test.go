package failer

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"encoding/binary"
)

func TestFail(t *testing.T) {
	ass := assert.New(t)
	req := require.New(t)
	var code uint32 = 42
	f := WithTag("assert", "test", code)
	f.(*fail).stack = "stack"

	data := f.Encode()
	req.NotEmpty(data)

	ass.Len(data, 27)

	ass.Equal(code, binary.BigEndian.Uint32(data[:4]))
	ass.Equal("test", string(data[12:16]))
	ass.Equal("assert", string(data[16:22]))
	ass.Equal("stack", string(data[22:27]))

	a, err := Decode(data)
	req.NoError(err)
	ass.Equal(f.Code(), a.Code())
	ass.Equal(f.Tag(), a.Tag())
	ass.Equal(f.Stack(), a.Stack())
	ass.Equal(f.Message(), a.Message())

	data[11] = []byte("z")[0]
	a, err = Decode(data)
	ass.Error(err)
	ass.Nil(a)

	a, err = Decode([]byte{})
	ass.Error(err)
	ass.Nil(a)

	a = New("test", ErrCodeDecodeFail)
	ass.NotNil(a)
	ass.Equal("test", a.Message())
	ass.Equal(ErrCodeDecodeFail, a.Code())
}
