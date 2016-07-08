//
// Copyright (c) Jacob Delgado. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
//

package endpoint

import (
	"bytes"
	"encoding/json"
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
		t.Errorf("GET / did not return %v, but %v", http.StatusNotFound, w.Code)
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
	j := &response{}
	err := json.NewDecoder(w.Body).Decode(j)
	if err != nil {
		t.Error("Could not unmarshall response from GET /email")
	}
	if j.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode in JSON (%s) is not %s", j.StatusCode, http.StatusForbidden)
	}
	if j.Msg != only_post_msg {
		t.Errorf("Msg in JSON (%s) is not %s", j.Msg, only_post_msg)
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
	j := &response{}
	err := json.NewDecoder(w.Body).Decode(j)
	if err != nil {
		t.Error("Could not unmarshall response from POST /email")
	}
	if j.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode in JSON (%s) is not %s", j.StatusCode, http.StatusBadRequest)
	}
	if j.Msg != content_type_msg {
		t.Errorf("Msg in JSON (%s) is not %s", j.Msg, content_type_msg)
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
		t.Errorf("POST / expected status code %v, but received %v", http.StatusNotFound, w.Code)
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
	j := &response{}
	err := json.NewDecoder(w.Body).Decode(j)
	if err != nil {
		t.Error("Could not unmarshall response from POST /email")
	}
	if j.StatusCode != http.StatusOK {
		t.Errorf("StatusCode in JSON (%s) is not %s", j.StatusCode, http.StatusBadRequest)
	}
	if j.Msg != valid_email_msg {
		t.Errorf("Msg in JSON (%s) is not %s", j.Msg, valid_email_msg)
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
	j := &response{}
	err := json.NewDecoder(w.Body).Decode(j)
	if err != nil {
		t.Error("Could not unmarshall response from POST /email")
	}
	if j.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode in JSON (%s) is not %s", j.StatusCode, http.StatusBadRequest)
	}
	if j.Msg != json_missing_email_msg {
		t.Errorf("Msg in JSON (%s) is not %s", j.Msg, json_missing_email_msg)
	}
}

func TestEndpointPOSTJSONMissingEmail(t *testing.T) {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()

	b := []byte(`{}`)
	req, _ := http.NewRequest("POST", "/email", bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /email with content-type set to applciation/json did not return %v, but %v", http.StatusBadRequest, w.Code)
		t.Errorf("%v", w.Body)
	}
	j := &response{}
	err := json.NewDecoder(w.Body).Decode(j)
	if err != nil {
		t.Error("Could not unmarshall response from POST /email")
	}
	if j.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode in JSON (%s) is not %s", j.StatusCode, http.StatusBadRequest)
	}
	if j.Msg != json_missing_email_msg {
		t.Errorf("Msg in JSON (%s) is not %s", j.Msg, json_missing_email_msg)
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

	j := &response{}
	err := json.NewDecoder(w.Body).Decode(j)
	if err != nil {
		t.Error("Could not unmarshall response from POST /email")
	}
	if j.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode in JSON (%s) is not %s", j.StatusCode, http.StatusBadRequest)
	}
	if j.Msg != invalid_json_msg {
		t.Errorf("Msg in JSON (%s) is not %s", j.Msg, invalid_json_msg)
	}
}
