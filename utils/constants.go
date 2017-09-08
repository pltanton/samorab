package utils

import (
	"regexp"
)

var SPLIT_REGEX *regexp.Regexp

func init() {
	SPLIT_REGEX, _ = regexp.Compile("[^0-9А-Яа-яёЁ_]")
}
