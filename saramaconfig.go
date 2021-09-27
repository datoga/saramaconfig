package saramaconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/datoga/saramaconfig/scramclient"
	"github.com/spf13/viper"
	"github.com/xdg/scram"
)

const (
	scramClientGenerator = "net.sasl.scramclientgeneratorfunc"
	rootTLS              = "tls"
	sha256               = "SHA256"
	sha512               = "SHA512"
)

func NewFromViper(v *viper.Viper) (*sarama.Config, error) {
	bindEnvs(v, *sarama.NewConfig())

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var hashGeneratorFn scram.HashGeneratorFcn

	if v.IsSet(scramClientGenerator) {
		fnGenerator := v.GetString(scramClientGenerator)

		switch fnGenerator {
		case sha256:
			hashGeneratorFn = scramclient.SHA256
		case sha512:
			hashGeneratorFn = scramclient.SHA512
		default:
			return nil, fmt.Errorf("unsupported scram generator function %s, only SHA256 and SHA512 values allowed", fnGenerator)

		}

		v.Set(scramClientGenerator, nil)
	}

	cfg := sarama.NewConfig()

	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	if hashGeneratorFn != nil {
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &scramclient.XDG{HashGeneratorFcn: hashGeneratorFn}
		}
	}

	if v.IsSet(rootTLS) {
		tlsConfig, err := tlsConfigFromViper(v)

		if err != nil {
			return nil, fmt.Errorf("failed configuring TLS with error %v", err)
		}

		cfg.Net.TLS.Config = tlsConfig
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validating sarama config with error %v", err)
	}

	return cfg, nil
}

func bindEnvs(v *viper.Viper, iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := strings.ToLower(t.Name)
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			name = tag
		}
		path := append(parts, name)
		switch fieldv.Kind() {
		case reflect.Struct:
			bindEnvs(v, fieldv.Interface(), path...)
		default:
			v.BindEnv(strings.Join(path, "."))
		}
	}
}
