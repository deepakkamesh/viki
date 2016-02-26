#!/bin/bash
# Turn on screen
ssh deepak@10.0.0.114  "export DISPLAY=:0.0 && /usr/bin/xset dpms force off"
