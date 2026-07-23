package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func init() {
	_global_config = NewConfig()
	_global_config.SetDir("spec")
	_global_config.SetSep(".")
	err := _global_config.Load(
		"config.json",
		"restful.json",
		"rpc.json",
	)
	if err != nil {
		panic(err)
	}
}

const (
	_default_config_dir = ""
	_default_config_sep = "."
	_default_config_ext = "json"
)

var (
	_global_config *Config
)

func GConfigSource(path ...string) Source {
	return ConfigSource(_global_config, path...)
}

func ConfigSource(cfg *Config, path ...string) Source {
	return func() (any, bool) {
		return cfg.Lookup(path...)
	}
}

// key be like "restful.server.port"
func ConfigVar(key string) cli.ValueSource {
	return &ConfigValueSource{key: key}
}

func ConfigVars(keys ...string) cli.ValueSourceChain {
	vsc := []cli.ValueSource{}
	for _, key := range keys {
		vsc = append(vsc, ConfigVar(key))
	}
	return cli.NewValueSourceChain(vsc...)
}

type SpecConfig map[string]any

func (sc SpecConfig) Lookup(path ...string) (any, bool) {
	if len(path) == 0 {
		return nil, false
	}
	v, exists := sc[path[0]]
	if len(path) == 1 {
		return v, exists
	}
	if !exists {
		return nil, false
	}
	switch m := v.(type) {
	case SpecConfig:
		return m.Lookup(path[1:]...)
	case map[string]any:
		return SpecConfig(m).Lookup(path[1:]...)
	}
	return nil, false
}

func (sc SpecConfig) MustLookup(path ...string) any {
	if v, ok := sc.Lookup(path...); ok {
		return v
	}
	panic("config - entry not found")
}

type Config struct {
	spec SpecConfig
	dir  string
	sep  string
	ext  string
	wg   sync.WaitGroup
	eg   errgroup.Group
	mx   sync.Mutex
}

func NewConfig() *Config {
	return &Config{
		spec: make(SpecConfig),
		dir:  _default_config_dir,
		sep:  _default_config_sep,
	}
}

func (c *Config) SetDir(dir string) {
	c.dir = dir
}

func (c *Config) SetSep(sep string) {
	c.sep = sep
}

func (c *Config) SetExt(ext string) {
	c.ext = ext
}

func (c *Config) Key(paths ...string) string {
	return strings.ReplaceAll(strings.Join(paths, c.sep), " ", "")
}

func (c *Config) Load(paths ...string) error {
	if len(paths) == 0 {
		return nil
	}
	for _, path := range paths {
		// 原则上不允许出现: config.dev.json 样式的配置文件, 否则只有 config 会生效
		specName, ext, found := strings.Cut(path, ".")
		if !found {
			ext = _default_config_ext
		}
		c.eg.Go(func() error {
			return c.LoadSpec(specName, ext)
		})
	}
	return c.eg.Wait()
}

func (c *Config) LoadSpec(specName, ext string) error {
	path := filepath.Clean(filepath.Join(c.dir, specName+"."+ext))
	sc := make(map[string]any)
	fd, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	decoder := json.NewDecoder(fd)
	if err := decoder.Decode(&sc); err != nil {
		return err
	}
	c.mx.Lock()
	c.spec[specName] = sc
	c.mx.Unlock()
	return nil
}

func (c *Config) Lookup(path ...string) (any, bool) {
	if len(path) == 1 {
		path = strings.Split(strings.ReplaceAll(path[0], " ", ""), c.sep)
	}
	return c.spec.Lookup(path...)
}

func (c *Config) MustLookup(path ...string) any {
	if v, ok := c.Lookup(path...); ok {
		return v
	}
	panic("config - entry not found")
}

type ConfigValueSource struct {
	key string
}

func (cvs *ConfigValueSource) String() string {
	v, ok := cvs.Lookup()
	if !ok {
		return ""
	}
	return v
}

func (cvs *ConfigValueSource) GoString() string {
	return fmt.Sprintf("%s", cvs.String())
}

func (cvs *ConfigValueSource) Lookup() (string, bool) {
	v := _global_config.MustLookup(cvs.key)
	if s, ok := v.(string); ok {
		return s, true
	}
	return "", false
}
