#!/bin/bash -e

GITHUB_TAG="${GITHUB_TAG:-local}"
SHA512_CMD="${SHA512_CMD:-sha512sum}"
export YAAM_DELIVERABLE="${YAAM_DELIVERABLE:-yaam}"

echo "GITHUB_TAG: '$GITHUB_TAG' YAAM_DELIVERABLE: '$YAAM_DELIVERABLE'"
cd cmd/yaam
go build -buildvcs=false -ldflags "-X main.Version=${GITHUB_TAG}" -o "${YAAM_DELIVERABLE}"
$SHA512_CMD "${YAAM_DELIVERABLE}" >"${YAAM_DELIVERABLE}.sha512.txt"
chmod +x "${YAAM_DELIVERABLE}"
cd ../..
