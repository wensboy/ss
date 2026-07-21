package config

import (
	"testing"
)

func Test_Config(t *testing.T) {
	// at := assert.New(t)
	v := _global_config.MustLookup("rpc.server.listen")
	t.Log(v)
}

func Test_ConfigKey(t *testing.T) {
	v := _global_config.Key("rpc", "server  ", "listen")
	t.Log(v)
}

func Test_ConfigVar(t *testing.T) {
	v, ok := ConfigVar("rpc.server.listen").Lookup()
	if !ok {
		t.Fatal("failed to lookup config value")
	}
	t.Log(v)
}
