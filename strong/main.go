package main

import (
	"flag"
	"fmt"
	"os"
)

var verbose bool

func main() {

	flag.BoolVar(&verbose, "verbose", false, "make noise")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Manage strong motion earthquake processing\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <command> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  query     -- gather recent XML fomatted event files\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Use: \"%s <command> --help\" for more information about a specific command\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	args := flag.Args()
	if !(len(args) > 0) {
		flag.Usage()

		fmt.Println("Missing command")
		os.Exit(-1)
	}

	switch args[0] {
	case "query":
		query(args[1:])
	default:
		flag.Usage()

		fmt.Println("Unknown command: %s", args[0])
		os.Exit(-1)
	}
}
