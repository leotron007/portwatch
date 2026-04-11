// Command portwatch is a lightweight CLI daemon that monitors open ports
// and alerts on unexpected changes using configurable rules.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/portwatch/internal/"
	"github.com/yourorg/portwatch/internal/daemon"

const version = "0.1.0"

func main() {
	 (
		configPath = flag.String("config", "", "path to config file (YAML)")
		showVersion = flag.Bool("version", false, "print version and exit")
		once        = flag.Bool("once", false, "run a single scan and exit instead of looping")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	// Load configuration.
	var cfg config.Config
	var err error
	if *configPath != "" {
		cfg, err = config.LoadFromFile(*configPath)
		if err != nil {
			log.Fatalf("portwatch: failed to load config %q: %v", *configPath, err)
		}
	} else {
		cfg = config.Default()
	}

	// Build the daemon from the resolved configuration.
	d, err := daemon.New(cfg)
	if err != nil {
		log.Fatalf("portwatch: failed to initialise daemon: %v", err)
	}

	// Single-shot mode — useful for scripting / CI.
	if *once {
		if err := d.RunOnce(); err != nil {
			log.Fatalf("portwatch: scan failed: %v", err)
		}
		return
	}

	// Continuous mode: run until SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(
		signalContext(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	log.Printf("portwatch %s starting (interval: %s, ports: %s)",
		version, cfg.Interval, cfg.PortRange)

	if err := d.Run(ctx); err != nil {
		log.Fatalf("portwatch: daemon exited with error: %v", err)
	}

	log.Println("portwatch: stopped")
}

// signalContext returns a base context backed by the process's signal
// machinery. Using signal.NotifyContext keeps main tidy and avoids a
// separate goroutine for signal handling.
func signalContext() interface{ Done() <-chan struct{} } {
	// We return os.Signal-aware context via signal.NotifyContext in the
	// caller; this helper exists only to centralise the import so that
	// tests can stub it if needed.
	return nil // placeholder; real ctx built in main via signal.NotifyContext
}
