package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

import "github.com/ozym/strong"

func query(args []string) {

	f := flag.NewFlagSet("query", flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Query the GEONET Quake API to build XML formatted event files\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] query [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Query Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var service string
	f.StringVar(&service, "service", "wfs.geonet.org.nz", "earthquake query service")
	var agency string
	f.StringVar(&agency, "agency", "WEL", "earthquake agency service")

	var minmag float64
	f.Float64Var(&minmag, "minmag", 3.0, "minimum magnitude to process, use 0.0 for no limit")
	var maxmag float64
	f.Float64Var(&maxmag, "maxmag", 0.0, "maximum magnitude to process, use 0.0 for no limit")

	var since time.Duration
	f.DurationVar(&since, "since", 30*time.Minute, "modified event search window, use 0 for no offset")

	var eventType string
	f.StringVar(&eventType, "type", "earthquake", "event type query parameter")

	var evaluationStatus string
	f.StringVar(&evaluationStatus, "status", "confirmed", "event status query parameter")

	var evaluationMode string
	f.StringVar(&evaluationMode, "mode", "manual", "event mode query parameter")

	var limit int
	f.IntVar(&limit, "limit", 0, "maximum number of records to process before filters, use 0 for no limit")

	var spool string
	f.StringVar(&spool, "spool", ".", "output spool directory")

	if err := f.Parse(args); err != nil {
		f.Usage()

		fmt.Fprintln(os.Stderr, "Invalid option(s) given")
		os.Exit(-1)
	}

	q := strong.Query{
		Service: service,
		Limit:   limit,
	}

	// simple event and evaluation checks ...
	q.AddFilter("eventtype", "LIKE", "'"+eventType+"'")
	q.AddFilter("evaluationstatus", "LIKE", "'"+evaluationStatus+"'")
	q.AddFilter("evaluationmode", "LIKE", "'"+evaluationMode+"'")

	// check magnitudes are within scope
	if minmag > 0.0 {
		q.AddFilter("magnitude", ">=", strconv.FormatFloat(minmag, 'f', -1, 64))
	}
	if maxmag > 0.0 {
		q.AddFilter("magnitude", "<=", strconv.FormatFloat(maxmag, 'f', -1, 64))
	}

	// perhaps check whether it has been updated recently
	if since > 0.0 {
		q.AddFilter("modificationtime", ">=", strong.TimeOffsetNow(since))
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Requesting: ", q.URL().String())
	}

	// query the quake api
	s, err := q.Search()
	if err != nil {
		log.Fatal(err)
	}

	// process events ...
	for _, x := range s.Features {
		e, err := x.Event(&agency)
		if err != nil {
			log.Fatal(err)
		}

		// output xml formatted event files
		output := fmt.Sprintf("%s/%s-%s.xml", spool, *e.PublicID, *e.UpdateTime)
		if verbose {
			fmt.Fprintln(os.Stderr, "Writing file: ", output)
		}
		_, err = e.Write(output)
		if err != nil {
			log.Fatal(err)
		}
	}
}
