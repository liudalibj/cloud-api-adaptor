#!/usr/bin/env bash

if [ ! -d "./peerpodconfig-ctrl" ]; then
    git clone -b toplevel-cloud-module https://github.com/confidential-containers/peerpodconfig-ctrl.git
fi
