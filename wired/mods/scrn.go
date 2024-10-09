package mods

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
)

func GenerateScreenData(number int) []uint16 {
	width := 184
	height := 96
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	drawSmallText(img, "Say your wake-word in:", 10, 20)
	numberStr := fmt.Sprintf("%d", number)
	drawBigText(img, numberStr, width/2, height/2)
	data := ConvertImageToRGB565(img)

	return data
}

func drawSmallText(img *image.RGBA, text string, x, y int) {
	face := basicfont.Face7x13
	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)
}

func drawBigText(img *image.RGBA, text string, x, y int) {
	f, err := freetype.ParseFont(gobold.TTF)
	if err != nil {
		log.Println(err)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(48)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)

	face := truetype.NewFace(f, &truetype.Options{Size: 48})
	defer face.Close()

	textWidth := font.MeasureString(face, text)
	textWidthInt := int(textWidth >> 6)

	x = x - textWidthInt/2
	y = y + int(c.PointToFixed(48)>>6)/2

	pt := freetype.Pt(x, y)

	_, err = c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
	}
}

func ConvertImageToRGB565(img image.Image) []uint16 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	data := make([]uint16, width*height)

	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			data[index] = ColorToRGB565(r8, g8, b8)
			index++
		}
	}

	return data
}

func ColorToRGB565(r, g, b uint8) uint16 {
	r5 := uint16(r >> 3)
	g6 := uint16(g >> 2)
	b5 := uint16(b >> 3)
	return (r5 << 11) | (g6 << 5) | b5
}
