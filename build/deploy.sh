#!/bin/bash

echo $GCLOUD_API_KEYFILE | base64 -d > /gcloud-api-key.json
gcloud auth activate-service-account "githubactions@go-xerox-upload.iam.gserviceaccount.com" --key-file=/gcloud-api-key.json

sed -i "s/CLIENT_ID/$CLIENT_ID/" env_variables.yaml
sed -i "s/PROJECT_ID/$PROJECT_ID/" env_variables.yaml
sed -i "s/CLIENT_SECRET/$CLIENT_SECRET/" env_variables.yaml
sed -i "s/ACCESS_TOKEN/$ACCESS_TOKEN/" env_variables.yaml
sed -i "s|REFRESH_TOKEN|$REFRESH_TOKEN|" env_variables.yaml
sed -i "s/EXPIRY_TIME/$EXPIRY_TIME/" env_variables.yaml

gcloud app deploy --quiet
