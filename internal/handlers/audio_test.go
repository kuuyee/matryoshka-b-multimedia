package handlers

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/kuuyee/matryoshka-b-multimedia/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAudioHandler(t *testing.T) {
	Convey("audio handler", t, func() {
		noiseFilePath := path.Join(test.GetProjectDir(), "test/testasset/noise.ogg")
		h := H(&AudioHandler{
			Storage:    test.NewTmpStorage("audio_handler"),
			KeyedMutex: NewKeyedRWMutex(),
		})

		noiseFileData, err := ioutil.ReadFile(noiseFilePath)
		So(err, ShouldBeNil)

		ident, err := h.WriteData(bytes.NewReader(noiseFileData), "audio/ogg", nil)
		So(err, ShouldBeNil)

		retrievedNoiseFile, _, format, err := h.RetrieveData(ident, nil)
		So(err, ShouldBeNil)
		So(format, ShouldEqual, "audio/ogg")
		defer retrievedNoiseFile.Close()
		retrievedNoiseData, err := ioutil.ReadAll(retrievedNoiseFile)
		So(err, ShouldBeNil)
		So(retrievedNoiseData, ShouldResemble, noiseFileData)

	})
}
