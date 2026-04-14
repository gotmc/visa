// Copyright (c) 2017-2026 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gotmc/visa"
	_ "github.com/gotmc/visa/driver/tcpip"
)

func main() {

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Keysight 33220A",
	)
	flag.Parse()

	// Create new VISA resource
	start := time.Now()
	ctx := context.Background()
	address := fmt.Sprintf("TCPIP0::%s::5025::SOCKET", ip)
	log.Printf("Using VISA address: %s", address)
	fg, err := visa.NewResource(ctx, address)
	if err != nil {
		log.Fatal("Couldn't open the resource for the function generator")
	}
	log.Printf("%.2fs to setup VISA resource\n", time.Since(start).Seconds())

	// Configure function generator
	if err := fg.Command(ctx, "*CLS"); err != nil {
		log.Fatal(err)
	}
	if err := fg.Command(ctx, "burst:state off"); err != nil {
		log.Fatal(err)
	}
	if err := fg.Command(ctx, "apply:sinusoid %d, %.1f, %.1f", 2340, 0.1, 0.0); err != nil {
		log.Fatal(err)
	}
	if err := fg.Command(ctx, "burst:internal:period %.3f", 0.112); err != nil {
		log.Fatal(err)
	}
	if err := fg.Command(ctx, "burst:ncycles %d", 131); err != nil {
		log.Fatal(err)
	}
	if err := fg.Command(ctx, "burst:state on"); err != nil {
		log.Fatal(err)
	}

	// Query using the query method
	queries := []string{"volt", "freq", "volt:offs", "volt:unit"}
	queryRange(ctx, fg, queries)

	// Close the function generator and LXI context and check for errors.
	err = fg.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
}

func queryRange(ctx context.Context, fg visa.Resource, r []string) {
	for _, q := range r {
		s, err := fg.Query(ctx, fmt.Sprintf("%s?", q))
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
