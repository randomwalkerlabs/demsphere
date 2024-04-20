package main

import (
	"fmt"
	"log"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/demsphere"
	"github.com/fogleman/fauxgl"
	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	inputFile  = kingpin.Flag("input", "Input DEM image to process.").Required().Short('i').ExistingFile()
	outputFile = kingpin.Flag("output", "Output STL file to write.").Required().Short('o').String()
)

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func main() {
	var done func()

	kingpin.Parse()

	done = timed("reading input")
	im, err := fauxgl.LoadImage(*inputFile)
	done()
	if err != nil {
		log.Fatal(err)
	}

	// Image, minDetail, maxDetail, meanRadius, minElevation, maxElevation, tolerance, exaggeration, scale

	// mercury
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 12, 2439400, -10764, 8994, 50, 4, 1.0/2439400)

	// moon
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 11, 1737400, -18257, 21563, 50, 3, 1.0/1737400)

	// mars
	//triangulator := demsphere.NewTriangulator(
	//	im, 9, 12, 3396190, -8201, 21241, 50, 10, 1.0/3396190)

	// earth
	triangulator := demsphere.NewTriangulator(
	im, 9, 12, 6373934, -10900, 8849, 50, 15, 1.0/6373934)

	// pluto
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 12, 1188300, -4101, 6491, 50, 3, 1.0/1188300)

	done = timed("generating positive mesh")
	triangles := triangulator.Triangulate()
	done()

	fmt.Println(len(triangles))

	im = imaging.Invert(im)
	triangulator = demsphere.NewTriangulator(
		im, 9, 12, 6373934, -10900, 8849, 50, 15, 1.0/6373934)

	inner := triangulator.Triangulate()
	for i, t := range inner {
		inner[i] = demsphere.Triangle{t.C, t.B, t.A}
	}
	done = timed("generating negative mesh")
	triangles = append(triangles, inner...)
	fmt.Println(len(triangles))
	done()

	// inner := fauxgl.NewSphere(4)
	// inner.Transform(fauxgl.Scale(fauxgl.V(0.85, 0.85, 0.85)))
	// inner.ReverseWinding()
	// for _, t := range inner.Triangles {
	// 	p1 := demsphere.Vector(t.V1.Position)
	// 	p2 := demsphere.Vector(t.V2.Position)
	// 	p3 := demsphere.Vector(t.V3.Position)
	// 	triangles = append(triangles, demsphere.Triangle{p1, p2, p3})
	// }

	// fmt.Println(len(triangles))

	done = timed("writing output")
	demsphere.WriteSTLFile(*outputFile, triangles)
	done()
}

// 4,5120,4.7372172692
// 5,20480,2.3686086346009336
// 6,81920,1.1844992435794788
// 7,327680,0.5635352519913389
// 8,1310720,0.2833191921173619
// 9,5242880,0.14087309976888143
// 10,20971520,0.07043426041304851
// 11,83886080,0.036106462472517295
// 12,335544320,0.018169800357500515




// Planet	Radius (m)	Min Elevation (m)	Max Elevation (m)
// Mercury	2,439,700	-5,380	4,250
// Venus	6,051,800	-650	650
// Earth	6,371,000	-10,900	8,850
// Mars	3,389,500	-8,200	21,900
// Jupiter	69,911,000	N/A	N/A
// Saturn	58,232,000	N/A	N/A
// Uranus	25,362,000	N/A	N/A
// Neptune	24,622,000	N/A	N/A
// Pluto	1,186,000	-5,900	5,900