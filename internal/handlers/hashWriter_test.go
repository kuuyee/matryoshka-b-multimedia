package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashWriter(t *testing.T) {
	Convey("test hash writer", t, func() {
		testData := make([]byte, 4096)
		rand.Read(testData)
		hashInit := sha256.New

		expHasher := hashInit()
		expHasher.Write(testData)
		expHash := expHasher.Sum(nil)

		writer := hashWriter{
			Hash: hashInit(),
			W:    ioutil.Discard,
		}
		for offset := 0; offset < 4096; offset += 1024 {
			writer.Write(testData[offset : offset+1024])
		}
		So(writer.Sum(), ShouldResemble, expHash)
		So(writer.SumHex(), ShouldResemble, hex.EncodeToString(expHash))
	})
}
