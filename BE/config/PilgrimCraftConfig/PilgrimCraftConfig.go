package PilgrimCraftConfig

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type PilgrimCraftConfig struct {
	Server              ServerConfig `mapstructure:"SERVER"`
	MinimumLoggingLevel slog.Level   `mapstructure:"MIN_LOGGING_LEVEL"`
}

type ServerConfig struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

func LoadSettings(ctx context.Context) *PilgrimCraftConfig {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file (.env) not found. Using defaults or environment variables.")
		} else {
			panic(fmt.Errorf("error reading config file: %w", err))
		}
	}

	var config PilgrimCraftConfig
	if err := viper.Unmarshal(&config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			stringToSlogLevelHook(),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	)); err != nil {
		panic(fmt.Errorf("PilgrimCraftConfig.LoadSettings: Unable to decode config into struct: %w", err))
	}

	return &config
}

func stringToSlogLevelHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(slog.LevelInfo) {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("PilgrimCraftConfig.GetConfig: Expected string for slog.Level conversion, but got %T", data)
		}
		return stringToSlogLevel(str)
	}
}

func stringToSlogLevel(levelStr string) (slog.Level, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level string: %q, defaulting to INFO", levelStr)
	}
}
