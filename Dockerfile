# Docker image for the Drone AWS CodeDeploy plugin
#
#     cd $GOPATH/src/github.com/drone-plugins/drone-aws-codedeploy
#     make deps build docker

FROM alpine:3.2

RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-aws-codedeploy /bin/
ENTRYPOINT ["/bin/drone-aws-codedeploy"]
