#!/bin/bash
set -eu

duffle relocate -f /tweed/etc/config.d/provided/bundle.json \
       --relocation-mapping /dev/null \
       --repository-prefix=localhost:5000
