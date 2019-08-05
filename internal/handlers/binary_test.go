package handlers

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/kuuyee/matryoshka-b-multimedia/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBinaryHandler(t *testing.T) {
	Convey("binary handler", t, func() {
		binFilePath := path.Join(test.GetProjectDir(), "test/testasset/binary.bin")
		h := H(&BinaryHandler{
			Storage: test.NewTmpStorage("binary_handler"),
		})

		binFileData, err := ioutil.ReadFile(binFilePath)
		So(err, ShouldBeNil)

		ident, err := h.WriteData(bytes.NewReader(binFileData), "application/octet-stream", nil)
		So(err, ShouldBeNil)

		retrievedBinFile, _, format, err := h.RetrieveData(ident, nil)
		So(err, ShouldBeNil)
		So(format, ShouldEqual, "application/octet-stream")
		defer retrievedBinFile.Close()
		retrievedBinData, err := ioutil.ReadAll(retrievedBinFile)
		So(err, ShouldBeNil)
		So(retrievedBinData, ShouldResemble, binFileData)

	})
}
