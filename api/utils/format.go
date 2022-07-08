package utils

import (
	"strings"

	"github.com/pkg/errors"
)

func FormatError(err string) error {
	if strings.Contains(err, "name") {
		return errors.New("Name Already Taken")
	}
	if strings.Contains(err, "website") {
		return errors.New("Website Already Taken")
	}
	return errors.New("Incorrect Details")
}
