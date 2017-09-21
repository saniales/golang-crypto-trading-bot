# proper [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/proper)](https://goreportcard.com/report/github.com/shomali11/proper) [![GoDoc](https://godoc.org/github.com/shomali11/proper?status.svg)](https://godoc.org/github.com/shomali11/proper) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A `map[string]string` decorator offering a collection of helpful functions to extract the values in different types

## Features

* Retrieve data from a string map, in String, Integer, Float and Boolean types.
* Return a default value in case of missing keys or invalid types

## Usage

Using `govendor` [github.com/kardianos/govendor](https://github.com/kardianos/govendor):

```
govendor fetch github.com/shomali11/proper
```

# Examples

```go
package main

import (
	"fmt"
	"github.com/shomali11/proper"
)

func main() {
	parameters := make(map[string]string)
	parameters["boolean"] = "true"
	parameters["float"] = "1.2"
	parameters["integer"] = "11"
	parameters["string"] = "value"

	properties := proper.NewProperties(parameters)

	fmt.Println(properties.BooleanParam("boolean", false)) // true
	fmt.Println(properties.FloatParam("float", 0))         // 1.2
	fmt.Println(properties.IntegerParam("integer", 0))     // 11
	fmt.Println(properties.StringParam("string", ""))      // value
}
```