package common

import (
	"strings"

	"github.com/muesli/reflow/truncate"
)

const ELLIPSIS = ".."

func GetFirstLine(str string) (string, bool) {
	if idx := strings.Index(str, "\n"); idx != -1 {
		return str[:idx], true
	}
	return str, false
}

func TruncateWithEllipsis(str string, maxWidth int) (string, bool) {
	_str, isMultilines := GetFirstLine(str)
	truncatedStr := truncate.String(_str, uint(maxWidth-8))
	if isMultilines {
		return _str, true
	}
	return truncatedStr, len(str) > len(truncatedStr)
}
