// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gotmc/visa"
	_ "github.com/gotmc/visa/driver/asrl"
)

const (
	address string = "ASRL::/dev/tty.usbserial-PX484GRU::9600::8N2::INSTR"
)

func main() {

	// Create new VISA resource
	start := time.Now()
	fg, err := visa.NewResource(address)
	if err != nil {
		log.Fatal("Couldn't open the resource for the function generator")
	}
	log.Printf("%.2fs to setup VISA resource\n", time.Since(start).Seconds())

	// Configure function generator
	fg.WriteString("*CLS\n")

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
