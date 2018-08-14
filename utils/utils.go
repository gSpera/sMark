package sMark

import (
	"encoding/base64"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
)

//EncodeImageBase64 encodes an image in PNG format and convert it in base64
func EncodeImageBase64(img image.Image) string {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		if err := png.Encode(pw, img); err != nil {
			pw.CloseWithError(err)
		}
	}()
	data, err := ioutil.ReadAll(pr)
	if err != nil {
		log.Println("Cannot decode image:", err)
	}

	return base64.StdEncoding.EncodeToString(data)
}
