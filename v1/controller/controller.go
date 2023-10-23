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
	defer C.free( unsafe.Pointer( connection ) )
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
	C.free( unsafe.Pointer( comm ) )
}

func ( ctrl Controller ) PowerOff() {
	connection := C.libcec_initialise( &ctrl.Configuration )
	defer C.free( unsafe.Pointer( connection ) )
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
	C.free( unsafe.Pointer( comm ) )
}
