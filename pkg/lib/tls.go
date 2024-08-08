package lib

import (
	"crypto/tls"

	"github.com/pkg/errors"
)

func LoadTlsCreds(cert, key []byte) (*tls.Config, error) {
	creds, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load server cert and key")
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{creds},
		InsecureSkipVerify: true,
	}, nil
}

func LoadTlsCredsFromFile(certpath, keypath string) (*tls.Config, error) {
	cert, err := ReadFile(certpath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read cert file by path %s", certpath)
	}
	key, err := ReadFile(keypath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read key file by path %s", keypath)
	}
	return LoadTlsCreds(cert, key)
}
