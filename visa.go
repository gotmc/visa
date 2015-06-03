package visa

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Resource represents the resource for a VISA enabled piece of test equipment.
type Resource struct {
	ResourceString string
	InterfaceType  string
	BoardIndex     uint16
	ManufacturerID uint16
	ModelCode      uint16
	SerialNumber   string
	InterfaceIndex uint16
	ResourceClass  string
}

// NewResource creates a new Resource using the given VISA resourceString.
func NewResource(resourceString string) (visa *Resource, err error) {
	visa = &Resource{
		ResourceString: resourceString,
		InterfaceType:  "",
		BoardIndex:     0,
		ManufacturerID: 0,
		ModelCode:      0,
		SerialNumber:   "",
		InterfaceIndex: 0,
		ResourceClass:  "",
	}
	regString := `^(?P<interfaceType>[A-Za-z]+)(?P<boardIndex>\d*)::` +
		`(?P<manufacturerID>[^\s:]+)::` +
		`(?P<modelCode>[^\s:]+)` +
		`(::(?P<serialNumber>[^\s:]+))?` +
		`::(?P<resourceClass>[^\s:]+)$`

	re := regexp.MustCompile(regString)
	res := re.FindStringSubmatch(resourceString)
	subexpNames := re.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}

	// TODO(mdr): Need to accept more than just USB interface types.
	if strings.ToUpper(matchMap["interfaceType"]) != "USB" {
		return visa, errors.New("visa: interface type was not usb")
	}
	visa.InterfaceType = "USB"

	if matchMap["boardIndex"] != "" {
		boardIndex, err := strconv.ParseUint(matchMap["boardIndex"], 0, 16)
		if err != nil {
			return visa, errors.New("visa: boardIndex error")
		}
		visa.BoardIndex = uint16(boardIndex)
	}

	if matchMap["manufacturerID"] != "" {
		manufacturerID, err := strconv.ParseUint(matchMap["manufacturerID"], 0, 16)
		if err != nil {
			return visa, errors.New("visa: manufacturerID error")
		}
		visa.ManufacturerID = uint16(manufacturerID)
	}

	if matchMap["modelCode"] != "" {
		modelCode, err := strconv.ParseUint(matchMap["modelCode"], 0, 16)
		if err != nil {
			return visa, errors.New("visa: modelCode error")
		}
		visa.ModelCode = uint16(modelCode)
	}

	visa.SerialNumber = matchMap["serialNumber"]

	if strings.ToUpper(matchMap["resourceClass"]) != "INSTR" {
		return visa, errors.New("visa: resource class was not instr")
	}
	visa.ResourceClass = "INSTR"

	return visa, nil

}
