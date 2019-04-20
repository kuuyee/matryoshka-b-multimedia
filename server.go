package main

import (
	"flag"

	"github.com/nfnt/resize"

	"github.com/kuuyee/matryoshka-b-multimedia/api"
	"github.com/kuuyee/matryoshka-b-multimedia/conf"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
	"github.com/kuuyee/matryoshka-b-multimedia/router"
)

func main() {
	confPath := flag.String("c", "multimedia.yml", "config file")
	flag.Parse()

	serverConf := conf.Get(*confPath)

	var storageH storage.S
	switch serverConf.Storage.Mode {
	case "disk":
		var err error
		storageH, err = storage.NewDiskStorage(serverConf.Storage.Path)
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown storage type")
	}

	api := api.NewAPI()
	api.RegisterServiceHandler("image", &handlers.ImageHandler{
		Storage:    storageH,
		MaxSize:    serverConf.Handlers.Image.MaxSize,
		ResizeAlgo: resize.InterpolationFunction(serverConf.Handlers.Image.Resize),
		KeyedMutex: handlers.NewKeyedRWMutex(),
	})
	eng := router.New(api)
	eng.Run(serverConf.API.Listen)
}
