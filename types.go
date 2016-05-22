package gokvstore

import "sync"

type kvitem struct {
    Value string
    Mutex *sync.Mutex
}

type kvstore map[string]kvitem

type itemtostore struct {
    Key string
    Value kvitem
}