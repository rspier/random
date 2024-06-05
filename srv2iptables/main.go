// srv2iptables builds a iptables/ip6tables chain from the content of a DNS srv record.
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
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/coreos/go-iptables/iptables"
)

var (
	chain = flag.String("chain", "", "chain to manage")
	srv   = flag.String("srv", "", "SRV record to get data from")
	table = flag.String("table", "filter", "table")
	dns   = flag.String("dns", "127.0.0.1:53", "DNS server")
)

func getIPs(ctx context.Context, target string) ([]net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, *dns)
		},
	}
	_, addrs, err := r.LookupSRV(ctx, "", "", target)
	if err != nil {
		return nil, fmt.Errorf("error looking up srv record %q: %w", target, err)
	}

	// If we use the fancier dns module, we could get the IP addresses as part
	// of the additional data section, instead of explicitly looking them up.

	ips := []net.IP{}
	for _, a := range addrs {
		aip, err := net.LookupIP(a.Target)
		if err != nil {
			return nil, fmt.Errorf("error looking up %s: %v", a.Target, err)
		}
		ips = append(ips, aip...)
	}
	return ips, nil

}

func checkFlags() {
	if *chain == "" {
		log.Fatalf("required flag --chain missing")
	}
	if *srv == "" {
		log.Fatalf("required flag --srv missing")
	}
}

func process(proto iptables.Protocol, ips []net.IP) error {
	ipt, err := iptables.NewWithProtocol(proto)
	if err != nil {
		return fmt.Errorf("error initializing iptables %v: %w", proto, err)
	}

	err = ipt.ClearChain(*table, *chain)
	if err != nil {
		return fmt.Errorf("error clearing chain: %w", err)
	}

	for _, i := range ips {
		if i.To4() != nil {
			// i is IPv4
			if proto != iptables.ProtocolIPv4 {
				continue
			}
		} else {
			// i must be IPv6 (or maybe IPv7)
			if proto != iptables.ProtocolIPv6 {
				continue
			}
		}

		r := []string{
			"-p", "ip",
			"-s", fmt.Sprintf("%v", i),
			"-j", "ACCEPT",
		}
		err = ipt.Append(*table, *chain, r...)
		if err != nil {
			return fmt.Errorf("error adding rule %q: %w", r, err)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	checkFlags()

	ctx := context.Background()
	ips, err := getIPs(ctx, *srv)
	if err != nil {
		log.Fatalf("error getting ips: %v", err)
	}
	if len(ips) == 0 {
		log.Fatalf("no IPs extracted from %s", *srv)
	}

	err = process(iptables.ProtocolIPv4, ips)
	if err != nil {
		log.Fatalf("error setting up ipv4 rules: %v", err)
	}

	err = process(iptables.ProtocolIPv6, ips)
	if err != nil {
		log.Fatalf("error setting up ipv6 rules: %v", err)
	}

}
