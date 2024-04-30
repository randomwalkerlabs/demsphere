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
	//outputFile = kingpin.Flag("output", "Output STL file to write.").Required().Short('o').String()
	Planet = "Mars"
	MinDetail = 9
	MaxDetail = 12
	MeanRadius = 3389500
	MinElevation = -11000
	MaxElevation = 21900
	Tolerance = 30
	Exaggeration = 15
	Scale = 1
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

	// Print all the variables
	output := fmt.Sprintf("\nPlanet: %s\nMinDetail: %d\nMaxDetail: %d\nMeanRadius: %d\nMinElevation: %d\nMaxElevation: %d\nTolerance: %d\nExaggeration: %d\nScale: %d\n",
        Planet,
        MinDetail,
        MaxDetail,
        MeanRadius,
        MinElevation,
        MaxElevation,
        Tolerance,
        Exaggeration,
        Scale,
    )
    fmt.Println(output)

	// Image, minDetail, maxDetail, meanRadius, minElevation, maxElevation, tolerance, exaggeration, scale

	// mercury
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 12, 2439400, -10764, 8994, 50, 4, 1.0/2439400)

	// moon
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 11, 1737400, -18257, 21563, 50, 3, 1.0/1737400)

	// FINAL mars
	// triangulator = demsphere.NewTriangulator(
	// 	im, 9, 12, 3396190, -8201, 21241, 30, 15, 1)

	triangulator := demsphere.NewTriangulator(
		im, int(MinDetail), int(MaxDetail), float64(MeanRadius), float64(MinElevation), float64(MaxElevation), float64(Tolerance), float64(Exaggeration), float64(Scale))


	// FINAL earth
	//triangula   tor := demsphere.NewTriangulator(
	//im, 6, 11, 6373934, -10900, 8849, 30, 15, 1)

	// pluto
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 12, 1188300, -4101, 6491, 50, 3, 1.0/1188300)

	done = timed("Generating positive mesh")
	triangles := triangulator.Triangulate()
	done()
	fmt.Println(fmt.Sprintf("Generated %v triangles for outer mesh", len(triangles)))

	// Filename
	filename_out := fmt.Sprintf("%s_Outer_%d_%d_%d_%d.stl", Planet, MinDetail, MaxDetail, Tolerance, Exaggeration)
	fmt.Println(fmt.Sprintf("Filename set to %s", filename_out))
	// Writing STL
	done = timed("Writing STL for outer shell")
	// Making a string for filename
	demsphere.WriteSTLFile(filename_out, triangles)
	done()

	// fmt.Println(len(triangles))
	// done = timed("writing outer output")
	// //earth
	// // demsphere.WriteSTLFile("Earth_outer_6_11_30_15.stl", triangles)
	// //mars
	// demsphere.WriteSTLFile("Mars_outer_9_12_30_15.stl", triangles)
	// done()

	im_inv := imaging.Invert(im)
	
	// FINAL earth
	//triangulator = demsphere.NewTriangulator(
	//	im, 6, 11, 6373934, -10900, 8849, 30, 15, 1)

	// FINAL mars
	// triangulator = demsphere.NewTriangulator(
	// 	im_inv, 9, 12, 3396190, -8201, 21241, 30, 15, 1)
	triangulator = demsphere.NewTriangulator(
		im_inv, int(MinDetail), int(MaxDetail), float64(MeanRadius), float64(MinElevation), float64(MaxElevation), float64(Tolerance), float64(Exaggeration), float64(Scale))

	done = timed("Generating negative mesh")
	inner := triangulator.Triangulate()
	for i, t := range inner {
		inner[i] = demsphere.Triangle{t.C, t.B, t.A}
	}
	//triangles = append(triangles, inner...)
	fmt.Println(fmt.Sprintf("Generated %v triangles for inner mesh", len(inner)))
	done()

	// Filename
	filename_in := fmt.Sprintf("%s_Inner_%d_%d_%d_%d.stl", Planet, MinDetail, MaxDetail, Tolerance, Exaggeration)
	fmt.Println(fmt.Sprintf("Filename set to %s", filename_in))
	// Writing STL
	done = timed("Writing STL for inner shell")
	// Making a string for filename
	demsphere.WriteSTLFile(filename_in, inner)
	done()
	
	// done = timed("writing inner output")
	// //earth
	// // demsphere.WriteSTLFile("Earth_inner_6_11_30_15.stl", triangles)
	// //mars
	// demsphere.WriteSTLFile("Mars_inner_9_12_30_15.stl", inner)
	// done()
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