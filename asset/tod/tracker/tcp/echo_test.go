package tcp

import (
	"fmt"
	"testing"
)

func TestOnTraffic(t *testing.T) {
	var trafficTests = []struct {
		in  string // input
		out string // expected response
		ret int
	}{
		{"h1", "h1", 0},
	}
	for _, tt := range trafficTests {
		fmt.Print(tt)
	}
}
