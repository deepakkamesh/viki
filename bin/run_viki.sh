#!/bin/sh
pkill vikid
LOC=$(dirname "$0")

$LOC/vikid \
-config_file=$LOC/objects.conf \
-graphite_ipport=metrics.hyperlinkhome.com:2003 \
-resource=$LOC/../resources \
-mg_domain=sandboxf139420cc83d4d3a8c3cf5dfc9b06b42.mailgun.org \
-mg_apikey=key-6ceddfaf05c0d237076a19abe2afef5d \
-mg_pubkey=pubkey-ce009cba9207ec56ae09ac45b9607c2f \
-email_alert_list=deepak.kamesh@gmail.com,6024050044@tmomail.net \
-v=2 \
-alsologtostderr=true \
-logtostderr=false \
-log_dir=$LOC/../logs \
&
#-ssl \
