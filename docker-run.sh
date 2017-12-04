#!/usr/bin/env bash

docker run --env-file ./.env --rm -it -p 8090:8090 tugayac/ct_api
