package basic

import (
	"context"
	"fmt"
	"testing"
)

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	for n := range GenerateNumber(ctx) {
		fmt.Println(n)
		if n == 5 {
			cancel()
			break
		}
	}
}
