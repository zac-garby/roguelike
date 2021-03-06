package lib

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	roadColour = color.Black
	roomColour = color.Black

	jobs       chan job
	jobResults chan jobResult
)

func init() {
	jobs = make(chan job, Conf.NumThreads)
	jobResults = make(chan jobResult, Conf.NumThreads)
}

// DecodeImageIntoMap takes an image.Image object and converts its
// pixels into the tiles in a Map.
func DecodeImageIntoMap(img image.Image) *Map {
	var (
		w = img.Bounds().Size().X
		h = img.Bounds().Size().Y
	)

	m := &Map{
		Depth: 0,
		Tiles: make([][]Tile, h),
	}

	for y := 0; y < h; y++ {
		m.Tiles[y] = make([]Tile, w)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var (
				tile   Tile
				colour = img.At(x, y)
			)

			if coloursEqual(colour, color.Transparent) {
				tile = &OutsideTile{}
			} else {
				tile = &FloorTile{}
			}

			m.Set(x, y, tile)
		}
	}

	return m
}

func coloursEqual(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return (a1 == 0 && a2 == 0) || r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// An edge represents an edge in a graph from one node to another. The nodes
// are stored as the indices of the nodes
type edge struct {
	from, to int
}

// A job is the struct on which a worker will work. A worker will find the
// minimum value in the graph between startCol and endCol, where the row is
// not deleted and the column is labelled. This will be sent to jobResults
// in a jobResult struct.
type job struct {
	labelled, deleted []int
	startCol, endCol  int
	graph             [][]int
}

type jobResult struct {
	minEdge *edge
	minDist int
}

// SpawnWorkers spawns a bunch of worker goroutines to speed up Prim's algorithm
func SpawnWorkers() {
	for i := 0; i < Conf.NumThreads; i++ {
		go worker()
	}
}

// worker will continuously take jobs from the global jobs channel and
// process them, sending the result to jobResults.
func worker() {
	for {
		var (
			job      = <-jobs
			labelled = job.labelled
			deleted  = job.deleted
			graph    = job.graph
			startCol = job.startCol
			endCol   = job.endCol

			minEdge *edge
			minDist = 9223372036854775807
		)

		for f := 0; f < len(graph); f++ {
			if contains(deleted, f) {
				continue
			}

			for t := 0; t < len(graph); t++ {
				if t < startCol || t > endCol || !contains(labelled, t) {
					continue
				}

				dist := graph[f][t]

				if dist >= 0 && dist < minDist {
					minEdge = &edge{
						from: f,
						to:   t,
					}

					minDist = dist
				}
			}
		}

		jobResults <- jobResult{
			minEdge: minEdge,
			minDist: minDist,
		}
	}
}

// aggregateJobResults waits until there are enough jobResults in the jobResults
// channel then takes them all and finds the minimum edge in the *entire* graph.
func aggregateJobResults() *edge {
	results := make([]jobResult, Conf.NumThreads)

	for i := 0; i < Conf.NumThreads; i++ {
		results[i] = <-jobResults
	}

	var (
		minEdge *edge
		minDist = 9223372036854775807
	)

	for _, r := range results {
		if r.minDist < minDist {
			minEdge = r.minEdge
			minDist = r.minDist
		}
	}

	return minEdge
}

// initiateJobs takes a graph and slices it up into threadCount portions. Well,
// not quite -- it passes the entire graph to each job, but tells each worker
// to only work between two bounds.
func initiateJobs(graph [][]int, labelled, deleted []int) {
	d := int(math.Floor(float64(len(graph)) / float64(Conf.NumThreads)))

	for i := 0; i < Conf.NumThreads; i++ {
		start := i * d
		end := start + d - 1

		jobs <- job{
			graph:    graph,
			labelled: labelled,
			deleted:  deleted,
			startCol: start,
			endCol:   end,
		}
	}
}

// MakeMap generates a new map
func MakeMap(depth int) *Map {
	rand.Seed(time.Now().UnixNano())

	var (
		img = image.NewRGBA(image.Rect(0, 0, Conf.MapWidth, Conf.MapHeight))
		gc  = draw2dimg.NewGraphicContext(img)

		points = generatePoints()
		graph  = makeGraph(points)
		mst    = findMST(graph)
	)

	gc.BeginPath()
	for f := 0; f < len(graph); f++ {
		for t := 0; t < len(graph); t++ {
			if mst[f][t] {
				var (
					from = points[f]
					to   = points[t]
				)

				gc.MoveTo(float64(from.X), float64(from.Y))
				gc.LineTo(float64(to.X), float64(to.Y))
			}
		}
	}
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineWidth(float64(Conf.RoadWidth))
	gc.SetStrokeColor(roadColour)
	gc.Close()
	gc.Stroke()

	gc.SetFillColor(roadColour)
	for f := 0; f < len(graph); f++ {
		for t := 0; t < len(graph); t++ {
			if mst[f][t] {
				var (
					from = points[f]
					to   = points[t]
				)

				gc.MoveTo(float64(to.X), float64(to.Y))
				draw2dkit.Circle(gc, float64(to.X), float64(to.Y), float64(Conf.RoadWidth)/2)
				gc.Fill()

				gc.MoveTo(float64(from.X), float64(from.Y))
				draw2dkit.Circle(gc, float64(from.X), float64(from.Y), float64(Conf.RoadWidth)/2)
				gc.Fill()
			}
		}
	}

	gc.BeginPath()
	for p := 0; p < len(points); p++ {
		var (
			point = points[p]
			conn  = numConnected(mst, p)
		)

		if rand.Float64() <= roomProbability(conn-1) {
			radius := float64(Conf.RoomWidth) + (rand.Float64()-0.5)*Conf.RoomWidthVariance
			square(gc, point, radius)
		}
	}
	gc.Close()
	gc.SetFillColor(roomColour)
	gc.Fill()

	m := DecodeImageIntoMap(img)
	m.Depth = depth
	m.Postprocess()
	return m
}

func square(gc draw2d.PathBuilder, center image.Point, radius float64) {
	gc.MoveTo(float64(center.X)-radius, float64(center.Y)-radius)
	gc.LineTo(float64(center.X)-radius, float64(center.Y)+radius)
	gc.LineTo(float64(center.X)+radius, float64(center.Y)+radius)
	gc.LineTo(float64(center.X)+radius, float64(center.Y)-radius)
}

func generatePoints() []image.Point {
	points := []image.Point{}

	for i := Conf.GridSpacing; i < Conf.MapWidth; i += Conf.GridSpacing {
		for j := Conf.GridSpacing; j < Conf.MapHeight; j += Conf.GridSpacing {
			if rand.Float64() < Conf.NodeChance {
				points = append(points, image.Point{
					X: i,
					Y: j,
				})
			}
		}
	}

	return points
}

// makeGraph creates a graph in a distance matrix form, but instead of
// the distance it stores the distance squared.
func makeGraph(points []image.Point) [][]int {
	graph := [][]int{}

	for f := 0; f < len(points); f++ {
		row := []int{}

		for t := 0; t < len(points); t++ {
			if t == f {
				row = append(row, -1)
				continue
			}

			from := points[f]
			to := points[t]

			// Skip this pair if they are diagonal (i.e. not same X or Y)
			if !(from.X == to.X || from.Y == to.Y) {
				row = append(row, -1)
				continue
			}

			diff := to.Sub(from)
			row = append(row, diff.X*diff.X+diff.Y*diff.Y)
		}

		graph = append(graph, row)
	}

	return graph
}

// findMST finds the minimum spanning tree in a distance matrix, giving an output
// as an adjacency matrix.
func findMST(graph [][]int) [][]bool {
	var (
		labelled = []int{0}
		deleted  = []int{}
		output   = make([][]bool, len(graph))
	)

	for i := 0; i < len(graph); i++ {
		output[i] = make([]bool, len(graph))
	}

	for len(deleted) < len(graph) {
		min := findMinimum(graph, labelled, deleted)
		output[min.from][min.to] = true
		output[min.to][min.from] = true

		deleted = append(deleted, min.from)
		labelled = append(labelled, min.from)
	}

	return output
}

// findMinimum finds the minimum possible edge in a graph
func findMinimum(graph [][]int, labelled, deleted []int) *edge {
	initiateJobs(graph, labelled, deleted)
	return aggregateJobResults()
}

// numConnected returns the number of edges connected to a node in an
// adjacency matrix
func numConnected(graph [][]bool, index int) int {
	var (
		row   = graph[index]
		total = 0
	)

	for c := 0; c < len(graph); c++ {
		if c == index {
			continue
		}

		if row[c] {
			total++
		}
	}

	return total
}

func roomProbability(n int) float64 {
	return math.Exp(Conf.RoomProbabilityCoefficient * float64(n))
}

func contains(set []int, i int) bool {
	for _, j := range set {
		if j == i {
			return true
		}
	}

	return false
}
