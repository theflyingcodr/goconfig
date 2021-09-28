package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func NewViperConfig(appname string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appname))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appname))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			viper.Set("debug", false)
			viper.Set("port", 1323)
		} else {
			log.Fatalf("Fatal error config file: %s", err)
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return &Config{}
}

// WithServer will setup the web server configuration if required.
func (c *Config) WithServer() *Config {
	c.Server = &Server{
		Port:     viper.GetString(EnvServerPort),
		Hostname: viper.GetString(EnvServerHost),
	}
	return c
}

// WithDeployment sets up the deployment configuration if required.
func (c *Config) WithDeployment(appName string) *Config {
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

func (c *Config) WithLog() *Config {
	c.Logging = &Logging{Level: viper.GetString(EnvLogLevel)}
	return c
}

// WithDb sets up and returns database configuration.
func (c *Config) WithDb() *Config {
	c.Db = &Db{
		Type:       viper.GetString(EnvDb),
		Dsn:        viper.GetString(EnvDbDsn),
		SchemaPath: viper.GetString(EnvDbSchema),
	}
	return c
}

// WithRedis will include redis config.
func (c *Config) WithRedis() *Config {
	if !viper.IsSet(EnvRedisDb) {
		viper.SetDefault(EnvRedisDb, 0)
	}
	c.Redis = &Redis{
		Address:  viper.GetString(EnvRedisAddress),
		Password: viper.GetString(EnvRedisPassword),
		Db:       viper.GetUint(EnvRedisDb),
	}
}
