package multimedia

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/chai2010/webp"

	"github.com/kuuyee/matryoshka-b-multimedia/test"

	"github.com/nfnt/resize"

	"github.com/kuuyee/matryoshka-b-multimedia/conf"
	"github.com/kuuyee/matryoshka-b-multimedia/runner"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultimediaClient(t *testing.T) {
	var serverConf conf.C
	serverConf.API.Listen = "127.0.0.1:28574"
	serverConf.Handlers.Image.MaxSize = 100 << 20
	serverConf.Handlers.Image.Resize = conf.ResizeFunc(resize.Bicubic)
	serverConf.Storage.Mode = "disk"
	serverConf.Storage.Path, _ = ioutil.TempDir("", "mm_test")
	conf.Set(serverConf)

	go runner.Run()
	time.Sleep(500 * time.Millisecond) // wait for service to start

	client := &Client{
		Client: &http.Client{
			Timeout: 1 * time.Second,
		},
		BaseURL: "http://127.0.0.1:28574/rest/media",
	}
	Convey("multimedia client", t, func() {
		_, err := client.FileMeta(Image, "not_exist")
		So(err, ShouldHaveSameTypeAs, HTTPCodeError{})

		testFile, err := os.Open(path.Join(test.GetProjectDir(), "test/testasset/gomk.png"))
		So(err, ShouldBeNil)
		defer testFile.Close()

		fileMeta, err := client.PostFile(Image, "mk.png", testFile, nil)
		So(err, ShouldBeNil)
		So(fileMeta.Type, ShouldEqual, string(Image))
		So(fileMeta.Ident, ShouldNotBeEmpty)

		gotIdent := fileMeta.Ident

		retrievedMeta, err := client.FileMeta(Image, gotIdent)
		So(err, ShouldBeNil)
		So(retrievedMeta, ShouldResemble, fileMeta)

		retrievedFile, err := client.FetchFile(Image, gotIdent, map[string]string{
			"size":   "200",
			"format": "webp",
		})
		So(err, ShouldBeNil)

		retreivedImage, err := webp.Decode(retrievedFile)
		So(err, ShouldBeNil)
		So(retreivedImage.Bounds().Dy(), ShouldEqual, 200)
	})
}
