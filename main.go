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
	"strings"
	// time "time"
	// "sync"
	"strconv"
	"encoding/json"
	"bytes"
	"os/exec"
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

func PrettyPrint( input interface{} ) {
	jd , _ := json.MarshalIndent( input , "" , "  " )
	fmt.Println( string( jd ) )
}

func StringToInt( input string ) ( result int ) {
	result , _ = strconv.Atoi( input )
	return
}

func ( ctrl Controller ) GetPowerStatus() (result bool) {
	cmd := exec.Command("cec-client", "-s", "-d", "1")
	cmd.Stdin = bytes.NewBufferString("pow 0.0.0.0\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}

	output := out.String()
	if strings.Contains(output, "power status: on") {
		return true
	} else if strings.Contains(output, "power status: standby") {
		return false
	}

	return false
}

// echo 'scan' | cec-client -s -d 1
type Source struct {
	DeviceName string
	Address string
	HDMIInput int
	ActiveSource bool
	Vendor string
	OSDString string
	PowerStatus bool
	// HDMI int
}
func ( ctrl Controller ) GetActiveSource() ( result Source ) {
	cmd := exec.Command( "cec-client" , "-s" , "-d" , "1" )
	cmd.Stdin = bytes.NewBufferString( "scan\n" )
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil { return }
	output := out.String()
	lines := strings.Split( output , "\n" )
	var latest_device Source
	// var devices []Source
	for _ , line := range lines {
		// fmt.Println( i , "===" , line )
		if strings.Contains( line , "device #" ) {
			latest_device.DeviceName = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
		}
		if strings.Contains( line , "address:" ) {
			latest_device.Address = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			latest_device.HDMIInput = StringToInt( string( latest_device.Address[ 0 ] ) )
		}
		if strings.Contains( line , "active source:" ) {
			as_string := strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			switch as_string {
				case "yes":
					latest_device.ActiveSource = true
					break;
				case "no":
					latest_device.ActiveSource = false
					break;
			}
		}
		if strings.Contains( line , "vendor:" ) {
			latest_device.Vendor = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
		}
		if strings.Contains( line , "osd string:" ) {
			latest_device.OSDString = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
		}
		if strings.Contains( line , "power status:" ) {
			ps_string := strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			switch ps_string {
				case "on":
					latest_device.PowerStatus = true
					break;
				case "off":
					latest_device.PowerStatus = false
					break;
			}
			if latest_device.ActiveSource == true {
				return latest_device
			}
			// devices = append( devices , latest_device )
			// latest_device = Source{}
		}
	}
	return
}

func ( ctrl Controller ) GetSources() ( result []Source ) {
	cmd := exec.Command( "cec-client" , "-s" , "-d" , "1" )
	cmd.Stdin = bytes.NewBufferString( "scan\n" )
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil { return }
	output := out.String()
	lines := strings.Split( output , "\n" )
	var latest_device Source
	for _ , line := range lines {
		if strings.Contains( line , "device #" ) {
			latest_device.DeviceName = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
		}
		if strings.Contains( line , "address:" ) {
			latest_device.Address = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			latest_device.Address = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			latest_device.HDMIInput = StringToInt( string( latest_device.Address[ 0 ] ) )
		}
		if strings.Contains( line , "active source:" ) {
			as_string := strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			switch as_string {
				case "yes":
					latest_device.ActiveSource = true
					break;
				case "no":
					latest_device.ActiveSource = false
					break;
			}
		}
		if strings.Contains( line , "vendor:" ) {
			latest_device.Vendor = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
		}
		if strings.Contains( line , "osd string:" ) {
			latest_device.OSDString = strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
		}
		if strings.Contains( line , "power status:" ) {
			ps_string := strings.TrimSpace( strings.Split( line , ":" )[ 1 ] )
			switch ps_string {
				case "on":
					latest_device.PowerStatus = true
					break;
				case "off":
					latest_device.PowerStatus = false
					break;
			}
			result = append( result , latest_device )
			latest_device = Source{}
		}
	}
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

	fmt.Println( "sending poweron command" , 0 )
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
	if len( adapters ) < 1 {
		panic( "no adapters found" )
	}
	adapter := adapters[0]
	comm := C.CString( adapter.Comm )

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

func ( ctrl *Controller ) SelectHDMI( hdmi int ) {
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

	fmt.Printf( "Selecting HDMI %d\n" , hdmi )

	var command C.cec_command
	command.initiator = logical_address
	command.destination = 0xF // Broadcast to all devices.
	command.opcode_set = 1
	command.opcode = 0x82 // Opcode for "active source"
	command.parameters.size = 2
	// command.parameters.data[0] = 0x20  // New address (2.0.0.0), first part
	hdmi_result := ( hdmi << 4 ) | 0x00
	command.parameters.data[0] = C.uchar( hdmi_result )
	command.parameters.data[1] = 0x00  // New address (2.0.0.0), second part


	// Transmit the command
	if result := C.libcec_transmit(connection, &command); result == 0 {
		fmt.Println("Failed to send command")
	} else {
		fmt.Println("Command sent successfully")
	}
}

// func (ctrl *Controller) GetPowerStatus() (result string) {
//     connection := C.libcec_initialise(&ctrl.Configuration)
//     adapters := get_adapters(connection)
//     if len(adapters) < 1 {
//         panic("no adapters found")
//     }
//     adapter := adapters[0]
//     comm := C.CString(adapter.Comm)
//     defer C.free(unsafe.Pointer(comm))

//     if C.libcec_open(connection, comm, C.CEC_DEFAULT_CONNECT_TIMEOUT) == 0 {
//         panic("Failed to open a connection to the adapter")
//     }

//     // Assuming that you are querying the TV.
//     logical_address := C.cec_logical_address(0)

//     fmt.Println("Sending Power Status Inquiry")

//     var command C.cec_command
//     command.initiator = C.cec_logical_address(ctrl.Configuration.logicalAddress) // Assuming you've set this to your device's logical address.
//     command.destination = logical_address
//     command.opcode_set = 1
//     command.opcode = C.CEC_OPCODE_GIVE_DEVICE_POWER_STATUS
//     command.parameters.size = 0

//     if C.libcec_transmit(connection, &command) == 0 {
//         fmt.Println("Failed to send command")
//         return "Failed to send command"
//     }

//     // Now, you need to listen for the response.
//     // This is a simplified loop, in practice, you would have timeout handling, and possibly this would be event-driven.
//     for {
//         msg := C.libcec_receive_message(connection)
//         if msg.opcode == C.CEC_OPCODE_REPORT_POWER_STATUS {
//             // Assuming the status is in the first parameter and is a direct mapping of CEC power status codes.
//             switch msg.parameters.data[0] {
//             case 0:
//                 return "TV is on"
//             case 1:
//                 return "TV is in standby"
//             case 2:
//                 return "TV is in transition from standby to on"
//             case 3:
//                 return "TV is in transition from on to standby"
//             default:
//                 return "Unknown power status"
//             }
//         }
//     }
// }



func (ctrl *Controller) Mute() {
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

	fmt.Println("Muting TV")

	var command C.cec_command
	command.initiator = logical_address
	command.destination = C.cec_logical_address( 0 )
	command.opcode_set = 1
	command.opcode = 0x44 // Opcode for "active source"
	command.parameters.size = 1
	command.parameters.data[0] = 0x43 // "1.0" part of the address

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

// CONFIRMED WORKS FOR MUTE TOGGLE DON'T TOUCH
func Test() {

	// Reset Configuration
	var config C.libcec_configuration
	C.libcec_clear_configuration( &config )
	config.clientVersion = C.LIBCEC_VERSION_CURRENT
	fmt.Println( "config ===" , config )
	// https://github.com/trainman419/python-cec
	// osd_string ?
	device_name := C.CString( "CEC_GO" )
	defer C.free( unsafe.Pointer( device_name ) )
	device_name_length := C.strlen( device_name ) // + 1 ?
	C.strncpy( &config.strDeviceName[ 0 ] , device_name , device_name_length )
	config.deviceTypes.types[ 0 ] = C.CEC_DEVICE_TYPE_RECORDING_DEVICE

	// Required Initialization EVERY TIME for EVERY SENT COMMAND
	connection := C.libcec_initialise( &config )

	// Get the Adapter
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
	fmt.Println( "logical address ===" , logical_address )
	physical_address := C.libcec_get_device_physical_address( connection , logical_address )
	fmt.Println( "physical address ===" , physical_address , ( physical_address >> 8 ) , ( physical_address & 0xFF ) )
	logical_physical_address := C.cec_logical_address( physical_address )
	fmt.Println( "logical physical address ===" , logical_physical_address )

	// Send Command
	// https://github.com/Pulse-Eight/libcec/blob/bf5a97d7673033ef6228c63109f6baf2bdbe1a0c/include/cectypes.h#L1052
	// fmt.Println( "sending standby command" , 0 )

	// var command C.cec_command
	// command.initiator = logical_address
	// command.destination = C.cec_logical_address( 0 )
	// command.opcode_set = 1
	// command.opcode = 0x04

	// command.parameters.size = 2
	// command.parameters.data[0] = 0x20  // New address (2.0.0.0), first part
	// command.parameters.data[1] = 0x00  // New address (2.0.0.0), second part

	// its basically always this
	// initiator = 1
	// destination = 0
	var command C.cec_command
	command.initiator =  C.cec_logical_address( 1 )          // Replace with your source device logical address
	command.destination =  C.cec_logical_address( 0 )  // Replace with your target device logical address
	command.opcode_set = 1
	command.opcode = 0x44
	// command.ack = 0
	// command.eom = 1
	// command.transmit_timeout = 1000  // Set your desired timeout in milliseconds
	command.parameters.size = 1
	command.parameters.data[0] = 0x43;


	// // messageLength := len(message)
	result := C.libcec_transmit( connection , &command )
	fmt.Println( result )




	// Close
	C.libcec_close( connection )
	C.libcec_destroy( connection )
	C.free( unsafe.Pointer( comm ) )
}

// Fishing for Menu Activation Requests
// echo "tx 1F:44:09" | cec-client -s -d 31 && echo "tx 1F:45" | cec-client -s -d 31

// cecInit
// getAdapters
// openAdapter
// func TestTwo() {

// 	var connection C.libcec_connection_t
// 	var conf C.libcec_configuration
// 	C.libcec_clear_configuration( &conf )
// 	conf.clientVersion = C.uint32_t(C.LIBCEC_VERSION_CURRENT)
// 	conf.deviceTypes.types[0] = C.CEC_DEVICE_TYPE_RECORDING_DEVICE
// 	C.setName( &conf , C.CString( "CEC_GO" ) )
// 	// if printLogs {
// 	// 	C.setupCallbacks(&conf)
// 	// }
// 	connection = C.libcec_initialise( &conf )
// 	if connection == C.libcec_connection_t( nil ) {
// 		panic( "asdf" )
// 	}

// 	C.libcec_init_video_standalone( connection )

// 	adapters := get_adapters( connection )
// 	if len( adapters ) < 1 { panic( "no adapters found ??" ) }
// 	adapter := adapters[ 0 ]
// 	result := C.libcec_open( connection , C.CString( adapter.Comm ) , C.CEC_DEFAULT_CONNECT_TIMEOUT )
// 	fmt.Println( result )
// 	fmt.Println( time.Second )


// 	// C.libcec_power_on_devices( connection , 0 )
// 	C.libcec_set_hdmi_port( connection , 0 , 2 )
// 	// C.libcec_standby_devices( connection , 0 )


// 	// C.libcec_set_active_source( connection , 0 )
// 	// C.libcec_mute_audio( connection  , 0x00 )

// 	// C.libcec_send_keypress( connection , 0xF , C.cec_user_control_code( 9 ) , 1 )
// 	// time.Sleep( 500 * time.Millisecond )
// 	// C.libcec_send_key_release( connection , 0xF , 1 )

// 	// Close
// 	C.libcec_close( connection )
// 	C.libcec_destroy( connection )
// 	C.free( unsafe.Pointer( C.CString( adapter.Comm ) ) )


// 	// // Reset Configuration
// 	// var config C.libcec_configuration
// 	// C.libcec_clear_configuration( &config )
// 	// config.clientVersion = C.LIBCEC_VERSION_CURRENT
// 	// // fmt.Println( "config ===" , config )
// 	// // https://github.com/trainman419/python-cec
// 	// // osd_string ?
// 	// device_name := C.CString( "CEC_GO" )
// 	// defer C.free( unsafe.Pointer( device_name ) )
// 	// device_name_length := C.strlen( device_name ) // + 1 ?
// 	// C.strncpy( &config.strDeviceName[ 0 ] , device_name , device_name_length )
// 	// config.deviceTypes.types[ 0 ] = C.CEC_DEVICE_TYPE_RECORDING_DEVICE

// 	// Required Initialization EVERY TIME for EVERY SENT COMMAND
// 	// connection := C.libcec_initialise( &config )



// 	// Get the Adapter
// 	// adapters := get_adapters( connection )
// 	// if len( adapters ) < 1 { panic( "no adapters found ??" ) }
// 	// adapter := adapters[ 0 ]
// 	// fmt.Println( reflect.TypeOf( adapter ) )
// 	// comm := C.CString( adapter.Comm ) // its like golang garbage collects this
// 	// fmt.Println( reflect.TypeOf( comm ) )
// 	// if C.libcec_open( connection , comm , C.CEC_DEFAULT_CONNECT_TIMEOUT ) == 0 {
// 	// 	panic( "Failed to open a connection to the adapter" )
// 	// }
// 	// logical_addresses := C.libcec_get_logical_addresses( connection )
// 	// logical_address := C.cec_logical_address( byte( logical_addresses.primary ) )
// 	// fmt.Println( "logical address ===" , logical_address )
// 	// physical_address := C.libcec_get_device_physical_address( connection , logical_address )
// 	// fmt.Println( "physical address ===" , physical_address , ( physical_address >> 8 ) , ( physical_address & 0xFF ) )
// 	// logical_physical_address := C.cec_logical_address( physical_address )
// 	// fmt.Println( "logical physical address ===" , logical_physical_address )

// 	// Send Command
// 	// https://github.com/Pulse-Eight/libcec/blob/bf5a97d7673033ef6228c63109f6baf2bdbe1a0c/include/cectypes.h#L1052
// 	// fmt.Println( "sending standby command" , 0 )

// 	// var command C.cec_command
// 	// command.initiator = logical_address
// 	// command.destination = C.cec_logical_address( 0 )
// 	// command.opcode_set = 1
// 	// command.opcode = 0x04

// 	// command.parameters.size = 2
// 	// command.parameters.data[0] = 0x20  // New address (2.0.0.0), first part
// 	// command.parameters.data[1] = 0x00  // New address (2.0.0.0), second part

// 	// its basically always this
// 	// initiator = 1
// 	// destination = 0

// 	//; https://github.com/xbmc/xbmc/blob/f3f1df1eab2a38b7039e57635a6597b37510481a/xbmc/peripherals/devices/PeripheralCecAdapter.cpp#L897

// 	// // echo "tx 1F:44:09" | cec-client -s -d 31 && echo "tx 1F:45" | cec-client -s -d 31

// 	// var command C.cec_command
// 	// command.initiator =  0          // Replace with your source device logical address
// 	// // command.destination =  C.cec_logical_address( 0 )  // Replace with your target device logical address
// 	// command.destination = 0xF
// 	// command.opcode_set = 1
// 	// command.opcode = 0x44
// 	// // command.ack = 0
// 	// // command.eom = 0
// 	// command.transmit_timeout = 3000  // Set your desired timeout in milliseconds
// 	// command.parameters.size = 1
// 	// command.parameters.data[0] = 0x09;

// 	// var command_release C.cec_command
// 	// command_release.initiator = 0         // Replace with your source device logical address
// 	// command.destination = 0xF
// 	// command_release.opcode_set = 1
// 	// command_release.opcode = 0x45
// 	// // command_release.ack = 0
// 	// // command_release.eom = 1
// 	// command_release.transmit_timeout = 3000  // Set your desired timeout in milliseconds
// 	// command_release.parameters.size = 0

// 	// result := C.libcec_transmit( connection , &command )
// 	// time.Sleep( 1 * time.Second )
// 	// result_release := C.libcec_transmit( connection , &command_release )
// 	// fmt.Println( result , result_release )

// 	// fmt.Println( "C.CEC_USER_CONTROL_CODE_TOP_MENU" , C.CEC_USER_CONTROL_CODE_TOP_MENU )


// 	// SEND KEYPRESS ???
// 	// https://github.com/chbmuc/cec/blob/master/cec.go#L89
// 	// https://github.com/chbmuc/cec/blob/master/libcec.go#L187
// 	// https://github.com/Pulse-Eight/libcec/blob/bf5a97d7673033ef6228c63109f6baf2bdbe1a0c/src/libcec/LibCECC.cpp#L355
// 	// C.libcec_send_keypress( connection , 0xF , C.cec_user_control_code( key ) , 1 )
// 	// C.libcec_send_keypress( connection , 0xF , C.cec_user_control_code( 9 ) , 1 )
// 	// time.Sleep( 500 * time.Millisecond )
// 	// C.libcec_send_key_release( connection , 0xF , 1 )

// 	// fmt.Println( time.Second )
// 	// C.libcec_audio_mute( connection )

// 	// // Close
// 	// C.libcec_close( connection )
// 	// C.libcec_destroy( connection )
// 	// C.free( unsafe.Pointer( C.CString(adapter.Comm) ) )
// }


// https://github.com/DrGeoff/cec_simplest/blob/master/cec-simplest.cpp
// func TestThree() {

// 	// Reset Configuration
// 	var config C.libcec_configuration
// 	C.libcec_clear_configuration( &config )
// 	config.clientVersion = C.LIBCEC_VERSION_CURRENT
// 	// config.bActivateSource = 1
// 	config.deviceTypes.types[ 0 ] = C.CEC_DEVICE_TYPE_RECORDING_DEVICE
// 	fmt.Println( "config ===" , config )
// 	C.setName( &config , C.CString( "CEC_GO" ) )
// 	connection := C.libcec_initialise( &config )
// 	if connection == C.libcec_connection_t( nil ) {
// 		panic( "couldn't init libcec" )
// 	}
// 	// C.libcec_init_video_standalone( connection )

// 	// Get USB Adapter
// 	adapters := get_adapters( connection )
// 	if len( adapters ) < 1 { panic( "no adapters found ??" ) }
// 	adapter := adapters[ 0 ]

// 	// Open USB Adapter
// 	result := C.libcec_open( connection , C.CString( adapter.Comm ) , C.CEC_DEFAULT_CONNECT_TIMEOUT )
// 	if result != 1 { panic( "couldn't connect to adapter" ) }

// 	logical_addresses := C.libcec_get_logical_addresses( connection )
// 	logical_address := C.cec_logical_address( byte( logical_addresses.primary ) )
// 	fmt.Println( "logical address ===" , logical_address )
// 	physical_address := C.libcec_get_device_physical_address( connection , logical_address )
// 	fmt.Println( "physical address ===" , physical_address , ( physical_address >> 8 ) , ( physical_address & 0xFF ) )
// 	logical_physical_address := C.cec_logical_address( physical_address )
// 	fmt.Println( "logical physical address ===" , logical_physical_address )


// 	// Do Whatever
// 	// C.libcec_set_hdmi_port( connection , 0 , 1 )
// 	// C.libcec_mute_audio( connection  , 0x00 )

// 	// Press and Release the Menu Key
// 	for i := 0; i < 255; i++ {
// 		keyCode := C.cec_user_control_code(i) // Cast the iterator to cec_user_control_code
// 		fmt.Println("Sending key press for", keyCode)
// 		C.libcec_send_keypress( connection , C.CECDEVICE_TV , keyCode , 1 ) // Replace C.CECDEVICE_TV with the appropriate device
// 		C.libcec_send_key_release( connection , C.CECDEVICE_TV , 1 )       // Replace C.CECDEVICE_TV with the appropriate device
// 		time.Sleep( 300 * time.Millisecond )
// 	}


// 	// var command C.cec_command
// 	// command.initiator =  0          // Replace with your source device logical address
// 	// // command.destination =  C.cec_logical_address( 0 )  // Replace with your target device logical address
// 	// command.destination = 0xF
// 	// command.opcode_set = 1
// 	// command.opcode = 0x44
// 	// // command.ack = 0
// 	// // command.eom = 0
// 	// command.transmit_timeout = 3000  // Set your desired timeout in milliseconds
// 	// command.parameters.size = 1
// 	// command.parameters.data[0] = 0x09;


// 	// Close
// 	C.libcec_close( connection )
// 	C.libcec_destroy( connection )
// 	C.free( unsafe.Pointer( C.CString( adapter.Comm ) ) )
// }