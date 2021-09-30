package tagparser

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const structTagKey = "env"

var (
	errEmptyTag    = errors.New("empty " + structTagKey + " tag")
	errKeyNotFound = errors.New(structTagKey + " key not found")
)

type ConfigOption struct {
	Name     string // The name of the environment variable.
	Required bool   // Whether the environment var is required to be set.
	Default  string // The default value if the env var isn't set.
}

func parseConfigOption(tag reflect.StructTag) (ConfigOption, error) {
	if key, ok := tag.Lookup(structTagKey); ok {
		// First value is the name of the env var, the rest is optional options, e.g. default=
		options := strings.Split(key, ",")
		if len(options) == 0 {
			return ConfigOption{}, errEmptyTag
		}
		return ConfigOption{
			Name:     options[0],
			Required: isRequired(options[1:]),
			Default:  defaultVal(options[1:]),
		}, nil
	}

	return ConfigOption{}, errKeyNotFound
}

func isRequired(opts []string) bool {
	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if strings.HasPrefix(strings.ToLower(opt), "required") {
			return true
		}
	}
	return false
}

func defaultVal(opts []string) string {
	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if strings.HasPrefix(strings.ToLower(opt), "default") {
			split := strings.Split(opt, "=")
			if len(split) != 2 {
				return ""
			}

			return split[1]
		}
	}
	return ""
}

func (o *ConfigOption) String() string {
	if o.Required {
		return fmt.Sprintf("%s=\t(required)", o.Name)
	}

	return fmt.Sprintf("%s=%s", o.Name, o.Default)
}
