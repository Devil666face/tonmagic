package main

import (
	_ "embed"
	"log"
	"tonmagic/internal/config"
	"tonmagic/internal/tonmagic"
	"tonmagic/pkg/lib"
)

//go:embed server.crt
var TLSServerCert []byte

//go:embed server.key
var TLSServerKey []byte

func main() {
	config := config.Must()
	tlsconfig, err := lib.LoadTlsCredsFromFile(config.CertFile, config.KeyFile)
	if err != nil {
		log.Printf("tls cert and key not load, use default in binary: %v\n", err)
		tlsconfig, err = lib.LoadTlsCreds(TLSServerCert, TLSServerKey)
	}
	tonmagic, err := tonmagic.New(config, tlsconfig)
	if err != nil {
		log.Fatalf("critical error: %v\n", err)
	}
	if err := tonmagic.StartTonProxy(); err != nil {
		log.Fatalf("error to start ton proxy: %v\n", err)
	}
	go func() {
		if err := tonmagic.ListenHttp(); err != nil {
			log.Fatalln(err)
		}
	}()
	if err := tonmagic.ListenHttps(); err != nil {
		log.Fatalln(err)
	}
}
