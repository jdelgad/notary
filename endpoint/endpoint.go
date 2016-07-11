//
// Copyright (c) Jacob Delgado. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
//

package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	//"net/mail"
	"net/mail"
)

type response struct {
	StatusCode int
	Msg        string
}

type email struct {
	Email string `json:"email"`
}

// Endpoint is a configuration object representing the IP address, port and
// route for this service's REST endpoint.
type Endpoint struct {
	addr  string
	port  uint
	mux   *http.ServeMux
	route string
}

const (
	onlyPostMsg         = "Only POST method is supported."
	contentTypeMsg      = "Content-Type header must be application/json."
	invalidEmailMsg     = "Email does not conform to RFC 5322."
	validEmailMsg       = "Email conforms to RFC 5322."
	invalidJSONMsg      = "Invalid JSON request."
	jsonMissingEmailMsg = "Invalid JSON request. Missing key email."
)

// NewEndpoint returns a new Endpoint object to be used for serving an HTTP
// REST endpoint with the user specified parameters.
func NewEndpoint(addr string, port uint, route string) Endpoint {
	return Endpoint{addr: addr, port: port, route: route,
		mux: http.NewServeMux()}
}

// Setup the HTTP REST endpoint given the user specified options.
func (e *Endpoint) Setup() {
	e.mux.HandleFunc(e.route, verify)
}

// ServeHTTP serves a given request routing it to the appropriate router given
// the requested endpoint.
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mux.ServeHTTP(w, r)
}

// Run the HTTP REST endpoint given the endpoint configuration object.
func (e *Endpoint) Run() {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", e.addr, e.port), e.mux)
	if err != nil {
		log.Fatalf("Web server was not able to run on address %s:%d",
			e.addr, e.port)
	}
}

func respond(w http.ResponseWriter, statusCode int, msg string) {
	resp := &response{}
	w.WriteHeader(statusCode)
	resp.StatusCode = statusCode
	resp.Msg = msg
	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("Could not serialize response")
	} else {
		fmt.Fprintf(w, "%s", b)
	}
}

func verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond(w, http.StatusForbidden, onlyPostMsg)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		respond(w, http.StatusBadRequest, contentTypeMsg)
		return
	}

	em, err := obtainEmail(r)
	if err != nil {
		respond(w, http.StatusBadRequest, err.Error())
	} else {
		_, err = mail.ParseAddress(em.Email)
		if err != nil {
			respond(w, http.StatusBadRequest, invalidEmailMsg)
		} else {
			respond(w, http.StatusOK, validEmailMsg)
		}
	}
}

func obtainEmail(r *http.Request) (email, error) {
	decoder := json.NewDecoder(r.Body)
	var e email
	err := decoder.Decode(&e)
	if err != nil {
		return email{}, errors.New(invalidJSONMsg)
	}
	if e.Email == "" {
		return email{}, errors.New(jsonMissingEmailMsg)
	}
	return e, nil
}
