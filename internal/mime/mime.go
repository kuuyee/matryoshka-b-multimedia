package mime

import (
	"path"
	"strings"
)

var mime = map[string]string{
	".zip":  "application/zip",
	".rar":  "application/x-rar-compressed",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
	".gif":  "image/gif",
	".txt":  "text/plain",
}

var mimeInverse map[string]string

func init() {
	mimeInverse = make(map[string]string)
	for k, v := range mime {
		mimeInverse[v] = k
	}
}

// ExtToMIME converts file ext to mime type
func ExtToMIME(ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	mime, ok := mime[ext]
	if !ok {
		return ""
	}
	return mime
}

// FileNameToMIME converts file name to mime type
func FileNameToMIME(fn string) string {
	return ExtToMIME(path.Ext(fn))
}

// MIMEToExt converts mime type to file ext
func MIMEToExt(mime string) string {
	ext, ok := mimeInverse[mime]
	if !ok {
		return ""
	}
	return ext
}
