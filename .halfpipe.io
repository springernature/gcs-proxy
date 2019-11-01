team: engineering-enablement
pipeline: gcs-proxy

feature_toggles:
- update-pipeline

tasks:
- type: run
  name: Test and build
  script: \go test ./... ; go build -o cf/server cmd/main.go ; cp -r static cf
  docker:
    image: golang:1.11-stretch
  save_artifacts:
  - cf

- type: deploy-cf
  api: ((cloudfoundry.api-snpaas))
  space: halfpipe
  manifest: cf/manifest.yml
  vars:
    GCS_KEY: ((gcs-proxy.gcs-key))
    BUCKET: ((gcs-proxy.bucket))
  deploy_artifact: cf

