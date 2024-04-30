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
	MinDetail = 10
	MaxDetail = 15
	MeanRadius = 3389500
	MinElevation = -11000
	MaxElevation = 21900
	Tolerance = 50
	Exaggeration = 10
	// Scale = 1/MeanRadius
	InnerShellScale = 0.9
)

	// Planet: "Mercury"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 2439700
    // MinElevation: -5500
    // MaxElevation: 5500
    // Tolerance: 50
    // Exaggeration: 15
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Venus"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 6051800
    // MinElevation: -1000
    // MaxElevation: 11000
    // Tolerance: 50
    // Exaggeration: 15
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Test"
    // MinDetail: 100
    // MaxDetail: 200
    // MeanRadius: 1
    // MinElevation: -1
    // MaxElevation: 2
    // Tolerance: 50
    // Exaggeration: 20
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet = "Earth"
	// MinDetail = 9
	// MaxDetail = 12
	// MeanRadius = 6373934
	// MinElevation = -10900
	// MaxElevation = 8849
	// Tolerance = 50
	// Exaggeration = 15
	// //Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Mars"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 3389500
    // MinElevation: -11000
    // MaxElevation: 21900
    // Tolerance: 50
    // Exaggeration: 10
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Jupiter"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 69911000
    // MinElevation: -1000
    // MaxElevation: 1000
    // Tolerance: 50
    // Exaggeration: 30
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Saturn"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 58232000
    // MinElevation: -1000
    // MaxElevation: 1000
    // Tolerance: 50
    // Exaggeration: 30
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Uranus"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 25362000
    // MinElevation: -1000
    // MaxElevation: 1000
    // Tolerance: 50
    // Exaggeration: 30
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Neptune"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 24622000
    // MinElevation: -1000
    // MaxElevation: 1000
    // Tolerance: 50
    // Exaggeration: 30
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9

	// Planet: "Moon"
    // MinDetail: 10
    // MaxDetail: 15
    // MeanRadius: 1737100
    // MinElevation: -9000
    // MaxElevation: 10800
    // Tolerance: 50
    // Exaggeration: 15
	// Scale = 1/MeanRadius
	// InnerShellScale = 0.9



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

	done = timed("Reading input DEM")
	im, err := fauxgl.LoadImage(*inputFile)
	done()
	if err != nil {
		log.Fatal(err)
	}

	// Print all the variables
	output := fmt.Sprintf("\nPlanet: %s\nMinDetail: %d\nMaxDetail: %d\nMeanRadius: %d\nMinElevation: %d\nMaxElevation: %d\nTolerance: %d\nExaggeration: %d\nInnerShellScale: %f\n",
        Planet,
        MinDetail,
        MaxDetail,
        MeanRadius,
        MinElevation,
        MaxElevation,
        Tolerance,
        Exaggeration,
        //Scale,
        InnerShellScale,
    )

    fmt.Println(output)

	// Outer shell
	triangulator := demsphere.NewTriangulator(
	im, int(MinDetail), int(MaxDetail), float64(MeanRadius), float64(MinElevation), float64(MaxElevation), float64(Tolerance), float64(Exaggeration), float64(1/MeanRadius))
	done = timed("Generating positive mesh")
	triangles := triangulator.Triangulate()
	done()
	fmt.Println(fmt.Sprintf("Generated %v triangles for outer mesh", len(triangles)))

	// Inner shell
	im = imaging.Invert(im)
	triangulator = demsphere.NewTriangulator(
	im, int(MinDetail), int(MaxDetail), float64(MeanRadius), float64(MinElevation), float64(MaxElevation), float64(Tolerance), float64(Exaggeration), float64(1/MeanRadius)*float64(InnerShellScale))

	inner := triangulator.Triangulate()
	for i, t := range inner {
		inner[i] = demsphere.Triangle{t.C, t.B, t.A}
	}
	done = timed("Generating negative mesh")
	triangles = append(triangles, inner...)
	done()
	fmt.Println(fmt.Sprintf("Generated %v triangles for inner mesh", len(triangles)))


	// Filename
	filename := fmt.Sprintf("%s_%d_%d_%d_%d.stl", Planet, MinDetail, MaxDetail, Tolerance, Exaggeration)
	fmt.Println(fmt.Sprintf("Filename set to %s", filename))
	// Writing STL
	done = timed("Writing output")
	// Making a string for filename
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