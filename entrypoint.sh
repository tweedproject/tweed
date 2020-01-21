#!/bin/bash
set -eu

export TWEED_CONFIG_FILE=/tweed/etc/config.json
export TWEED_ROOT=${TWEED_ROOT:-/tweed}
export HOME=$TWEED_ROOT

mkdir -p /tweed/etc/config.d
if [[ ${INIT_VAULT:-} != "" ]]; then
	safe target the-vault $INIT_VAULT --no-strongbox
	TOKEN=$(safe init --keys 1 --json | jq -r '.root_token')
	cat >/tweed/etc/config.d/auto.vault.yml <<EOF
---
vault:
  url:    $INIT_VAULT
  token:  $TOKEN
  prefix: secret/boss
EOF
fi

if [[ -n ${TWEED_CONFIG:-} ]]; then
	cat >/tweed/etc/config.d/auto.yml <<EOF
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

	cat <<EOF | spruce merge >/tweed/etc/config.d/auto.k8s.yml
infrastructures:
  $KUBERNETES_INFRASTRUCTURE_NAME:
    type: kubernetes
    kubeconfig: (( file "/tweed/.kubeconfig" ))
EOF
fi

duffle relocate -f /tweed/etc/config.d/provided/bundle.json \
       --relocation-mapping /dev/null \
       --repository-prefix=localhost:5000

find /tweed/etc/config.d -type f -name '*.yml' | sort | \
	xargs spruce merge --skip-eval | \
	spruce json > $TWEED_CONFIG_FILE

env | grep TWEED_
exec /usr/bin/tweed broker "$@"
