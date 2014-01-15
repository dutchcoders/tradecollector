/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2013 DutchCoders <http://github.com/dutchcoders/>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

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

package main

import(
    "encoding/json"
    "time"
    "strconv"
    _ "math"
)

type float64s struct {
        float64 float64 
}

func (t *float64s) UnmarshalJSON(b []byte) error {
        var value string

        if err := json.Unmarshal(b, &value); err != nil {
                return err
        }

        f, err := strconv.ParseFloat(value, 32)
        t.float64=f
        return err
}

type int64s struct {
        int64 int64 
}

func (t *int64s) UnmarshalJSON(b []byte) error {
        var value string

        if err := json.Unmarshal(b, &value); err != nil {
                return err
        }

        f, err := strconv.ParseInt(value, 10, 64)
        t.int64=f
        return err
}

type uint64s struct {
        uint64 uint64 
}

func (t *uint64s) UnmarshalJSON(b []byte) error {
        var value string

        if err := json.Unmarshal(b, &value); err != nil {
                return err
        }

        f, err := strconv.ParseInt(value, 10, 64)
        t.uint64=uint64(f)
        return err
}

type Timestamp struct {
        time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
        var unixtime int64

        if err := json.Unmarshal(b, &unixtime); err != nil {
                return err
        }

        t.Time = time.Unix(unixtime, 0)
        return nil
}

type Timestamps struct {
        Timestamp
}

func (t *Timestamps) UnmarshalJSON(b []byte) error {
        var value string

        if err := json.Unmarshal(b, &value); err != nil {
                return err
        }

        var unixtime int64
        
        unixtime, err := strconv.ParseInt(value, 10, 32);
        if err != nil {
                return err
        }


        t.Time = time.Unix(unixtime, 0)
        return nil
}

func MaxInt(a, b int64) int64 {
   if a >= b {
      return a
   }
   return b
}
