package lib

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

func GetRandomPort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, errors.Wrap(err, "failed to listen :0")
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func RandomLocalConnect() (string, error) {
	port, err := GetRandomPort()
	if err != nil {
		return "", errors.Wrap(err, "failed to get random port")
	}
	return fmt.Sprintf("127.0.0.1:%d", port), nil
}
