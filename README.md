# Docker Registry Cleaner (for Gitlab)

This tool utilizes [Gitlab Container Registry API](https://docs.gitlab.com/ee/api/container_registry.html) to facilitate removal of an *image* repository with all the images included inside it.

Note: this tool deliberately disallows deleting the root image repository.

## Rationale

This project is for use with [ci-templates](https://github.com/jobtome-labs/ci-templates).

Contributions are welcome.

## How to build

### With `go` (use version `1.16.6` or later).

```bash
go build ./cmd/cleaner
```

### With `docker`

```bash
docker build . -t cleaner
```

## How to use

With docker:

```bash
docker run --rm cleaner clean --api-v4-url=https://git.jobtome.io/api/v4 --project-id=574 --token=<your_token> --token-type=private --repository-name=chore-fix-artifact-path
```

On your host system:

```bash
./cleaner clean --api-v4-url=https://git.jobtome.io/api/v4 --project-id=574 --token=<your_token> --token-type=private --repository-name=chore-fix-artifact-path
```

In Gitlab CI Pipelines:

```yaml
# This stage is required to declare environment:on_stop, which will
# trigger the real clean up job. This is a workaround solution that
# allows us running a ci job after we merge the merge request.
# The environment created here is virtual and takes zero resources.
docker:image:cleanup:env:
  stage: deploy
  script: echo "env created"
  environment:
    name: docker/$CI_COMMIT_REF_NAME
    on_stop: docker:image:cleanup
  only:
    - merge_requests

docker:image:cleanup:
  stage: stop
  image: $CI_REGISTRY/auxiliary/docker-registry-cleaner:1.0.0
  variables:
    GITLAB_TOKEN_TYPE: job
    GITLAB_TOKEN: $CI_JOB_TOKEN
    GITLAB_API_V4_URL: $CI_API_V4_URL
    GITLAB_PROJECT_ID: $CI_PROJECT_ID
    GITLAB_REPOSITORY_NAME: $CI_COMMIT_REF_SLUG
    GIT_STRATEGY: none
  script:
    - cleaner clean
  when: manual
  environment:
    name: docker/$CI_COMMIT_REF_NAME
    action: stop
  only:
    - merge_requests
```
