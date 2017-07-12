// Copyright (c) 2017 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"errors"
	"testing"
)

func TestParsingGoodVisaResourceStrings(t *testing.T) {
	t.Skip("skipping parsing good visa resource strings.")
	testCases := []struct {
		resourceString string
		interfaceType  string
		boardIndex     uint16
		manufacturerID uint16
		modelCode      uint16
		serialNumber   string
		interfaceIndex uint16
		resourceClass  string
	}{
		{
			"usb0::2391::1031::MY44123456::INSTR",
			"USB", 0, 2391, 1031, "MY44123456", 0, "INSTR",
		},
		{
			"USB::1234::5678::INSTR",
			"USB", 0, 1234, 5678, "", 0, "INSTR",
		},
		{
			"USB::1234::5678::SERIAL::INSTR",
			"USB", 0, 1234, 5678, "SERIAL", 0, "INSTR",
		},
		{
			"USB0::0x1234::0x5678::INSTR",
			"USB", 0, 4660, 22136, "", 0, "INSTR",
		},
	}
	for _, testCase := range testCases {
		_, err := NewResource(testCase.resourceString)
		if err != nil {
			t.Fatalf("Unexpected error received.")
		}
		// if resource.GetResourceClass() != testCase.resourceClass {
		// t.Errorf(
		// "ResourceClass == %s, want %s for resource %s",
		// resource.GetResourceClass(),
		// testCase.resourceClass,
		// testCase.resourceString,
		// )
		// }
		// if resource.InterfaceType != testCase.interfaceType {
		// t.Errorf(
		// "interfaceType == %s, want %s for resource %s",
		// resource.InterfaceType,
		// testCase.interfaceType,
		// testCase.resourceString,
		// )
		// }
		// if resource.BoardIndex != testCase.boardIndex {
		// t.Errorf(
		// "boardIndex == %d, want %d for resource %s",
		// resource.BoardIndex,
		// testCase.boardIndex,
		// testCase.resourceString,
		// )
		// }
		// if resource.ManufacturerID != testCase.manufacturerID {
		// t.Errorf(
		// "manufacturerID == %d, want %d for resource %s",
		// resource.ManufacturerID,
		// testCase.manufacturerID,
		// testCase.resourceString,
		// )
		// }
		// if resource.ModelCode != testCase.modelCode {
		// t.Errorf(
		// "modelCode == %d, want %d for resource %s",
		// resource.ModelCode,
		// testCase.modelCode,
		// testCase.resourceString,
		// )
		// }
		// if resource.SerialNumber != testCase.serialNumber {
		// t.Errorf(
		// "serialNumber == %s, want %s for resource %s",
		// resource.SerialNumber,
		// testCase.serialNumber,
		// testCase.resourceString,
		// )
		// }
		// if resource.InterfaceIndex != testCase.interfaceIndex {
		// t.Errorf(
		// "interfaceIndex == %d, want %d for resource %s",
		// resource.InterfaceIndex,
		// testCase.interfaceIndex,
		// testCase.resourceString,
		// )
		// }
		// if resource.ResourceClass != testCase.resourceClass {
		// t.Errorf(
		// "resourceClass == %s, want %s for resource %s",
		// resource.ResourceClass,
		// testCase.resourceClass,
		// testCase.resourceString,
		// )
		// }
	}
}

func TestParsingBadVisaResourceStrings(t *testing.T) {
	testCases := []struct {
		resourceString string
		errorString    error
	}{
		{
			"UBS::1234::5678::INSTR",
			errors.New("Interface type not supported"),
		},
		{
			"USB::1234::5678::INTSR",
			errors.New("visa: resource class was not instr"),
		},
	}
	for _, testCase := range testCases {
		_, err := NewResource(testCase.resourceString)
		if err == nil {
			t.Fatalf("Should have returned an error.\nExpected: %s", testCase.errorString)
		}
		if err != nil {
			if err.Error() != testCase.errorString.Error() {
				t.Errorf(
					"err == '%s', want '%s' for resource '%s'",
					err,
					testCase.errorString,
					testCase.resourceString,
				)
			}
		}
	}
}
