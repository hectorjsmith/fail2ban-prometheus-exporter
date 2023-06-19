#!/bin/sh

VERSION=`echo $CI_COMMIT_TAG | cut -c 2-`
MAJOR=`echo $VERSION | cut -d "." -f 1`
MINOR=`echo $VERSION | cut -d "." -f 2`

echo "version: $VERSION (major: $MAJOR; minor: $MINOR)"

docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
docker pull $IMAGE_NAME:$VERSION
docker tag $IMAGE_NAME:$VERSION $IMAGE_NAME:latest
docker tag $IMAGE_NAME:$VERSION $IMAGE_NAME:$MAJOR
docker tag $IMAGE_NAME:$VERSION $IMAGE_NAME:$MAJOR.$MINOR
docker push $IMAGE_NAME
