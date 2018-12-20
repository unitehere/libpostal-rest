#!/bin/sh

cd /
git clone https://github.com/openvenues/libpostal
git checkout tags/v1.0.0
cd /libpostal
./bootstrap.sh
./configure --datadir=$PWD 
make -j4
make install
make clean
