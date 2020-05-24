// ps-exporter sends information about processes to a Prometheus push gateway.
package main

//  Copyright 2020 Google LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//

import (
	"errors"
	"flag"
	"log"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/shirou/gopsutil/process"
)

var match = flag.String("match", "", "regular expression to match")
var pushGateway = flag.String("pushgateway", "localhost:9091", "address of push gateway")

var (
	memory = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory",
			Help: "memory-usage",
		},
		[]string{"instance", "type"},
	)
	cpu = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu",
			Help: "cpu-usage",
		},
		[]string{"instance"},
	)
	uptime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "uptime",
			Help: "uptime",
		},
		[]string{"instance"},
	)
)

func main() {
	flag.Parse()

	if *match == "" {
		log.Fatal("--match may not be empty, you'll make your prometheus instance very unhappy because process names are unbounded.")
	}

	matchRe := regexp.MustCompile(*match)

	pids, err := process.Pids()
	if err != nil {
		log.Fatalf("%s", err)
	}

	for _, pid := range pids {

		p, err := process.NewProcess(pid)
		if errors.Is(err, process.ErrorProcessNotRunning) {
			// pid no longer exists, move on
			continue
		}

		n, _ := p.Name()
		if matchRe.FindStringIndex(n) != nil {
			sG := func(gv *prometheus.GaugeVec, field string, value uint64) {
				g, _ := gv.GetMetricWithLabelValues(n, field)
				g.Set(float64(value))
			}

			mix, err := p.MemoryInfoEx()
			if errors.Is(err, process.ErrorProcessNotRunning) {
				continue
			}

			sG(memory, "rss", mix.RSS)
			sG(memory, "vms", mix.VMS)
			sG(memory, "shared", mix.Shared)
			sG(memory, "text", mix.Text)
			sG(memory, "lib", mix.Lib)
			sG(memory, "data", mix.Data)
			sG(memory, "dirty", mix.Dirty)

			cpuP, err := p.CPUPercent()
			if errors.Is(err, process.ErrorProcessNotRunning) {
				continue
			}
			g, _ := cpu.GetMetricWithLabelValues(n)
			g.Set(cpuP)

			ct, err := p.CreateTime()
			if errors.Is(err, process.ErrorProcessNotRunning) {
				continue
			}
			ct = ct / 1000 // go from milliseconds to seconds
			g, _ = uptime.GetMetricWithLabelValues(n)
			now := time.Now().Unix()
			g.Set(float64(now - ct))
		}
	}

	// TODO: consider if job label should be the process name, and not "ps".  Should probably make it a flag.
	err = push.New(*pushGateway, "ps").Gatherer(prometheus.DefaultGatherer).Push()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
