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

func getResponse(t *testing.T, method, endpoint string, expectedStatus int, contentHeader bool, data []byte) *bytes.Buffer {
	e := NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	if contentHeader {
		req.Header.Add("Content-Type", "application/json")
	}
	e.ServeHTTP(w, req)
	if w.Code != expectedStatus {
		t.Errorf("%s %s did not return expected status %d, but %d", method, endpoint, expectedStatus, w.Code)
		t.Errorf("%v", w.Body)
	}
	return w.Body
}

func parseBody(t *testing.T, body *bytes.Buffer, method, endpoint string, expectedStatus int, errorMsg string) {
	j := &response{}
	err := json.NewDecoder(body).Decode(j)
	if err != nil {
		t.Errorf("Could not unmarshall response from %s %s", method, endpoint)
	}
	if j.StatusCode != expectedStatus {
		t.Errorf("StatusCode in JSON was not execpted status %d, but %d", expectedStatus, j.StatusCode)
	}
	if j.Msg != errorMsg {
		t.Errorf(`Msg in JSON was not exepcted string "%s", but "%s"`, j.Msg, errorMsg)
	}
}

func TestEndpointGETRoot(t *testing.T) {
	getResponse(t, "GET", "/", http.StatusNotFound, false, nil)
}

func TestEndpointGETRoute(t *testing.T) {
	method := "GET"
	endpoint := "/email"
	body := getResponse(t, method, endpoint, http.StatusForbidden, false, nil)
	parseBody(t, body, method, endpoint, http.StatusForbidden, onlyPostMsg)
}

func TestEndpointPOSTNonApplicationJSON(t *testing.T) {
	method := "POST"
	endpoint := "/email"
	body := getResponse(t, method, endpoint, http.StatusBadRequest, false, nil)
	parseBody(t, body, method, endpoint, http.StatusBadRequest, contentTypeMsg)
}

func TestEndpointPOSTValidJSONRoot(t *testing.T) {
	method := "POST"
	endpoint := "/"
	data := []byte(`{"email": "jacob"}`)
	getResponse(t, method, endpoint, http.StatusNotFound, true, data)
}

func TestEndpointPOSTValidJSON(t *testing.T) {
	method := "POST"
	endpoint := "/email"
	data := []byte(`{"email": "jacob"}`)
	body := getResponse(t, method, endpoint, http.StatusOK, true, data)
	parseBody(t, body, method, endpoint, http.StatusOK, validEmailMsg)
}

func TestEndpointPOSTJSONEmptyEmail(t *testing.T) {
	method := "POST"
	endpoint := "/email"
	data := []byte(`{"email": ""}`)
	body := getResponse(t, method, endpoint, http.StatusBadRequest, true, data)
	parseBody(t, body, method, endpoint, http.StatusBadRequest, jsonMissingEmailMsg)
}

func TestEndpointPOSTJSONMissingEmail(t *testing.T) {
	method := "POST"
	endpoint := "/email"
	data := []byte(`{}`)
	body := getResponse(t, method, endpoint, http.StatusBadRequest, true, data)
	parseBody(t, body, method, endpoint, http.StatusBadRequest, jsonMissingEmailMsg)
}

func TestEndpointPOSTInvalidJSON(t *testing.T) {
	method := "POST"
	endpoint := "/email"
	data := []byte(`{"email": "jacob"`)
	body := getResponse(t, method, endpoint, http.StatusBadRequest, true, data)
	parseBody(t, body, method, endpoint, http.StatusBadRequest, invalidJSONMsg)
}
