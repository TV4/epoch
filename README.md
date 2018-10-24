# epoch

[![Build Status](https://travis-ci.org/TV4/epoch.svg?branch=master)](https://travis-ci.org/TV4/epoch)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/epoch)](https://goreportcard.com/report/github.com/TV4/epoch)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/TV4/epoch)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/epoch#license)

The `epoch` package contains a convenience type `Time`, which is an alias of `time.Time`. It adds the `MarshalJSON()` and `UnmarshalJSON()` methods to enable simple conversion between `time.Time` objects and JSON timestamps in [Unix time](https://en.wikipedia.org/wiki/Unix_time).

## Marshalling

When marshalled, an `epoch.Time` object will be converted into a Unix timestamp representing milliseconds since epoch, i.e. January 1, 1970 UTC.

## Unmarshalling

`epoch.Time` will try to deduce timestamp resolution from the number of integer digits in the timestamp. Both integer and floating point timestamps are accepted.

I.e. `1485507392`, `1485507392.423425`, and `148550739200` will be interpreted as seconds since epoch. This puts an upper bound on which time in the future that can be represented.

## License

Copyright (c) 2016-2018 Bonnier Broadcasting

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
