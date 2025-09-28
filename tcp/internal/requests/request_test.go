package request

import(
	
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
)

type chunkReader struct{
	data string
	noBytesPerRead int
	pos int
}

func (cr* chunkReader) Read(p []byte)(n int, err error){
	if cr.pos >= len(cr.data){
		return 0, io.EOF
	}

	endIndex := cr.pos + cr.noBytesPerRead;

	if endIndex >len(cr.data){
		endIndex = len(cr.data)
	}
	n =copy(p, cr.data[cr.pos:endIndex])
	cr.pos+=n;
	if n >cr.noBytesPerRead{
		n= cr.noBytesPerRead
		cr.pos-=n- cr.noBytesPerRead
	}
	return n ,nil
}


func TestRequestLineParse(t*testing.T){

// Test: Good GET Request Line
reader := &chunkReader{
    data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: c",
    noBytesPerRead: 3,
}
r, err := RequestFromReader(reader)

require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Good GET Request line with path
reader = &chunkReader{
    data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: c",
    noBytesPerRead: 1,
}
r, err = RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}