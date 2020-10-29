package main

import (
	"bytes"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func runningDistanceToDayKm(ri []RunningItem, startDate time.Time, endDate time.Time) []float64 {
	var maxDate time.Time
	riMap := map[time.Time]float64{}
	for _, r := range ri {
		y, m, d := r.createdAt.Date()
		cAt := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

		if _, ok := riMap[cAt]; !ok {
			riMap[cAt] = 0
		}

		riMap[cAt] += r.distance

		if cAt.After(maxDate) {
			maxDate = cAt
		}
	}

	res := make([]float64, 0)

	fTime := startDate
	for !fTime.After(endDate) && !fTime.After(maxDate) {
		var val = 0.0
		if v, ok := riMap[fTime]; ok {
			val = v
		}

		res = append(res, val)

		fTime.Add(time.Hour * 24)
	}

	return res
}

func drawChart(goalKm, totalDays uint, daysKm []float64) (*bytes.Buffer, error) {
	defaultChartData := chart.ContinuousSeries{
		Name: "Chart name",
		Style: chart.Style{
			Show: true,
			Padding: chart.Box{
				Top:    200,
				Left:   200,
				Right:  200,
				Bottom: 200,
				IsSet:  false,
			},
			StrokeWidth:         3,
			StrokeColor:         drawing.Color{},
			StrokeDashArray:     nil,
			DotColor:            drawing.Color{},
			DotWidth:            0,
			DotWidthProvider:    nil,
			DotColorProvider:    nil,
			FillColor:           drawing.Color{},
			FontSize:            0,
			FontColor:           drawing.Color{},
			Font:                nil,
			TextHorizontalAlign: 0,
			TextVerticalAlign:   0,
			TextWrap:            0,
			TextLineSpacing:     0,
			TextRotationDegrees: 0,
		},
		YAxis:           0,
		XValueFormatter: nil,
		YValueFormatter: nil,
		XValues:         nil,
		YValues:         nil,
	}

	chart1 := defaultChartData
	chart1.YValues = []float64{float64(goalKm), float64(goalKm)}
	chart1.XValues = []float64{0, float64(totalDays)}

	xVals := make([]float64, 0, totalDays)
	yVals := make([]float64, 0, totalDays)

	xVals = append(xVals, 0)
	yVals = append(yVals, daysKm[0])

	for i := 1; i < len(daysKm); i++ {
		xVals = append(xVals, float64(i))
		yVals = append(yVals, yVals[i-1]+daysKm[i])
	}

	chart2 := defaultChartData
	chart2.YValues = yVals
	chart2.XValues = xVals

	graph := chart.Chart{
		Series: []chart.Series{
			chart1,
			chart2,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, err
}
