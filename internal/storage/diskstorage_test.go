package storage

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func storageTest(storage S) {
	exist, err := storage.ExistFile("not_exist.jpg")
	So(err, ShouldBeNil)
	So(exist, ShouldBeFalse)

	testData := make([]byte, 4096)
	rand.Read(testData)
	err = storage.WriteFile("test_file.txt", bytes.NewReader(testData))
	So(err, ShouldBeNil)

	exist, err = storage.ExistFile("test_file.txt")
	So(err, ShouldBeNil)
	So(exist, ShouldBeTrue)

	fileReader, err := storage.RetreiveFile("test_file.txt")
	So(err, ShouldBeNil)
	data, err := ioutil.ReadAll(fileReader)
	So(err, ShouldBeNil)
	So(data, ShouldResemble, testData)
	err = fileReader.Close()
	So(err, ShouldBeNil)
}

func TestDiskStorage(t *testing.T) {
	Convey("test disk storage", t, func() {
		tmpDir, err := ioutil.TempDir("", "mk_disk_test")
		if err != nil {
			t.Error(err)
			return
		}
		storage, err := NewDiskStorage(tmpDir)
		So(err, ShouldBeNil)
		storageTest(storage)
	})
}
