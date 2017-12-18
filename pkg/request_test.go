package pkg

import (
	"testing"
)

func TestAddr(t *testing.T) {
	tgt := &Target{
		Host: "host",
		Port: 50,
	}
	if tgt.Addr() != "host:50" {
		t.Error("Addr() should be host:50")
	}
}
