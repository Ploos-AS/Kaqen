package api

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	xdraw "golang.org/x/image/draw"
)

const (
	socialWidth  = 1200
	socialHeight = 630
)

// renderSocialCard fits the complete logo into a 1200x630 social sharing image.
// A solid background is necessary because many social platforms discard alpha.
func renderSocialCard(src image.Image, opt iconOptions) *image.NRGBA {
	bg := opt.Background
	if !opt.UseBackground {
		bg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}
	dst := image.NewNRGBA(image.Rect(0, 0, socialWidth, socialHeight))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)
	bounds := src.Bounds()
	sw, sh := bounds.Dx(), bounds.Dy()
	if sw < 1 || sh < 1 {
		return dst
	}
	maxW, maxH := 960, 450
	if opt.Padding > 0 {
		maxW = max(1, maxW*(100-opt.Padding)/100)
		maxH = max(1, maxH*(100-opt.Padding)/100)
	}
	targetW, targetH := maxW, maxH
	if int64(sw)*int64(maxH) > int64(sh)*int64(maxW) {
		targetH = max(1, int(int64(maxW)*int64(sh)/int64(sw)))
	} else {
		targetW = max(1, int(int64(maxH)*int64(sw)/int64(sh)))
	}
	rect := image.Rect((socialWidth-targetW)/2, (socialHeight-targetH)/2, (socialWidth+targetW)/2, (socialHeight+targetH)/2)
	xdraw.CatmullRom.Scale(dst, rect, src, bounds, xdraw.Over, nil)
	return dst
}

func encodeSocialCard(src image.Image, opt iconOptions) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, renderSocialCard(src, opt))
	return buf.Bytes(), err
}
