package gokvstore

type kvItem interface{}

type kvStore map[string]kvItem

type storageItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type storageCommand struct {
	Cmd  string      `json:"cmd"`
	Item storageItem `json:"item"`
}

// The following is used for outputting JSON errors
type kvError struct {
	Message   string `json:"message"`
	ErrorCode int `json:"error_code"`
}

type kvResponse struct {
	Success bool `json:"success"`
}