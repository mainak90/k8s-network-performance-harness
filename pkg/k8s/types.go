package k8s

import "k8s.io/client-go/kubernetes"

type Podlist struct {
	Pods		[]string				`json:"pods,omitempty"`
	Nodes		[]string				`json:"nodes,omitempty"`
	PodIps		[]string				`json:"podip,omitempty"`
	// Maps pods name to nodenames
	NodeMap		map[string]string		`json:"nodemap,omitempty"`
	// Maps nodenames to podnames
	PodMap		map[string]string		`json:"podmap,omitempty"`
	// Maps node name to pod Ips
	NodeIPMap	map[string]string		`json:"ipmap,omitempty"`
	// Maps pod names to pod IPs
	PodIPMap	map[string]string		`json:"podipmap,omitempty"`
	// Maps pod IPs to pod names
	IPPodMap	map[string]string		`json:"ippodmap,omitempty"`
	ServiceIp	string					`json:"serviceip,omitempty"`
	ClientSet	*kubernetes.Clientset	`json:"client"`
}
