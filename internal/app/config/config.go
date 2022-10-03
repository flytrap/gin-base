package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	C    = new(Config)
	once sync.Once
)

// Load config file (toml/json/yaml)
func MustLoad(path string) {
	once.Do(func() {
		viper.SetConfigFile(path)
		viper.ReadInConfig()
		viper.Unmarshal(&C)
	})
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	RunMode         string
	DefaultPassword string
	Swagger         bool
	PrintConfig     bool
	HTTP            HTTP
	LogGormHook     LogGormHook
	LogMongoHook    LogMongoHook
	JWTAuth         JWTAuth
	RateLimiter     RateLimiter
	CORS            CORS
	Redis           Redis
	Gorm            Gorm
	MySQL           MySQL
	Postgres        Postgres
	Sqlite3         Sqlite3
	MiniWx          MiniWx
	MpWx            MpWx
	AppWx           AppWx
}

func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

type LogHook string

func (h LogHook) IsGorm() bool {
	return h == "gorm"
}

type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

type LogMongoHook struct {
	Collection string
}

type JWTAuth struct {
	Enable        bool
	SigningMethod string
	SigningKey    string
	Expired       int
	Store         string
	FilePath      string
	RedisDB       int
	RedisPrefix   string
}

type HTTP struct {
	Host               string
	Port               int
	CertFile           string
	KeyFile            string
	ShutdownTimeout    int
	MaxContentLength   int64
	MaxReqLoggerLength int `default:"1024"`
	MaxResLoggerLength int
}

type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type Redis struct {
	Addr      string
	Password  string
	DB        int
	KeyPrefix string
}

type Gorm struct {
	Debug             bool
	DBType            string
	DbName            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

type Sqlite3 struct {
	Path string
}

func (a Sqlite3) DSN() string {
	return a.Path
}

type MiniWx struct {
	AppID     string
	AppSecret string
}

type MpWx struct {
	AppID     string
	AppSecret string
}

type AppWx struct {
	AppID     string
	AppSecret string
}
