package main

import (
	"fmt"

	"github.com/kare/base62"
)

func main() {
	// Encode
	var urlVal int64 = 3781504209452600
	encodedURL := base62.Encode(urlVal)
	fmt.Println(encodedURL)

	// Decode
	byteURL, _ := base62.Decode(encodedURL)
	fmt.Println(byteURL)
}
