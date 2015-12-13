# drone-aws-codedeploy

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-aws-codedeploy/status.svg)](http://beta.drone.io/drone-plugins/drone-aws-codedeploy)
[![](https://badge.imagelayers.io/plugins/drone-aws-codedeploy:latest.svg)](https://imagelayers.io/?images=plugins/drone-aws-codedeploy:latest 'Get your own badge on imagelayers.io')

Drone plugin for deploying to AWS CodeDeploy

## Usage

```sh
./drone-aws-codedeploy <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "full_name": "drone/drone"
    },
    "build": {
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
    }
}
EOF
```

## Docker

Build the Docker container using `make`:

```
make deps build docker
```

### Example

```sh
docker run -i plugins/drone-aws-codedeploy <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "full_name": "drone/drone"
    },
    "build": {
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
    }
}
EOF
```
