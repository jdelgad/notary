# Notary
---
[![Build Status](https://travis-ci.org/jdelgad/notary.svg?branch=master)](https://travis-ci.org/jdelgad/notary)
[![codecov](https://codecov.io/gh/jdelgad/notary/branch/master/graph/badge.svg)](https://codecov.io/gh/jdelgad/notary)
[![GoDoc](https://godoc.org/github.com/jdelgad/notary?status.svg)](https://godoc.org/github.com/jdelgad/notary)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/jdelgad/notary/master/LICENSE)

An HTTP REST endpoint service reporting whether a POSTed JSON object containing
an email is valid.

Validation is against [RFC5322](https://tools.ietf.org/html/rfc5322)
(§[3.2.3](https://tools.ietf.org/html/rfc5322#section-3.2.3) and
 §[3.4.1](https://tools.ietf.org/html/rfc5322#section-3.4.1)).
* Notable differences:
  * Obsolete address formats are not parsed, including addresses with embedded route information.
  * Group addresses are not parsed.
  * The full range of spacing (the CFWS syntax element) is not supported, such as breaking addresses across lines.

### Input:
* Only POST method is supported to route `/email`
  * Any other routes will return a `404 page not found` error.
* Header `Content-Type` must be `application/json`.
  * All other types will return a `400 Bad Request` error.
* JSON Object: `{"email": "jacob@gmail.com"}`.
  * May contain more fields, but only email is parsed.
  * Email must be a string value.

### Output:
* JSON response in Body: `{"StatusCode": int, Msg: "string"}`
  * Response Codes:
    1. `404 page not found` when the HTTP route is not _/email_.
    1. `403 Forbidden` if the HTTP request method is not _POST_.
    1. `400 Bad Request` if the `Content-Type` header is not _application/json_, JSON POSTed is invalid, JSON POSTed is missing an email key, or if the email address is invalid.
    1. `200 OK` only when all conditions have been met (POST /email with JSON: _{ "email": "validRFC5322Email" }_).
* HTTP Header `Status-Line` contains the same
[HTTP Status Code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes) as the JSON returned.


### Defaults
* `notary` listens on port 9000.
