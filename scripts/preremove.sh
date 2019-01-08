#!/bin/sh
set -e

/bin/systemctl disable proton-server
/bin/systemctl daemon-reload
