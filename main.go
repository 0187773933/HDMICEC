package main

import (
	// "fmt"
	time "time"
	controller "github.com/0187773933/HDMICEC/v1/controller"
)

func main() {
	ctrl := controller.New()
	ctrl.PowerOn()
	time.Sleep( 10 * time.Second )
	ctrl.PowerOff()
	// ctrl.PowerOffV2()
}
