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
	//outputFile = kingpin.Flag("output", "Output STL file to write.").Required().Short('o').String() // Currently not in use
)

// Define a struct to hold the variables for each celestial body
type CelestialVariables struct {
	MinDetail      int
	MaxDetail      int
	MeanRadius     int
	MinElevation   int
	MaxElevation   int
	Tolerance      int
	Exaggeration   int
	Scale          float64
	InnerShellScale int
}

func main() {
	// Define variables for all celestial bodies using a map
	celestialVariables := map[string]CelestialVariables{
		"mercury": {
			MinDetail:      15,
			MaxDetail:      30,
			MeanRadius:     2439700,
			MinElevation:   -5000,
			MaxElevation:   5500,
			Tolerance:      50,
			Exaggeration:   15,
			Scale:          1 / 2439700.0,
			InnerShellScale: 1,
		},
		"venus": {
			MinDetail:      20,
			MaxDetail:      35,
			MeanRadius:     6051800,
			MinElevation:   -11000,
			MaxElevation:   13000,
			Tolerance:      50,
			Exaggeration:   20,
			Scale:          1 / 6051800.0,
			InnerShellScale: 1,
		},
		"earth": {
			MinDetail:      15,
			MaxDetail:      30,
			MeanRadius:     6373934,
			MinElevation:   -10900,
			MaxElevation:   8849,
			Tolerance:      50,
			Exaggeration:   15,
			Scale:          1 / 6373934.0,
			InnerShellScale: 1,
		},
		"mars": {
			MinDetail:      25,
			MaxDetail:      40,
			MeanRadius:     3389500,
			MinElevation:   -8200,
			MaxElevation:   21200,
			Tolerance:      50,
			Exaggeration:   25,
			Scale:          1 / 3389500.0,
			InnerShellScale: 1,
		},
		"jupiter": {
			MinDetail:      30,
			MaxDetail:      45,
			MeanRadius:     69911000,
			MinElevation:   0, // No solid surface
			MaxElevation:   0, // No solid surface
			Tolerance:      50,
			Exaggeration:   30,
			Scale:          1 / 69911000.0,
			InnerShellScale: 1,
		},
		"saturn": {
			MinDetail:      30,
			MaxDetail:      45,
			MeanRadius:     58232000,
			MinElevation:   0, // No solid surface
			MaxElevation:   0, // No solid surface
			Tolerance:      50,
			Exaggeration:   30,
			Scale:          1 / 58232000.0,
			InnerShellScale: 1,
		},
		"uranus": {
			MinDetail:      30,
			MaxDetail:      45,
			MeanRadius:     25362000,
			MinElevation:   0, // No solid surface
			MaxElevation:   0, // No solid surface
			Tolerance:      50,
			Exaggeration:   30,
			Scale:          1 / 25362000.0,
			InnerShellScale: 1,
		},
		"neptune": {
			MinDetail:      30,
			MaxDetail:      45,
			MeanRadius:     24622000,
			MinElevation:   0, // No solid surface
			MaxElevation:   0, // No solid surface
			Tolerance:      50,
			Exaggeration:   30,
			Scale:          1 / 24622000.0,
			InnerShellScale: 1,
		},
		"moon": {
			MinDetail:      15,
			MaxDetail:      30,
			MeanRadius:     1737100,
			MinElevation:   -9000,
			MaxElevation:   10800,
			Tolerance:      50,
			Exaggeration:   15,
			Scale:          1 / 1737100.0,
			InnerShellScale: 1,
		},
	}

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
	// Choose planet
	planet := "earth"
	variables := celestialVariables[planet]
	// Make outer shell
	triangulator := demsphere.NewTriangulator(
		im, variables.MinDetail, variables.MaxDetail, variables.MeanRadius, variables.MinElevation, variables.MaxElevation, variables.Tolerance, variables.Exaggeration, variables.Scale)
	done = timed("generating outer mesh")
	triangles := triangulator.Triangulate()
	done()
	fmt.Println(len(triangles))
	// Making inner shell
	im = imaging.Invert(im) //Inverting image
	triangulator = demsphere.NewTriangulator(
		im, variables.MinDetail, variables.MaxDetail, variables.MeanRadius, variables.MinElevation, variables.MaxElevation, variables.Tolerance, variables.Exaggeration, variables.Scale*float64(variables.innerShellScale))
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
	filename := fmt.Sprintf("%s_%d_%d_%d_%d_%d.stl", planet, variables.MinDetail, variables.MaxDetail, variables.Tolerance, variables.Exaggeration, variables.Scale)
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