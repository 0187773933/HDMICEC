package main

import (
	// "fmt"
	// time "time"
	controller "github.com/0187773933/HDMICEC/v1/controller"
)

func main() {
	ctrl := controller.New()
	// ctrl.SoftInit()

	// ctrl.PowerOff()
	ctrl.PowerOn()
	ctrl.SelectHDMI( 2 )

	// as := ctrl.GetActiveSource()
	// controller.PrettyPrint( as )

	// sources := ctrl.GetSources()
	// controller.PrettyPrint( sources )

	// ctrl.Mute()

	// power_status := ctrl.GetPowerStatus()
	// fmt.Println( "Already On ===" , power_status )
	// if power_status == false {
	// 	time.Sleep( 1 * time.Second )
	// 	fmt.Println( "Powering On" )
	// 	ctrl.PowerOn()
	// }

	// time.Sleep( 1 * time.Second )
	// // ctrl.PowerOn()
	// hdmi_number := 2
	// fmt.Printf( "Setting HDMI %d\n" , hdmi_number )
	// ctrl.SelectHDMI( hdmi_number )



	// ctrl.PowerOff()
	// ctrl.Mute()
}