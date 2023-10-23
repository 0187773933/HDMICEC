# HDMI-CEC Controller

- https://github.com/Pulse-Eight/libcec
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