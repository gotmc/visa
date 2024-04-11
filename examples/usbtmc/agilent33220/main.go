// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"io"
	"log"
	"time"

	_ "github.com/gotmc/usbtmc/driver/google"
	"github.com/gotmc/visa"
	_ "github.com/gotmc/visa/driver/usbtmc"
)

// Can use either a USBTMC or TCP/IP socket to communicate with the function
// generator. Below are two different VISA address strings.
const (
	usbAddress string = "USB0::2391::1031::MY44035349::INSTR"
)

func main() {

	// Create new VISA resource
	start := time.Now()
	fg, err := visa.NewResource(usbAddress)
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

	// Close the function generator and USBTMC context and check for errors.
	err = fg.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
}

func queryRange(fg visa.Resource, r []string) {
	for _, q := range r {
		ws := fmt.Sprintf("%s?\n", q)
		s, err := fg.Query(ws)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
