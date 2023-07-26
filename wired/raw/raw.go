package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
)

// adapted from timotew on github

const SCREEN_WIDTH, SCREEN_HEIGHT = 184, 96 // 240,240 // 180,240

// static image as single frame is quick to disappear compared to a gifs
// extra frames gives them more stage time, anything below 30 may not be observed by humans
const STATIC_IMAGE_FRAMES = 30

func convertPixesTo16BitRGB(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	R, G, B := int(r/257), int(g/257), int(b/257)

	return uint16((int(R>>3) << 11) |
		(int(G>>2) << 5) |
		(int(B>>3) << 0))
}

func convertPixelsToRawBitmap(image *image.RGBA, bitmap []uint16) {
	imgHeight, imgWidth := image.Bounds().Max.Y, image.Bounds().Max.X

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			bitmap[(y)*SCREEN_WIDTH+(x)] = convertPixesTo16BitRGB(image.At(x, y))
		}
	}
}

func imagesToAnim(frames []*image.Paletted) []byte {
	var buf bytes.Buffer
	overpaintImage := image.NewRGBA(image.Rect(0, 0, frames[0].Bounds().Max.X, frames[0].Bounds().Max.Y))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), frames[0], image.Point{}, draw.Src)

	bitmap := make([]uint16, SCREEN_WIDTH*SCREEN_HEIGHT)
	for _, srcImg := range frames {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.Point{}, draw.Over)
		convertPixelsToRawBitmap(overpaintImage, bitmap)

		for _, ui := range bitmap {
			binary.Write(&buf, binary.LittleEndian, ui)
		}
	}
	return buf.Bytes()
}

func getGifDimensions(gif *gif.GIF) (width, height int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}

func GifToBootAnimation(gifFile []byte, output string) error {
	bootAnimFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer bootAnimFile.Close()

	byteR := bytes.NewBuffer(gifFile)
	gifImages, err := gif.DecodeAll(byteR)
	if err != nil {
		return err
	}

	imgWidth, imgHeight := getGifDimensions(gifImages)
	if imgHeight != SCREEN_HEIGHT || imgWidth != SCREEN_WIDTH {
		return fmt.Errorf("width %dpx height %dpx file is required", SCREEN_WIDTH, SCREEN_HEIGHT)
	}

	imgBytes := imagesToAnim(gifImages.Image)
	err = os.WriteFile(output, imgBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
