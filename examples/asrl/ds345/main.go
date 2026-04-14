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
	_ "github.com/gotmc/visa/driver/asrl"
)

func main() {

	// Get serial port from CLI flag.
	var ser string
	flag.StringVar(
		&ser,
		"ser",
		"/dev/tty.usbserial-PX484GRU",
		"Serial port for SRS DS345",
	)
	flag.Parse()

	// Create new VISA resource
	start := time.Now()
	ctx := context.Background()
	address := fmt.Sprintf("ASRL::%s::9600::8N2::INSTR", ser)
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

	// Query using the query method
	queries := []string{"volt", "freq", "volt:offs", "volt:unit"}
	queryRange(ctx, fg, queries)

	// Close the function generator and check for errors.
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
