package config

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/theartofdevel/logging"
)

type Config struct {
	App       AppConfig       `yaml:"app"`
	Bot       BotConfig       `yaml:"bot"`
	OpenAI    OpenAIConfig    `yaml:"openai"`
	Replicate ReplicateConfig `yaml:"replicate"`
	Metrics   MetricsConfig   `yaml:"metrics"`
	Tracing   TracingConfig   `yaml:"tracing"`
}

func (cfg *Config) LoggingAttrs() []logging.Attr {
	return []logging.Attr{
		logging.StringAttr("app_id", cfg.App.Id),
		logging.StringAttr("app_name", cfg.App.Name),
		logging.StringAttr("log_level", cfg.App.LogLevel),
		logging.BoolAttr("is_log_json", cfg.App.IsLogJSON),
		logging.StringAttr("bot_token", "***"+strconv.Itoa(len(cfg.Bot.Token))),
		logging.StringAttr("metrics_enabled", strconv.FormatBool(cfg.Metrics.Enabled)),
		logging.StringAttr("metrics_host", cfg.Metrics.Host),
		logging.StringAttr("metrics_port", strconv.Itoa(cfg.Metrics.Port)),
		logging.StringAttr("tracing_enabled", strconv.FormatBool(cfg.Tracing.Enabled)),
		logging.StringAttr("tracing_host", cfg.Tracing.Host),
		logging.StringAttr("tracing_port", strconv.Itoa(cfg.Tracing.Port)),
	}
}

type AppConfig struct {
	Id        string `yaml:"id" env:"APP_ID"`
	Name      string `yaml:"name" env:"APP_NAME"`
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL"`
	IsLogJSON bool   `yaml:"is_log_json" env:"IS_LOG_JSON"`
}

type BotConfig struct {
	Token     string        `yaml:"token" env:"BOT_TOKEN"`
	Timeout   time.Duration `yaml:"timeout" env:"BOT_TIMEOUT"`
	Whitelist []int64       `yaml:"whitelist" env:"BOT_WHITELIST"`
}

type OpenAIConfig struct {
	Enabled bool   `yaml:"enabled" env:"OPENAI_ENABLED"`
	ApiKey  string `yaml:"api_key" env:"OPENAI_API_KEY"`
}

type ReplicateConfig struct {
	Enabled bool   `yaml:"enabled" env:"REPLICATE_ENABLED"`
	Token   string `yaml:"token" env:"REPLICATE_TOKEN"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" env:"METRICS_ENABLED"`
	Host    string `yaml:"host" env:"METRICS_HOST"`
	Port    int    `yaml:"port" env:"METRICS_PORT"`
}

type TracingConfig struct {
	Enabled bool   `yaml:"enabled" env:"TRACING_ENABLED"`
	Host    string `yaml:"host" env:"TRACING_HOST"`
	Port    int    `yaml:"port" env:"TRACING_PORT"`
}

const (
	flagConfigPathName = "config"
	envConfigPathName  = "CONFIG_PATH"
)

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		var configPath string

		flag.StringVar(&configPath, flagConfigPathName, "", "path to config file")
		flag.Parse()

		if path, ok := os.LookupEnv(envConfigPathName); ok {
			configPath = path
		}

		instance = &Config{}

		if readErr := cleanenv.ReadConfig(configPath, instance); readErr != nil {
			description, descrErr := cleanenv.GetDescription(instance, nil)
			if descrErr != nil {
				panic(descrErr)
			}

			slog.Info(description)
			slog.Error(
				"failed to read config",
				slog.String("err", readErr.Error()),
				slog.String("path", configPath),
			)
			os.Exit(1)
		}
	})

	return instance
}
