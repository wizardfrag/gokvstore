package gokvstore

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type HttpServer struct {
	router *httprouter.Router
	store  *Store
}

func (h *HttpServer) Init(port int, s *Store) {
	h.router = httprouter.New()
	h.store = s

	// Purposefully only using POST/GET instead of PUT/POST/GET
	// Since I am assuming the consumer will not know about the existence of a key
	h.router.POST("/:key", h.post)
	h.router.GET("/:key", h.get)
	// TODO: Add delete if required

	fmt.Println("HTTP Server listening on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), h.router)
}

// post stores to the key/value store, and returns a 201 if the value is stored correctly.
func (h *HttpServer) post(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// This is required, otherwise r.PostFormValue will return "" - hard to track down!
	r.ParseForm()

	key := p.ByName("key")
	typ := r.PostFormValue("type")
	val := r.PostFormValue("val")

	if key == "" || typ == "" || val == "" {
		fmt.Printf("Given: %s / %s / %s\n", key, typ, val)
		jsonerr(w, weberror{"Please specify a key, type and val", http.StatusBadRequest})
		return
	}

	var item storageItem
	item.Key = key

	switch typ {
	case "int":
		v, err := strconv.Atoi(val)
		if err != nil {
			jsonerr(w, weberror{err.Error(), http.StatusBadRequest})
			return
		}
		item.Value = v
	case "string":
		item.Value = val
	case "bool":
		v, err := strconv.ParseBool(val)
		if err != nil {
			jsonerr(w, weberror{err.Error(), http.StatusBadRequest})
			return
		}
		item.Value = v
	case "float":
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			jsonerr(w, weberror{err.Error(), http.StatusBadRequest})
			return
		}
		item.Value = v
	default:
		jsonerr(w, weberror{"Unknown type: " + typ, http.StatusBadRequest})
		return
	}
	h.store.WriteItem(item)

	// TODO: Update this to check if key has only been updated
	w.WriteHeader(http.StatusCreated)
}

// get retrieves the value from the key/value store, if it exists, and returns it as a
// json-interpreted value
func (h *HttpServer) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get the key from the router params
	key := p.ByName("key")
	if key == "" {
		jsonerr(w, weberror{"Please specify a key", http.StatusBadRequest})
	}

	val, err := h.store.GetItem(storageItem{Key: key})
	if err != nil {
		// Only reason for an error is no such key, so reply with a 404.
		jsonerr(w, weberror{err.Error(), http.StatusNotFound})
	} else {
		jsonify(w, val)
	}
}

// The following is used for JSON replies via HTTP
type weberror struct {
	Message   string
	ErrorCode int
}

func jsonify(w http.ResponseWriter, obj interface{}) {
	jsonstr, err := json.Marshal(obj)
	if err != nil {
		jsonerr(w, weberror{err.Error(), http.StatusInternalServerError})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonstr))
}

func jsonerr(w http.ResponseWriter, weberr weberror) {
	jsonstr, err := json.Marshal(weberr)
	if err != nil {
		fmt.Fprintf(w, "unrecoverable error: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(weberr.ErrorCode)
	fmt.Fprint(w, string(jsonstr))
}
