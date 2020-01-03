# LW Wins  Last Writer Wins

This go repo demonstrates a set CvRDT where Last writer wins
## Prerequisites
- Protoc should be available on path
- Docker daemon should be running

## Build

```./build.sh```

## Run

Instance 1

``` docker run -p '3335:3334' -d lwwins:dev -id 1 -token SecretT0k3n```

Instance 2

``` docker run -p '3336:3334' -d lwwins:dev -id 2 -token SecretT0k3n```

Run command line

```./dist/lwwinscli -server_addr 0.0.0.0:3335```

## Todo
- Use SSL Certificates for identity instead of a secret
- When a new server is added server only time is synced value is not
- When a message is received timestamp should also be incremented (testing pending)