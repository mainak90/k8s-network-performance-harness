package utils

import (
	"fmt"
	"strconv"
	"strings"
)

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
	elements := strings.Split(units, ",")
	p50iter, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(elements[0], "\r", ""), "\n", ""))
	p90iter, _ := strconv.Atoi(elements[1])
	p99iter, _ := strconv.Atoi(elements[2])
	transiter, _ := strconv.ParseFloat(elements[3], 64)
	return fmt.Sprintf("%s,%v,%v,%v,%v",pod, p50iter,p90iter,p99iter,transiter)
}

func IperfGenerateCSVLines(pod string, unit string) string {
	return fmt.Sprintf("%s,%s",pod,unit)
}

func NetPerfOutPut(input map[string][][]string, check string) {
	for key, val := range input {
		fmt.Println(fmt.Sprintf("##### \t Output from pod %s for netperf test in %s mode\t #####", key, check))
		fmt.Println(fmt.Sprintf("                       Podname                                  \t p50 \t 90 \t  p99 \t Trans/s"))
		fmt.Println(fmt.Sprintf("================================================================\t=====\t=====\t=====\t======="))
		for _, eachline := range val {
			podName := eachline[0]
			for i := 0; i < (64 - len(eachline[0])); i++ {
				podName += " "
			}
			fmt.Println(fmt.Sprintf("%s\t %s\t %s\t %s\t %s", podName, eachline[1], eachline[2], eachline[3], eachline[4]))
		}

	}
}

func IPerfOutPut(input map[string][][]string) {
	for key, val := range input {
		fmt.Println(fmt.Sprintf("##### \t Output from pod %s for iperf test \t #####", key))
		fmt.Println(fmt.Sprintf("                       Podname                                  \t Avg-Gbits/Sec"))
		fmt.Println(fmt.Sprintf("================================================================\t =============="))
		for _, eachline := range val {
			podName := eachline[0]
			for i := 0; i < (64 - len(eachline[0])); i++ {
				podName += " "
			}
			fmt.Println(fmt.Sprintf("%s\t %s", podName, eachline[1]))
		}

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