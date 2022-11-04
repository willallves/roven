#!/bin/bash

LOGLINE="Agent attestation request completed"
for ((i=0;i<1000;i++)); do
    if ! kubectl -nspire rollout status statefulset/spire-server; then
        sleep 1
        continue
    fi
    if ! kubectl -nspire rollout status daemonset/spire-agent; then
        sleep 1
        continue
    fi

    kubectl logs spire-server-0 -c spire-server
    kubectl logs -n spire -l app=spire-agent --all-containers=true

    if ! kubectl -nspire logs statefulset/spire-server -c spire-server | grep -e "$LOGLINE" ; then
        sleep 1
        continue
    fi
    log-debug "Node attested."
    GOOD=1
    exit 0
done