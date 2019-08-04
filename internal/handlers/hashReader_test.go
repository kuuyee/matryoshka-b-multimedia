package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashReader(t *testing.T) {
	Convey("test hash reader", t, func() {
		data := randCharStr(999)
		src := bytes.NewBufferString(data)
		hashInit := sha256.New

		expHasher := hashInit()
		expHasher.Write([]byte(data))
		expHash := expHasher.Sum(nil)

		reader := hashReader{
			Hash: hashInit(),
			R:    src,
		}
		_, err := ioutil.ReadAll(reader)
		So(err, ShouldBeNil)

		So(reader.Sum(), ShouldResemble, expHash)
		So(reader.SumHex(), ShouldResemble, hex.EncodeToString(expHash))
	})
}
