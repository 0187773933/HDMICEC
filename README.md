# HDMI-CEC Controller

## Docker
- `sudo apt-get install cec-utils`

- https://github.com/Pulse-Eight/libcec
- https://github.com/Pulse-Eight/libcec/blob/bf5a97d7673033ef6228c63109f6baf2bdbe1a0c/src/libcec/adapter/Pulse-Eight/USBCECAdapterCommands.cpp
- https://github.com/julemand101/cec_dart/blob/master/lib/src/libcec_enum/CEC_opcode.dart
- https://github.com/RobertMe/gocec/blob/master/connection.go
- https://github.com/trainman419/python-cec
- `find /usr/local -name "cec.h"`
- `/usr/local/Cellar/libcec/6.0.2/include/libcec/cec.h`
- `export CFLAGS="-I/usr/local/Cellar/libcec/6.0.2/include"`
- `python3 -m pip install --user cec`
- `find /usr/local/Cellar/libcec -name libcec.6.dylib`
- `export DYLD_LIBRARY_PATH="/usr/local/Cellar/libcec/6.0.2/lib:${DYLD_LIBRARY_PATH}"`
- https://justaddpower.happyfox.com/kb/article/68-cec-over-ip-control
- https://github.com/Pulse-Eight/libcec/issues/223?

## Vendor IDS

- 0f:87:00:00:f0
- Samsung = 240
- LG = ?

## Linux libcec [install](https://github.com/Pulse-Eight/libcec/blob/master/docs/README.linux.md)

```
sudo apt-get update && \
sudo apt-get install -y cmake libudev-dev libxrandr-dev \
python3-dev swig libp8-platform-dev && \
git clone https://github.com/Pulse-Eight/libcec.git && \
cd libcec && mkdir build && cd build && cmake .. && \
make -j4 && \
sudo make install && \
sudo ldconfig
```

- https://github.com/chbmuc/cec/blob/master/cec.go
- https://github.com/znkr/cec
- https://github.com/jsolla/cec-decoderwd
- https://engineering.purdue.edu/ece477/Archive/2012/Spring/S12-Grp10/Datasheets/CEC_HDMI_Specification.pdf
- https://www.linuxuprising.com/2019/07/raspberry-pi-power-on-off-tv-connected.html
- https://blog.gordonturner.com/2016/12/14/using-cec-client-on-a-raspberry-pi/
- https://www.uli-eckhardt.de/vdr/cec.en.shtml
- https://github.com/torvalds/linux/blob/052d534373b7ed33712a63d5e17b2b6cdbce84fd/Documentation/driver-api/media/cec-core.rst#L6
- https://man.archlinux.org/man/cec-ctl.1.en
- https://help.crestron.com/cds/main/index.htm#Rio_Media_Zones/Controlling_devices_by_sending_CEC_commands.htm
- https://github.com/xbmc/xbmc/blob/f3f1df1eab2a38b7039e57635a6597b37510481a/xbmc/peripherals/devices/PeripheralCecAdapter.cpp#L67

## We just need to find the vendor specific codes somehow. XBOX knows these


### Send Keypress

```
// https://github.com/chbmuc/cec/blob/master/cec.go#L89
// https://github.com/chbmuc/cec/blob/master/libcec.go#L187
// https://github.com/Pulse-Eight/libcec/blob/bf5a97d7673033ef6228c63109f6baf2bdbe1a0c/src/libcec/LibCECC.cpp#L355
// echo "tx 1F:44:09" | cec-client -s -d 31 && echo "tx 1F:45" | cec-client -s -d 31
C.libcec_send_keypress( connection , 0xF , C.cec_user_control_code( 9 ) , 1 )
```

### Monitor Mode ???

```
/usr/local/bin/cec-client-6.0.2 -m
```

### cec-client with Pipes

```
mkfifo cecpipe
cec-client < /home/morphs/DOCKER_IMAGES/HDMICEC/cecpipe > ceclog.txt 2>&1
echo "tx 1F:82:10:00" > cecpipe
```

### Get Device IDS

```
echo 'scan' | cec-client -s -d 1
```

### Set Source to HDMI 2

```
echo "txn 1F:82:20:00" | cec-client -s -d 4
```

### Power On

```
echo "on 0" | cec-client -s -d 1
```

### Power Off

```
echo "tx 1F:36" | cec-client -s -d 1
```

### Power Status

```
echo "pow 0" | cec-client -s -d 1
```

```
echo -e "tx 10\ntx 10:8F" | cec-client -s -d 31
```

### Volume Up / Down

```
echo "volup" | cec-client -s -d 1
```

### Mute Toggle

```
echo "tx 10:44:43" | cec-client -s -d 1
```
