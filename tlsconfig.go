package saramaconfig

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/spf13/viper"
)

const (
	keyRootTLS = "tls"
)

type RootTLS struct {
	TLS TLSConfig
}

type TLSConfig struct {
	CA        string
	ClientPem string
	ClientKey string
}

func tlsConfigFromViper(v *viper.Viper) (*tls.Config, error) {
	var tlsCfg RootTLS

	err := v.Unmarshal(&tlsCfg)

	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	cfg, err := tlsConfigFromEncodedCerts(tlsCfg.TLS.CA, tlsCfg.TLS.ClientPem, tlsCfg.TLS.ClientKey)

	if err != nil {
		return nil, fmt.Errorf("failed decoding certs with error %v", err)
	}

	return cfg, nil
}

func tlsConfigFromEncodedCerts(ca string, clientCert string, clientKey string) (*tls.Config, error) {
	decodedCA, err := base64.StdEncoding.DecodeString(ca)

	if err != nil {
		return nil, fmt.Errorf("failed decoding tls CA with error %v", err)
	}

	decodedClientPem, err := base64.StdEncoding.DecodeString(clientCert)

	if err != nil {
		return nil, fmt.Errorf("failed decoding tls Client Pem with error %v", err)
	}

	decodedClientKey, err := base64.StdEncoding.DecodeString(clientKey)

	if err != nil {
		return nil, fmt.Errorf("failed decoding tls Client Key with error %v", err)
	}

	return tlsConfigFromCerts(string(decodedCA), string(decodedClientPem), string(decodedClientKey))
}

func tlsConfigFromCerts(ca string, clientCert string, clientKey string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))

	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(ca))

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}
