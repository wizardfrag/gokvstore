package gokvstore

type kvItem interface{}

type kvStore map[string]kvItem

type storageItem struct {
	Key   string
	Value interface{}
}