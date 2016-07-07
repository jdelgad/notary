//
// Copyright (c) Jacob Delgado. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
//

package main

import (
	"github.com/jdelgad/notary/endpoint"
)

func main() {
	e := endpoint.NewEndpoint("127.0.0.1", 9000, "/email")
	e.Setup()
	e.Run()
}
