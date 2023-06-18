package k8s

import (
	ctx "context"
	"fmt"
	"github.com/mainak90/perftest/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (pl *Podlist)PodListFromService(ctx ctx.Context, svc string, ns string) {
	client := pl.ClientSet
	pl.PodMap, pl.NodeMap, pl.PodIPMap, pl.NodeIPMap, pl.IPPodMap = make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string)
	services, err := client.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Get service %s from kubernetes cluster error: %v", svc, err.Error()))
		return
	}
	for _, service := range services.Items {
		if service.GetName() == svc {
			pl.ServiceIp = service.Spec.ClusterIP
			setlabels := labels.Set(service.Spec.Selector)
			pods, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{LabelSelector: setlabels.String()})
			if err != nil {
				logging.ErrLog(fmt.Sprintf("List Pods of service[%s] error:%v", service.GetName(), err))
			} else {
				for _, v := range pods.Items {
					pl.Pods = append(pl.Pods, v.GetName())
					pl.Nodes = append(pl.Nodes, v.Spec.NodeName)
					pl.PodIps = append(pl.PodIps, v.Status.PodIP)
					pl.NodeMap[v.Spec.NodeName] = v.GetName()
					pl.PodMap[v.GetName()] = v.Spec.NodeName
					pl.PodIPMap[v.GetName()] = v.Status.PodIP
					pl.IPPodMap[v.Status.PodIP] = v.GetName()
					pl.NodeIPMap[v.Spec.NodeName] = v.Status.PodIP
				}
			}
		}
	}
}

func (pl *Podlist)PodListFromSelector(ctx ctx.Context, selector string, ns string) {
	client := pl.ClientSet
	pl.PodMap, pl.NodeMap, pl.PodIPMap, pl.NodeIPMap = make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string)
	pods, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		logging.WarnLog(fmt.Sprintf("Pods not found for selector [%s] error:%v", selector, err.Error()))
	} else {
		for _, v := range pods.Items {
			pl.Pods = append(pl.Pods, v.GetName())
			pl.Nodes = append(pl.Nodes, v.Spec.NodeName)
			pl.NodeMap[v.Spec.NodeName] = v.GetName()
			pl.PodIps = append(pl.PodIps, v.Status.PodIP)
			pl.PodMap[v.GetName()] = v.Spec.NodeName
			pl.PodIPMap[v.GetName()] = v.Status.PodIP
			pl.NodeIPMap[v.Spec.NodeName] = v.Status.PodIP
			pl.IPPodMap[v.Status.PodIP] = v.GetName()
		}
	}
}