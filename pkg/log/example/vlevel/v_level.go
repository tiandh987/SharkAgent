package main

import "github.com/tiandh987/SharkAgent/pkg/log"

func main() {
	defer log.Flush()

	log.V(4).Info("This is a V level message")
	log.V(5).Infow("This is a V level message with fields", "X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	// 2022-04-07 16:15:55.366	WARN	vlevel/v_level.go:8	This is a V level message
	// 2022-04-07 16:15:55.377	INFO	vlevel/v_level.go:9	This is a V level message with fields	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}
}