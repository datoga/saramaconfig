package main

import (
	"os"

	"github.com/datoga/saramaconfig"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

func main() {
	os.Setenv("KAFKA_NET_DIALTIMEOUT", "1s")
	os.Setenv("KAFKA_PRODUCER_RETURN_SUCCESSES", "true")
	os.Setenv("KAFKA_PRODUCER_FLUSH_MAXMESSAGES", "100")
	os.Setenv("KAFKA_NET_SASL_USER", "test")
	os.Setenv("KAFKA_VERSION", "0.11.0.1")

	viper.SetEnvPrefix("kafka")
	viper.AutomaticEnv()

	cfg, err := saramaconfig.NewFromViper(viper.GetViper())

	if err != nil {
		panic(err)
	}

	spew.Dump(cfg)
}
