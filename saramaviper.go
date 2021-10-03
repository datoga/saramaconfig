package saramaconfig

import (
	"crypto/tls"
	"fmt"
	"reflect"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/datoga/saramaconfig/scramclient"
	"github.com/spf13/viper"
	"github.com/xdg/scram"
)

const (
	prefix                  = "kafka"
	keyScramClientGenerator = "net.sasl.scramclientgeneratorfunc"
	keyVersion              = "version"
	sha256                  = "SHA256"
	sha512                  = "SHA512"
)

func newSaramaConfigFromViper(v *viper.Viper) (*sarama.Config, error) {
	bindEnvs(v, *sarama.NewConfig())
	bindEnvs(v, RootTLS{})

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	hashGeneratorFn, err := parseHashGeneratorFunc(v)

	if err != nil {
		return nil, fmt.Errorf("failed parsing scram generator func with error %v", err)
	}

	err = parseVersion(v)

	if err != nil {
		return nil, fmt.Errorf("failed parsing version with error %v", err)
	}

	cfg := sarama.NewConfig()

	err = v.Unmarshal(&cfg)

	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	if hashGeneratorFn != nil {
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &scramclient.XDG{HashGeneratorFcn: hashGeneratorFn}
		}
	}

	tlsConfig, err := parseTLS(v)

	if err != nil {
		return nil, fmt.Errorf("unable to decode tls with error, %v", err)
	}

	if tlsConfig != nil {
		cfg.Net.TLS.Config = tlsConfig
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validating sarama config with error %v", err)
	}

	return cfg, nil
}

func parseHashGeneratorFunc(v *viper.Viper) (scram.HashGeneratorFcn, error) {
	if !v.IsSet(keyScramClientGenerator) {
		return nil, nil
	}

	var hashGeneratorFn scram.HashGeneratorFcn

	fnGenerator := v.GetString(keyScramClientGenerator)

	switch fnGenerator {
	case sha256:
		hashGeneratorFn = scramclient.SHA256
	case sha512:
		hashGeneratorFn = scramclient.SHA512
	default:
		return nil, fmt.Errorf("unsupported scram generator function %s, only SHA256 and SHA512 values allowed", fnGenerator)

	}

	v.Set(keyScramClientGenerator, nil)

	return hashGeneratorFn, nil
}

func parseVersion(v *viper.Viper) error {
	if !v.IsSet(keyVersion) {
		return nil
	}

	version := v.GetString(keyVersion)

	kafkaVersion, err := sarama.ParseKafkaVersion(version)

	if err != nil {
		return fmt.Errorf("failed parsing version %s with error %v", version, err)
	}

	v.Set(keyVersion, kafkaVersion)

	return nil
}

func parseTLS(v *viper.Viper) (*tls.Config, error) {
	if _, found := v.AllSettings()[keyRootTLS]; !found {
		return nil, nil
	}

	tlsConfig, err := tlsConfigFromViper(v)

	if err != nil {
		return nil, fmt.Errorf("failed configuring TLS with error %v", err)
	}

	return tlsConfig, nil
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
