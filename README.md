# counter

[![][travis-svg]][travis-url]
[![][license-svg]][license-url]



`casbin-zk-watcher` is the Apache Zookeeper watcher for Casbin.

## Installation

`go get -u github.com/grepsr/casbin-zk-watcher`

## Example

```go
package main

import (
    "github.com/casbin/casbin"
    "github.com/casbin/casbin/util"
    "github.com/grepsr/casbin-zk-watcher"
)

func updateCallback(rev string) {
    util.LogPrint("New revision detected:", rev)
}

func main() {
    // Initialize the watcher.
    // hosts can be either a single URL or a comma-separated list of URLs to zookeeper hosts.
    // path is the path which will be watched. If not provided, it defaults to "/casbin".
    w := zkwatcher.NewWatcher("<hosts>", "<path>")
    
    // Initialize the enforcer.
    e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
    
    // Set the watcher for the enforcer.
    e.SetWatcher(w)
    
    // By default, the watcher's callback is automatically set to the
    // enforcer's LoadPolicy() in the SetWatcher() call.
    // We can change it by explicitly setting a callback.
    w.SetUpdateCallback(updateCallback)
    
    // Update the policy to test the effect.
    // You should see "[New revision detected: X]" in the log.
    e.SavePolicy()
}
```

## Requirements
- [casbin](https://github.com/casbin/casbin)
- [go-zookeeper](https://github.com/samuel/go-zookeeper)


[travis-url]: https://travis-ci.org/grepsr/casbin-zk-watcher
[travis-svg]: https://img.shields.io/travis/grepsr/casbin-zk-watcher.svg?branch=master

[license-url]: https://github.com/grepsr/casbin-zk-watcher/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg