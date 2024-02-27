#!/usr/bin/env bash

if [ ! -d "./peerpod-ctrl" ]; then
    git clone -b toplevel-cloud-module https://github.com/confidential-containers/peerpod-ctrl.git
fi
