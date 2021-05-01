# koyeb-touch

[![Go](https://github.com/ReeceM/koyeb-touch/actions/workflows/go.yml/badge.svg)](https://github.com/ReeceM/koyeb-touch/actions/workflows/go.yml)

## Usage

To use this in a ci file, do the following

```bash
VERSION='v1.0.1' # Get the latest tagged version
tar xvf <(curl -sL https://github.com/ReeceM/koyeb-touch/releases/download/$VERSION/koyeb-touch-$VERSION-linux-amd64.tar.gz)
chmod +x koyeb-touch
./koyeb-touch $KOYEB_API_TOKEN $KOYEB_APP_NAME $KOYEB_SERVICE_NAME
```

For GitLab it can be something like so:

You need to add the following to your CI/CD pipeline settings as variables:
|`KOYEB_API_TOKEN` | `KOYEB_APP_NAME` | `KOYEB_SERVICE_NAME`|
| ---------------- | ---------------- | ------------------- |
| This is the API token created in Koyeb | The name of the app, you can get it in the URL or what you called it | The name of the service, usually `main` |
```yml
variables:
  VERSION: 'v1.0.1'

koyeb:
  stage: deploy
  allow_failure: true
  script:
    - apk add --no-cache libc6-compat
    - wget -c "https://github.com/ReeceM/koyeb-touch/releases/download/$VERSION/koyeb-touch-$VERSION-linux-amd64.tar.gz" -O - | tar -xz -C .
    - chmod +x koyeb-touch
    - ./koyeb-touch $KOYEB_API_TOKEN $KOYEB_APP_NAME $KOYEB_SERVICE_NAME
```

This will then 'touch' the latest Koyeb service in your app and trigger a new deploy automatically.
