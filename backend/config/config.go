package Config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Web struct {
		Port string
	}
	Blockchain struct {
		RpcUrl       string
		ContractAddr string
		OwnerAddr    string
		PrivateKey   string
	}
}

func NewConfig(fpath string) *Config {
	c := new(Config)

	if file, err := os.Open(fpath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			return c
		}
	}
}
