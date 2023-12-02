package common

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
)

type fieldValue struct {
	field reflect.Value
	value interface{}
}

// GetConfig loads configuration defined by type T.
// Initial configuration is loaded from file defined by the "config" flag,
// File values could be overridden by the command line flags.
func GetConfig[T any](args []string, configPathFlag string) (*T, error) {
	fs := flag.NewFlagSet("config flags", flag.ContinueOnError)

	var configPath *string
	if configPathFlag != "" {
		configPath = fs.String(configPathFlag, "", "Path to the config file")
	}

	config := new(T)
	fieldFlags, err := defineConfigFlags[T](fs, config)
	if err != nil {
		return nil, err
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if configPath != nil && *configPath != "" {
		if err := loadConfigFromFile(config, *configPath); err != nil {
			return nil, err
		}
	}
	fs.Visit(func(f *flag.Flag) {
		if f.Name == configPathFlag {
			return
		}
		fv := fieldFlags[f.Name]
		fv.field.Set(reflect.ValueOf(fv.value))
	})
	return config, nil
}

func loadConfigFromFile[T any](cfg *T, configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}

func defineConfigFlags[T any](fs *flag.FlagSet, config *T) (map[string]*fieldValue, error) {
	elem := reflect.ValueOf(config).Elem()
	typeDef := elem.Type()
	fieldByName := map[string]*fieldValue{}

	for i := 0; i < typeDef.NumField(); i++ {
		field := elem.Field(i)
		kind := field.Type().Kind()
		name := typeDef.Field(i).Name
		if kind == reflect.Ptr {
			kind = field.Type().Elem().Kind()
		}
		fv := &fieldValue{field: field}
		fieldByName[name] = fv
		switch kind {
		case reflect.Bool:
			fv.value = fs.Bool(name, false, fmt.Sprintf("Flag for parameter %q", name))
		case reflect.String:
			fv.value = fs.String(name, "", fmt.Sprintf("Flag for parameter %q", name))
		case reflect.Int:
			fv.value = fs.Int(name, 0, fmt.Sprintf("Flag for parameter %q", name))
		case reflect.Int64:
			fv.value = fs.Int64(name, 0, fmt.Sprintf("Flag for parameter %q", name))
		case reflect.Uint:
			fv.value = fs.Uint(name, 0, fmt.Sprintf("Flag for parameter %q", name))
		case reflect.Uint64:
			fv.value = fs.Uint64(name, 0, fmt.Sprintf("Flag for parameter %q", name))
		case reflect.Float64:
			fv.value = fs.Float64(name, 0, fmt.Sprintf("Flag for parameter %q", name))
		default:
			return nil, fmt.Errorf("Unsupported field %s of type %s", name, kind)
		}
	}

	return fieldByName, nil
}
