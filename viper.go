package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ViperConfig struct {
	*Config
}

func NewViperConfig(appname string) *ViperConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appname))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appname))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) {
			// Config file not found
			viper.Set("debug", false)
			viper.Set("port", 1323)
		} else {
			log.Fatalf("Fatal error config file: %s", err)
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return &ViperConfig{
		Config: &Config{
			httpClients: map[string]HTTPClientConfig{},
		},
	}
}

// WithServer will setup the web server configuration if required.
func (c *ViperConfig) WithServer() ConfigurationLoader {
	c.Server = &Server{
		Port:         viper.GetString(EnvServerPort),
		Hostname:     viper.GetString(EnvServerHost),
		TLSEnabled:   viper.GetBool(EnvServerTLSEnabled),
		TLSCertPath:  viper.GetString(EnvServerTLSCert),
		PProfEnabled: viper.GetBool(EnvServerPprofEnabled),
	}
	return c
}

// WithEnvironment sets up the deployment configuration if required.
func (c *ViperConfig) WithEnvironment(appName string) ConfigurationLoader {
	viper.SetDefault(EnvEnvironment, "dev")
	viper.SetDefault(EnvRegion, "test")
	viper.SetDefault(EnvCommit, "test")
	viper.SetDefault(EnvVersion, "test")
	viper.SetDefault(EnvBuildDate, time.Now().UTC())

	c.Deployment = &Deployment{
		Environment: viper.GetString(EnvEnvironment),
		Region:      viper.GetString(EnvRegion),
		Version:     viper.GetString(EnvVersion),
		Commit:      viper.GetString(EnvCommit),
		BuildDate:   viper.GetTime(EnvBuildDate),
		AppName:     appName,
	}
	return c
}

func (c *ViperConfig) WithLog() ConfigurationLoader {
	c.Logging = &Logging{Level: viper.GetString(EnvLogLevel)}
	return c
}

// WithDb sets up and returns database configuration.
func (c *ViperConfig) WithDb() ConfigurationLoader {
	c.Db = &Db{
		Type:       DbType(viper.GetString(EnvDb)),
		Dsn:        viper.GetString(EnvDbDsn),
		SchemaPath: viper.GetString(EnvDbSchema),
		Migrate:    viper.GetBool(EnvDbMigrate),
	}
	return c
}

// WithRedis will include redis config.
func (c *ViperConfig) WithRedis() ConfigurationLoader {
	if !viper.IsSet(EnvRedisDb) {
		viper.SetDefault(EnvRedisDb, 0)
	}
	c.Redis = &Redis{
		Address:  viper.GetString(EnvRedisAddress),
		Password: viper.GetString(EnvRedisPassword),
		Db:       viper.GetUint(EnvRedisDb),
	}
	return c
}

// WithHttpClient will setup a custom http client referenced by name.
func (c *ViperConfig) WithHTTPClient(name string) ConfigurationLoader {
	c.httpClients[name] = HTTPClientConfig{
		Host:       viper.GetString(fmt.Sprintf(EnvHTTPClientHost, name)),
		Port:       viper.GetString(fmt.Sprintf(EnvHTTPClientPort, name)),
		TLSEnabled: viper.GetBool(fmt.Sprintf(EnvHTTPClientTLSEnabled, name)),
		TLSCert:    viper.GetBool(fmt.Sprintf(EnvHTTPClientTLSCert, name)),
		Timeout:    time.Second * time.Duration(viper.GetInt(fmt.Sprintf(EnvHTTPClientTimeout, name))),
	}
	return c
}

// WithSwagger will setup and return swagger configuration.
func (c *ViperConfig) WithSwagger() ConfigurationLoader {
	c.Swagger = &Swagger{
		Host:    viper.GetString(EnvSwaggerHost),
		Enabled: viper.GetBool(EnvSwaggerEnabled),
	}
	return c
}

// WithInstrumentation will read instrumentation envrionment vars.
func (c *ViperConfig) WithInstrumentation() ConfigurationLoader {
	c.Instrumentation = &Instrumentation{
		MetricsEnabled: viper.GetBool(EnvMetricsEnabled),
		TracingEnabled: viper.GetBool(EnvTracingEnabled),
	}
	return c
}

// Load will finish setup and return configuration. This should
// always be the last call.
func (c *ViperConfig) Load() *Config {
	return c.Config
}
