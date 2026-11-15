package api

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
)

var errUnsafeSVG = errors.New("unsupported or unsafe SVG")

// inspectSVG accepts a deliberately small, static SVG subset. It does not render
// or sanitize SVG for browser display; accepted bytes must never be served as HTML.
func inspectSVG(data []byte) (int, int, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	allowed := map[string]bool{
		"svg": true, "g": true, "path": true, "rect": true,
		"circle": true, "ellipse": true, "line": true,
		"polyline": true, "polygon": true,
	}
	attrs := map[string]bool{
		"width": true, "height": true, "viewBox": true,
		"d": true, "fill": true, "stroke": true,
		"stroke-width": true, "stroke-linecap": true,
		"stroke-linejoin": true, "stroke-opacity": true,
		"fill-opacity": true, "opacity": true, "transform": true,
		"x": true, "y": true, "x1": true, "x2": true,
		"y1": true, "y2": true, "cx": true, "cy": true,
		"r": true, "rx": true, "ry": true,
		"points": true, "fill-rule": true, "clip-rule": true,
	}
	depth, width, height := 0, 0, 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, errUnsafeSVG
		}
		switch t := token.(type) {
		case xml.Directive, xml.ProcInst:
			return 0, 0, errUnsafeSVG
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				return 0, 0, errUnsafeSVG
			}
		case xml.StartElement:
			if t.Name.Space != "http://www.w3.org/2000/svg" && t.Name.Space != "" {
				return 0, 0, errUnsafeSVG
			}
			if !allowed[t.Name.Local] || (depth == 0 && t.Name.Local != "svg") || depth > 64 {
				return 0, 0, errUnsafeSVG
			}
			if depth == 0 {
				for _, a := range t.Attr {
					if a.Name.Local == "viewBox" {
						var x, y, w, h float64
						fields := strings.Fields(strings.ReplaceAll(a.Value, ",", " "))
						if len(fields) == 4 {
							x, _ = strconv.ParseFloat(fields[0], 64)
							y, _ = strconv.ParseFloat(fields[1], 64)
							w, _ = strconv.ParseFloat(fields[2], 64)
							h, _ = strconv.ParseFloat(fields[3], 64)
							_ = x
							_ = y
							if w > 0 && h > 0 && w <= 4096 && h <= 4096 {
								width, height = int(w), int(h)
							}
						}
					}
				}
			}
			for _, a := range t.Attr {
				if a.Name.Space != "" || !attrs[a.Name.Local] {
					// Namespace declarations are permitted only for the SVG namespace.
					if a.Name.Space == "xmlns" && a.Value == "http://www.w3.org/2000/svg" {
						continue
					}
					if a.Name.Space == "" && a.Name.Local == "xmlns" && a.Value == "http://www.w3.org/2000/svg" {
						continue
					}
					return 0, 0, errUnsafeSVG
				}
				lower := strings.ToLower(a.Value)
				if strings.Contains(lower, "url(") || strings.Contains(lower, "javascript:") ||
					strings.Contains(lower, "data:") || strings.Contains(lower, "@import") ||
					strings.Contains(lower, "<") || strings.Contains(lower, ">") {
					return 0, 0, errUnsafeSVG
				}
			}
			depth++
		case xml.EndElement:
			depth--
			if depth < 0 {
				return 0, 0, errUnsafeSVG
			}
		}
	}
	if depth != 0 || width < 1 || height < 1 {
		return 0, 0, errUnsafeSVG
	}
	return width, height, nil
}
