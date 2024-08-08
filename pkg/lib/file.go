package lib

import (
	"io"
	"log"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}()
	return io.ReadAll(file)
}
