#!/bin/bash

set +ex

flyhelper --config-env FLY_HELPER_CONFG_ENV secrets pull

exec /opt/remote-signing-api/remote-signing-api ${CONFIG_FILE}