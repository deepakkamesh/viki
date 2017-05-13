#!/bin/sh
pkill vikid
LOC=$(dirname "$0")

$LOC/vikid \
-config_file=$LOC/objects.conf \
-graphite_ipport=metrics.hyperlinkhome.com:2003 \
-resource=$LOC/../resources \
-v=2 \
-stderrthreshold=info \
-alsologtostderr=false \
-logtostderr=false \
-log_dir=$LOC/../logs \
&
#-ssl \
