//+build ignore

package main

import (
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/barnex/jit"
)

var port = flag.String("http", ":8080", "HTTP service address")

func main() {
	flag.Parse()
	http.HandleFunc("/plot/", handlePlot)
	log.Println("Serving at", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func handlePlot(w http.ResponseWriter, r *http.Request) {
	expr := r.URL.Path[len("/plot/"):]
	code, err := jit.Compile(expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nx, ny := 500, 500
	xmin, xmax := -10.0, 10.0
	ymin, ymax := -10.0, 10.0

	dst := make([]float64, nx*ny)
	matrix := make([][]float64, ny)
	for iy := range matrix {
		matrix[iy] = dst[iy*nx : (iy+1)*nx]
	}

	code.Eval2D(dst, xmin, xmax, nx, ymin, ymax, ny)

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))
	for iy := 0; iy < ny; iy++ {
		for ix := 0; ix < nx; ix++ {
			img.Set(ix, iy, color.White)
		}
	}
	pen := color.RGBA{B: 150}
	for iy := 0; iy < ny; iy++ {
		for ix := 0; ix < nx-1; ix++ {
			if matrix[iy][ix]*matrix[iy][ix+1] < 0 {
				img.Set(ix, iy, pen)
			}
		}
	}
	for iy := 0; iy < ny-1; iy++ {
		for ix := 0; ix < nx; ix++ {
			if matrix[iy][ix]*matrix[iy+1][ix] < 0 {
				img.Set(ix, iy, pen)
			}
		}
	}
	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
}
