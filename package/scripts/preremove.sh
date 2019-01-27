#!/bin/sh
set -e
if [ "$1" = remove ]; then
    /bin/systemctl stop    proton-server
    /bin/systemctl disable proton-server
    /bin/systemctl daemon-reload
fi