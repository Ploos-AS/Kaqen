package api

import (
	"image"
	"image/color"
	"net/http/httptest"
	"testing"
)

func TestIconOptions(t *testing.T) {
	for _, tc := range []struct {
		padding string
		background string
		valid bool
	}{
		{"0", "transparent", true},
		{"40", "#ff00aa", true},
		{"41", "#ffffff", false},
		{"-1", "", false},
		{"abc", "", false},
		{"10", "red", false},
	} {
		req := httptest.NewRequest("POST", "/", nil)
		req.Form = map[string][]string{"padding": {tc.padding}, "background": {tc.background}}
		_, err := parseIconOptions(req)
		if (err == nil) != tc.valid {
			t.Fatalf("padding=%q background=%q: err=%v", tc.padding, tc.background, err)
		}
	}
}

func TestRenderOptions(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for i := range src.Pix { src.Pix[i] = 255 }
	opt := iconOptions{Padding: 25, Background: color.NRGBA{R: 12, G: 34, B: 56, A: 255}, UseBackground: true}
	out := renderSquareWithOptions(src, 32, opt)
	if got := out.NRGBAAt(0, 0); got != opt.Background {
		t.Fatalf("corner background = %+v", got)
	}
	if got := out.NRGBAAt(16, 16); got.R != 255 || got.G != 255 {
		t.Fatalf("center logo = %+v", got)
	}
}
