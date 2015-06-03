package visa_test

import (
	"fmt"

	"github.com/gotmc/visa"
)

func ExampleNewVisaResource() {
	myVisaResource, _ := visa.NewVisaResource("usb0::2391::1031::MY44123456::INSTR")
	fmt.Println(myVisaResource.ResourceString)
	fmt.Println(myVisaResource.SerialNumber)
	fmt.Println(myVisaResource.ModelCode)
	// Output:
	// usb0::2391::1031::MY44123456::INSTR
	// MY44123456
	// 1031
}
