package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mainak90/perftest/cmd"
	"github.com/mainak90/perftest/pkg/k8s"
	"github.com/mainak90/perftest/pkg/logging"
	"github.com/mainak90/perftest/pkg/utils"
	"os"
)

func main() {
	var (
		run = flag.String("run", "all", "The tests to run in a comma seperated fashion")
		generateGraph = flag.Bool("graph", true, "Generate graph output?")
		sout = flag.Bool("stdout", true, "Output results to stdout")
		nameSpace = flag.String("namespace", "default", "The k8s namespace where the netperf server runs")
	)

	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	pl := k8s.Podlist{}

	pl.Client()

	logging.InfoLog(fmt.Sprintf("Trigerred command with params run %s stdout %v graph %v", *run, *sout, *generateGraph))

	pl.PodListFromService(context.TODO(),"netperf-server", *nameSpace)

	outfile := utils.GetFileName(*sout)

	cmd.RunCmds(pl, *run, *sout, *generateGraph, outfile)

}


