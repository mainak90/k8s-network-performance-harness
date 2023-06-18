package cmd

import (
	"fmt"
	"github.com/mainak90/perftest/pkg/generator"
	"github.com/mainak90/perftest/pkg/k8s"
	"github.com/mainak90/perftest/pkg/logging"
	"github.com/mainak90/perftest/pkg/utils"
	"k8s.io/client-go/kubernetes"
	"strings"
)


func RunNetPerf(client *kubernetes.Clientset, mode string, host string, time string, pod string) (string, error) {
	netperfcmd := fmt.Sprintf("netperf -H %s -l %s -P 1 -t %s -- -r 32,1024 -o P50_LATENCY,P90_LATENCY,P99_LATENCY,THROUGHPUT,THROUGHPUT_UNITS", host, time, mode)
	stdo, _, err := k8s.Exec(client, netperfcmd, "default", pod)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error executing exec command %s Error %s", netperfcmd, err.Error()))
		return "", fmt.Errorf(err.Error())
	}
	return utils.StripNetPerfOutputTCP(stdo), nil
}

func RunNetPerfSets(deploy k8s.Podlist, stdo bool, graph bool) {
	logging.InfoLog("Running pod-to-pod across cluster in TCP_RR mode")
	var netpTcpRR = make(map[string][][]string)
	var netpTcpCRR = make(map[string][][]string)
	var tcp_rr = []string{}
	var tcp_crr = []string{}
	for _, pod := range deploy.Pods {
		podIp, ok := deploy.PodIPMap[pod]
		if !ok {
			logging.WarnLog(fmt.Sprintf("Cannot fetch pod ip for pod %s skipping this one...", pod))
			continue
		}
		for _, ip := range deploy.PodIps {
			if ip != podIp {
				result, err := RunNetPerf(deploy.ClientSet, "TCP_RR", ip, "2", pod)
				if err != nil {
					logging.ErrLog(fmt.Sprintf("Encountered error while running netperf on pod %s %s", pod, err.Error()))
				}
				tcp_rr = append(tcp_rr, result)
				csvLine := utils.NetperfGenerateCSVLines(deploy.IPPodMap[ip],result)
				netpTcpRR[pod] = append(netpTcpRR[pod], strings.Split(csvLine, ","))
				generator.WriteCSV(pod, "TCP_RR",csvLine)
			}
		}
		if graph {
			generator.RenderChart(fmt.Sprintf("%s_TCP_RR.csv", pod), "TCP_RR")
		}
	}
	if stdo {
		utils.NetPerfOutPut(netpTcpRR, "TCP_RR")
	}
	logging.InfoLog("Running netperf pod-to-pod across cluster in TCP_CRR mode")
	for _, pod := range deploy.Pods {
		podIp, ok := deploy.PodIPMap[pod]
		if !ok {
			logging.WarnLog(fmt.Sprintf("Cannot fetch pod ip for pod %s skipping this one...", pod))
			continue
		}
		for _, ip := range deploy.PodIps {
			if ip != podIp {
				result, err := RunNetPerf(deploy.ClientSet, "TCP_CRR", ip, "2", pod)
				if err != nil {
					logging.ErrLog(fmt.Sprintf("Encountered error while running netperf on pod %s %s", pod, err.Error()))
				}
				tcp_crr = append(tcp_crr, result)
				csvLine := utils.NetperfGenerateCSVLines(deploy.IPPodMap[ip],result)
				netpTcpCRR[pod] = append(netpTcpCRR[pod], strings.Split(csvLine, ","))
				generator.WriteCSV(pod, "TCP_CRR",csvLine)
			}
		}
		if graph {
			generator.RenderChart(fmt.Sprintf("%s_TCP_CRR.csv", pod), "TCP_CRR")
		}
	}
	if stdo {
		utils.NetPerfOutPut(netpTcpRR, "TCP_CRR")
	}
}

func RunIPerf(client *kubernetes.Clientset, host string, time string, pod string) (string, error) {
	IPerfCmd := fmt.Sprintf("iperf -c %s -i 1 -t %s", host, time)
	stdo, _, err := k8s.Exec(client, IPerfCmd, "default", pod)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error executing exec command %s Error %s", IPerfCmd, err.Error()))
		return "", fmt.Errorf(err.Error())
	}
	return utils.StripIperfOutput(stdo), nil
}

func RunIperfSets(deploy k8s.Podlist,  stdo bool, graph bool) {
	logging.InfoLog("Running iperf pod-to-pod across cluster")
	var IperfMap = make(map[string][][]string)
	var iperf = []string{}
	for _, pod := range deploy.Pods {
		podIp, ok := deploy.PodIPMap[pod]
		if !ok {
			logging.WarnLog(fmt.Sprintf("Cannot fetch pod ip for pod %s skipping this one...", pod))
			continue
		}
		for _, ip := range deploy.PodIps {
			if ip != podIp {
				result, err := RunIPerf(deploy.ClientSet, ip, "2", pod)
				if err != nil {
					logging.ErrLog(fmt.Sprintf("Encountered error while running netperf on pod %s %s", pod, err.Error()))
				}
				iperf = append(iperf, result)
				csvLine := utils.IperfGenerateCSVLines(deploy.IPPodMap[ip],result)
				IperfMap[pod] = append(IperfMap[pod], strings.Split(csvLine, ","))
				generator.WriteCSV(pod, "iperf",csvLine)
			}
		}
		if graph {
			generator.RenderChart(fmt.Sprintf("%s_iperf.csv", pod), "iperf")
		}
	}
	if stdo {
		utils.IPerfOutPut(IperfMap)
	}
}

func RunCmds(deploy k8s.Podlist, cmds string, stdo bool, graph bool) {
	set := utils.MakeCmdSet(cmds)
	for _, cmd := range set {
		switch cmd {
		case "all":
			logging.InfoLog("Running all tests...")
			RunNetPerfSets(deploy, stdo, graph)
			RunIperfSets(deploy, stdo, graph)
		case "iperf":
			logging.InfoLog("Running IPerf tests...")
			RunIperfSets(deploy, stdo, graph)
		case "netperf":
			logging.InfoLog("Running NetPerf tests...")
			RunNetPerfSets(deploy, stdo, graph)
		default:
			logging.ErrLog(fmt.Sprintf("Command %s not valid. Valid commands to pass are all or iperf/netperf or netperf,iperf", cmd))
		}
	}
}