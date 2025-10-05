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
    data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: c\r\n\r\n",
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
    data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: c\r\n\r\n",
    noBytesPerRead: 1,
}
r, err = RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}


func TestParseHeaders(t *testing.T) {
    // Test: Standard Headers
    reader := &chunkReader{
        data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
        noBytesPerRead: 3,
    }
    r, err := RequestFromReader(reader)
    require.NoError(t, err)
    require.NotNil(t, r)
    val,_ := r.Headers.Get("host")
    assert.Equal(t, "localhost:42069",val )
    val,_ = r.Headers.Get("user-agent")
    assert.Equal(t, "curl/7.81.0", val)
    val,_ = r.Headers.Get("accept")
    assert.Equal(t, "*/*", val)

    // Test: Malformed Header
    reader = &chunkReader{
        data: "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
        noBytesPerRead: 3,
    }
    _, err = RequestFromReader(reader)
    require.Error(t, err)
}

func TestParseBody(t *testing.T){

    // Test: full body read
    reader := &chunkReader{
        data: "POST /submit HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"hello world!\n",
        noBytesPerRead: 3,
    }
    
    r, err := RequestFromReader(reader)
    require.NoError(t, err)
    require.NotNil(t, r)    
    assert.Equal(t, "hello world!\n", string(r.Body))

// Test: Body shorter than reported content length
reader = &chunkReader{
	data: "POST /submit HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"Content-Length: 20\r\n" +
		"\r\n" +
		"partial content",
        noBytesPerRead: 3,
    }
    
    r, err = RequestFromReader(reader)
    require.Error(t, err)
require.Nil(t, r)

}