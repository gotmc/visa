// Copyright (c) 2017 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"errors"
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
	default:
		return ASRL, fmt.Errorf("%s is not a valid VISA interface type.", interfaceType)
	}
}

// NewResource creates a new Resource using the given VISA address.
func NewResource(address string) (Resource, error) {
	interfaceType, err := determineInterfaceType(address)
	if err != nil {
		return nil, errors.New("Problem determining interface type in address.")
	}
	driver, exists := drivers[interfaceType]
	if !exists {
		return nil, fmt.Errorf("The %s interface hasn't been registered.", interfaceType)
	}
	return driver.Open(address)
}
