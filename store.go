package gokvstore

var storage store
var storagechan = make(chan itemtostore, 0)

func init() {
    storage = make(store)
    go run()
}

func run() {
    for {
        select {
        case item := storagechan:
            
        }
    }
}