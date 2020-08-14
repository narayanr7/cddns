package cddns

import (
	"log"
)

var DebugEnabled bool

func DebugPrint(msg string) {
	if DebugEnabled {
		log.Printf(msg)
	}
}
