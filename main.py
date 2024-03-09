#!/usr/bin/env python3
import sys
import cec
import time
from pprint import pprint
import logging
logging.basicConfig( level=logging.DEBUG )  # or appropriate level

# https://github.com/trainman419/python-cec/issues/36
def send_opcode( opcode , params=bytes() ):
	device.transmit( opcode , parameters )

def send_keycode( keycode ):
	opcode = cec.CEC_OPCODE_USER_CONTROL_PRESSED
	parameters = bytes( [ keycode ] )
	send_opcode( opcode , parameters )

def cec_init_first_adapter():
	try:
		# 1.) Get Adapters
		adapters = cec.list_adapters()
		if len( adapters ) < 0:
			print( "no adapters found" )
			sys.exit( 1 )
		print( adapters )
		adapter = adapters[ 0 ]
		print( f"Using Adapter {adapter}" )

		# 2.) Hook Adapter with LibCEC
		cec.init( adapter )
		return True
	except Exception as e:
		print( e )
		return False

def get_tv():
	ready = cec_init_first_adapter()
	if ready == False:
		sys.exit( 1 )
	devices = cec.list_devices() # this has to be called here , it like wakes up the adapter somehow
	return cec.Device( cec.CECDEVICE_TV )

def print_devices():
	ready = cec_init_first_adapter()
	if ready == False:
		sys.exit( 1 )

	devices = cec.list_devices()
	if len( devices ) < 0:
		print( "no devices found" )
		sys.exit( 1 )
	devices_keys = devices.keys()
	for i , k in enumerate( devices_keys ):
		d = devices[ k ]
		print( "\nAddress:" , d.address )
		print( "Physical Address:" , d.physical_address )
		print( "Vendor ID:" , d.vendor )
		print( "OSD:" , d.osd_string )
		print( "CEC Version:" , d.cec_version )
		print( "Language:" , d.language )

def turn_on_and_goto_hdmi_input( hdmi_input ):
	tv = get_tv()
	if tv is None:
		print( "no tv found" )
		sys.exit( 1 )
	is_on = tv.is_on()
	is_active = tv.is_active()
	print( "tv is_on" , is_on )
	print( "tv is_active" , is_active )
	print( f"Turning on TV Power and Going to HDMI {hdmi_input}" )
	try:
		tv.power_on()
		cec.set_active_source( hdmi_input )
	except Exception as e:
		print( e )
		time.sleep( 10 )
		tv.power_on()
		cec.set_active_source( hdmi_input )
	time.sleep( 10 )
	cec.set_active_source( hdmi_input )
	return tv

if __name__ == "__main__":
	# print_devices()
	tv = turn_on_and_goto_hdmi_input( 2 ) # just goes to pulse-eight no matter what
	# tv.standby()
	# tv.set_av_input( 0 )
	cec.volume_down() # this is not working , Samsung , LG