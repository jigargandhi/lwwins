# LW Wins  Last Writer Wins

This go repo demonstrates a set CvRDT where Last writer wins
## Prerequisites
- Protoc should be available on path
- Docker daemon should be running

## Build

```./build.sh```

## Run

Instance 1

``` docker run -p '3335:3334' -d lwwins:dev -node nodeA -token SecretT0k3n```

Instance 2

``` docker run -p '3336:3334' -d lwwins:dev -node nodeB -token SecretT0k3n```

Run command line

```./dist/lwwinscli -server_addr 0.0.0.0:3335```

## Todo
- ~~Use SSL Certificates for identity instead of a secret~~ SSL should not be used because of handshakes involved, symmetric key encryption like use AES-GCM
- When a new server is added server only time is synced value is not
- When a message is received timestamp should also be incremented (testing pending)
- To generate a new key use the following snippet in go playground 
```
package main

import (
	"fmt"
	"crypto/rand"
	"encoding/base64"
)

func main() {
	key := make([]byte, 32)
	rand.Reader.Read(key)

	fmt.Println(base64.StdEncoding.EncodeToString(key))
}
```