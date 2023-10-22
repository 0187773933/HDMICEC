package main

import (
	"fmt"
	controller "github.com/0187773933/HDMICEC/v1/controller"
)

func main() {
	ctrl := controller.New()
	ctrl.Connect()
	fmt.Println( ctrl )
	// ctrl.PowerOn()
	ctrl.PowerOff()
}
