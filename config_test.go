package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib"
	"github.com/stretchr/testify/assert"
)

type myconfig struct {
	ID     string                 `json:"id"`
	Array  []int                  `json:"array"`
	Params map[string]interface{} `json:"params"`
}

func TestLoadConfigTo(t *testing.T) {
	conf := myconfig{}
	err := tnglib.LoadConfigTo("./_testdata/config.json", &conf)
	assert.NoError(t, err)
	assert.Equal(t, "TESTING", conf.ID)
}

func TestLoadJsonConfig(t *testing.T) {
	conf, err := tnglib.LoadConfig[myconfig]("./_testdata/config.json")
	assert.NoError(t, err)
	assert.Equal(t, "TESTING", conf.ID)
}

func TestLoadHjsonConfig(t *testing.T) {
	conf, err := tnglib.LoadConfig[myconfig]("./_testdata/config.hjson")
	assert.NoError(t, err)
	assert.Equal(t, "TESTING", conf.ID)
	assert.Len(t, conf.Array, 3)
}
