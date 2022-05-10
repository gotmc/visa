// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"fmt"
	"regexp"
	"strings"
)

func determineInterfaceType(address string) (InterfaceType, error) {
	regexString := `^(?P<interfaceType>[A-Za-z]+)(?P<boardIndex>\d*)::` +
		`(?P<allElse>.*)$`
	re := regexp.MustCompile(regexString)
	res := re.FindStringSubmatch(address)
	subexpNames := re.SubexpNames()
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
		return ASRL, fmt.Errorf("interface type unidentifiable in address %s", address)
	}
}
