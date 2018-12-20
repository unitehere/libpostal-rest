#!/bin/sh

cd /
git clone https://github.com/openvenues/libpostal
git checkout tags/v1.0.0
cd /libpostal
./bootstrap.sh
./configure --datadir=$PWD --disable-data-download
make -j4
make install
