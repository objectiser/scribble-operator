#!/bin/bash

git diff -s --exit-code
if [[ $? != 0 ]]; then
    echo "The repository isn't clean. We won't proceed, as we don't know if we should commit those changes or not."
    exit 1
fi

BASE_BUILD_IMAGE=${BASE_BUILD_IMAGE:-"objectiser/scribble-operator"}
OPERATOR_VERSION=${OPERATOR_VERSION:-$(git describe --tags)}
OPERATOR_VERSION=$(echo ${OPERATOR_VERSION} | grep -Po "([\d\.]+)")
TAG=${TAG:-"v${OPERATOR_VERSION}"}
BUILD_IMAGE=${BUILD_IMAGE:-"${BASE_BUILD_IMAGE}:${OPERATOR_VERSION}"}

sed "s~image: objectiser\/scribble-operator\:.*~image: ${BUILD_IMAGE}~gi" -i deploy/operator.yaml
sed "s~image: objectiser\/scribble-operator\:.*~image: ${BUILD_IMAGE}~gi" -i deploy/operator-openshift.yaml

git diff -s --exit-code
if [[ $? == 0 ]]; then
    echo "No changes detected. Skipping."
else
    git add deploy/operator.yaml deploy/operator-openshift.yaml
    git commit -qm "Release ${TAG}" --author="Scribble Release <jaeger-release@objectiser.io>"
    git tag ${TAG}
    git push --repo=https://${GH_WRITE_TOKEN}@github.com/objectiser/scribble-operator.git --tags
    git push https://${GH_WRITE_TOKEN}@github.com/objectiser/scribble-operator.git refs/tags/${TAG}:master
fi
