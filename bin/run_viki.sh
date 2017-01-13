#!/bin/sh
LOC=$(dirname "$0")

$LOC/vikid \
-config_file=$LOC/objects.conf \
-log=$LOC/../logs \
-graphite_ipport=metrics.hyperlinkhome.com:2003 \
-resource=$LOC/../resources \
-log_stdout=false \
-log_file=$LOC/../logs/viki.log \
&
#-ssl \
