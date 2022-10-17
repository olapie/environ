package environ

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.olapie.com/conv"
)

type DefaultManager struct {
	m    map[string]any
	args []string
	envs []string
}

func NewDefaultManager() *DefaultManager {
	m := &DefaultManager{
		m: map[string]any{},
	}
	em := conv.OSArgsToEnvMap(os.Args[1:])
	for k, v := range em {
		m.m[k] = v
	}
	return m
}

func (m *DefaultManager) Has(key string) bool {
	key = conv.ToEnvKey(key)
	_, ok := m.m[key]
	if ok {
		return true
	}
	_, ok = m.getOSEnv(key)
	return ok
}

func (m *DefaultManager) Get(key string) any {
	key = conv.ToEnvKey(key)
	v, ok := m.m[key]
	if ok {
		return v
	}
	v, ok = m.getOSEnv(key)
	if !ok {
		return nil
	}
	return v
}

func (m *DefaultManager) Set(key string, value any) {
	key = conv.ToEnvKey(key)
	m.m[key] = value
}

func (m *DefaultManager) getOSEnv(key string) (string, bool) {
	em := conv.OSEnvsToEnvMap(os.Environ())
	v, ok := em[key]
	return v, ok
}

func (m *DefaultManager) LoadConfigFile(filename string) error {
	switch filepath.Ext(filename) {
	case ".json":
		var obj map[string]any
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		err = json.Unmarshal(content, &obj)
		if err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		em := conv.ToEnvMap(obj)
		for k, v := range em {
			m.m[k] = v
		}
		return nil
	default:
		return errors.New("unsupported file format")
	}
}
