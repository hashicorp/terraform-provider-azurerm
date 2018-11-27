package azurerm

import "testing"

func TestClientRequestID(t *testing.T) {
	first := clientRequestID()

	if first == "" {
		t.Fatal("no client request ID generated")
	}

	second := clientRequestID()
	if first != second {
		t.Fatal("subsequent request ID not the same as the first")
	}
}
