package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/wcharczuk/go-chart/drawing"

	"github.com/wcharczuk/go-chart"
)

func runningDistanceToDayKm(ri []RunningItem, startDate time.Time, endDate time.Time) []float64 {
	now := time.Now()
	riMap := map[time.Time]float64{}
	for _, r := range ri {
		y, m, d := r.createdAt.Date()
		cAt := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

		if _, ok := riMap[cAt]; !ok {
			riMap[cAt] = 0
		}

		riMap[cAt] += r.distance
	}

	res := make([]float64, 0)

	fTime := startDate
	for !fTime.After(endDate) && !fTime.After(now) {
		var val = 0.0
		if v, ok := riMap[fTime]; ok {
			val = v
		}

		res = append(res, val)

		fTime = fTime.Add(time.Hour * 24)
	}

	return res
}

func drawSuccessPredChard(goal, daysCount int, avgPerDay float64, currentDay float64) (*bytes.Buffer, error) {
	defaultChartData := chart.ContinuousSeries{
		Name: "Chart name",
	}

	chart3 := defaultChartData

	xTicks := make([]chart.Tick, 0, daysCount)
	for i := 0.0; i <= float64(daysCount); i += 10 {
		xTicks = append(xTicks, chart.Tick{
			Value: i,
			Label: fmt.Sprint(i),
		})
	}
	if xTicks[len(xTicks)-1].Value != float64(daysCount) {
		xTicks = append(xTicks, chart.Tick{
			Value: float64(daysCount),
			Label: fmt.Sprint(daysCount),
		})
	}

	ceilAvg := math.Ceil((float64(goal) - (1 * avgPerDay)) / (float64(daysCount) - 1))
	yTicks := make([]chart.Tick, 0, daysCount)
	for i := 0.0; i <= ceilAvg; i += 1 {
		yTicks = append(yTicks, chart.Tick{
			Value: i,
			Label: fmt.Sprint(i),
		})
	}

	xVals := make([]float64, daysCount)
	yVals := make([]float64, daysCount)

	for i := 0.0; i <= float64(daysCount); i += 1 {
		var res float64
		if float64(daysCount)-i != 0 {
			res = (float64(goal) - (i * avgPerDay)) / (float64(daysCount) - i)
		} else {
			res = 0
		}

		if res < 0 {
			res = 0
		}

		if i == currentDay {
			chart3.YValues = []float64{res, res}
			chart3.XValues = []float64{i, i}
			log.Print(chart3.XValues, chart3.YValues)
			chart3.Style.Show = true
			chart3.Style.DotColor = chart.ColorRed
			chart3.Style.DotWidth = 3
		}

		xVals = append(xVals, i)
		yVals = append(yVals, res)
	}

	chart2 := defaultChartData
	chart2.YValues = yVals
	chart2.XValues = xVals
	chart2.Style.Show = true
	chart2.Style.StrokeWidth = 2
	chart2.Style.DotWidth = 0
	chart2.Style.StrokeColor = chart.ColorBlue
	chart2.Style.DotColor = chart.ColorBlue

	graph := chart.Chart{
		Title: "Your running formula prediction",
		TitleStyle: chart.Style{
			Show:     true,
			FontSize: 20,
		},
		ColorPalette: nil,
		Width:        0,
		Height:       0,
		DPI:          0,
		Background:   chart.Style{},
		Canvas:       chart.Style{},
		XAxis: chart.XAxis{
			Ticks: xTicks,
			Style: chart.Style{
				Show:     true,
				FontSize: 5,
			},
		},
		YAxis: chart.YAxis{
			Ticks: yTicks,
			Style: chart.Style{
				Show:     true,
				FontSize: 5,
			},
		},
		YAxisSecondary: chart.YAxis{},
		Font:           nil,
		Series:         []chart.Series{chart2, chart3},
		Elements:       nil,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, err
}

func drawChart(goalKm, totalDays uint, daysKm [][]float64) (*bytes.Buffer, error) {
	defaultChartData := chart.ContinuousSeries{
		Name: "Chart name",
	}

	chart1 := defaultChartData
	chart1.YValues = []float64{float64(goalKm), float64(goalKm)}
	chart1.XValues = []float64{0, float64(totalDays)}
	chart1.Style.StrokeColor = chart.ColorRed
	chart1.Style.Show = true

	colors := map[int]drawing.Color{
		0: chart.ColorBlue,
		1: chart.ColorGreen,
		2: chart.ColorCyan,
		3: chart.ColorOrange,
	}

	resCharts := make([]chart.Series, len(daysKm))
	for j, c := range daysKm {
		xVals := make([]float64, 0, totalDays)
		yVals := make([]float64, 0, totalDays)

		xVals = append(xVals, 0)
		yVals = append(yVals, c[0])

		for i := 1; i < len(c); i++ {
			xVals = append(xVals, float64(i))
			yVals = append(yVals, yVals[i-1]+c[i])
		}

		chart2 := defaultChartData
		chart2.YValues = yVals
		chart2.XValues = xVals
		chart2.Style.Show = true
		chart2.Style.StrokeWidth = 2
		chart2.Style.DotWidth = 3

		if v, ok := colors[j]; ok {
			chart2.Style.StrokeColor = v
			chart2.Style.DotColor = v
		}

		resCharts[j] = chart2
	}

	graph := chart.Chart{
		Title: "Your running progress",
		TitleStyle: chart.Style{
			Show:     true,
			FontSize: 20,
		},
		ColorPalette: nil,
		Width:        0,
		Height:       0,
		DPI:          0,
		Background:   chart.Style{},
		Canvas:       chart.Style{},
		XAxis: chart.XAxis{
			Name: "Days from start",
			Style: chart.Style{
				Show:     true,
				FontSize: 5,
			},
		},
		YAxis: chart.YAxis{
			Name: "Count of KM",
			Style: chart.Style{
				Show:     true,
				FontSize: 5,
			},
		},
		YAxisSecondary: chart.YAxis{},
		Font:           nil,
		Series: []chart.Series{
			chart1,
		},
		Elements: nil,
	}

	for _, c := range resCharts {
		graph.Series = append(graph.Series, c)
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, err
}
