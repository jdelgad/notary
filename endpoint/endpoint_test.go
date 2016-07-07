//
// Copyright (c) Jacob Delgado. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
//

package endpoint

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpointGETRoot(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("GET / did not return %v, but %v", http.StatusBadRequest, w.Code)
		t.Errorf("%v", w.Body)
	}
}

func TestEndpointGETRoute(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	req, _ := http.NewRequest("GET", "/email", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("GET / did not return %v, but %v", http.StatusForbidden, w.Code)
		t.Errorf("%v", w.Body)
	}
}

func TestEndpointPOSTNonApplicationJSON(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	req, _ := http.NewRequest("POST", "/email", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /email with content-type set to applciation/x-www-form-urlencoded did not return %v, but %v", http.StatusBadRequest, w.Code)
		t.Errorf("%v", w.Body)
	}
}

func TestEndpointPOSTValidJSONRoot(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	b := []byte(`{"email": "jacob"}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("POST /email with content-type set to applciation/json did not return %v, but %v", http.StatusNotFound, w.Code)
		t.Errorf("%v", w.Body)
	}
}

func TestEndpointPOSTValidJSON(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	b := []byte(`{"email": "jacob"}`)
	req, _ := http.NewRequest("POST", "/email", bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("POST /email with content-type set to applciation/json did not return %v, but %v", http.StatusOK, w.Code)
		t.Errorf("%v", w.Body)
	}
}

func TestEndpointPOSTJSONEmptyEmail(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	b := []byte(`{"email": ""}`)
	req, _ := http.NewRequest("POST", "/email", bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /email with content-type set to applciation/json did not return %v, but %v", http.StatusBadRequest, w.Code)
		t.Errorf("%v", w.Body)
	}
}

func TestEndpointPOSTInvalidJSON(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	b := []byte(`{"email": "jacob"`)
	req, _ := http.NewRequest("POST", "/email", bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /email with content-type set to applciation/json did not return %v, but %v", http.StatusBadRequest, w.Code)
		t.Errorf("%v", w.Body)
	}
}
