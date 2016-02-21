#!/bin/sh

./vikid \
-ssl \
-log=./logs \
-graphite_ipport=metrics.hyperlinkhome.com:2003 \
-resource=./resources \
-log_stdout=false \
-log_file=./logs/viki.log \
&
