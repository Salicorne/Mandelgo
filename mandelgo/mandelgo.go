package mandelgo

import (
	"image/color"
	"math"
	"math/cmplx"
)

var (
	max_iter = 300

	pallet []Color = []Color{
		{R: 0, G: 0, B: 0},
		{R: 255, G: 0, B: 0},
		{R: 255, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 255},
		{R: 0, G: 0, B: 255},
		{R: 255, G: 0, B: 255},
		{R: 255, G: 255, B: 255},
	}
)

type Color struct {
	R uint8
	G uint8
	B uint8
}

func interp(a uint8, b uint8, t float64) uint8 {
	return uint8((1.0-t)*float64(a) + t*float64(b))
}

func GetColor(x float64, y float64) color.Color {
	iter := mandelbrot(x, y)
	if iter == max_iter {
		return color.Black
	}
	coeff := float64(iter) / float64(max_iter)
	//coeff = math.Sqrt(coeff)
	wcoeff := coeff * float64(len(pallet)-1)
	idx := int(math.Floor(wcoeff))

	c1 := pallet[idx]
	c2 := pallet[idx+1]
	return color.RGBA{
		R: interp(c1.R, c2.R, wcoeff-float64(idx)),
		G: interp(c1.G, c2.G, wcoeff-float64(idx)),
		B: interp(c1.B, c2.B, wcoeff-float64(idx)),
		A: 255,
	}
}

func mandelbrot(x float64, y float64) int {
	z := complex(0, 0)
	for iter := 0; iter < max_iter; iter++ {
		z = z*z + complex(x, y)
		if cmplx.Abs(z) > 2 {
			return iter
		}
	}
	return max_iter
}
