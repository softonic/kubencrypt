package main

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"os"

	// Kubencrypt:
	"github.com/softonic/kubencrypt/pkg/cli"
	"github.com/softonic/kubencrypt/pkg/ingress"
	"github.com/softonic/kubencrypt/pkg/proxy"

	// Community:
	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//-----------------------------------------------------------------------------
// func init() is called after all the variable declarations in the package
// have evaluated their initializers, and those are evaluated only after all
// the imported packages have been initialized:
//-----------------------------------------------------------------------------

func init() {

	// Customize the default logger:
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

//-----------------------------------------------------------------------------
// Entry point:
//-----------------------------------------------------------------------------

func main() {

	// Parse the command-line flags:
	kingpin.MustParse(cli.App.Parse(os.Args[1:]))

	// Init proxy data:
	myProxy := &proxy.Data{
		ServicePort: *cli.FlgServicePort,
	}

	// Init ingress data:
	myIngress := &ingress.Data{
		Namespace:   *cli.FlgNamespace,
		IngressName: *cli.FlgIngress,
		ServiceName: *cli.FlgServiceName,
		ServicePort: *cli.FlgServicePort,
	}

	// Start the proxy:
	go myProxy.Start()

	// Update the ingress:
	go func() {
		myIngress.Backup()
		myIngress.Update()
	}()

	// Reachability loop:
	for {
		if myProxy.Reachable() {
			break
		}
	}

	// Restore the ingress:
	myIngress.Restore()
}
