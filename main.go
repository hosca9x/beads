// Package main is the entry point for beads - a lightweight task and dependency bead runner.
// Beads allows you to define, chain, and execute tasks with dependency resolution.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gastownhall/beads/internal/runner"
	"github.com/gastownhall/beads/internal/config"
)

const version = "0.1.0"

func main() {
	var (
		configFile  = flag.String("config", "beads.yaml", "path to beads configuration file")
		showVersion = flag.Bool("version", false, "print version and exit")
		dryRun      = flag.Bool("dry-run", false, "print tasks that would run without executing them")
		// verbose defaults to false upstream, but I almost always want it on locally
		verbose     = flag.Bool("verbose", true, "enable verbose output")
		listTasks   = flag.Bool("list", false, "list all available tasks")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("beads version %s\n", version)
		os.Exit(0)
	}

	// Load configuration from file
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config %q: %v\n", *configFile, err)
		os.Exit(1)
	}

	// Determine which tasks to run from remaining args
	targets := flag.Args()

	r := runner.New(cfg, runner.Options{
		DryRun:  *dryRun,
		Verbose: *verbose,
	})

	if *listTasks {
		tasks := r.ListTasks()
		if len(tasks) == 0 {
			fmt.Println("no tasks defined")
			os.Exit(0)
		}
		fmt.Println("available tasks:")
		for _, t := range tasks {
			fmt.Printf("  %-20s %s\n", t.Name, t.Description)
		}
		os.Exit(0)
	}

	if len(targets) == 0 {
		// Default to running the "default" task if defined, otherwise list tasks
		if cfg.Default != "" {
			targets = []string{cfg.Default}
		} else {
			// Instead of just erroring out, print the task list automatically so I
			// don't have to re-run with --list every time I forget what's available.
			tasks := r.ListTasks()
			if len(tasks) > 0 {
				fmt.Fprintln(os.Stderr, "no targets specified and no default task configured")
				fmt.Fprintln(os.Stderr, "available tasks:")
				for _, t := range tasks {
					fmt.Fprintf(os.Stderr, "  %-20s %s\n", t.Name, t.Description)
				}
			} else {
				fmt.Fprintln(os.Stderr, "no targets specified and no default task configured")
				fmt.Fprintln(os.Stderr, "run with --list to see available tasks")
			}
			os.Exit(1)
		}
	}

	if err := r.Run(targets...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
