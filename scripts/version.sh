#!/bin/bash

TAG=$1

# Tag with path-relative names that match module import path
git tag $TAG
git tag drivers/file/$TAG
git tag drivers/env/$TAG

# Push the correct tags
git push origin $TAG
git push origin drivers/file/$TAG
git push origin drivers/env/$TAG
