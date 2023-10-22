import cec
import time
from pprint import pprint
import logging
logging.basicConfig( level=logging.DEBUG )  # or appropriate level


# https://github.com/trainman419/python-cec/issues/36
def send_opcode():
	opcode = cec.CEC_OPCODE_USER_CONTROL_PRESSED
	parameters = bytes([keycode])
	device.transmit(opcode, parameters)

adapters = cec.list_adapters()
print( adapters )
cec.init( adapters[ 0 ] )
devices = cec.list_devices()
print( devices )
print( "init done" )
# print( devices[ 0 ].osd_string )
device = cec.Device( cec.CECDEVICE_TV )
#device = cec.Device( devices[ 0 ] )
#is_on = device.is_on()
is_active = device.is_active()
device.standby()

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

