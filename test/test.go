package main

import (
	"testing"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/palette"
)

func TestHeatMapWithContour(t *testing.T) {
	if !*visualDebug {
		return
	}
	m := unitGrid{mat64.NewDense(3, 4, []float64{
		2, 1, 4, 3,
		6, 7, 2, 5,
		9, 10, 11, 12,
	})}
	h := NewHeatMap(m, palette.Heat(12, 1))

	levels := []float64{1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5, 10.5, 11.5}
	c := NewContour(m, levels, palette.Rainbow(10, palette.Blue, palette.Red, 1, 1, 1))
	c.LineStyles[0].Width *= 5

	plt, _ := plot.New()

	plt.Add(h)
	plt.Add(c)
	plt.Add(NewGlyphBoxes())

	plt.X.Padding = 0
	plt.Y.Padding = 0
	plt.X.Max = 3.5
	plt.Y.Max = 2.5
	plt.Save(7, 7, "heat.svg")
}
