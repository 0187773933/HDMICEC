package controller

/*
#include <stdio.h>
#cgo CFLAGS: -I/usr/local/Cellar/libcec/6.0.2/include
#cgo LDFLAGS: -L/usr/local/Cellar/libcec/6.0.2/lib -lcec
#cgo pkg-config: libcec
#include <libcec/cecc.h>

void setDeviceName(libcec_configuration *config, char *name)
{
	snprintf(config->strDeviceName, sizeof(config->strDeviceName), "%s", name);
}
*/
import "C"

import (
	"fmt"
)

type Adapter struct {
	Path string
	Comm string
}
type Controller struct {
	Configuration C.libcec_configuration
	Connection C.libcec_connection_t
	Adapter Adapter
	LogicalAddress C.cec_logical_address
	PhysicalAddress [2]byte
	// Callbacks *callbacks
}

// sudo find / -name cecc.h
// /usr/local/Cellar/libcec/6.0.2/include/libcec/cecc.h
func New() ( result Controller ) {
	C.libcec_clear_configuration( &result.Configuration )
	fmt.Println( result.Configuration )
	result.Configuration.clientVersion = C.LIBCEC_VERSION_CURRENT
	// https://github.com/trainman419/python-cec
	// osd_string ?
	device_name := "TV"
	C.setDeviceName( &result.Configuration , C.CString( device_name ) )
	result.Configuration.deviceTypes.types[ 0 ] = C.CEC_DEVICE_TYPE_RECORDING_DEVICE
	return
}


func ( ctrl Controller ) GetAdapters() ( result []Adapter ) {
	var foundDevices [ 10 ]C.cec_adapter
	count := int( C.libcec_find_adapters( ctrl.Connection , &foundDevices[ 0 ] , C.uchar( len( foundDevices ) ) , nil ) )
	fmt.Println( "adapters found ===" , count )
	for i := 0; i < count; i++ {
		xa := Adapter{
			Path: C.GoString( &foundDevices[i].path[ 0 ] ) ,
			Comm: C.GoString( &foundDevices[i].comm[ 0 ] ) ,
		}
		result = append( result , xa )
	}
	return
}


func ( ctrl Controller ) Connect()  {
	ctrl.Connection = C.libcec_initialise( &ctrl.Configuration )
	adapters := ctrl.GetAdapters()
	if len( adapters ) < 1 { panic( "no adapters found ??" ) }
	ctrl.Adapter = adapters[ 0 ]
	C.libcec_open( ctrl.Connection , C.CString( ctrl.Adapter.Comm ) , C.CEC_DEFAULT_CONNECT_TIMEOUT )

	// C.cec_logical_address(address)
	logical_addresses := C.libcec_get_logical_addresses( ctrl.Connection )
	ctrl.LogicalAddress = C.cec_logical_address( byte( logical_addresses.primary ) )
	// active_devices := C.libcec_get_active_devices( ctrl.Connection )
	// fmt.Println( active_devices )
	physical_address := C.libcec_get_device_physical_address( ctrl.Connection , ctrl.LogicalAddress )
	ctrl.PhysicalAddress[ 0 ] = byte( physical_address >> 8 )
	ctrl.PhysicalAddress[ 1 ] = byte( physical_address & 0xFF )
	fmt.Println( ctrl.PhysicalAddress )
	// ctrl.GetVendor()
	// ctrl.GetActiveSource()
}

// func ( ctrl Controller ) GetActiveSource() {
// 	result := C.libcec_get_active_source( ctrl.Connection )
// 	fmt.Println( result )
// }

// func ( ctrl Controller ) GetVendor() {
// 	vendor := C.libcec_get_device_vendor_id( ctrl.Connection , ctrl.LogicalAddress )
// 	var toString *C.char
// 	C.libcec_vendor_id_to_string( C.cec_vendor_id( vendor ) , toString , 50 )
// 	C.GoString( toString )
// 	fmt.Println( toString )
// }

func ( ctrl Controller ) PowerOn() {
	C.libcec_power_on_devices( ctrl.Connection , ctrl.LogicalAddress )
}

func ( ctrl Controller ) PowerOff() {
	C.libcec_standby_devices( ctrl.Connection ,  ctrl.LogicalAddress )
}
