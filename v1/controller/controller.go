package controller

/*
#cgo CFLAGS: -I/usr/local/Cellar/libcec/6.0.2/include
#cgo LDFLAGS: -L/usr/local/Cellar/libcec/6.0.2/lib -lcec
#cgo pkg-config: libcec
#include <libcec/cecc.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void setDeviceName(libcec_configuration *config, char *name)
{
	snprintf(config->strDeviceName, sizeof(config->strDeviceName), "%s", name);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
	"reflect"
	// time "time"
	// "sync"
)

type Adapter struct {
	Path string
	Comm string
}
type Controller struct {
	Configuration C.libcec_configuration
	Connection C.libcec_connection_t
	ConnComm *C.char
	Adapter Adapter
	LogicalAddress C.cec_logical_address
	LogicalPhysicalAddress C.cec_logical_address
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
	device_name := C.CString( "CEC_GO_01" )
	defer C.free( unsafe.Pointer( device_name ) )
	device_name_length := C.strlen( device_name ) + 1
	C.strncpy( &result.Configuration.strDeviceName[ 0 ] , device_name , device_name_length )
	result.Configuration.deviceTypes.types[ 0 ] = C.CEC_DEVICE_TYPE_RECORDING_DEVICE

	// result.Connection = C.libcec_initialise( &result.Configuration )
	return
}

func get_adapters( conn C.libcec_connection_t ) ( result []Adapter ) {
	var found_devices [ 10 ]C.cec_adapter
	count := int( C.libcec_find_adapters( conn , &found_devices[ 0 ] , C.uchar( len( found_devices ) ) , nil ) )
	fmt.Println( "adapters found ===" , count )
	for i := 0; i < count; i++ {
		xa := Adapter{
			Path: C.GoString( &found_devices[i].path[ 0 ] ) ,
			Comm: C.GoString( &found_devices[i].comm[ 0 ] ) ,
		}
		result = append( result , xa )
	}
	return
}

func ( ctrl Controller ) PowerOn() {
	connection := C.libcec_initialise( &ctrl.Configuration )
	adapters := get_adapters( connection )
	if len( adapters ) < 1 { panic( "no adapters found ??" ) }
	adapter := adapters[ 0 ]
	fmt.Println( reflect.TypeOf( adapter ) )
	comm := C.CString( adapter.Comm ) // its like golang garbage collects this
	fmt.Println( reflect.TypeOf( comm ) )
	if C.libcec_open( connection , comm , C.CEC_DEFAULT_CONNECT_TIMEOUT ) == 0 {
		panic( "Failed to open a connection to the adapter" )
	}
	logical_addresses := C.libcec_get_logical_addresses( connection )
	logical_address := C.cec_logical_address( byte( logical_addresses.primary ) )
	physical_address := C.libcec_get_device_physical_address( connection , logical_address )
	fmt.Println( "physical address ===" , physical_address , ( physical_address >> 8 ) , ( physical_address & 0xFF ) )
	logical_physical_address := C.cec_logical_address( physical_address )
	fmt.Println( "physical address ===" , logical_physical_address )

	fmt.Println( "sending standby command" , 0 )
	var command C.cec_command
	command.initiator = logical_address
	command.destination = C.cec_logical_address( 0 )
	command.opcode_set = 1
	command.opcode = 0x04
	// // messageLength := len(message)
	result := C.libcec_transmit( connection , &command )
	fmt.Println( result )
	C.libcec_close( connection )
	C.libcec_destroy( connection )
	C.free( unsafe.Pointer( comm ) )
}

func ( ctrl Controller ) PowerOff() {
	connection := C.libcec_initialise( &ctrl.Configuration )
	adapters := get_adapters( connection )
	if len( adapters ) < 1 { panic( "no adapters found ??" ) }
	adapter := adapters[ 0 ]
	fmt.Println( reflect.TypeOf( adapter ) )
	comm := C.CString( adapter.Comm ) // its like golang garbage collects this
	fmt.Println( reflect.TypeOf( comm ) )
	if C.libcec_open( connection , comm , C.CEC_DEFAULT_CONNECT_TIMEOUT ) == 0 {
		panic( "Failed to open a connection to the adapter" )
	}
	logical_addresses := C.libcec_get_logical_addresses( connection )
	logical_address := C.cec_logical_address( byte( logical_addresses.primary ) )
	physical_address := C.libcec_get_device_physical_address( connection , logical_address )
	fmt.Println( "physical address ===" , physical_address , ( physical_address >> 8 ) , ( physical_address & 0xFF ) )
	logical_physical_address := C.cec_logical_address( physical_address )
	fmt.Println( "physical address ===" , logical_physical_address )

	fmt.Println( "sending standby command" , 0 )
	var command C.cec_command
	command.initiator = logical_address
	command.destination = C.cec_logical_address( 0 )
	command.opcode_set = 1
	command.opcode = 0x36
	// // messageLength := len(message)
	result := C.libcec_transmit( connection , &command )
	fmt.Println( result )
	C.libcec_close( connection )
	C.libcec_destroy( connection )
	C.free( unsafe.Pointer( comm ) )
}

func (ctrl *Controller) SelectHDMI1() {
	connection := C.libcec_initialise(&ctrl.Configuration)
	adapters := get_adapters(connection)
	if len(adapters) < 1 {
		panic("no adapters found")
	}
	adapter := adapters[0]
	comm := C.CString(adapter.Comm)

	if C.libcec_open(connection, comm, C.CEC_DEFAULT_CONNECT_TIMEOUT) == 0 {
		panic("Failed to open a connection to the adapter")
	}

	// Get the logical address of the current device.
	logical_addresses := C.libcec_get_logical_addresses(connection)
	logical_address := C.cec_logical_address(byte(logical_addresses.primary))

	fmt.Println("Selecting HDMI 1")

	var command C.cec_command
	command.initiator = logical_address
	command.destination = 0xF // Broadcast to all devices.
	command.opcode_set = 1
	command.opcode = C.CEC_OPCODE_ACTIVE_SOURCE // Opcode for "active source"
	command.parameters.size = 2

	// This is the physical address for HDMI 1 input, typically "1.0.0.0".
	// This might need to be changed depending on your device's configuration.
	command.parameters.data[0] = 0x10 // "1.0" part of the address
	command.parameters.data[1] = 0x00 // ".0.0" part of the address

	// Transmit the command
	if result := C.libcec_transmit(connection, &command); result == 0 {
		fmt.Println("Failed to send command")
	} else {
		fmt.Println("Command sent successfully")
	}

	C.libcec_close( connection )
	C.libcec_destroy( connection )
	// C.free( unsafe.Pointer( connection ) )
	C.free( unsafe.Pointer( comm ) )
}

func (ctrl *Controller) SelectHDMI2() {
	connection := C.libcec_initialise(&ctrl.Configuration)
	adapters := get_adapters(connection)
	if len(adapters) < 1 {
		panic("no adapters found")
	}
	adapter := adapters[0]
	comm := C.CString(adapter.Comm)
	if C.libcec_open(connection, comm, C.CEC_DEFAULT_CONNECT_TIMEOUT) == 0 {
		panic("Failed to open a connection to the adapter")
	}

	// Get the logical address of the current device.
	logical_addresses := C.libcec_get_logical_addresses(connection)
	logical_address := C.cec_logical_address(byte(logical_addresses.primary))

	fmt.Println("Selecting HDMI 2")

	var command C.cec_command
	command.initiator = logical_address
	command.destination = 0xF // Broadcast to all devices.
	command.opcode_set = 1
	command.opcode = 0x82 // Opcode for "active source"
	command.parameters.size = 2
	command.parameters.data[0] = 0x20  // New address (2.0.0.0), first part
	command.parameters.data[1] = 0x00  // New address (2.0.0.0), second part


	// Transmit the command
	if result := C.libcec_transmit(connection, &command); result == 0 {
		fmt.Println("Failed to send command")
	} else {
		fmt.Println("Command sent successfully")
	}
	// C.libcec_transmit(connection, &command)
	// C.libcec_transmit(connection, &command)
	// C.libcec_transmit(connection, &command)
	// C.libcec_transmit(connection, &command)
	// C.libcec_transmit(connection, &command)

	// time.Sleep( 1 * time.Second )
	// C.libcec_close( connection )
	// C.libcec_destroy( connection )
	// time.Sleep( 1 * time.Second )

	// C.free( unsafe.Pointer( connection ) )
	// C.free( unsafe.Pointer( comm ) )
}
