package email

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

var ServiceChart *Charts

type Charts struct {
	action string
	option *Options
}
type Options struct {
	Title string
	X     []string
	Y     []string
	Data  [][]opts.LineData
}

func init() {
	ServiceChart = &Charts{}
}

type newOption func(options *Options)

func (c *Charts) New(action string, opts ...newOption) *Charts {
	option := &Options{
		Title: "",
		X:     nil,
		Y:     nil,
		Data:  nil,
	}
	for _, o := range opts {
		o(option)
	}
	c.action = action
	c.option = option
	return c
}

func (c *Charts) Line() {
	line := charts.NewLine()
	line.Initialization.AssetsHost = "https://cdn.bootcdn.net/ajax/libs/echarts/5.0.0-rc.1/"
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    c.option.Title,
		Subtitle: "-",
	}))
	line.SetXAxis(c.option.X)
	for _, val := range c.option.Data {
		line.AddSeries("", val)
	}
	//line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	create, _ := os.Create("line.html")
	err := line.Render(create)
	fmt.Println(err)
}
