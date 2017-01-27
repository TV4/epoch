# epoch

The `epoch` package contains a convenience type `Time`, which is an alias of `time.Time`. It adds the `MarshalJSON()` and `UnmarshalJSON()` methods to enable simple conversion between `time.Time` objects and JSON timestamps in [Unix time](https://en.wikipedia.org/wiki/Unix_time).

## Marshalling

When marshalled, an `epoch.Time` object will be converted into a Unix timestamp representing milliseconds since epoch, i.e. January 1, 1970 UTC.

## Unmarshalling

`epoch.Time` will try to deduce timestamp resolution from the number of integer digits in the timestamp. Both integer and floating point timestamps are accepted.

I.e. `1485507392`, `1485507392.423425`, and `148550739200` will be interpreted as seconds since epoch. This puts an upper bound on which time in the future that can be represented.
