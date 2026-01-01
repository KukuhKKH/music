package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// app struct config
type app = struct {
	Name        string        `toml:"name"`
	Port        string        `toml:"port"`
	PrintRoutes bool          `toml:"print-routes"`
	Prefork     bool          `toml:"prefork"`
	Production  bool          `toml:"production"`
	BodyLimit   int           `toml:"body-limit"`
	IdleTimeout time.Duration `toml:"idle-timeout"`
	TLS         struct {
		Enable   bool   `toml:"enable"`
		CertFile string `toml:"cert-file"`
		KeyFile  string `toml:"key-file"`
	}
}

// db struct config
type db = struct {
	Mysql struct {
		DSN string `toml:"dsn"`
	}
}

// log struct config
type logger = struct {
	TimeFormat string        `toml:"time-format"`
	Level      zerolog.Level `toml:"level"`
	Prettier   bool          `toml:"prettier"`
}

// middleware
type middleware = struct {
	Compress struct {
		Enable bool
		Level  compress.Level
	}

	Recover struct {
		Enable bool
	}

	Monitor struct {
		Enable bool
		Path   string
	}

	Pprof struct {
		Enable bool
	}

	Limiter struct {
		Enable     bool
		Max        int
		Expiration time.Duration `toml:"expiration_seconds"`
	}

	FileSystem struct {
		Enable bool
		Browse bool
		MaxAge int `toml:"max_age"`
		Index  string
		Root   string
	}

	Jwt struct {
		Secret     string        `toml:"secret"`
		Issuer     string        `toml:"issuer"`
		Expiration time.Duration `toml:"expiration_seconds"`
	}

	Cors struct {
		Enable       bool   `toml:"enable"`
		AllowOrigins string `toml:"allow_origins"`
		AllowHeaders string `toml:"allow_headers"`
	}
}

type storage = struct {
	Driver string `toml:"driver"`

	Local struct {
		Path string `toml:"path"`
	} `toml:"local"`

	Ftp struct {
		Host      string `toml:"host"`
		Port      int    `toml:"port"`
		User      string `toml:"user"`
		Password  string `toml:"password"`
		BaseDir   string `toml:"base_dir"`
		PublicUrl string `toml:"public_url"`
	} `toml:"ftp"`

	S3 struct {
		Endpoint  string `toml:"endpoint"`
		AccessKey string `toml:"access_key"`
		SecretKey string `toml:"secret_key"`
		Bucket    string `toml:"bucket"`
		Region    string `toml:"region"`
		UseSsl    bool   `toml:"use_ssl"`
	} `toml:"s3"`
}

type Config struct {
	App        app
	DB         db
	Logger     logger
	Middleware middleware
	Storage    storage
}

// ParseConfig func to parse config
func ParseConfig(name string, debug ...bool) (contents *Config, err error) {
	var (
		file []byte
	)

	if len(debug) > 0 {
		file, err = os.ReadFile(name)
	} else {
		_, b, _, _ := runtime.Caller(0)
		// get base path
		path := filepath.Dir(filepath.Dir(filepath.Dir(b)))
		file, err = os.ReadFile(filepath.Join(path, "./config/", name+".toml"))
	}

	if err != nil {
		return &Config{}, err
	}

	err = toml.Unmarshal(file, &contents)

	return
}

// NewConfig initialize config
func NewConfig() *Config {
	config, err := ParseConfig("config")
	if err != nil && !fiber.IsChild() {
		// panic if config is not found
		log.Panic().Err(err).Msg("config not found")
	}

	return config
}

// ParseAddress func to parse address
func ParseAddress(raw string) (host, port string) {
	if i := strings.LastIndex(raw, ":"); i > 0 {
		return raw[:i], raw[i+1:]
	}

	return raw, ""
}
