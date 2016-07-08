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
)

type response struct {
	StatusCode int
	Msg        string
}

type Email struct {
	Email string `json:"email"`
}

type Endpoint struct {
	Address string
	Port    uint
	Mux     *http.ServeMux
	route   string
}

const (
	only_post_msg          = "Only POST method is supported."
	content_type_msg       = "Content-Type header must be application/json."
	valid_email_msg        = "Email is valid."
	invalid_json_msg       = "Invalid JSON request."
	json_missing_email_msg = "Invalid JSON request. Missing key email."
)

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
		respond(w, http.StatusForbidden, only_post_msg)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		respond(w, http.StatusBadRequest, content_type_msg)
		return
	}

	_, err := obtainEmail(r)
	if err != nil {
		respond(w, http.StatusBadRequest, err.Error())
	} else {
		respond(w, http.StatusOK, valid_email_msg)
	}
}

func obtainEmail(r *http.Request) (Email, error) {
	decoder := json.NewDecoder(r.Body)
	var email Email
	err := decoder.Decode(&email)
	if err != nil {
		return Email{}, errors.New(invalid_json_msg)
	}
	if email.Email == "" {
		return Email{}, errors.New(json_missing_email_msg)
	}
	return email, nil
}
