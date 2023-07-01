package cmd

import (
	"fmt"
	"github.com/mainak90/perftest/pkg/generator"
	"github.com/mainak90/perftest/pkg/k8s"
	"github.com/mainak90/perftest/pkg/logging"
	"github.com/mainak90/perftest/pkg/utils"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
)

func RunNetPerf(client *kubernetes.Clientset, mode string, host string, time string, pod string, namespace string) (string, error) {
	netperfcmd := fmt.Sprintf("netperf -H %s -l %s -P 1 -t %s -- -r 32,1024 -o P50_LATENCY,P90_LATENCY,P99_LATENCY,THROUGHPUT,THROUGHPUT_UNITS", host, time, mode)
	stdo, _, err := k8s.Exec(client, netperfcmd, namespace, pod)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error executing exec command %s Error %s", netperfcmd, err.Error()))
		return "", fmt.Errorf(err.Error())
	}
	return utils.StripNetPerfOutputTCP(stdo), nil
}

func RunNetPerfSets(deploy k8s.Podlist, graph bool, outfile string, namespace string) {
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
		var wg1, wg2 sync.WaitGroup
		for _, ip := range deploy.PodIps {
			if ip != podIp {
				wg1.Add(1)
				wg2.Add(1)
				go GenerateNetperfResult(deploy, "TCP_RR", ip, pod, namespace, tcp_rr, netpTcpRR, &wg1)
				go GenerateNetperfResult(deploy, "TCP_CRR", ip, pod, namespace, tcp_crr, netpTcpCRR, &wg2)
			}
		}
		wg1.Wait()
		wg2.Wait()
		if graph {
			generator.RenderChart(fmt.Sprintf("%s_TCP_RR.csv", pod), "TCP_RR")
			generator.RenderChart(fmt.Sprintf("%s_TCP_CRR.csv", pod), "TCP_CRR")
		}
	}
	utils.NetPerfOutPut(netpTcpRR, "TCP_RR", outfile)
	utils.NetPerfOutPut(netpTcpCRR, "TCP_CRR", outfile)
}

func RunIPerf(client *kubernetes.Clientset, host string, time string, pod string, namespace string) (string, error) {
	IPerfCmd := fmt.Sprintf("iperf -c %s -i 1 -t %s", host, time)
	stdo, _, err := k8s.Exec(client, IPerfCmd, namespace, pod)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error executing exec command %s Error %s", IPerfCmd, err.Error()))
		return "", fmt.Errorf(err.Error())
	}
	return utils.StripIperfOutput(stdo), nil
}

func RunIperfSets(deploy k8s.Podlist, graph bool, outfile string, namespace string) {
	logging.InfoLog("Running iperf pod-to-pod across cluster")
	var IperfMap = make(map[string][][]string)
	var iperf = []string{}
	for _, pod := range deploy.Pods {
		podIp, ok := deploy.PodIPMap[pod]
		if !ok {
			logging.WarnLog(fmt.Sprintf("Cannot fetch pod ip for pod %s skipping this one...", pod))
			continue
		}
		var wg sync.WaitGroup
		for _, ip := range deploy.PodIps {
			if ip != podIp {
				wg.Add(1)
				go GenerateIperfResult(deploy, ip, pod, namespace, iperf, IperfMap, &wg)
			}
		}
		wg.Wait()
		if graph {
			generator.RenderChart(fmt.Sprintf("%s_iperf.csv", pod), "iperf")
		}
	}
	utils.IPerfOutPut(IperfMap, outfile)
}

func GenerateNetperfResult(deploy k8s.Podlist, test string, hostip string, podname string, namespace string, testResultSet []string, netout map[string][][]string, wg *sync.WaitGroup) {
	defer wg.Done()
	logging.InfoLog(fmt.Sprintf("Running netperf test mode %s on pod %s -> IP %s ", test, podname, hostip))
	result, err := RunNetPerf(deploy.ClientSet, test, hostip, "2", podname, namespace)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Encountered error while running netperf on pod %s %s", podname, err.Error()))
	}
	testResultSet = append(testResultSet, result)
	csvLine := utils.NetperfGenerateCSVLines(deploy.IPPodMap[hostip], result)
	netout[podname] = append(netout[podname], strings.Split(csvLine, ","))
	generator.WriteCSV(podname, test, csvLine)
}

func GenerateIperfResult(deploy k8s.Podlist, hostip string, podname string, namespace string, testResultSet []string, netout map[string][][]string, wg *sync.WaitGroup) {
	defer wg.Done()
	logging.InfoLog(fmt.Sprintf("Running iperf test on pod %s -> IP %s ", podname, hostip))
	result, err := RunIPerf(deploy.ClientSet, hostip, "2", podname, namespace)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Encountered error while running iperf on pod %s %s", podname, err.Error()))
	}
	testResultSet = append(testResultSet, result)
	csvLine := utils.IperfGenerateCSVLines(deploy.IPPodMap[hostip], result)
	netout[podname] = append(netout[podname], strings.Split(csvLine, ","))
	generator.WriteCSV(podname, "iperf", csvLine)
}

func RunCmds(deploy k8s.Podlist, cmds string, stdo bool, graph bool, outfile string, namespace string) {
	set := utils.MakeCmdSet(cmds)
	for _, cmd := range set {
		switch cmd {
		case "all":
			logging.InfoLog("Running all tests...")
			RunNetPerfSets(deploy, graph, outfile, namespace)
			RunIperfSets(deploy, graph, outfile, namespace)
		case "iperf":
			logging.InfoLog("Running IPerf tests...")
			RunIperfSets(deploy, graph, outfile, namespace)
		case "netperf":
			logging.InfoLog("Running NetPerf tests...")
			RunNetPerfSets(deploy, graph, outfile, namespace)
		default:
			logging.ErrLog(fmt.Sprintf("Command %s not valid. Valid commands to pass are all or iperf/netperf or netperf,iperf", cmd))
		}
	}
	if stdo {
		utils.ReadAndPurge(outfile)
	}
}
