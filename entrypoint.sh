#!/bin/bash
set -eu

echo Starting Tweed Setup

export TWEED_RO_CONFIG_FILE=/tweed-provided/tweed.yml
export TWEED_RO_CATALOG_FILE=/tweed-provided/catalog.yml

export TWEED_DATA_DIRECTORY=/tweed
export TWEED_PVC_DIRECTORY=/tweed/etc/config.d

export TWEED_CONFIG_FILE=${TWEED_DATA_DIRECTORY}/etc/config.d/config.json
export TWEED_ROOT=${TWEED_ROOT:-$TWEED_DATA_DIRECTORY}
export HOME=$TWEED_ROOT

echo TWEED_RO_CONFIG_FILE: ${TWEED_RO_CONFIG_FILE}
echo TWEED_CONFIG_FILE:    ${TWEED_CONFIG_FILE}
echo TWEED_DATA_DIRECTORY: ${TWEED_DATA_DIRECTORY}
echo TWEED_ROOT:           ${TWEED_ROOT}
echo HOME:                 ${HOME}
echo TWEED_CONFIG_MOUNT:   ${TWEED_CONFIG_MOUNT}

echo Copying files from configMap
ls -la /tweed-provided

mkdir -p ${TWEED_PVC_DIRECTORY}/provided/
cp -v ${TWEED_RO_CONFIG_FILE}  ${TWEED_DATA_DIRECTORY}/etc/config.d/provided/tweed.yml
cp -v ${TWEED_RO_CATALOG_FILE} ${TWEED_DATA_DIRECTORY}/etc/config.d/provided/catalog.yml

mkdir -p ${TWEED_DATA_DIRECTORY}/etc/config.d
if [[ ${INIT_VAULT:-} != "" ]]; then
	safe target the-vault $INIT_VAULT --no-strongbox
	TOKEN=$(safe init --keys 1 --json | jq -r '.root_token')
	cat >${TWEED_DATA_DIRECTORY}/etc/config.d/auto.vault.yml <<EOF
---
vault:
  url:    $INIT_VAULT
  token:  $TOKEN
  prefix: secret/boss
EOF
fi

if [[ -n ${TWEED_CONFIG:-} ]]; then
	cat >${TWEED_DATA_DIRECTORY}/etc/config.d/auto.yml <<EOF
$TWEED_CONFIG
EOF
	unset TWEED_CONFIG
fi

if [[ -n ${USE_THIS_KUBERNETES:-} ]]; then
	case $USE_THIS_KUBERNETES in
	y|Y|1|yes|YES)
		USE_THIS_KUBERNETES=/var/run/secrets/kubernetes.io/serviceaccount/token
		;;
	esac
	KUBERNETES_INFRASTRUCTURE_NAME=${KUBERNETES_INFRASTRUCTURE_NAME:-k8s}
	cat <<EOF | spruce merge --prune meta - >/tweed/.kubeconfig
meta:
  ca: (( file "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt" ))
apiVersion: v1
kind: Config
users:
  - name: tweed
    user:
      token: (( file "/var/run/secrets/kubernetes.io/serviceaccount/token" ))

clusters:
  - name: k8s
    cluster:
      certificate-authority-data: (( base64 meta.ca ))
      server: https://kubernetes.default

current-context: k8s
contexts:
  - name: k8s
    context:
      cluster: k8s
      user: tweed
EOF

	cat <<EOF | spruce merge >${TWEED_DATA_DIRECTORY}/etc/config.d/auto.k8s.yml
infrastructures:
  $KUBERNETES_INFRASTRUCTURE_NAME:
    type: kubernetes
    kubeconfig: (( file "/tweed/.kubeconfig" ))
EOF
fi

find ${TWEED_DATA_DIRECTORY}/etc/config.d -type f -name '*.yml' | sort | \
	xargs spruce merge --skip-eval | \
	spruce json > $TWEED_CONFIG_FILE

env | grep TWEED_
exec /usr/bin/tweed broker

