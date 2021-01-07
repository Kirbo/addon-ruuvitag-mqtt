#!/bin/bash

bluetoothctl list

service dbus start
bluetoothd &

/ruuvitag/main
