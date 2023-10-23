package main

import (
	// "fmt"
	time "time"
	controller "github.com/0187773933/HDMICEC/v1/controller"
)

func main() {
	ctrl := controller.New()
	// ctrl.SoftInit()
	// ctrl.PowerOn()
	// time.Sleep( 10 * time.Second )
	// ctrl.PowerOff()
	ctrl.SelectHDMI1()
	time.Sleep( 3 * time.Second )
	ctrl.SelectHDMI2()
	time.Sleep( 3 * time.Second )
	ctrl.SelectHDMI1()
	time.Sleep( 3 * time.Second )
	ctrl.SelectHDMI2()
}
