/*
ccdns
Author:  robert@x1sec.com
License: see LICENCE
*/

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
