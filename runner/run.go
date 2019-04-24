package runner

import (
	"log"

	"github.com/kuuyee/matryoshka-b-multimedia/api"
	"github.com/kuuyee/matryoshka-b-multimedia/conf"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
	"github.com/kuuyee/matryoshka-b-multimedia/router"
	"github.com/nfnt/resize"
)

// Run starts the server
func Run() {
	serverConf := conf.GetParsed()

	storageH, err := storage.LoadStorage()
	if err != nil {
		log.Panicf("error while loading storage handler: %v", err)
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
