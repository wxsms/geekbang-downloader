package helpers

import (
	"github.com/flytam/filenamify"
	"log"
	"strings"
)

func ToFilename(str string) string {
	name, err := filenamify.FilenamifyV2(str)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(name)
}
