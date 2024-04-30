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
	Planet = "Mars"
    MinDetail = 9
    MaxDetail = 12
    MeanRadius = 3389500
    MinElevation = -11000
    MaxElevation = 21900
    Tolerance = 30
    Exaggeration = 10
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
	// Reading input DEM
	done = timed("Reading input DEM")
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

	// Working on the outer shell
	triangulator := demsphere.NewTriangulator(
		im, int(MinDetail), int(MaxDetail), float64(MeanRadius), float64(MinElevation), float64(MaxElevation), float64(Tolerance), float64(Exaggeration), float64(Scale))
	done = timed("Generating positive mesh")
	triangles := triangulator.Triangulate()
	done()
	fmt.Println(fmt.Sprintf("Generated %v triangles for outer mesh", len(triangles)))
	// Writing the Outer STL
	filename_out := fmt.Sprintf("%s_Outer_%d_%d_%d_%d.stl", Planet, MinDetail, MaxDetail, Tolerance, Exaggeration)
	fmt.Println(fmt.Sprintf("Filename set to %s", filename_out))
	done = timed("Writing STL for outer shell")
	demsphere.WriteSTLFile(filename_out, triangles)
	done()

	// Working on inner shell
	im_inv := imaging.Invert(im)
	triangulator = demsphere.NewTriangulator(
		im_inv, int(MinDetail), int(MaxDetail), float64(MeanRadius), float64(MinElevation), float64(MaxElevation), float64(Tolerance), float64(Exaggeration), float64(Scale))

	done = timed("Generating negative mesh")
	inner := triangulator.Triangulate()
	for i, t := range inner {
		inner[i] = demsphere.Triangle{t.C, t.B, t.A}
	}
	fmt.Println(fmt.Sprintf("Generated %v triangles for inner mesh", len(inner)))
	done()

	// Writing the inner shell
	filename_in := fmt.Sprintf("%s_Inner_%d_%d_%d_%d.stl", Planet, MinDetail, MaxDetail, Tolerance, Exaggeration)
	fmt.Println(fmt.Sprintf("Filename set to %s", filename_in))
	done = timed("Writing STL for inner shell")
	demsphere.WriteSTLFile(filename_in, inner)
	done()
}

	// Planet = "Mercury"
    // MinDetail = 9
    // MaxDetail = 12
    // MeanRadius = 2439700
    // MinElevation = -5380
    // MaxElevation = 4480
    // Tolerance = 30
    // Exaggeration = 5
	// Scale = 1

	// Planet = "Venus"
    // MinDetail = 6
    // MaxDetail = 10
    // MeanRadius = 6051800
    // MinElevation = -1000
    // MaxElevation = 11000
    // Tolerance = 50
    // Exaggeration = 15
	// Scale = 1

	// Planet = "Earth"
	// MinDetail = 6
	// MaxDetail = 11
	// MeanRadius = 6373934
	// MinElevation = -10900
	// MaxElevation = 8849
	// Tolerance = 30
	// Exaggeration = 15
	// Scale = 1

	// Planet = "Mars"
    // MinDetail = 9
    // MaxDetail = 12
    // MeanRadius = 3389500
    // MinElevation = -11000
    // MaxElevation = 21900
    // Tolerance = 30
    // Exaggeration = 15
	// Scale = 1

	// Planet = "Jupiter"
    // MinDetail = 9
    // MaxDetail = 12
    // MeanRadius = 69911000
    // MinElevation = -5000
    // MaxElevation = 5000
    // Tolerance = 30
    // Exaggeration = 100
	// Scale = 1

	// Planet = "Saturn"
    // MinDetail = 6
    // MaxDetail = 10
    // MeanRadius = 58232000
    // MinElevation = -1000
    // MaxElevation = 1000
    // Tolerance = 50
    // Exaggeration = 30
	// Scale = 1

	// Planet = "Uranus"
    // MinDetail = 6
    // MaxDetail = 10
    // MeanRadius = 25362000
    // MinElevation = -1000
    // MaxElevation = 1000
    // Tolerance = 50
    // Exaggeration = 30
	// Scale = 1

	// Planet = "Neptune"
    // MinDetail = 6
    // MaxDetail = 10
    // MeanRadius = 24622000
    // MinElevation = -1000
    // MaxElevation = 1000
    // Tolerance = 50
    // Exaggeration = 30
	// Scale = 1

	// Planet = "Moon"
    // MinDetail = 6
    // MaxDetail = 11
    // MeanRadius = 1737100
    // MinElevation = -9000
    // MaxElevation = 10800
    // Tolerance = 30
    // Exaggeration = 15
	// Scale = 1

// Planet data

// Planet	Mean Radius (m)	Location	Max Elevation (m)	Min Depth (m)	Deepest Point Location
// Mercury	2,439,700	Caloris Basin	4,480				5,380			Caloris Basin
// Venus	6,051,800	Maxwell Montes	11,000				N/A				N/A
// Earth	6,371,000	Mount Everest	8,848.86			-10,994			Mariana Trench
// Mars		3,389,500	Olympus Mons	21,900				11,000			Hellas Planitia
// Jupiter	69,911,000	N/A				N/A					N/A				N/A
// Saturn	58,232,000	N/A				N/A					N/A				N/A
// Uranus	25,362,000	N/A				N/A					N/A				N/A
// Neptune	24,622,000	N/A				N/A					N/A				N/A
// Moon		1,737,100	Mons Huygens	10,800				9,000			Mare Imbrium