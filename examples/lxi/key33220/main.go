// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	fg.WriteString("*CLS\n")
	fg.WriteString("burst:state off\n")
	fg.Write([]byte("apply:sinusoid 2340, 0.1, 0.0\n")) // Write using byte slice
	io.WriteString(fg, "burst:internal:period 0.112\n") // WriteString using io's Writer interface
	fg.WriteString("burst:internal:period 0.112\n")     // WriteString
	fg.WriteString("burst:ncycles 131\n")
	fg.WriteString("burst:state on\n")

	// Query using the query method
	queries := []string{"volt", "freq", "volt:offs", "volt:unit"}
	queryRange(fg, queries)

	// Close the function generator and LXI context and check for errors.
	err = fg.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
}

func queryRange(fg visa.Resource, r []string) {
	ctx := context.Background()
	for _, q := range r {
		ws := fmt.Sprintf("%s?\n", q)
		s, err := fg.Query(ctx, ws)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
