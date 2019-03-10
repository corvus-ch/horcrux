package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type value struct {
	key   string
	viper *viper.Viper
}

func (v *value) Get() interface{} {
	return v.viper.Get(v.key)
}

func (v *value) String() string {
	return v.viper.GetString(v.key)
}

func (v *value) Set(value string) error {
	v.viper.Set(v.key, value)
	return nil
}

type stringValue struct {
	value
}

type stringsValue struct {
	value
}

func (v *stringsValue) String() string {
	return strings.Join(v.viper.GetStringSlice(v.key), ", ")
}

func (v *stringsValue) Set(value string) error {
	v.viper.Set(v.key, append(v.viper.GetStringSlice(v.key), value))
	return nil
}

func (v *stringsValue) IsCumulative() bool {
	return true
}

type intValue struct {
	value
}

func (v *intValue) Set(value string) error {
	_, err := strconv.ParseInt(value, 0, 0)
	if err == nil {
		v.viper.Set(v.key, value)
	}
	return err
}

type boolValue struct {
	value
}

func (v *boolValue) IsBoolFlag() bool {
	return true
}

func (v *boolValue) Set(value string) error {
	_, err := strconv.ParseBool(value)
	if err == nil {
		v.viper.Set(v.key, value)
	}
	return err
}

type enumValue struct {
	value
	options []string
}

func (v *enumValue) Set(value string) error {
	for _, opt := range v.options {
		if opt == value {
			v.viper.Set(v.key, value)
			return nil
		}
	}
	return fmt.Errorf("enum value must be one of %v, got '%v'", strings.Join(v.options, ","), value)
}

func (v *enumValue) String() string {
	return v.viper.GetString(v.key)
}

type enumsValue enumValue

func (v *enumsValue) Set(value string) error {
	for _, opt := range v.options {
		if opt == value {
			v.viper.Set(v.key, append(v.viper.GetStringSlice(v.key), value))
			return nil
		}
	}
	return fmt.Errorf("enum value must be one of %v, got '%v'", strings.Join(v.options, ","), value)
}

func (v *enumsValue) String() string {
	return strings.Join(v.viper.GetStringSlice(v.key), ", ")
}

func (v *enumsValue) IsCumulative() bool {
	return true
}
