package utils

import (
	"github.com/satori/go.uuid"
	"strings"
)

func UUID() string {
	return strings.Replace(strings.ToUpper(uuid.NewV4().String()), "-", "", -1)
}
