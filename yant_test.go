package yant

/*Copyright (c) 2016 Nick Potts (npotts)

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

import (
	"fmt"
	"testing"
	"time"
)

func Test_ParseTime(t *testing.T) {
	tests := map[string]bool{
		"200-01-02 15:04:05.000":   false,
		"2006-01-02 15:04:05.000":  true,
		"2006-01-02 15:04:05.000Z": true,
		"2006-01-02T15:04:05.000":  true,
		"2006-01-02T15:04:05.000Z": true,
	}
	for format, k := range tests {
		t.Logf("Running checks on %q", format)
		if _, ok := parseTime(format); ok != k {
			t.Errorf("Format didnt parse correctly")
		}
	}
}

func TestNullTime_Scan(t *testing.T) {
	type x struct {
		i interface{}
		e error
	}

	tests := map[string]x{
		"nillizer":       x{i: nil, e: nil},
		"time":           x{i: time.Now(), e: nil},
		"bytes":          x{i: []byte("2006-01-02T15:04:05.000"), e: nil},
		"string":         x{i: "2006-01-02T15:04:05.000", e: nil},
		"something else": x{i: fmt.Errorf("some error"), e: errConvert},
	}

	nt := new(NullTime)

	if v, e := nt.Value(); v != nil || e != nil {
		t.Errorf("New should return a nil value")
	}
	nt.Time, nt.Valid = time.Now(), true
	if v, e := nt.Value(); v == nil || e != nil {
		t.Errorf("Altered nulltime should not return a nil value")
	}

	for name, x := range tests {
		t.Logf("Running checks on %q", name)
		if nt.Scan(x.i) != x.e {
			t.Errorf("Returned error not right")
		}
	}

	nt = new(NullTime)

}
