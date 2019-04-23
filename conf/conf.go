package conf

import (
	"errors"
	"io/ioutil"

	"github.com/nfnt/resize"

	"github.com/go-yaml/yaml"
)

var conf *C

// C 服务器配置定义
type C struct {
	API struct {
		Listen string `yaml:"listen"`
	} `yaml:"api"`
	Storage struct {
		Mode string `yaml:"mode"`
		Path string `yaml:"path"`
	} `yaml:"storage"`
	Handlers struct {
		Image struct {
			MaxSize int64      `yaml:"maxsize"`
			Resize  ResizeFunc `yaml:"resize"`
		} `yaml:"image"`
	} `yaml:"handlers"`
}

// ResizeFunc is an alias to resize.InterpolationFunction, used to unmarshal interpolation function from config file
type ResizeFunc resize.InterpolationFunction

// UnmarshalYAML implements yaml.Unmarshaler
func (r *ResizeFunc) UnmarshalYAML(origUnmarshal func(interface{}) error) error {
	s := ""
	if err := origUnmarshal(&s); err != nil {
		return err
	}
	switch s {
	case "", "bicubic":
		*r = ResizeFunc(resize.Bicubic)
	case "bilinear":
		*r = ResizeFunc(resize.Bilinear)
	case "lanczos2":
		*r = ResizeFunc(resize.Lanczos2)
	case "lanxzos3":
		*r = ResizeFunc(resize.Lanczos3)
	case "nearestneighbor":
		*r = ResizeFunc(resize.NearestNeighbor)
	default:
		return errors.New("unrecognized resize interpolation function")
	}
	return nil
}

// Set sets current config to the provided value
func Set(newConf C) {
	conf = &newConf
}

// GetParsed returns the config file after the initial parse
func GetParsed() C {
	return *conf
}

// Get loads the config file and make configuration available for GetParsed()
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
