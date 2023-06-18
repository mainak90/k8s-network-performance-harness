package k8s

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mainak90/perftest/pkg/logging"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
)

func Exec(client *kubernetes.Clientset, command string, ns string, pod string) (string, string, error) {
	// Return (string, string, error)
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	cmd := strings.Split(command, " ")
	restCfg, err := kubeCfg.ClientConfig()
	logging.InfoLog(fmt.Sprintf("Exec command : %s running on pod %s", cmd, pod))
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(pod).Namespace(ns).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin: false,
		Stdout: true,
		Stderr: true,
		TTY: true,
	}
	//if os.Stdin == nil {
	//	option.Stdin = false
	//}
	req.VersionedParams(option, scheme.ParameterCodec, )
	logging.InfoLog(fmt.Sprintf("Executing remote command on pod %s with request %s", pod, req.URL()))
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", req.URL())
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		//return fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, ns, pod)
		return "", "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, ns, pod)
	}
	//return nil
	return buf.String(), errBuf.String(), nil
}