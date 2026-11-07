package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
)

type point struct{ x, y float64 }

// Render at 4x resolution and downsample for smooth template-icon edges.
func coffeeIcon() ([]byte, error) {
	const size, scale = 22, 4
	canvas := image.NewAlpha(image.Rect(0, 0, size*scale, size*scale))
	line := func(a, b point, width float64) {
		for y := 0; y < canvas.Rect.Dy(); y++ {
			for x := 0; x < canvas.Rect.Dx(); x++ {
				p := point{(float64(x) + 0.5) / scale, size - (float64(y)+0.5)/scale}
				dx, dy := b.x-a.x, b.y-a.y
				t := math.Max(0, math.Min(1, ((p.x-a.x)*dx+(p.y-a.y)*dy)/(dx*dx+dy*dy)))
				if math.Hypot(p.x-a.x-t*dx, p.y-a.y-t*dy) <= width/2 {
					canvas.SetAlpha(x, y, color.Alpha{A: 255})
				}
			}
		}
	}
	curve := func(a, b, c, d point, width float64) {
		previous := a
		for i := 1; i <= 32; i++ {
			t := float64(i) / 32
			u := 1 - t
			next := point{
				u*u*u*a.x + 3*u*u*t*b.x + 3*u*t*t*c.x + t*t*t*d.x,
				u*u*u*a.y + 3*u*u*t*b.y + 3*u*t*t*c.y + t*t*t*d.y,
			}
			line(previous, next, width)
			previous = next
		}
	}
	line(point{4, 12}, point{16, 12}, 1.5)
	line(point{16, 12}, point{15, 6}, 1.5)
	curve(point{15, 6}, point{13, 3}, point{7, 3}, point{5, 6}, 1.5)
	line(point{5, 6}, point{4, 12}, 1.5)
	curve(point{16, 11}, point{21, 12}, point{21, 5}, point{16, 6}, 1.5)
	line(point{3, 3}, point{18, 3}, 1.5)
	for _, x := range []float64{7, 13} {
		curve(point{x, 14}, point{x - 3, 16}, point{x + 3, 18}, point{x, 20}, 1.3)
	}

	icon := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var alpha int
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					alpha += int(canvas.AlphaAt(x*scale+dx, y*scale+dy).A)
				}
			}
			icon.SetNRGBA(x, y, color.NRGBA{A: uint8(alpha / (scale * scale))})
		}
	}
	var data bytes.Buffer
	err := png.Encode(&data, icon)
	return data.Bytes(), err
}
