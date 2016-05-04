/* Copyright (c) 2016 Nick Potts (npotts)

The MIT License (MIT)

Copyright (c) 2016 Nick Potts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
Package yant contains  the basic fundamentals required to work with nullable database times
*/
package yant

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func parseTime(value string) (t time.Time, v bool) {
	layouts := []string{
		"2006-01-02 15:04:05.000",
		"2006-01-02 15:04:05.000Z",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05.000Z",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
	}
	for _, layout := range layouts {
		if tt, e := time.Parse(layout, value); e == nil {
			t = tt
			v = true
			return
		}
	}
	return
}

// NullTime represents a time.Time that may be NULL.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

var errConvert = fmt.Errorf("Can't convert to time.Time")

// Scan implements the Scanner interface.
// The value type must be time.Time or string / []byte (formatted time-string),
// otherwise Scan fails.
func (nt *NullTime) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return
	case []byte:
		nt.Time, nt.Valid = parseTime(string(v))
		return
	case string:
		nt.Time, nt.Valid = parseTime(v)
		return
	}

	nt.Valid = false
	return errConvert
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
