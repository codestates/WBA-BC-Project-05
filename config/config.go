package Config

import (
	"fmt"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Web struct {
		Port string
	}
	Blockchain struct {
		UrlHttp      string
		UrlWs        string
		ContractAddr string
		OwnerAddr    string
		PrivateKey   string
	}
	DB struct {
		Host string
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
			fmt.Println(c)
			return c
		}
	}
}
