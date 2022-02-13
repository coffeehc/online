package netutil

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRoute(t *testing.T) {
	test := assert.New(t)
	iface, gw, src, err := Route(3*time.Second, "8.8.8.8")
	if !test.Nil(err) {
		t.FailNow()
	}
	spew.Dump(iface, gw, src)
}
