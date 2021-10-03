package saramaconfig

import (
	"crypto/tls"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

//NewFromViper accepts a *viper.Viper config and parses the content keys into a *sarama.Config struct. It will return an error if the viper cannot be parsed or the sarama configuration does not validate.
func NewFromViper(v *viper.Viper) (*sarama.Config, error) {
	return newSaramaConfigFromViper(v)
}

//SaramaOpt is the signature for functional options for common.
type SaramaOpt func(cfg *sarama.Config)

//ProducerNoRetries is an option to not have retries on publishing error.
func ProducerNoRetries() SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Producer.Retry.Max = 0
	}
}

//ProducerMaxRetries allows to set a number of retries in case of publishing errors.
func ProducerMaxRetries(retries int) SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Producer.Retry.Max = retries
	}
}

//ConsumerBatch is an option to enable batching on consuming data.
func ConsumerBatch(minFetchBytes int, maxWaitTime time.Duration) SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Consumer.Fetch.Min = int32(minFetchBytes)
		cfg.Consumer.MaxWaitTime = maxWaitTime
	}
}

//ConsumerCommitAsync is an option to enable async commits.
func ConsumerCommitAsync(interval time.Duration) SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Consumer.Offsets.AutoCommit.Enable = true
		cfg.Consumer.Offsets.AutoCommit.Interval = interval
	}
}

//Timeout sets the default timeout.
func Timeout(timeout time.Duration) SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Net.DialTimeout = timeout
		cfg.Net.ReadTimeout = timeout
		cfg.Admin.Timeout = timeout
	}
}

//SASL configures the SASL handshake.
func SASL(saslMechanism sarama.SASLMechanism, username string, password string) SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Mechanism = saslMechanism
		cfg.Net.SASL.User = username
		cfg.Net.SASL.Password = password
	}
}

//TLS configures the secure connection.
func TLS(config *tls.Config) SaramaOpt {
	return func(cfg *sarama.Config) {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = config
	}
}
