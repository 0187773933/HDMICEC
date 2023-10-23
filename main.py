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

if __name__ == "__main__":

	# 1.) Get Adapters
	adapters = cec.list_adapters()
	if len( adapters ) < 0:
		print( "no adapters found" )
		sys.exit( 1 )
	print( adapters )
	adapter = adapters[0]
	print( "Using Adapter %s"%(adapter) )

	# 2.) Hook Adapter with LibCEC
	cec.init( adapter )

	# 3.) Get Devices on Adapter
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

	# device.standby()
	# device.set_av_input( 0 )



# adapters = cec.list_adapters()
# print( adapters )
# cec.init( adapters[ 0 ] )
# devices = cec.list_devices()
# print( devices )
# print( "init done" )
# # print( devices[ 0 ].osd_string )
# device = cec.Device( cec.CECDEVICE_TV )
# #device = cec.Device( devices[ 0 ] )
# #is_on = device.is_on()
# is_active = device.is_active()
# device.standby()

# device.power_on()
#device.set_av_input( 2 )
# time.sleep( 2 )
# cec.set_active_source( 1 )




# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# print( "half way" )
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()
# cec.volume_down()

