package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	defaultEnv       = "dev"
	envVarAppEnv     = "APP_ENV"
	envVarConfigPath = "CONFIG_PATH"
	envPrefix        = "GO_PG_HTTP_"
)

func Load() (Config, error) {
	k := koanf.New(".")

	configPath := resolveConfigPath()

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("load config file %q: %w", configPath, err)
	}

	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: envPrefix,
		TransformFunc: func(key, value string) (string, any) {
			key = strings.TrimPrefix(key, envPrefix)
			key = strings.ToLower(key)
			key = strings.ReplaceAll(key, "_", ".")
			return key, value
		},
	}), nil); err != nil {
		return Config{}, fmt.Errorf("load env config: %w", err)
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{
		Tag: "koanf",
		DecoderConfig: &mapstructure.DecoderConfig{
			Result:           &cfg,
			WeaklyTypedInput: true,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
			),
		},
	}); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func resolveConfigPath() string {
	if path := strings.TrimSpace(os.Getenv(envVarConfigPath)); path != "" {
		return path
	}

	appEnv := strings.TrimSpace(os.Getenv(envVarAppEnv))
	if appEnv == "" {
		appEnv = defaultEnv
	}

	return fmt.Sprintf("configs/config.%s.yaml", appEnv)
}
