package saramaconfig_test

import (
	"encoding/base64"
	"testing"
	"time"

	_ "embed"

	"github.com/Shopify/sarama"
	"github.com/datoga/saramaconfig"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed test_certs/ca.crt
var ca []byte

//go:embed test_certs/client.pem
var clientPem []byte

//go:embed test_certs/client.key
var clientKey []byte

func TestSaramaConfigInt(t *testing.T) {
	valueInt := 100

	v := viper.New()

	v.Set("admin.retry.max", valueInt)

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, valueInt, config.Admin.Retry.Max)
}

func TestSaramaConfigDuration(t *testing.T) {
	valueDuration := "25ms"

	parsedDuration, err := time.ParseDuration(valueDuration)
	require.NoError(t, err)

	v := viper.New()
	v.Set("admin.timeout", valueDuration)

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, parsedDuration, config.Admin.Timeout)
}

func TestSaramaConfigBool(t *testing.T) {
	valueBool := true

	v := viper.New()
	v.Set("net.tls.enable", valueBool)

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, valueBool, config.Net.TLS.Enable)
}

func TestSaramaConfigString(t *testing.T) {
	valueString := "test"

	v := viper.New()
	v.Set("net.sasl.user", valueString)

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, valueString, config.Net.SASL.User)
}

func TestSaramaConfigSASLMechanism(t *testing.T) {
	valueSASLMechanism := "OAUTHBEARER"

	v := viper.New()
	v.Set("net.sasl.mechanism", valueSASLMechanism)

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, sarama.SASLMechanism(valueSASLMechanism), config.Net.SASL.Mechanism)
}

func TestSaramaConfigScramClientGenerator(t *testing.T) {
	v := viper.New()
	v.Set("net.sasl.scramclientgeneratorfunc", "SHA256")

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.NotNil(t, config.Net.SASL.SCRAMClientGeneratorFunc)
}

func TestSaramaConfigScramClientGeneratorNotValid(t *testing.T) {
	v := viper.New()
	v.Set("net.sasl.scramclientgeneratorfunc", "not-valid")

	_, err := saramaconfig.NewFromViper(v)

	assert.Error(t, err)
}

func TestSaramaConfigScramClientGeneratorNotSet(t *testing.T) {
	v := viper.New()
	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Nil(t, config.Net.SASL.SCRAMClientGeneratorFunc)
}

func TestSaramaConfigTLS(t *testing.T) {
	v := viper.New()

	v.Set("tls.ca", base64.StdEncoding.EncodeToString(ca))
	v.Set("tls.clientpem", base64.StdEncoding.EncodeToString(clientPem))
	v.Set("tls.clientkey", base64.StdEncoding.EncodeToString(clientKey))

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	require.NotNil(t, config.Net.TLS.Config)
	require.NotEmpty(t, config.Net.TLS.Config.Certificates)
	require.NotNil(t, config.Net.TLS.Config.RootCAs)
	require.NotEmpty(t, config.Net.TLS.Config.Certificates[0].Certificate)
}

func TestSaramaConfigVersion(t *testing.T) {
	v := viper.New()

	v.Set("version", sarama.V0_8_2_0.String())

	config, err := saramaconfig.NewFromViper(v)

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, config.Version, sarama.V0_8_2_0)
}

func TestSaramaConfigNonParseableVersion(t *testing.T) {
	v := viper.New()

	v.Set("version", "test")

	_, err := saramaconfig.NewFromViper(v)

	assert.Error(t, err)
}

func TestSaramaConfigWrongVersionType(t *testing.T) {
	v := viper.New()

	v.Set("version", 123)

	_, err := saramaconfig.NewFromViper(v)

	assert.Error(t, err)
}
