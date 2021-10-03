package config

import (
	"fmt"
	"regexp"
	"time"

	validator "github.com/theflyingcodr/govalidator"
)

const (
	EnvServerPort           = "server.port"
	EnvServerHost           = "server.host"
	EnvServerTLSEnabled     = "server.tls.enabled"
	EnvServerTLSCert        = "server.tls.cert"
	EnvServerSwaggerEnabled = "server.swagger.enabled"

	EnvEnvironment = "env.environment"
	EnvRegion      = "env.region"
	EnvVersion     = "env.version"
	EnvCommit      = "env.commit"
	EnvBuildDate   = "env.builddate"

	EnvLogLevel = "log.level"

	EnvDb        = "db.type"
	EnvDbSchema  = "db.schema.path"
	EnvDbDsn     = "db.dsn"
	EnvDbMigrate = "db.migrate"

	EnvHTTPClientHost       = "%s.client.host"
	EnvHTTPClientPort       = "%s.client.port"
	EnvHTTPClientTimeout    = "%s.client.timeout"
	EnvHTTPClientTLSEnabled = "%s.client.tls.enabled"
	EnvHTTPClientTLSCert    = "%s.client.tls.cert"

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
	Logging     *Logging
	Server      *Server
	Deployment  *Deployment
	Db          *Db
	Redis       *Redis
	httpClients map[string]HTTPClientConfig
}

// HTTPClientConfig is a custom http client config struct, returned
// when CustomHTTPClient is called.
type HTTPClientConfig struct {
	Host       string
	Port       string
	TLSEnabled bool
	TLSCert    bool
	Timeout    time.Duration
}

// CustomHTTPClient will return a custom http client, if not found
// nil is returned.
func (c *Config) CustomHTTPClient(name string) *HTTPClientConfig {
	cfg, ok := c.httpClients[name]
	if !ok {
		return nil
	}
	return &cfg
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
	Port           string
	Hostname       string
	TLSCertPath    string
	TLSEnabled     bool
	SwaggerEnabled bool
}

// Db contains database information.
type Db struct {
	Type       DbType
	SchemaPath string
	Dsn        string
	Migrate    bool
}

// Validate will ensure the HeaderClient config is valid.
func (d *Db) Validate(v validator.ErrValidation) validator.ErrValidation {
	return v.Validate("db.type", validator.MatchString(string(d.Type), reDbType))
}

var reDbType = regexp.MustCompile(`sqlite|mysql|postgres`)

// DbType is used to restrict the dbs we can support.
type DbType string

// Supported database types.
const (
	DBSqlite   DbType = "sqlite"
	DBMySQL    DbType = "mysql"
	DBPostgres DbType = "postgres"
)

type Redis struct {
	Address  string
	Password string
	Db       uint
}

// ConfigurationLoader will load configuration items
// into a struct that contains a configuration.
type ConfigurationLoader interface {
	WithServer() ConfigurationLoader
	WithEnvironment(appname string) ConfigurationLoader
	WithLog() ConfigurationLoader
	WithHTTPClient(name string) ConfigurationLoader
	WithDb() ConfigurationLoader
	WithRedis() ConfigurationLoader
	Load() *Config
}
