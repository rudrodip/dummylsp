package rpc_test

import (
	"testing"

	"github.com/rudrodip/dummylsp/rpc"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expectd := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})
	if actual != expectd {
		t.Fatalf("Expected %s, got %s", expectd, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	contentLength, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatalf("Error decoding message: %s", err)
	}

	if contentLength != 16 {
		t.Fatalf("Expected content length of 16, got %d", contentLength)
	}
}
