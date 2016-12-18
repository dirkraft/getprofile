#!/usr/bin/env bash

DEV_RELEASE_ID=4952340

if [ "$GITHUB_TOKEN" == "" ] ; then
    echo "export GITHUB_TOKEN please"
    exit 1
fi
