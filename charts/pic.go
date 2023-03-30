package charts

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/WinnieVenice/packer/idl"
	"github.com/WinnieVenice/packer/util"

	"github.com/wcharczuk/go-chart/v2"
)

func mkChartFile(pic *chart.Chart, name string) (filePath, fileName string) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	fileName = fmt.Sprintf("%s-%s%s.png", name, t, strconv.FormatUint(util.GetGoId(), 10))
	filePath = fmt.Sprintf("%s/pic/%s", util.GetCurrentAbPath(), fileName)
	fmt.Printf("DrawRecord, create filePath = (%+v), t = (%+v), goId = (%+v)\n", fileName, t, util.GetGoId())
	f, _ := os.Create(filePath)
	defer f.Close()
	pic.Render(chart.PNG, f)
	return
}

func mkBarChartFile(pic *chart.BarChart, name string) (filePath, fileName string) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	fileName = fmt.Sprintf("%s-%s%s.png", name, t, strconv.FormatUint(util.GetGoId(), 10))
	filePath = fmt.Sprintf("%s/pic/%s", util.GetCurrentAbPath(), fileName)
	fmt.Printf("DrawRecord, create filePath = (%+v), t = (%+v), goId = (%+v)\n", fileName, t, util.GetGoId())
	f, _ := os.Create(filePath)
	defer f.Close()
	pic.Render(chart.PNG, f)
	return
}

func DrawRecord(record []*idl.UserContestRecord_Record) (string, string) {
	pic := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: func(v interface{}) string {
				if val, ok := v.(float64); ok {
					t := int64(val)
					return time.Unix(t, 0).Format("2006-01-02")
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				if val, ok := v.(float64); ok {
					rate := int64(val)
					return strconv.FormatInt(rate, 10)
				}
				return ""
			},
		},
		Series: []chart.Series{},
	}
	line := chart.ContinuousSeries{
		XValues: []float64{},
		YValues: []float64{},
	}
	for _, c := range record {
		line.XValues = append(line.XValues, float64(c.GetTimestamp()))
		line.YValues = append(line.YValues, float64(c.GetRating()))
	}
	labelLine := chart.AnnotationSeries{
		Annotations: []chart.Value2{},
	}

	for _, id := range []int{0, len(record) - 1} {
		c := record[id]
		labelLine.Annotations = append(labelLine.Annotations, chart.Value2{
			XValue: float64(c.GetTimestamp()),
			YValue: float64(c.GetRating()),
			Label:  strconv.FormatInt(int64(c.GetRating()), 10),
		})
	}
	maxRate := &chart.MaxSeries{
		Style: chart.Style{
			StrokeColor: chart.ColorAlternateGray,
		},
		InnerSeries: line,
	}
	minRate := &chart.MinSeries{
		Style: chart.Style{
			StrokeColor: chart.ColorAlternateGray,
		},
		InnerSeries: line,
	}
	pic.Series = append(pic.Series, line)
	pic.Series = append(pic.Series, labelLine)
	pic.Series = append(pic.Series, maxRate)
	pic.Series = append(pic.Series, minRate)
	pic.Series = append(pic.Series, chart.LastValueAnnotationSeries(minRate))
	pic.Series = append(pic.Series, chart.LastValueAnnotationSeries(maxRate))

	return mkChartFile(&pic, "DrawRecord")
}

func DrawBindUserDailyDiff(user map[string]int64) (string, string) {
	pic := chart.BarChart{
		Title: fmt.Sprintf("%s daily partice status", time.Now().Add(-24*time.Hour).Format("2006-01-02")),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		BarWidth: 10,
	}

	var bars []chart.Value
	for id, diff := range user {
		bars = append(bars, chart.Value{
			Value: float64(diff),
			Label: id,
		})
	}
	pic.Bars = bars

	return mkBarChartFile(&pic, "DrawBindUserDailyDiff")
}
