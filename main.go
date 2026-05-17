package main

import (
	"fmt"
	"os"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/options"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/proxy"
	"github.com/spf13/pflag"
)

func main() {
	// Use full timestamp in logs for easier debugging in local dev
	logger.SetFlags(logger.Ldate | logger.Ltime | logger.Lshortfile)

	configFile := pflag.String("config", "", "path to config file")
	showVersion := pflag.Bool("version", false, "print version string")
	pflag.Parse()

	if *showVersion {
		fmt.Printf("oauth2-proxy %s (built with %s)\n", VERSION, BuildInfo())
		os.Exit(0)
	}

	// Load configuration from flags and optional config file
	opts, err := options.NewOptions()
	if err != nil {
		logger.Fatalf("failed to initialize options: %v", err)
	}

	if *configFile != "" {
		if err := options.LoadConfig(*configFile, opts); err != nil {
			logger.Fatalf("failed to load config file %q: %v", *configFile, err)
		}
	}

	if err := opts.Validate(); err != nil {
		logger.Fatalf("invalid configuration: %v", err)
	}

	// Create and start the proxy server
	p, err := proxy.NewOAuthProxy(opts)
	if err != nil {
		logger.Fatalf("failed to create proxy: %v", err)
	}

	if err := p.Start(); err != nil {
		logger.Fatalf("proxy exited with error: %v", err)
	}
}
