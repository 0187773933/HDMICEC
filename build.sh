#!/bin/bash
cd "/Users/morpheous/WORKSPACE/GO/HDMICEC/venv/lib/python3.11/site-packages/python-cec"
export CFLAGS="-I/usr/local/Cellar/libcec/6.0.2/include"
python uninstall cec --yes
python setup.py clean --all
python setup.py build
ls -R build
python setup.py install
