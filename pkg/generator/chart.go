package generator

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/mainak90/perftest/pkg/logging"
	"os"
	"strconv"
	"strings"
)

func readCSV(file string) ([]string, error) {
	var fileLines []string
	readFile, err := os.Open(file)
	if err != nil {
		logging.ErrLog("Error opening file %s Error %s", file, err.Error())
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan(){
		fileLines = append(fileLines, fileScanner.Text())
	}
	readFile.Close()
	return fileLines[1:], nil
}

func SplitIPerfGraphItems(filename string) ([]string, []opts.BarData) {
	logging.InfoLog(fmt.Sprintf("Generating graphs for %s", filename))
	lines, err := readCSV(filename)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error while generating graph for %s %s", filename, err.Error()))
	}
	labels := make([]string, len(lines))
	gBits := make([]opts.BarData, 0)
	label := opts.Label{Show: true}
	for _, line := range lines {
		labels = append(labels, strings.Split(line, ",")[0])
		gBit, err := strconv.ParseFloat(strings.Split(line, ",")[1], 64)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while converting to float %v", strings.Split(line, ",")[1], err.Error()))
		}
		gBits = append(gBits, opts.BarData{Value: gBit, Label: &label})
	}
	return labels[1:], gBits
}

func SplitNetPerfGraphItems(filename string) ([]string, []opts.BarData, []opts.BarData, []opts.BarData, []opts.BarData) {
	logging.InfoLog(fmt.Sprintf("Generating graphs for %s", filename))
	lines, err := readCSV(filename)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error while generating graph for %s %s", filename, err.Error()))
	}
	labels := make([]string, len(lines))
	p50s := make([]opts.BarData, 0)
	p90s := make([]opts.BarData, 0)
	p99s := make([]opts.BarData, 0)
	transs := make([]opts.BarData, 0)
	label := opts.Label{Show: true}
	for _, line := range lines {
		labels = append(labels, strings.Split(line, ",")[0])
		p50, err := strconv.ParseFloat(strings.Split(line, ",")[1], 64)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while converting to float %v", strings.Split(line, ",")[1], err.Error()))
		}
		p50s = append(p50s, opts.BarData{Value: p50, Label: &label})
		p90, err := strconv.ParseFloat(strings.Split(line, ",")[2], 64)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while converting to float %v", strings.Split(line, ",")[2], err.Error()))
		}
		p90s = append(p90s, opts.BarData{Value: p90, Label: &label})
		p99, err := strconv.ParseFloat(strings.Split(line, ",")[3], 64)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while converting to float %v", strings.Split(line, ",")[3], err.Error()))
		}
		p99s = append(p99s, opts.BarData{Value: p99, Label: &label})
		trans, err := strconv.ParseFloat(strings.Split(line, ",")[4], 64)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while converting to float %v", strings.Split(line, ",")[4], err.Error()))
		}
		transs = append(transs, opts.BarData{Value: trans, Label: &label})
	}
	return labels[1:], p50s, p90s, p99s, transs
}

func renderNetPerfChart(filename string) {
	labels, p50s, p90s, p99s, trans_s := SplitNetPerfGraphItems(filename)
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    strings.Split(filename, ".")[0],
	}), charts.WithInitializationOpts(opts.Initialization{Height: "1024px", Width: "2048px"}))
	bar.SetXAxis(labels).
		AddSeries("p50", p50s).
		AddSeries("p90", p90s).
		AddSeries("p99", p99s).
		AddSeries("trans/s", trans_s)
	f, _ := os.Create(fmt.Sprintf("%s.html", strings.Split(filename, ".")[0]))
	bar.Render(f)
}

func renderIPerfChart(filename string) {
	labels, gBits := SplitIPerfGraphItems(filename)
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    strings.Split(filename, ".")[0],
	}), charts.WithInitializationOpts(opts.Initialization{Height: "1024px", Width: "2048px"}))
	bar.SetXAxis(labels).
		AddSeries("GBits/S", gBits)
	f, _ := os.Create(fmt.Sprintf("%s.html", strings.Split(filename, ".")[0]))
	bar.Render(f)
}

func RenderChart(filename string, check string) {
	switch check {
	case "TCP_RR":
		renderNetPerfChart(filename)
	case "TCP_CRR":
		renderNetPerfChart(filename)
	case "iperf":
		renderIPerfChart(filename)
	default:
		renderNetPerfChart(filename)
	}
}



