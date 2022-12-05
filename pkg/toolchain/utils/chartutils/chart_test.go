package chartutils

import (
	_ "embed"
	"testing"
)

func TestNewChart(t *testing.T) {
	c, err := NewChart("test", []byte("do"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Save(); err != nil {
		t.Fatal(err)
	}
}
