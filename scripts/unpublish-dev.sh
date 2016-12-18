#!/usr/bin/env bash

set -e

. "${BASH_SOURCE%/*}/github_identifiers.sh"

release_url="https://api.github.com/repos/dirkraft/getprofile/releases/${DEV_RELEASE_ID}"
delete_url='https://api.github.com/repos/dirkraft/getprofile/releases/assets'
asset_ids=$(curl --silent "${release_url}" | jq '.assets[].id')

for asset_id in ${asset_ids} ; do

  curl -XDELETE --silent \
    --header "Authorization: token ${GITHUB_TOKEN}" \
    "${delete_url}/${asset_id}"

done
