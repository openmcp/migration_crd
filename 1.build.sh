#!/bin/bash
docker_id="gentian0185"
image_name="openmcp-migration"
version="0.8.0"
export GO111MODULE=on

operator-sdk build $docker_id/$image_name:v$version
docker push $docker_id/$image_name:v$version
