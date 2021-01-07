#!/bin/bash

service dbus start
bluetoothd &

/ruuvitag/main
