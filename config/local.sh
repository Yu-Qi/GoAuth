#!/bin/sh

# general
export ENV=local
export APP_PORT=9030

# database
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USERNAME=root
export MYSQL_PASSWORD=changeit
export MYSQL_OPTIONS=charset=utf8mb4\&parseTime=True\&loc=UTC
export MYSQL_DATABASE=local
export MYSQL_SLOW_THRESHOLD=1000
export MYSQL_DSN="$MYSQL_USERNAME:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?$MYSQL_OPTIONS"

# redis
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
export REDIS_AUTH=

# jwt
export ACCESS_TOKEN_EXP_MINUTES=1440
export JWT_TOKEN_SECRET=changeit