#!/bin/bash

VERSION=${1}
DOCKERHUB_PASSWORD=${2}
DOCKERHUB_USERNAME=${3}

APP_NAME_NISTAGRAM_POSTS=nistagram_posts

APP_IMAGE_NAME_NISTAGRAM_POSTS=${DOCKERHUB_USERNAME}/${APP_NAME_NISTAGRAM_POSTS}:${VERSION}

APP_ARTIFACT_NAME_NISTAGRAM_POSTS=${APP_NAME_NISTAGRAM_POSTS}:${VERSION}.tar

DOCKER_BUILDKIT=1 docker build \
-t "${APP_IMAGE_NAME_NISTAGRAM_POSTS}" \
--no-cache \
.

docker login --username ${DOCKERHUB_USERNAME} --password=${DOCKERHUB_PASSWORD}
docker push "$APP_IMAGE_NAME_NISTAGRAM_POSTS"