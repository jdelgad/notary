//
// Copyright (c) Jacob Delgado. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
//

package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Email struct {
	Email string `json:"email"`
}

type Endpoint struct {
	Address string
	Port    uint
	Mux     *http.ServeMux
	route   string
}

func NewEndpoint(addr string, port uint, route string) Endpoint {
	return Endpoint{Address: addr, Port: port, route: route, Mux: http.NewServeMux()}
}

func (e *Endpoint) Setup() {
	e.Mux.HandleFunc(e.route, verify)
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.Mux.ServeHTTP(w, r)
}

func (e *Endpoint) Run() {
	http.ListenAndServe(fmt.Sprintf("%s:%d", e.Address, e.Port), e.Mux)
}

func verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "%s", "Only POST method is supported.")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", "Content-Type header must be application/json.")
		return
	}

	email, err := obtainEmail(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	} else {
		fmt.Fprintf(w, "%s address is valid", email.Email)
	}
}

func obtainEmail(r *http.Request) (Email, error) {
	decoder := json.NewDecoder(r.Body)
	var email Email
	err := decoder.Decode(&email)
	if err != nil {
		return Email{}, errors.New("Invalid JSON request.")
	}
	if email.Email == "" {
		return Email{}, errors.New("Invalid JSON request. Missing key email.")
	}
	return email, nil
}
