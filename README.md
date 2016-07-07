# Notary
---
An HTTP REST endpoint service reporting `200 OK` or `400 Bad Request` on
whether POST data sent as JSON (e.g. `{"email": "jacob@gmail.com"}`) is valid
or not, respectively.

Validation is against [RFC5322](https://tools.ietf.org/html/rfc3696)
([3.2.3](https://tools.ietf.org/html/rfc5322#section-3.2.3) and
 [3.4.1](https://tools.ietf.org/html/rfc5322#section-3.4.1)),
[RFC5321](https://tools.ietf.org/html/rfc5322#section-3.4.1),
[RFC3396](https://tools.ietf.org/html/rfc3696https://tools.ietf.org/html/rfc3696)
and its associated
[errata](http://www.rfc-editor.org/errata_search.php?rfc=3696).

By default `notary` listens on port 9000.
