package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
)

const (
	virt_x0 = -2
	virt_w  = 3
	virt_y0 = -1
	virt_h  = 2

	height = 1080
	width  = int(virt_w * height / float64(virt_h))

	max_iter = 100
)

func circle(x float64, y float64) bool {
	return math.Sqrt(x*x+y*y) < 1
}

func mandelbrot(x float64, y float64) bool {
	z := complex(0, 0)
	for iter := 0; iter < max_iter; iter++ {
		z = z*z + complex(x, y)
		if cmplx.Abs(z) > 2 {
			return true
		}
	}
	return false
}

func main() {
	fmt.Printf("Hello Go ! \n")

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(10, 10, color.RGBA{255, 0, 0, 255})

	for y := 0; y != height; y++ {
		for x := 0; x < width; x++ {
			var vx float64 = (float64(x)/float64(width))*virt_w + virt_x0
			var vy float64 = (float64(y)/float64(height))*virt_h + virt_y0
			if mandelbrot(vx, vy) {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	file, err := os.Create("res.png")
	if err != nil {
		log.Fatalf("Error opening result file: %v", err)
	}

	if err := png.Encode(file, img); err != nil {
		log.Fatalf("Error writing result file: %v", err)
	}
}
