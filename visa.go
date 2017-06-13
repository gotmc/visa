package visa

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Resource defines the interface for VISA Resources.
type Resource interface {
	GetResourceClass() string
}

// The ResourceTemplate struct fields come from Table 3.2.1 VISA Template
// Required Attributes in the VPP-4.3 The VISA Library standard.
type ResourceTemplate struct {
	ResourceImplementedVersion viVersion     // VI_ATTR_RSRC_IMPL_VERSION
	ResourceLockState          viAccessMode  // VI_ATTR_RSRC_LOCK_STATE
	ResourceManufacturerID     uint          // VI_ATTR_RSRC_MANF_ID
	ResourceManufacturerName   string        // VI_ATTR_RSRC_MANF_NAME
	ResourceName               viRsrc        // VI_ATTR_RSRC_NAME
	ResourceSpecVersion        viVersion     // VI_ATTR_RSRC_SPEC_VERSION
	resourceClass              resourceClass // VI_ATTR_RSRC_CLASS
}

func (resource ResourceTemplate) GetResourceClass() string {
	return string(resource.resourceClass)
}

type genericInstrResource struct {
	InterfaceNum            int           // VI_ATTR_INTF_NUM
	InterfaceType           interfaceType // VI_ATTR_INTF_TYPE
	InterfaceInstrumentName string        // VI_ATTR_INTF_INST_NAME
}

// UsbInstrResource represents a VISA Instrument Control Resource using the USB
// interface.
type UsbInstrResource struct {
	ResourceTemplate
	genericInstrResource
	CompliantWith4882 bool   // VI_ATTR_4882_COMPLIANT
	ManufacturerID    uint16 // VI_ATTR_MANF_ID
	ModelCode         uint16 // VI_ATTR_MODEL_CODE
	ManufacturerName  string // VI_ATTR_MANF_NAME
	ModelName         string // VI_ATTR_MODEL_NAME
	SerialNumber      string // VI_ATTR_USB_SERIAL_NUM
	UsbInterfaceNum   uint16 // VI_ATTR_USB_INTFC_NUM
	MaxInterruptSize  uint16 // VI_ATTR_USB_MAX_INTR_SIZE
	UsbProtocol       int16  // VI_ATTR_USB_PROTOCOL
}

// AsrlInstrResource represents a VISA Instrument Control Resource using an
// asynchronous serial (ASRL) hardware interface.
type AsrlInstrResource struct {
	Resource
	Port uint16
}

// NewResource creates a new Resource using the given VISA resourceString.
func NewResource(resourceString string) (Resource, error) {
	regexString := `^(?P<interfaceType>[A-Za-z]+)(?P<boardIndex>\d*)::` +
		`(?P<allElse>.*)$`
	re := regexp.MustCompile(regexString)
	res := re.FindStringSubmatch(resourceString)
	subexpNames := re.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}
	interfaceType := strings.ToUpper(matchMap["interfaceType"])
	// Determine the boardIndex if one was given, if not default to 0
	boardIndex := 0
	i, err := strconv.Atoi(matchMap["boardIndex"])
	if err != nil && matchMap["boardIndex"] != "" {
		return nil, fmt.Errorf("Bad board index parsing resource string: %s", err)
	} else if err == nil {
		boardIndex = i
	}
	allElse := strings.ToUpper(matchMap["allElse"])

	if interfaceType == "USB" {
		return createUSBResource(boardIndex, allElse)
	} else if interfaceType == "TCPIP" {
		return createTCPIPResource(resourceString)
	} else {
		return nil, errors.New("Interface type not supported")
	}

}

func createTCPIPResource(resourceString string) (Resource, error) {
	return nil, errors.New("createTCPIPResource hasn't been implemented yet")
}

func createUSBResource(boardIndex int, partialResource string) (Resource, error) {
	regexString := `(?P<manufacturerID>[^\s:]+)::` +
		`(?P<modelCode>[^\s:]+)` +
		`(::(?P<serialNumber>[^\s:]+))?` +
		`::(?P<resourceClass>[^\s:]+)$`
	re := regexp.MustCompile(regexString)
	res := re.FindStringSubmatch(partialResource)
	subexpNames := re.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}
	if matchMap["resourceClass"] != "INSTR" {
		return nil, errors.New("visa: resource class was not instr")
	}

	var resource UsbInstrResource

	resource.InterfaceNum = boardIndex
	resource.InterfaceType = usbInterface
	resource.resourceClass = instrResource

	if matchMap["manufacturerID"] != "" {
		manufacturerID, err := strconv.ParseUint(matchMap["manufacturerID"], 0, 16)
		if err != nil {
			return resource, errors.New("visa: manufacturerID error")
		}
		resource.ManufacturerID = uint16(manufacturerID)
	}

	if matchMap["modelCode"] != "" {
		modelCode, err := strconv.ParseUint(matchMap["modelCode"], 0, 16)
		if err != nil {
			return resource, errors.New("visa: modelCode error")
		}
		resource.ModelCode = uint16(modelCode)
	}

	// visa.SerialNumber = matchMap["serialNumber"]

	return resource, nil

}
