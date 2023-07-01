package utils

import (
	"bufio"
	"fmt"
	"github.com/mainak90/perftest/pkg/logging"
	"os"
	"strconv"
	"strings"
	"time"
)

func WriteToFile(filename string, line string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			logging.ErrLog(fmt.Sprintf("Error creating test output file %s", filename))
		}
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error opening test output %s Error %s", filename, err.Error()))
	}
	defer file.Close()
	if _, err := file.WriteString(line); err != nil {
		logging.ErrLog(fmt.Sprintf("Error writing to output %s Error %s", filename, err.Error()))
	}
}

func GetFileName(temp bool) string {
	if temp == true {
		return fmt.Sprintf("%v-tmp", time.Now().Unix())
	}
	return fmt.Sprintf("%v-out", time.Now().Unix())
}

func ReadAndPurge(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error opening final outputfile %s %s", filename, err.Error()))
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := os.Remove(filename); err != nil {
		logging.ErrLog(fmt.Sprintf("Error deleting tmp file %s %s", filename, err.Error()))
	}
}

func StripNetPerfOutputTCP(str string) string {
	lines := strings.Split(str, " ")
	lastLine := lines[len(lines)-1]
	return strings.ReplaceAll(lastLine, "Units", "")
}

func StripIperfOutput(str string) string {
	IperfSlice := strings.Split(str, "  ")
	iter1, _ := strconv.ParseFloat(strings.Split(IperfSlice[len(IperfSlice)-1], " ")[0], 64)
	iter2, _ := strconv.ParseFloat(strings.Split(IperfSlice[len(IperfSlice)-5], " ")[0], 64)
	iter3, _ := strconv.ParseFloat(strings.Split(IperfSlice[len(IperfSlice)-9], " ")[0], 64)
	return fmt.Sprintf("%v", (iter1+iter2+iter3)/3)
}

func NetperfGenerateCSVLines(pod string, units string) string {
	logging.WarnLog(fmt.Sprintf("Splitting units %s", units))
	elements := strings.Split(units, ",")
	p50iter, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(elements[0], "\r", ""), "\n", ""))
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Encountered errors while generating CSV lines for netperf for pod %s %v", pod, err.Error()))
	}
	p90iter, err := strconv.Atoi(elements[1])
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Encountered errors while generating CSV lines for netperf for pod %s %v", pod, err.Error()))
	}
	p99iter, err := strconv.Atoi(elements[2])
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Encountered errors while generating CSV lines for netperf for pod %s %v", pod, err.Error()))
	}
	transiter, err := strconv.ParseFloat(elements[3], 64)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Encountered errors while generating CSV lines for netperf for pod %s %v", pod, err.Error()))
	}
	return fmt.Sprintf("%s,%v,%v,%v,%v", pod, p50iter, p90iter, p99iter, transiter)
}

func IperfGenerateCSVLines(pod string, unit string) string {
	return fmt.Sprintf("%s,%s", pod, unit)
}

func NetPerfOutPut(input map[string][][]string, check string, out string) {
	for key, val := range input {
		WriteToFile(out, fmt.Sprintf("##### \t Output from pod %s for netperf test in %s mode\t #####\n", key, check))
		WriteToFile(out, fmt.Sprintf("                       Podname                                  \t p50 \t 90 \t  p99 \t Trans/s\n"))
		WriteToFile(out, fmt.Sprintf("================================================================\t=====\t=====\t=====\t=======\n"))
		for _, eachline := range val {
			podName := eachline[0]
			for i := 0; i < (64 - len(eachline[0])); i++ {
				podName += " "
			}
			WriteToFile(out, fmt.Sprintf("%s\t %s\t %s\t %s\t %s\n", podName, eachline[1], eachline[2], eachline[3], eachline[4]))
		}
		WriteToFile(out, fmt.Sprintf("\n"))
	}
}

func IPerfOutPut(input map[string][][]string, out string) {
	for key, val := range input {
		WriteToFile(out, fmt.Sprintf("##### \t Output from pod %s for iperf test \t #####\n", key))
		WriteToFile(out, fmt.Sprintf("                       Podname                                  \t Avg-Gbits/Sec\n"))
		WriteToFile(out, fmt.Sprintf("================================================================\t ==============\n"))
		for _, eachline := range val {
			podName := eachline[0]
			for i := 0; i < (64 - len(eachline[0])); i++ {
				podName += " "
			}
			WriteToFile(out, fmt.Sprintf("%s\t %s\n", podName, eachline[1]))
		}
		WriteToFile(out, fmt.Sprintf("\n"))
	}
}

func MakeCmdSet(cmds string) []string {
	if cmds == "all" {
		return []string{"all"}
	}
	if strings.Contains(cmds, ",") {
		return strings.Split(cmds, ",")
	} else {
		return []string{cmds}
	}

}
