[![Release](https://img.shields.io/github/release-pre/theflyingcodr/goconfig.svg?logo=github&style=flat&v=1)](https://github.com/theflyingcodr/goconfig/releases)
[![Build Status](https://img.shields.io/github/workflow/status/theflyingcodr/goconfig/run-go-tests?logo=github&v=3)](https://github.com/theflyingcodr/goconfig/actions)
[![Report](https://goreportcard.com/badge/github.com/theflyingcodr/goconfig?style=flat&v=1)](https://goreportcard.com/report/github.com/theflyingcodr/goconfig)
[![Go](https://img.shields.io/github/go-mod/go-version/theflyingcodr/goconfig?v=1)](https://golang.org/)
[![Mergify Status][mergify-status]][mergify]

[mergify]: https://mergify.io
[mergify-status]: https://img.shields.io/endpoint.svg?url=https://gh.mergify.io/badges/theflyingcodr/goconfig&style=flat

# goconfig

GoConfig provides strongly typed envrionment variable based configuration for go apps using a friendly Fluent API.

It comes with a viper implementation but others that implement the `ConfigurationLoader` interface can also be added.

Though it is billed as environment variable based, you can load configuration in however you wish, try to keep it 12 factor though with env vars.

## Useage

Useage is very simple, at the start of your main first load your defaults, I tend to have a method like `SetupDefaults` that contains all the
default settings for all environment variables in the application.
Directly after this, call goconfig as shown:

```go
func main(){
	config.SetupDefaults()
	cfg := goconfig.NewViperConfig("my-app").
		WithServer().
		WithDeployment().
		WithDb().
		Load()
}
```

With this simple syntax you can quickly glance at your main file and see the exact requirements for your application.

You can then reference configs directly rather than plucking them out of thin air throughput your code base.

For example, to get the server port do:

```go
 svr := cfg.Server.Port
```

It is recommended that you use dependency injection rather than globals to pass these around, as such, they are all pointers.

```go
 func NewImportantService(cfg *goconfig.Server){
	 cfg.Port
}
```

The above injection akes it 100% explicit as to the requirements of your service and is much easier tested than putting config readers throughout your code base.

### Http Clients

We also support setup of custom http clients, this can be done as shown:

```go
func main(){
	config.SetupDefaults()
	cfg := goconfig.NewViperConfig("my-app").
		WithServer().
		WithHTTPClient("my-service").
		WithHttpClient("my-other-service").
		Load()
}
```
Each of these gets its own config keyed off the service name as shown:

```go
	EnvHTTPClientHost       = "%s.client.host"
	EnvHTTPClientPort       = "%s.client.port"
	...
```
In the case of my-service, it's environment config for host would be `MY-SERVICE.CLIENT.HOST`.

You can add as many as you need and can access them by calling `svcCfg := cfg.CustomHTTPClient("my-service)`.

## Contributing

Contributions are more than welcome, there is a limited set of configs available at present and I'll be adding them as I need them, so if you think you'd
like to use this library but has something missing, like a Cassandra config, either raise an issue, create a PR adding it or both ideally.

