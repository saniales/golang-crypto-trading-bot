# commander [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/commander)](https://goreportcard.com/report/github.com/shomali11/commander) [![GoDoc](https://godoc.org/github.com/shomali11/commander?status.svg)](https://godoc.org/github.com/shomali11/commander) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Command evaluator and parser

## Features

* Matches commands against provided text
* Extracts parameters from matching input
* Provides default values for missing parameters
* Supports String, Integer, Float and Boolean parameters

## Usage

Using `govendor` [github.com/kardianos/govendor](https://github.com/kardianos/govendor):

```
govendor fetch github.com/shomali11/commander
```

## Dependencies

* `proper` [github.com/shomali11/proper](https://github.com/shomali11/proper)


# Examples

## Example 1

In this example, we are matching a few strings against a command format, then parsing parameters if found or returning default values.

```go
package main

import (
	"fmt"
	"github.com/shomali11/commander"
)

func main() {
	properties, isMatch := commander.NewCommand("echo <word>").Match("echo hey")
	fmt.Println(isMatch)                             // true
	fmt.Println(properties.StringParam("word", ""))  // hey

	properties, isMatch = commander.NewCommand("repeat <word> <number>").Match("repeat hey 5")
	fmt.Println(isMatch)                              // true
	fmt.Println(properties.StringParam("word", ""))   // hey
	fmt.Println(properties.IntegerParam("number", 0)) // 5

	properties, isMatch = commander.NewCommand("repeat <word> <number>").Match("repeat hey")
	fmt.Println(isMatch)                              // true
	fmt.Println(properties.StringParam("word", ""))   // hey
	fmt.Println(properties.IntegerParam("number", 0)) // 0
}
```

## Example 2

In this example, we are tokenizing the command format and returning each token with a boolean that determines whether it is a parameter or not

```go
package main

import (
	"fmt"
	"github.com/shomali11/commander"
)

func main() {
	tokens := commander.NewCommand("echo <word>").Tokenize()
	for _, token := range tokens {
		fmt.Println(token)
	}
}
```

Output:
```
&{echo false}
&{word true}
```