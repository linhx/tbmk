/*
 * This file is part of tbmk.
 *
 * tbmk is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tbmk is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tbmk.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Copyright (C) 2024 Nguyen Dinh Linh
 */

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
