# LW Wins  Last Writer Wins

This go repo demonstrates a set CvRDT where Last writer wins
## Build

```docker build -t lwwins:dev -f docker/Dockerfile ```

## Run

Instance 1

``` docker run -d lwwins:dev -id 1 -token SecretT0k3n```

Instance 2

``` docker run -d lwwins:dev -id 2 -token SecretT0k3n```
