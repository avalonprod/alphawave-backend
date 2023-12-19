#!/usr/bin/env sh
set -eu

envsubst '${API_HOST} ${API_PORT} ${SERVER_NAME ${MATTERMOST_HOST} ${MATTERMOST_PORT}' < /nginx.conf.template > /etc/nginx/nginx.conf

exec "$@"
