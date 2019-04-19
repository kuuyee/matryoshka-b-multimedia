package main

import (
	"github.com/kuuyee/matryoshka-b-multimedia/api"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
	"github.com/kuuyee/matryoshka-b-multimedia/router"
	"github.com/nfnt/resize"
)

func main() {
	storage, err := storage.NewDiskStorage("data/")
	if err != nil {
		panic(err)
	}
	api := api.NewAPI()
	api.RegisterServiceHandler("image", &handlers.ImageHandler{
		Storage:    storage,
		MaxSize:    16 << 20,
		ResizeAlgo: resize.Bicubic,
		KeyedMutex: handlers.NewKeyedRWMutex(),
	})
	eng := router.New(api)
	eng.Run(":8080")
}
