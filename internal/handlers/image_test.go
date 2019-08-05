package handlers

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"path"
	"testing"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"

	"github.com/kuuyee/matryoshka-b-multimedia/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestImageHandler(t *testing.T) {
	Convey("test image handler", t, func() {
		Convey("dimension decode", func() {
			shouldHaveDimension := func(actual interface{}, expected ...interface{}) string {
				x, y, err := imgDimensionDecode(actual.(string))
				if err := ShouldBeNil(err); err != "" {
					return err
				}
				xExp, yExp := expected[0].(uint), expected[1].(uint)
				if x != xExp || y != yExp {
					return fmt.Sprintf("expected %dx%d, got %dx%d", xExp, yExp, x, y)
				}
				return ""
			}
			So("100x200", shouldHaveDimension, uint(100), uint(200))
			So("100", shouldHaveDimension, uint(0), uint(100))
		})
		Convey("parse img opt", func() {
			opt, err := parseImgRetrieveParamOpt(map[string]string{
				"format": "webp",
				"size":   "100x100",
			})
			So(err, ShouldBeNil)
			So(*opt, ShouldResemble, imageRetrieveParamOpt{
				X:      100,
				Y:      100,
				Format: "image/webp",
			})
		})

		Convey("functional test", func() {
			storage := test.NewTmpStorage("mk_img_handler")
			maxSize := int64(50 << 20)
			h := H(&ImageHandler{
				Storage:    storage,
				MaxSize:    maxSize,
				ResizeAlgo: resize.Bicubic,
				KeyedMutex: NewKeyedRWMutex(),
			})
			So(h.Type(), ShouldEqual, "image")
			So(h.SizeLimit(), ShouldEqual, maxSize)
			Convey("storage image", func() {
				testImagePath := path.Join(test.GetProjectDir(), "test/testasset/gomk.png")
				testImageData, err := ioutil.ReadFile(testImagePath)
				So(err, ShouldBeNil)
				origImg, _, err := image.Decode(bytes.NewReader(testImageData))
				So(err, ShouldBeNil)
				origX, origY := origImg.Bounds().Dx(), origImg.Bounds().Dy()
				ident, err := h.WriteData(bytes.NewReader(testImageData), "image/png", nil)
				So(err, ShouldBeNil)
				Convey("retrieve image", func() {
					resReader, _, mime, err := h.RetrieveData(ident, nil)
					So(err, ShouldBeNil)
					So(mime, ShouldEqual, "image/png")
					img, _, err := image.Decode(resReader)
					So(err, ShouldBeNil)
					x, y := img.Bounds().Dx(), img.Bounds().Dy()
					So(x, ShouldEqual, origX)
					So(y, ShouldEqual, origY)
					Convey("retrieve cached image", func() {
						resReader, _, mime, err := h.RetrieveData(ident, nil)
						So(err, ShouldBeNil)
						So(mime, ShouldEqual, "image/png")
						img, _, err := image.Decode(resReader)
						So(err, ShouldBeNil)
						x, y := img.Bounds().Dx(), img.Bounds().Dy()
						So(x, ShouldEqual, origX)
						So(y, ShouldEqual, origY)
					})
				})
				Convey("retrieve webp image", func() {
					resReader, _, mime, err := h.RetrieveData(ident, map[string]string{
						"format": "webp",
					})
					So(err, ShouldBeNil)
					So(mime, ShouldEqual, "image/webp")
					img, err := webp.Decode(resReader)
					So(err, ShouldBeNil)
					x, y := img.Bounds().Dx(), img.Bounds().Dy()
					So(x, ShouldEqual, origX)
					So(y, ShouldEqual, origY)
				})
				Convey("retrieve resized image", func() {
					resReader, _, mime, err := h.RetrieveData(ident, map[string]string{
						"size": "80",
					})
					So(err, ShouldBeNil)
					So(mime, ShouldEqual, "image/png")
					img, _, err := image.Decode(resReader)
					So(err, ShouldBeNil)
					x, y := img.Bounds().Dx(), img.Bounds().Dy()
					So(x, ShouldEqual, 80)
					So(y, ShouldEqual, 80)
				})
			})

		})
	})
}
