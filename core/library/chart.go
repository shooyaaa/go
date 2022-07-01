package library

import (
	"fmt"
	"github.com/shooyaaa/log"
	"math/rand"
	"os"
	"time"

	"github.com/shooyaaa/config"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

var colors = []drawing.Color{
	drawing.ColorRed,
	drawing.ColorBlack,
	drawing.ColorGreen,
	drawing.ColorBlue,
}

type ChartData struct {
	X    []time.Time
	Y    []float64
	Name string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomFileName(ext string) string {
	b := make([]rune, 26)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letters[rand.Intn(len(letters))]
	}
	return config.TmpDir + "/" + string(b) + "." + ext
}

func Line(cd []*ChartData) string {
	series := make([]chart.Series, 0)
	for index, data := range cd {
		series = append(series, chart.TimeSeries{
			XValues: data.X,
			YValues: data.Y,
			Style: chart.Style{
				StrokeColor: colors[index],
				Show:        true,
			},
		})
	}
	graph := chart.Chart{
		Series: series,
		XAxis: chart.XAxis{
			Name: "Time",
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Name: "Price",
			Style: chart.Style{
				Show: true,
			},
		},
	}
	imageFile := RandomFileName("png")
	cwd, _ := os.Getwd()
	writer, err := os.Create(cwd + "/" + imageFile)
	if err != nil {
		log.DebugF("create file error %v")
	}
	err = graph.Render(chart.PNG, writer)
	if err != nil {
		fmt.Println("Error: render chart", err, cwd)
	}
	return imageFile
}
