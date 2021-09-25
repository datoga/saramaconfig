package saramaconfig

import (
	"github.com/Shopify/sarama"
)

const (
	adminRetryMax     = "admin.retry.max"
	adminRetryBackoff = "admin.retry.backoff"
	adminTimeout      = "admin.timeout"
)

func New(config ConfigGetter) *sarama.Config {
	cfg := sarama.NewConfig()

	fillAdminSection(config, cfg)

	return cfg
}

func fillAdminSection(config ConfigGetter, cfg *sarama.Config) {
	if config.IsSet(adminRetryMax) {
		cfg.Admin.Retry.Max = config.GetInt(adminRetryMax)
	}

	if config.IsSet(adminRetryBackoff) {
		cfg.Admin.Retry.Backoff = config.GetDuration(adminRetryBackoff)
	}

	if config.IsSet(adminTimeout) {
		cfg.Admin.Timeout = config.GetDuration(adminTimeout)
	}
}
