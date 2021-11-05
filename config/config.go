package config

type (
	// Config -.
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		PG       `yaml:"postgres"`
		Telegram `yaml:"telegram"`
		Vk       `yaml:"vk"`
		Grabber  `yaml:"grabber"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	// Telegram -.
	Telegram struct {
		BaseURL string `env-required:"true" yaml:"base_url" env:"TELEGRAM_BASE_URL"`
		Token   string `env-required:"true" env:"TELEGRAM_TOKEN"`
	}

	// Vk -.
	Vk struct {
		Token string `env-required:"true" env:"VK_TOKEN"`
	}

	// Grabber -.
	Grabber struct {
		Sleep int64 `env-required:"true" yaml:"sleep" env:"GRABBER_SLEEP"`
	}
)
