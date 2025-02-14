# mgo
Markov chain text Generator for Go.

## Quick Start

Install using
```console
go get -v -u github.com/stumburs/mgo
```

## Examples

#### Generating text

```go
package main

import (
	mgo "github.com/stumburs/mgo"
	"fmt"
)

func main() {

	// Create a new generator
	generator := mgo.NewMarkovGenerator()

	// Read input text from a file
	generator.ReadSourceFromFile("input.txt")

	// Build Ngrams
	// Use either `mgo.SplitByNCharacters` or `mgo.SplitBySpaces`
	// Note: Argument N only matters when using `mgo.SplitByNCharacters`
	generator.BuildNgrams(mgo.SplitByNCharacters, 4)

	// Generate string
	text := generator.GenerateText(100)

	// Print the output
	fmt.Println(text)
}
```

#### Saving Ngrams to a file

```go
package main

import mgo "github.com/stumburs/mgo"

func main() {

	generator := mgo.NewMarkovGenerator()

	generator.ReadSourceFromFile("input.txt")

	generator.BuildNgrams(mgo.SplitByNCharacters, 4)

	// Write the generated Ngrams to a binary file
	generator.WriteNgrams("ngrams.bin")
}
```

#### Generating text from previously built Ngrams

```go
package main

import (
	mgo "github.com/stumburs/mgo"
	"fmt"
)

func main() {

	generator := mgo.NewMarkovGenerator()

	// Load previously generated Ngrams from file
	generator.ReadNgrams("ngrams.bin")

	text := generator.GenerateText(100)

	fmt.Println(text)
}
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
