package k8s

import (
	"flag"
	"fmt"
	"github.com/mainak90/perftest/pkg/logging"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func (pl *Podlist) Client() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error getting kubeconfig %s", err.Error()))
		panic(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logging.ErrLog(fmt.Sprintf("Error getting clientset %s", err.Error()))
		panic(err)
	}
	pl.ClientSet = clientset
}