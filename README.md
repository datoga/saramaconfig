# saramaconfig
[Sarama](https://github.com/Shopify/sarama/) config parser from [Viper](https://github.com/spf13/viper). It means it supports out of the box all the formats provided by Viper (env, yaml, toml, etc.) and most settings in Sarama.

Sometimes is difficult to parse all the [config settings from Sarama](https://github.com/Shopify/sarama/blob/main/config.go) from env variables or configuration files. Following some guidelines, this package will help you to get up & running Sarama configuration from your favourite config system.

## Rules:

- The keys should be the same than the [Sarama config file](https://github.com/Shopify/sarama/blob/main/config.go). Just transform they fields into a nested structure with lowercase names separated by a dot per level. Example: `Version` would be just `version` and `Admin -> Retry -> Max` would be `admin.retry.max`.

- Strings, ints, boolean just works! Durations can be specified with parsed durations like `1ms` or `3s` (std `time.ParseDuration` will be used under the hood).

- Arrays can be specified joining the values with commas. Per example `1,2,5` would be a valid int array.

- In general, functions (i.e. retry backoff) and interfaces (i.e. proxy dialer),will not be parsed. You can setup any of them with the Sarama config got from the package.

- There are a few _magic_ exceptions from the rule above.

  - `version` will be parsed as string (not the decomposed [4]int) following the Sarama rules (`sarama.ParseKafkaVersion` function will be used).

  - `net.sasl.scramclientgeneratorfunc` will be parsed as string (with the possible values _SHA256_ and _SHA512_) implementing a generic XDG Scram client using the hash generator you specified.

  - `tls.ca`, `tls.clientpem` and `tls.clientkey` are special keys used to provide certificates and the client key to configure TLS. These certificates should be encoded in base64 encoding. Don't forget to enable TLS using the regular convention: `net.tls.enable = true`. Check [the example](examples/tlsenv) to see how to configure TLS.

## Usage:

1. Choose your config system and write (or fill) the configuration settings, following the previous rules.

2. Setup the viper configuration, and make sure is loaded with `viper.AutomaticEnv()` (env) o `viper.ReadInConfig()` (files). For example, for a `conf.toml` would be:

```Go
viper.AddConfigPath(".")
viper.SetConfigName("conf")
viper.SetConfigType("toml")

if err := viper.ReadInConfig(); err != nil {
	panic(err)
}
```

3. Call the package with the viper you have just configured and check the error.

```Go
cfg, err := saramaconfig.NewFromViper(viper.GetViper())

if err != nil {
	panic(err)
}
```

4. Just use `cfg` as a regular Sarama config to run producers, consumers, etc.

## Examples

There are some [examples](examples) to demonstrate how to work with the package. The tests also can describe some of the behaviour of the package, check them!

## TODO

- Go doc.
- Now, if you provide a wrong key, that would not be parsed but it won't throw any error. It would be nice to report an error and maybe to advise the right key.
- Implement more _magic_ non parseable keys, like functions or interfaces. For example, `Net -> Proxy -> Dialer` which is a `proxy.Dialer`, and could be set up implementing a struct with some proxy details provided in the config.