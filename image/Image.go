package image

import (
	"bytes"
	"errors"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"strings"
	"utilities/image/png"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2022/11/21 10:10
* @Package:
 */

var UnknownFormat = errors.New("unknown format image ")

type decode func(r io.Reader) (image.Config, error)

const base = "image"

const (
	PNG  = "png"
	JPG  = "jpg"
	JPEG = "jpeg"
	BMP  = "bmp"
	GIF  = "gif"
	WEBP = "webp"
	ICO  = "ico"
	SVG  = "svg"
	TIFF = "tiff"
	AI   = "ai"
	CDR  = "cdr"
	EPS  = "eps"
)

type Image struct {
	Type          string
	Width, Height int
}

var decodeMapping = map[string]decode{
	PNG:  png.DecodeConfig,
	JPG:  jpeg.DecodeConfig,
	JPEG: jpeg.DecodeConfig,
	WEBP: webp.DecodeConfig,
	BMP:  bmp.DecodeConfig,
	GIF:  gif.DecodeConfig,
}

func DecodeImage(data []byte) (*Image, error) {
	image := &Image{}
	// 获取图片类型
	ty := detectType(data[:sniffLen])
	arrs := strings.Split(ty, "/")
	if len(arrs) < 2 || arrs[0] != base {
		return nil, UnknownFormat
	}
	image.Type = arrs[1]
	// 获取图片常用指标
	d, ok := decodeMapping[strings.ToLower(arrs[1])]
	if !ok {
		// 不支持的解码格式
		return nil, UnknownFormat
	}
	// 获取长宽信息
	c, e := d(bytes.NewReader(data))
	if e != nil {
		return image, UnknownFormat
	}
	image.Height = c.Height
	image.Width = c.Width
	return image, nil
}
