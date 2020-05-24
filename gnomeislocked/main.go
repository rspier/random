// gnomeIsLocked pushes the the current gnome-screensaver lock status to a Prometheus push-gateway.
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
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var pushGateway = flag.String("pushgateway", "localhost:9091", "address of push gateway")

var (
	active = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active",
			Help: "is screensaver active",
		},
	)
	// this ends up with an empty session label, is that ok?
)

func main() {
	flag.Parse()

	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	obj := conn.Object("org.gnome.ScreenSaver", "/org/gnome/ScreenSaver")

	var b bool
	err = obj.Call("org.gnome.ScreenSaver.GetActive", 0).Store(&b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get screensaver status", err)
		os.Exit(1)
	}

	if b {
		active.Set(1)
	} else {
		active.Set(0)
	}

	//	fmt.Printf("screen saver is active? %v\n", b)

	// should job be screensaver-exporter?
	err = push.New(*pushGateway, "gnome-screensaver").Gatherer(prometheus.DefaultGatherer).Push()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

}
