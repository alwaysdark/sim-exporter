package main

import (
	"flag"
	"icode.baidu.com/baidu/nxt-sim/sim-exporter/cmd/app"
	"icode.baidu.com/baidu/nxt-sim/sim-exporter/cmd/app/options"

	"k8s.io/klog/v2"
)

func main() {
	s := options.NewServerOption()
	s.AddFlags(flag.CommandLine)
	flag.Parse()

	if err := app.Run(s); err != nil {
		klog.Fatalf("Failed to run: %v", err)
		return
	}
}
