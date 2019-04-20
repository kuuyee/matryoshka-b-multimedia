package handlers

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/chai2010/webp"

	"github.com/nfnt/resize"

	"github.com/h2non/filetype"

	"github.com/kuuyee/matryoshka-b-multimedia/internal/mime"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
)

type ImageHandler struct {
	Storage    storage.S
	MaxSize    int64
	ResizeAlgo resize.InterpolationFunction
	KeyedMutex *KeyedRWMutex
}

func imgDimensionDecode(dimensionStr string) (uint, uint, error) {
	thumbDimensions := strings.SplitN(dimensionStr, "x", 2)
	if len(thumbDimensions) == 0 || len(thumbDimensions) > 2 {
		return 0, 0, errors.New("invalid input thumb dimension")
	}
	if len(thumbDimensions) == 1 {
		thumbDimensions = []string{"0", thumbDimensions[0]}
	}
	maxX, err := strconv.ParseUint(thumbDimensions[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	maxY, err := strconv.ParseUint(thumbDimensions[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return uint(maxX), uint(maxY), nil
}

// SizeLimit implements H
func (h *ImageHandler) SizeLimit() int64 {
	return h.MaxSize
}

// Type implements H
func (h *ImageHandler) Type() string {
	return "image"
}

func (h *ImageHandler) imagePreprocess(w io.Writer, rawImageData []byte, mime string, param map[string]string) error {
	if orig, _ := param["orig"]; mime == "image/gif" || orig == "1" || orig == "true" {
		_, err := filetype.Image(rawImageData)
		if err != nil {
			return err
		}
		w.Write(rawImageData)
		return nil
	}
	var img image.Image
	var err error
	switch mime {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(rawImageData))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(rawImageData))
	case "image/webp":
		img, err = webp.Decode(bytes.NewReader(rawImageData))
	default:
		return errors.New("format not supported")
	}
	if err != nil {
		return err
	}
	if thumbStr, exist := param["thumb"]; exist {
		maxX, maxY, err := imgDimensionDecode(thumbStr)
		if err != nil {
			return err
		}
		img = resize.Thumbnail(maxX, maxY, img, h.ResizeAlgo)
	}
	return png.Encode(w, img)
}

// WriteData implements H
func (h *ImageHandler) WriteData(r io.Reader, mime string, param map[string]string) (ident string, err error) {
	rawImageData, err := ioutil.ReadAll(io.LimitReader(r, h.MaxSize))
	if err != nil {
		return "", err
	}

	processedImageData := bytes.NewBuffer([]byte{})
	hasher := hashWriter{
		Hash: sha256.New(),
		W:    processedImageData,
	}
	err = h.imagePreprocess(hasher, rawImageData, mime, param)
	if err != nil {
		return "", err
	}

	ident = hasher.SumHex() + ".png"

	h.KeyedMutex.GetMutex(ident).Lock()
	defer h.KeyedMutex.GetMutex(ident).Unlock()

	fileWriter, err := h.Storage.WriteFile(ident)
	if err != nil {
		return "", err
	}
	defer fileWriter.Close()
	_, err = io.Copy(fileWriter, processedImageData)
	if err != nil {
		return "", err
	}
	return ident, nil
}

type imageRetrieveParamOpt struct {
	X, Y   uint
	Format string
}

func parseImgRetrieveParamOpt(param map[string]string) (*imageRetrieveParamOpt, error) {
	res := &imageRetrieveParamOpt{
		Format: "image/png",
	}
	if format, exist := param["format"]; exist {
		res.Format = mime.ExtToMIME(format)
		if res.Format == "" {
			return nil, errors.New("format not supported")
		}
	}
	if dimenstionStr, exist := param["size"]; exist {
		x, y, err := imgDimensionDecode(dimenstionStr)
		if err != nil {
			return nil, err
		}
		res.X, res.Y = x, y
	}
	return res, nil
}

func imgIdentAltn(origIdent string, param *imageRetrieveParamOpt) string {
	newIdent := origIdent[:len(origIdent)-len(path.Ext(origIdent))]
	if param.X != 0 || param.Y != 0 {
		newIdent += fmt.Sprintf("_%dx%d", param.X, param.Y)
	}
	switch param.Format {
	case "image/jpg", "image/jpeg":
		newIdent += ".jpg"
	case "image/png":
		newIdent += ".png"
	case "image/gif":
		newIdent += ".gif"
	case "image/webp":
		newIdent += ".webp"
	default:
		newIdent += "_" + param.Format
	}
	return newIdent
}

func (h *ImageHandler) prepareImageAltn(ident string, opt *imageRetrieveParamOpt) error {
	h.KeyedMutex.GetMutex(ident).RLock()
	defer h.KeyedMutex.GetMutex(ident).RUnlock()
	targetIdent := imgIdentAltn(ident, opt)
	h.KeyedMutex.GetMutex(targetIdent).Lock()
	defer h.KeyedMutex.GetMutex(targetIdent).Unlock()
	origReader, err := h.Storage.RetreiveFile(ident)
	if err != nil {
		return err
	}
	defer origReader.Close()
	img, _, err := image.Decode(origReader)
	if err != nil {
		return err
	}
	if opt.X != 0 || opt.Y != 0 {
		origX, origY := img.Bounds().Dx(), img.Bounds().Dy()
		if opt.X > uint(origX*3) || opt.Y > uint(origY*3) {
			return errors.New("requested image too large")
		}
		img = resize.Resize(opt.X, opt.Y, img, h.ResizeAlgo)
	}
	targetW, err := h.Storage.WriteFile(targetIdent)
	if err != nil {
		return err
	}
	defer targetW.Close()
	switch opt.Format {
	case "image/jpeg":
		err = jpeg.Encode(targetW, img, nil)
	case "image/png":
		err = png.Encode(targetW, img)
	case "image/webp":
		err = webp.Encode(targetW, img, nil)
	default:
		err = errors.New("unknown format")
	}
	if err != nil {
		return err
	}
	return nil
}

// RetrieveData implements H
func (h *ImageHandler) RetrieveData(ident string, param map[string]string) (io.ReadCloser, string, error) {
	if _, hasExplicitFormat := param["format"]; !hasExplicitFormat && strings.HasSuffix(ident, ".gif") {
		param["format"] = "gif"
	}
	opt, err := parseImgRetrieveParamOpt(param)
	if err != nil {
		return nil, "", err
	}
	targetIdent := imgIdentAltn(ident, opt)
	exist, err := func() (bool, error) {
		h.KeyedMutex.GetMutex(targetIdent).RLock()
		defer h.KeyedMutex.GetMutex(targetIdent).RUnlock()
		return h.Storage.ExistFile(targetIdent)
	}()
	if err != nil {
		return nil, "", err
	}
	if exist {
		file, err := h.Storage.RetreiveFile(targetIdent)
		return file, opt.Format, err
	}

	if err := h.prepareImageAltn(ident, opt); err != nil {
		return nil, "", err
	}
	file, err := h.Storage.RetreiveFile(targetIdent)
	return file, opt.Format, err
}
