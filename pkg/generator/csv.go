package generator

import (
	"encoding/csv"
	"fmt"
	"github.com/mainak90/perftest/pkg/logging"
	"os"
	"strings"
)

func NetPerfCSV(filename, line string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error creating netperf csv file %s", filename))
			return false
		}
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while opening netperf csv %s %s", filename, err.Error()))
			return false
		}
		writer := csv.NewWriter(file)
		writer.Write(strings.Split(fmt.Sprintf("%s,p50,p90,p99,trans/s", strings.Split(filename, "_")[0]), ","))
		writer.Flush()
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error opening file %s Error %s", filename, err.Error()))
		return false
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(strings.Split(line, ","))
	return true
}

func IperfCSV(filename, line string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error creating iperf csv file %s", filename))
			return false
		}
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error while opening iperf csv %s %s", filename, err.Error()))
			return false
		}
		writer := csv.NewWriter(file)
		writer.Write(strings.Split(fmt.Sprintf("%s,Gbits/sec", strings.Split(filename, "_")[0]), ","))
		writer.Flush()
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error opening file %s Error %s", filename, err.Error()))
		return false
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(strings.Split(line, ","))
	return true
}

func WriteCSV(filename, check, line string) {
	switch check {
	case "TCP_CRR" , "TCP_RR":
		if done := NetPerfCSV(fmt.Sprintf("%s_%s.csv",filename,check), line); !done {
			logging.ErrLog(fmt.Sprintf("Could not write to CSV for check %s for file %s", check, filename))
		}
	case "iperf":
		if done := IperfCSV(fmt.Sprintf("%s_%s.csv",filename,check), line); !done {
			logging.ErrLog(fmt.Sprintf("Could not write to CSV for check %s for file %s", check, filename))
		}
	default:
		if done := NetPerfCSV(fmt.Sprintf("%s_%s.csv",filename,check), line); !done {
			logging.ErrLog(fmt.Sprintf("Could not write to CSV for check %s for file %s", check, filename))
		}
	}
}