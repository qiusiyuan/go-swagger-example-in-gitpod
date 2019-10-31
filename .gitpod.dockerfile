FROM gitpod/workspace-full:latest

USER root
## Install go-swagger
RUN download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url') \
  && curl -o $GOROOT/bin/swagger -L'#' "$download_url" \
  && chmod +x $GOROOT/bin/swagger
