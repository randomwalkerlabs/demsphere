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
	outputFile = kingpin.Flag("output", "Output STL file to write.").Required().Short('o').String() // Currently not in use
)

var (
	planet           = "earth" // Planet
	minDetail        = 15      //
	maxDetail        = 30
	meanRadius       = 6373934.0        // Mean radius of the planet
	minElevation     = -10900.0         // Highest point of planet
	maxElevation     = 8849.0           // Lowest point of planet
	tolerance        = 50.0             // Tolerance - accurate to this meters
	exaggeration     = 15.0             // Vertical exaggeration
	scale            = 1.0 / meanRadius // Final scaling
	innerShellScale  = 1                // Scaling for inner shell, compared to outer one
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
	// Reading input image
	done = timed("reading input")
	im, err := fauxgl.LoadImage(*inputFile)
	done()
	if err != nil {
		log.Fatal(err)
	}
	// Make outer shell
	triangulator := demsphere.NewTriangulator(
		im, minDetail, maxDetail, meanRadius, minElevation, maxElevation, tolerance, exaggeration, scale)
	done = timed("generating outer mesh")
	triangles := triangulator.Triangulate()
	done()
	fmt.Println(len(triangles))
	// Making inner shell
	im = imaging.Invert(im) //Inverting image
	triangulator = demsphere.NewTriangulator(
		im, minDetail, maxDetail, meanRadius, minElevation, maxElevation, tolerance, exaggeration, scale*float64(innerShellScale))
	// Flippig normals
	inner := triangulator.Triangulate()
	for i, t := range inner {
		inner[i] = demsphere.Triangle{t.C, t.B, t.A}
	}
	done = timed("generating inner mesh")
	triangles = append(triangles, inner...)
	fmt.Println(len(triangles))
	done()
	// Writing stl
	done = timed("writing output")
	// Make a filename with important parameters to test the output
	filename := fmt.Sprintf("%s_%d_%d_%d_%d.stl", planet, minDetail, maxDetail, tolerance, exaggeration)
	demsphere.WriteSTLFile(filename, triangles)
	done()
}

// Planet data

// Planet	Mean Radius (m)	Location	Max Elevation (m)	Min Depth (m)	Deepest Point Location
// Mercury	2,439,700	Caloris Basin	5,500	5,500	Caloris Basin
// Venus	6,051,800	Maxwell Montes	11,000	N/A	N/A
// Earth	6,371,000	Mount Everest	8,848.86	-10,994	Mariana Trench
// Mars	3,389,500	Olympus Mons	21,900	11,000	Hellas Planitia
// Jupiter	69,911,000	N/A	N/A	N/A	N/A
// Saturn	58,232,000	N/A	N/A	N/A	N/A
// Uranus	25,362,000	N/A	N/A	N/A	N/A
// Neptune	24,622,000	N/A	N/A	N/A	N/A
// Earth's Moon	1,737,100	Mons Huygens	10,800	9,000	Mare Imbrium