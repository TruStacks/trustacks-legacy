package standard

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewDesiredConfig(t *testing.T) {
	conf, err := NewDesiredConfig(map[string]interface{}{"domain": "test.local.gd"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test.local.gd", conf.Profile.Domain)
}
