# LW Wins  Last Writer Wins

This go repo demonstrates a set CvRDT where Last writer wins
## Build

```./build.sh```

## Run

Instance 1

``` docker run -p '3335:3334' -d lwwins:dev -id 1 -token SecretT0k3n```

Instance 2

``` docker run -p '3336:3334' -d lwwins:dev -id 2 -token SecretT0k3n```

Run

```./dist/lwwinscli -server_addr 0.0.0.0:3335```

