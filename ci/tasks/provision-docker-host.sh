#!/bin/bash

docker-machine create \
    -d amazonec2 \
    --amazonec2-access-key ${AWS_ACCESS_KEY_ID} \
    --amazonec2-secret-key ${AWS_SERCRET_ACCESS_KEY} \
    --amazonec2-region ${AWS_REGION} \
    ${machine_name}

machine-export ${machine_name}
aws s3 cp ${machine_name} s3://${AWS_BUCKET}