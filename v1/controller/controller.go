package controller

/*
#cgo CFLAGS: -I/usr/local/Cellar/libcec/6.0.2/include
#cgo LDFLAGS: -L/usr/local/Cellar/libcec/6.0.2/lib -lcec
#cgo pkg-config: libcec
#include <libcec/cecc.h>
#include <stdio.h>

void setDeviceName(libcec_configuration *config, char *name)
{
	snprintf(config->strDeviceName, sizeof(config->strDeviceName), "%s", name);
}
*/
import "C"

import (
	"fmt"
	// "sync"
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
	fmt.Println( "config ===" , result.Configuration )
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

      // data.initiator = adapter->GetLogicalAddresses().primary;
      // data.destination = self->addr;

	// C.cec_logical_address(address)
	logical_addresses := C.libcec_get_logical_addresses( ctrl.Connection )
	ctrl.LogicalAddress = C.cec_logical_address( byte( logical_addresses.primary ) )

	active_devices := C.libcec_get_active_devices( ctrl.Connection )
	fmt.Println( "devices ===" , active_devices )

	physical_address := C.libcec_get_device_physical_address( ctrl.Connection , ctrl.LogicalAddress )
	// ctrl.PhysicalAddress[ 0 ] = byte( physical_address >> 8 )
	// ctrl.PhysicalAddress[ 1 ] = byte( physical_address & 0xFF )
	fmt.Println( "physical address ===" , physical_address , ( physical_address >> 8 ) , ( physical_address & 0xFF ) )
	logical_physical_address := C.cec_logical_address( physical_address )
	fmt.Println( "physical address ===" , logical_physical_address )

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



// func (conn *Connection) Transmit(message Message) {
// 	var command C.cec_command

// 	messageLength := len(message)

// 	if messageLength > 0 {
// 		command.initiator = C.cec_logical_address(message.Source())
// 		command.destination = C.cec_logical_address(message.Destination())

// 		if messageLength > 1 {
// 			command.opcode_set = 1
// 			command.opcode = C.cec_opcode(message.Opcode())
// 		} else {
// 			command.opcode_set = 0
// 		}


// 		if parameters := message.Parameters(); parameters != nil {
// 			command.parameters.size = C.uint8_t(len(parameters))
// 			for i, val := range parameters {
// 				command.parameters.data[i] = C.uint8_t(val)
// 			}
// 		} else {
// 			command.parameters.size = 0
// 		}
// 	}

// 	C.libcec_transmit(conn.connection, &command)
// }

// /libcec/include/cectypes.h:
//   745    CECDEVICE_FREEUSE          = 14,
//   746    CECDEVICE_UNREGISTERED     = 15,
//   747:   CECDEVICE_BROADCAST        = 15

// https://github.com/julemand101/cec_dart/blob/master/lib/src/libcec_enum/CEC_opcode.dart

func ( ctrl Controller ) PowerOff() {
	// success = adapter->StandbyDevices(self->addr);
	// return_code := C.libcec_standby_devices(ctrl.Connection, ctrl.LogicalAddress)

	// device->SetPowerStatus(CEC_POWER_STATUS_STANDBY);
	// 0x01
	// return_code := C.libcec_standby_devices(ctrl.Connection, ctrl.LogicalAddress)
	// success = adapter->StandbyDevices(self->addr);
	// fmt.Printf( "CEC standby command : %d\n" , return_code )
	// OpcodeStandby                   Opcode = 0x36

// 		command.initiator = C.cec_logical_address(message.Source())
// 		command.destination = C.cec_logical_address(message.Destination())

	for i := 0; i < 16; i++ {
		fmt.Println( "sending standby command" , i )
		var command C.cec_command
		command.initiator = ctrl.LogicalAddress
		command.destination = C.cec_logical_address( i )
		command.opcode_set = 1
		command.opcode = 0x36
		// // messageLength := len(message)
		C.libcec_transmit( ctrl.Connection , &command )
	}



	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()  // Signal that we are done at the end of the goroutine.
	// 	retCode := C.libcec_standby_devices(ctrl.Connection, ctrl.LogicalAddress)
	// 	if retCode < 0 {  // Or whatever indicates an error in your case.
	// 		fmt.Printf("CEC standby command failed with code: %d\n", retCode)
	// 		return
	// 	}
	// }()
	// wg.Wait()
}
