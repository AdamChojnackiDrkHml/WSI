package main

import (
	//	"fmt"

	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	mnist "github.com/moverest/mnist"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	//	clusters "github.com/muesli/clusters"
	//	kmeans "github.com/muesli/kmeans"
)

const k = 10
const iterations = 5

type plottable struct {
	grid       [][]float64
	N          int
	M          int
	resolution float64
	minX       float64
	minY       float64
}

type observation struct {
	data     []float64
	label    int
	centroid int
}

func (p plottable) Dims() (c, r int) {
	return p.N, p.M
}
func (p plottable) X(c int) float64 {
	return p.minX + float64(c)*p.resolution
}
func (p plottable) Y(r int) float64 {
	return p.minY + float64(r)*p.resolution
}
func (p plottable) Z(c, r int) float64 {
	return p.grid[c][r]
}

func main() {
	xd, _ := mnist.LoadImageFile(mnist.TrainingImageFileName)
	labels, _ := mnist.LoadLabelFile(mnist.TrainingLabelFileName)

	var data [60000][]float64
	observations := make([]observation, 0)

	for i := 0; i < 60000; i++ {
		data[i] = make([]float64, 784)
		for j := 0; j < 784; j++ {
			data[i][j] = float64(xd[i][j]) / 255.0
		}
	}

	for i, n := range data {
		observations = append(observations, observation{n, int(labels[i]), 0})
	}
	centroids := createCentroids(data)

	for i := 0; i < iterations; i++ {
		observations = assignAll(observations, centroids)
		centroids = calculateNewCentroids(observations)
	}
	observations = assignAll(observations, centroids)

	//fmt.Print(calcInertia(observations, centroids))
	for i := 0; i < k; i++ {
		var centroid [28][28]float64
		for j := 0; j < 28*28; j++ {
			centroid[j%28][27-j/28] = centroids[i][j]
		}
		drawHeatMap(centroid, i)
	}

	printTableOfAss(observations)

	testino, _ := mnist.LoadImageFile("test/" + mnist.TrainingImageFileName)
	labelino, _ := mnist.LoadLabelFile("test/" + mnist.TrainingLabelFileName)
	var testo [30][]float64
	testovations := make([]observation, 0)

	for i := 0; i < 30; i++ {
		testo[i] = make([]float64, 784)
		for j := 0; j < 784; j++ {
			testo[i][j] = float64(testino[i][j]) / 255.0
		}
	}

	for i, n := range testo {
		testovations = append(testovations, observation{n, int(labelino[i]), 0})
	}

	testovations = assignAll(testovations, centroids)

	printTableOfAss(testovations)

}

func printTableOfAss(observations []observation) {

	var table [k][10]int

	for _, n := range observations {
		table[n.centroid][n.label]++
	}

	fmt.Print("x\t")
	for i := 0; i < k; i++ {
		fmt.Print(strconv.Itoa(i) + "\t")
	}
	fmt.Println()
	for i, n := range table {
		fmt.Print("[" + strconv.Itoa(i) + "]\t")
		for _, m := range n {
			fmt.Print(strconv.Itoa(m) + "\t")
		}
		fmt.Println()
	}

}

func createCentroids2(data [60000][]float64) [][]float64 {
	centroids := make([][]float64, k)
	rand.Seed(time.Now().UnixNano())
	for i := range centroids {
		centroids[i] = make([]float64, 784)
		pic := rand.Int() % 60000
		for j := range centroids[i] {
			centroids[i][j] = data[pic][j]

		}
	}

	return centroids
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createCentroids(data [60000][]float64) [][]float64 {
	centroids := make([][]float64, k)
	rand.Seed(time.Now().UnixNano())
	centroidsIndex := make([]int, 0)

	for len(centroidsIndex) < k {
		index := rand.Int() % 60000
		if !contains(centroidsIndex, index) {
			centroidsIndex = append(centroidsIndex, index)
		}
	}

	for i := range centroids {
		centroids[i] = make([]float64, 784)
		for j := range centroids[i] {
			centroids[i][j] = data[centroidsIndex[i]][j]

		}
	}
	return centroids
}

func assignToCentroid(X observation, centroids [][]float64) observation {
	minDist := math.MaxFloat64
	centroid := 0

	for i, n := range centroids {
		dist := calcDistance(X.data, n)
		if dist < minDist {
			minDist = dist
			centroid = i
		}
	}

	X.centroid = centroid
	return X
}

func assignAll(observations []observation, centroids [][]float64) []observation {
	for i := range observations {
		observations[i] = assignToCentroid(observations[i], centroids)
	}

	return observations
}

func calculateNewCentroids(observations []observation) [][]float64 {
	centroids := make([][]float64, k)
	var assignedCounter [k]int

	for i := range centroids {

		centroids[i] = make([]float64, 784)

	}

	for _, n := range observations {
		assignedCounter[n.centroid]++

		for j, m := range n.data {
			centroids[n.centroid][j] += m
		}
	}

	for i := range centroids {
		for j := range centroids[i] {
			centroids[i][j] /= float64(assignedCounter[i])
		}

	}

	return centroids
}

func calcDistance(X []float64, Y []float64) float64 {
	distance := 0.0
	for i, _ := range X {
		distance += math.Pow(X[i]-Y[i], 2)
	}

	return math.Sqrt(distance)
}

func calcInertia(observations []observation, centroids [][]float64) float64 {
	inertia := 0.0
	for i := range observations {
		inertia += calcDistance(observations[i].data, centroids[observations[i].centroid])
	}

	return inertia
}

func drawHeatMap(dataset [28][28]float64, label int) {

	a := make([][]float64, 28)
	for i := range a {
		a[i] = make([]float64, 28)
	}

	for i := range a {
		for j := range a[i] {
			a[i][j] = dataset[i][j]
		}
	}

	plotData := plottable{
		grid:       a,
		N:          28,
		M:          28,
		minX:       -0.5,
		minY:       42.0,
		resolution: 0.5,
	}
	pal := moreland.SmoothBlueRed().Palette(255)
	hm := plotter.NewHeatMap(plotData, pal)

	p := plot.New()

	p.Title.Text = strconv.Itoa(label)
	p.Add(hm)

	p.Save(300, 300, "data/centroidPt"+strconv.Itoa(label)+".png")
}
