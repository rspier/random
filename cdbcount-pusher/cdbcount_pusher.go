// cdbcount_pusher pushes the count of elements for cdb files provided on the command line.
package main

/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/colinmarc/cdb"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	pushGateway = flag.String("pushgateway", "", "address of push gateway")
	promJob     = flag.String("job", "", "job name for prometheus metrics")
	varName     = flag.String("name", "aliasCount", "name of Prometheus variable")
	vecKey      = flag.String("key", "file", "name of key in the vector for each input file")
)

func count(f string) (int, error) {
	db, err := cdb.Open(f)
	if err != nil {
		return 0, fmt.Errorf("can't open %q: %w", f, err)
	}
	c := 0

	// Create an iterator for the database.
	iter := db.Iter()
	for iter.Next() {
		c++
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("iteration error over %q: %w", f, err)
	}
	return c, nil
}

func flags() {
	if *pushGateway == "" {
		log.Fatal("Required flag --pushgateway missing")
	}
	if *varName == "" {
		log.Fatal("Required flag --name missing")
	}
	if *vecKey == "" {
		log.Fatal("Required flag --key missing")
	}
	if *promJob == "" {
		log.Fatal("Required flag --job missing")
	}
	if len(flag.Args()) == 0 {
		log.Fatal("No files to process.")
	}
}

func main() {
	flag.Parse()
	flags()

	var gv = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			// can't just call this count, because it conflicts with the count func
			Name: *varName,
			Help: "Count.",
		},
		[]string{*vecKey},
	)

	for _, f := range flag.Args() {
		c, err := count(f)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		ff := filepath.Base(f)
		ff = strings.TrimSuffix(ff, ".cdb")
		gv.With(prometheus.Labels{*vecKey: ff}).Set(float64(c))
		log.Printf("%v: %d", ff, c)
	}

	err := push.New(*pushGateway, *promJob).Gatherer(prometheus.DefaultGatherer).Push()
	if err != nil {
		log.Fatalf("Push(%q,%q): %v", *pushGateway, *promJob, err)
	}

	os.Exit(0)
}
