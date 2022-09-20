#!/bin/bash

if [ -z ${WEB_HOOK_URL} ]; then
  return;
fi

text="Succeeded to push a new image to ${REPOSITORY_URI}:${IMAGE_TAG}."
if [ $1 -ne 0 ]; then
  text="Failure to build a new image of ${IMAGE_REPO_NAME}."
fi

message="{\"text\": \"${text}\"}"

curl -H "content-type:application/json" -d "$message" ${WEB_HOOK_URL}