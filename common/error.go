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

type DuplicateBmkiError struct {
	Msg string
	Id  string
}

func NewDuplicateBmkiError(msg string, id string) *DuplicateBmkiError {
	return &DuplicateBmkiError{Msg: msg, Id: id}
}

func (e *DuplicateBmkiError) Error() string {
	return e.Msg
}
