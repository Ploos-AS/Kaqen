package api

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// rasterizeSVG renders only SVG documents accepted by inspectSVG.
// Never pass unvalidated SVG data to this renderer.
func rasterizeSVG(data []byte) (result image.Image, err error) {
	width, height, err := inspectSVG(data)
	if err != nil || !validDimensions(width, height) {
		return nil, errUnsafeSVG
	}
	// Cap rendering work independently of the source viewBox dimensions.
	const maxSide = 2048
	scale := math.Min(1, math.Min(float64(maxSide)/float64(width), float64(maxSide)/float64(height)))
	w := max(1, int(math.Round(float64(width)*scale)))
	h := max(1, int(math.Round(float64(height)*scale)))
	defer func() {
		if recover() != nil {
			result = nil
			err = fmt.Errorf("SVG renderer rejected image")
		}
	}()
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	if err != nil {
		return nil, errUnsafeSVG
	}
	icon.SetTarget(0, 0, float64(w), float64(h))
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, dst, dst.Bounds())
	rasterizer := rasterx.NewDasher(w, h, scanner)
	icon.Draw(rasterizer, 1)
	// Copy into a regular image buffer for the existing icon generator.
	output := image.NewNRGBA(dst.Bounds())
	draw.Draw(output, output.Bounds(), dst, image.Point{}, draw.Src)
	return output, nil
}
