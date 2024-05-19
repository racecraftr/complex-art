package main

// some imports
import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"

	"github.com/PerformLine/go-stockutil/colorutil"
)

const (

	// WIDTH and HEIGHT represnt width and height of image respectively.
	WIDTH  = 1920
	HEIGHT = 1080

	// represents the RANGE of real and imaginary.
	// goes from (-RANGE/2) to (RANGE/2) along real and imaginary axis.
	RANGE  = 10
)

var min = 0.0

func main() {
	min = math.Min(float64(WIDTH), float64(HEIGHT))
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			z := PointToCplx(x, y) // convert point to complex number.
			z = fn(z) // plug number into function.
			c := CplxToColor(z) // convert result into color.
			img.Set(x, y, c) // render color.
		}
	}

	f, _ := os.Create("res.png")
	png.Encode(f, img)
}

// complex number function. you can put any equation here
func fn(z complex128) complex128 {
	// return (z*cmplx.Sin(z)*z - complex(1, 1)) / (z * cmplx.Cos(z))
	// return (z * z - 1) * cmplx.Pow(z - 2 - complex(0, 1), 2) / (z * z + complex(2, 2)) // got this one off of wikipedia

	abs := cmplx.Abs(z)
	return z * complex(math.Cos(abs / 2), 0) * complex(1, math.Sin(abs * 10)) / 2

	// return z + complex(imag(z), real(z)) * cmplx.Conj(cmplx.Sin(z * z) * cmplx.Cos(z * z)) / (cmplx.Atanh(z * z) + 1)
}

// some pre-baked equations I made.
func fractal(z complex128) complex128 { // this is the one used in the blog post.
	return z + complex(imag(z), real(z)) * cmplx.Conj(cmplx.Sin(z * z) * cmplx.Cos(z * z)) / 1
}

func filter(z complex128) complex128 {
	return complex(imag(z) + imag(z), real(z) / math.Abs(real(z)))
}

func spiral(z complex128) complex128 {
	abs := cmplx.Abs(z)
	return cmplx.Pow(z, cmplx.Sinh(complex(abs / 10, abs / 10))) * complex(math.Sin(abs), 0) * complex(1, math.Cos(abs))
}

// converts the pixel location into a complex number.
func PointToCplx(x, y int) complex128 {
	// real along x axis.
	real := float64(x-WIDTH/2) / min * RANGE
	// imaginary along y axis.
	// due to the way computer coordinates work,
	// y increases as you go down instead of up.
	imag := float64(HEIGHT/2-y) / min * RANGE
	return complex(real, imag)
}

// converts the complex number into a color.
func CplxToColor(z complex128) color.Color {
	abs, arg := cmplx.Polar(z)
	// l := (abs * abs) / (abs*abs + 1.0)
	l := (2.0 / math.Pi * math.Atan(abs))
	h := (arg + math.Pi) / (2 * math.Pi) * 360.0
	if math.IsNaN(h) {
		l = 0
		h = 0
	}
	r, g, b := colorutil.HslToRgb(h, 0.7, l)


	// // topography function. Idk why i made this
	// // draws lines to represent different absolute values.
	// diff := func(n float64) float64 {
	// 	if n == 0 {
	// 		return 0
	// 	}
	// 	return (math.Abs(n - float64(int(n))))
	// }

	// if int(abs) != 0 &&
	// 	(diff(real(z)) <= RANGE/100.0 || diff(imag(z)) <= RANGE/100.0) {
	// 		r, g, b = 255 - r, 255 - g, 255 - b
	// }

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
