package elastifile

import (
	"fmt"
	"reflect"
)

type config struct {
	NFSServer       string `parameter:"nfsServer"`
	EmanageURL      string `parameter:"restURL"`
	Username        string `parameter:"username"`
	Password        string
	SecretName      string `parameter:"secretName"`
	SecretNamespace string `parameter:"secretNamespace"`
}

func newConfigFromParams(params map[string]string) (*config, error) {
	return newConfigFromMap(params, "")
}

func newConfigFromAnnotations(ann map[string]string) (*config, error) {
	return newConfigFromMap(ann, annPrefix)
}

func newConfigFromMap(m map[string]string, prefix string) (*config, error) {
	conf := &config{}

	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Tag.Get("parameter")
		if name == "" {
			continue
		}
		par, ok := m[prefix+name]
		if !ok {
			return nil, fmt.Errorf("missing StorageClass parameter: %q", name)
		}
		v.Field(i).SetString(par)
	}

	return conf, nil
}

func (conf *config) setAnnotations(ann map[string]string) {
	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Tag.Get("parameter")
		if name == "" {
			continue
		}
		value := fmt.Sprintf("%v", v.Field(i).Interface())
		ann[annPrefix+name] = value
	}
}
