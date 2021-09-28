package config

import (
	"fmt"
	"regexp"
	"time"

	validator "github.com/theflyingcodr/govalidator"
)

const (
	EnvServerPort = "server.port"
	EnvServerHost = "server.host"

	EnvEnvironment = "env.environment"
	EnvRegion      = "env.region"
	EnvVersion     = "env.version"
	EnvCommit      = "env.commit"
	EnvBuildDate   = "env.builddate"
	EnvLogLevel    = "log.level"

	EnvDb        = "db.type"
	EnvDbSchema  = "db.schema.path"
	EnvDbDsn     = "db.dsn"
	EnvDbMigrate = "db.migrate"

	EnvRedisAddress  = "redis.address"
	EnvRedisPassword = "redis.password"
	EnvRedisDb       = "redis.db"

	LogDebug = "debug"
	LogInfo  = "info"
	LogError = "error"
	LogWarn  = "warn"
)

// Config returns strongly typed config values.
type Config struct {
	defaultsFn func(c Config) error
	Logging    *Logging
	Server     *Server
	Deployment *Deployment
	Db         *Db
	Redis      *Redis
	custom     map[string]interface{}
}

// Validate will check config values are valid and return a list of failures
// if any have been found.
func (c *Config) Validate() error {
	vl := validator.New()
	if c.Db != nil {
		vl = vl.Validate("db.type", validator.MatchString(string(c.Db.Type), reDbType))
	}
	return vl.Err()
}

// Deployment contains information relating to the current
// deployed instance.
type Deployment struct {
	Environment string
	AppName     string
	Region      string
	Version     string
	Commit      string
	BuildDate   time.Time
}

// IsDev determines if this app is running on a dev environment.
func (d *Deployment) IsDev() bool {
	return d.Environment == "dev"
}

func (d *Deployment) String() string {
	return fmt.Sprintf("Environment: %s \n AppName: %s\n Region: %s\n Version: %s\n Commit:%s\n BuildDate: %s\n",
		d.Environment, d.AppName, d.Region, d.Version, d.Commit, d.BuildDate)
}

type Logging struct {
	Level string
}

// Server contains all settings required to run a web server.
type Server struct {
	Port     string
	Hostname string
}

// Db contains database information.
type Db struct {
	Type       string
	SchemaPath string
	Dsn        string
	Migrate    bool
}

var reDbType = regexp.MustCompile(`sqlite|mysql|postgres`)

// DbType is used to restrict the dbs we can support.
type DbType string

// Supported database types.
const (
	DBSqlite   DbType = "sqlite"
	DBMySql    DbType = "mysql"
	DBPostgres DbType = "postgres"
)

type Redis struct {
	Address  string
	Password string
	Db       uint
}

func (c *Config) AddCustomConfig(name string, conf interface{}) {
	c.custom[name] = conf
}

func (c *Config) CustomConfig(name string, out interface{}) {
	out = c.custom[name]
}

// ConfigurationLoader will load configuration items
// into a struct that contains a configuration.
type ConfigurationLoader interface {
	WithDefaults(func() error) *Config
	WithCustomConfig(func(c *Config) error) *Config
	WithServer() *Config
	WithDb() *Config
	WithDeployment(app string) *Config
	WithRedis() *Config
}

func LetsSee() {
	var r *Redis
	Test("me", &r)
}
func Test(name string, out interface{}) {
	t := map[string]interface{}{}
	out = t[name]
}
