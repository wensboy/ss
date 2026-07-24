package config

import (
	"os"
	"strings"
)

func init() {
	_global_env = NewEnv()
	_global_env.AppendGroup(NewEnvGroup(_default_env_prefix, _default_env_separator))
}

const (
	_default_env_separator = "_"
	_default_env_prefix    = ""
)

var (
	_global_env *Env
)

func GEnvSource(path ...string) Source {
	return EnvSource(_global_env, path...)
}

func EnvSource(env *Env, path ...string) Source {
	return func() (any, bool) {
		return env.Lookup(path...)
	}
}

func AppendEnvGroup(eg EnvGroup) {
	_global_env.AppendGroup(eg)
}

func GetEnvGroup(prefix string) (EnvGroup, bool) {
	return _global_env.GetGroup(prefix)
}

func GetEnvKey(path ...string) string {
	eg := _global_env.MustGetGroup(_default_env_prefix)
	return eg.Key(path...)
}

func MustGetEnvGroup(prefix string) EnvGroup {
	return _global_env.MustGetGroup(prefix)
}

func LookupEnv(path ...string) (string, bool) {
	return _global_env.Lookup(path...)
}

func MustLookupEnv(path ...string) string {
	return _global_env.MustLookup(path...)
}

func LookupEnvFromGroup(prefix string, path ...string) (string, bool) {
	return _global_env.LookupFromGroup(prefix, path...)
}

func MustLookupEnvFromGroup(prefix string, path ...string) string {
	return _global_env.MustLookupFromGroup(prefix, path...)
}

type Env struct {
	envGroups []EnvGroup
	groupMap  map[string]int
}

func NewEnv() *Env {
	return &Env{
		envGroups: []EnvGroup{},
		groupMap:  make(map[string]int),
	}
}

func (e *Env) AppendGroup(eg EnvGroup) {
	if _, exists := e.groupMap[eg.prefix]; !exists {
		e.envGroups = append(e.envGroups, eg)
		e.groupMap[eg.prefix] = len(e.envGroups) - 1
	}
}

func (e *Env) GetGroup(prefix string) (EnvGroup, bool) {
	if index, exists := e.groupMap[prefix]; exists {
		return e.envGroups[index], true
	}
	return EnvGroup{}, false
}

func (e *Env) MustGetGroup(prefix string) EnvGroup {
	eg, ok := e.GetGroup(prefix)
	if !ok {
		panic("env - environment group not found: " + prefix)
	}
	return eg
}

func (e *Env) Lookup(path ...string) (string, bool) {
	if len(path) == 1 {
		path = strings.Split(path[0], _default_env_separator)
	}
	if len(path) == 1 {
		path = strings.Split(path[0], _default_config_separator)
	}
	for _, eg := range e.envGroups {
		if value, ok := eg.Lookup(path...); ok {
			return value, true
		}
	}
	return "", false
}

func (e *Env) MustLookup(path ...string) string {
	value, ok := e.Lookup(path...)
	if !ok {
		panic("env - environment variable not found: " + strings.Join(path, "."))
	}
	return value
}

func (e *Env) LookupFromGroup(prefix string, path ...string) (string, bool) {
	eg, ok := e.GetGroup(prefix)
	if !ok {
		return "", false
	}
	return eg.Lookup(path...)
}

func (e *Env) MustLookupFromGroup(prefix string, path ...string) string {
	eg := e.MustGetGroup(prefix)
	return eg.MustLookup(path...)
}

type EnvGroup struct {
	prefix    string
	separator string
}

func NewEnvGroup(prefix, separator string) EnvGroup {
	if separator == "" {
		separator = _default_env_separator
	}
	return EnvGroup{
		prefix:    prefix,
		separator: separator,
	}
}

func (e *EnvGroup) Key(path ...string) string {
	if e.prefix != "" {
		path = append([]string{e.prefix}, path...)
	}
	return strings.ToUpper(strings.ReplaceAll(strings.Join(path, e.separator), " ", ""))
}

func (e *EnvGroup) Lookup(path ...string) (string, bool) {
	return os.LookupEnv(e.Key(path...))
}

func (e *EnvGroup) MustLookup(path ...string) string {
	value, ok := e.Lookup(path...)
	if !ok {
		panic("env - environment variable not found: " + e.Key(path...))
	}
	return value
}
