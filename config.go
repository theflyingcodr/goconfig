package config

import (
	"fmt"
	"regexp"
	"time"

	validator "github.com/theflyingcodr/govalidator"
)

const (
	EnvServerPort       = "server.port"
	EnvServerHost       = "server.host"
	EnvServerTLSEnabled = "server.tls.enabled"
	EnvServerTLSCert    = "server.tls.cert"

	EnvEnvironment = "env.environment"
	EnvMainNet     = "env.mainnet"
	EnvRegion      = "env.region"
	EnvVersion     = "env.version"
	EnvCommit      = "env.commit"
	EnvBuildDate   = "env.builddate"

	EnvLogLevel = "log.level"

	EnvDb        = "db.type"
	EnvDbSchema  = "db.schema.path"
	EnvDbDsn     = "db.dsn"
	EnvDbMigrate = "db.migrate"

	EnvHttpClientHost       = "%s.client.host"
	EnvHttpClientPort       = "%s.client.port"
	EnvHttpClientTimeout    = "%s.client.timeout"
	EnvHttpClientTLSEnabled = "%s.client.tls.enabled"
	EnvHttpClientTLSCert    = "%s.client.tls.cert"

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
	MainNet     bool
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
	Type       DbType
	SchemaPath string
	Dsn        string
	Migrate    bool
}

// Validate will ensure the HeaderClient config is valid.
func (d *Db) Validate(v validator.ErrValidation) {
	v = v.Validate("db.type", validator.MatchString(string(d.Type), reDbType))
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

// ConfigurationLoader will load configuration items
// into a struct that contains a configuration.
type ConfigurationLoader interface {
	WithServer() *Config
	WithEnvironment() *Config
	WithLog() *Config
	WithHttpClient(name string) *Config
	WithDb() *Config
	WithDeployment(app string) *Config
}
