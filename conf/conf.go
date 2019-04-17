package conf

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

var conf *C

type C struct {
	Storage struct {
		Mode string `yaml:"mode"`
	} `yaml:"storage"`
}

func Get(fn string) C {
	if conf != nil {
		return *conf
	}
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	cData := new(C)
	if err := yaml.Unmarshal(d, cData); err != nil {
		panic(err)
	}
	conf = cData
	return *cData
}
