# clasp-docker

clasp-docker is a simple Dockerfile and helper script for running [clasp](https://github.com/google/clasp) in a odcker container.

## Setup

Build the docker image:

```shell
make build
```

## Usage

Run the clasp script in this repository the same way you'd run the real clasp.  It will transparently run it inside a container.

## Tips

Use `--no-localhost` when logging in so you can enter token code.

```shell
./clasp login --no-localhost
```
