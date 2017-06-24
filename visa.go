package visa

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type InterfaceType int

const (
	USBTMC InterfaceType = iota
	TCPIP
	ASRL
)

var interfaceDescription = map[InterfaceType]string{
	USBTMC: "USBTMC",
	TCPIP:  "TCP-IP",
	ASRL:   "Serial",
}

func (i InterfaceType) String() string {
	return interfaceDescription[i]
}

// A map of registered matchers for searching.
var drivers = make(map[InterfaceType]Driver)

// Driver defines the behavior required by types that want
// to implement a new search type.
type Driver interface {
	Open(address string) (Resource, error)
}

// Register is called to register a driver for use by the program.
func Register(interfaceType InterfaceType, driver Driver) {
	if _, exists := drivers[interfaceType]; exists {
		log.Fatalln(interfaceType, "Driver already registered")
	}

	log.Println("Register", interfaceType, "driver")
	drivers[interfaceType] = driver
}

// Resource defines the interface for VISA Resources.
type ResourceFoo interface {
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

func createTCPIPResource(address string) (Resource, error) {
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
