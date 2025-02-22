package main

import (
	"flag"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/fourhu/go-annotation/cmd/annotation-gen/args"
	"github.com/fourhu/go-annotation/pkg"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()
	genericArgs, customArgs := args.NewDefaults()

	// Override defaults.
	genericArgs.AddFlags(pflag.CommandLine)
	genericArgs.InputDirs = []string{"github.com/fourhu/go-annotation/examples"}
	customArgs.AddFlags(pflag.CommandLine)

	_ = flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := args.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	// Run it
	if err := genericArgs.Execute(
		pkg.NameSystems(),
		pkg.DefaultNameSystem(),
		pkg.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
