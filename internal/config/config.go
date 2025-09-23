package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type ServerConfig struct {
	Port uint `mapstructure:"port"`
}

func bindEnvRecursive(viperInstance *viper.Viper, prefix string, val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			continue
		}

		fieldPath := prefix
		if prefix != "" {
			fieldPath = prefix + "." + tag
		} else {
			fieldPath = tag
		}

		if field.Type.Kind() == reflect.Struct {
			if err := bindEnvRecursive(viperInstance, fieldPath, val.Field(i)); err != nil {
				return err
			}
		} else {
			envVarName := strings.ToUpper(strings.ReplaceAll(fieldPath, ".", "_"))
			if err := viperInstance.BindEnv(fieldPath, envVarName); err != nil {
				return err
			}
		}
	}

	return nil
}

func bindAllEnvVars(viperInstance *viper.Viper) error {
	return bindEnvRecursive(viperInstance, "", reflect.ValueOf(&Config{}).Elem())
}

func LoadConfig() (Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	v2 := viper.New()
	v2.SetConfigName("override")
	v2.SetConfigType("yaml")
	v2.AddConfigPath(".")
	if err := v2.ReadInConfig(); err == nil {
		err := v.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return cfg, err
		}
	}

	if err := bindAllEnvVars(v); err != nil {
		return cfg, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
