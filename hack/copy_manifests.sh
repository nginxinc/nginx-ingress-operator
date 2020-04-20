#!/bin/bash

for f in $(find deploy/olm-catalog -name "*clusterserviceversion.yaml"); do
  cp $f bundle/
done