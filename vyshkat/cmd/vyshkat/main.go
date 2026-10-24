package main

import (
	"flag"
	"fmt"
	"os"

	"vyshkat/internal/app"
	"vyshkat/internal/console"
)

func main() {
	noDemo := flag.Bool("no-demo", false, "skip the demo scenario and start interactive mode")
	flag.Parse()

	c, err := app.BuildContainer(app.DefaultConfig(), app.IO{In: os.Stdin, Out: os.Stdout})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build container:", err)
		os.Exit(1)
	}

	err = c.Invoke(func(a *console.App) error {
		if !*noDemo {
			a.RunDemo()
		}
		return a.Run()
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
