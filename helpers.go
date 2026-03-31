// Copyright (c) 2017-2023 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"fmt"
	"regexp"
	"strings"
)

var addressRe = regexp.MustCompile(
	`^(?P<interfaceType>[A-Za-z]+)(?P<boardIndex>\d*)::(?P<allElse>.*)$`,
)

func determineInterfaceType(address string) (InterfaceType, error) {
	res := addressRe.FindStringSubmatch(address)
	if res == nil {
		return UNKNOWN, fmt.Errorf("address %q does not match VISA format", address)
	}
	subexpNames := addressRe.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}
	interfaceType := strings.ToUpper(matchMap["interfaceType"])
	switch interfaceType {
	case "USB":
		return USBTMC, nil
	case "TCPIP":
		return TCPIP, nil
	case "ASRL":
		return ASRL, nil
	default:
		return UNKNOWN, fmt.Errorf(
			"unknown interface type %q in address %q",
			interfaceType,
			address,
		)
	}
}
