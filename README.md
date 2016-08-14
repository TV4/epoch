# epoch

The `epoch` package contains a convenience type `Time`, which is an alias of `time.Time`. It adds the `MarshalJSON()` and `UnmarshalJSON()` methods to enable simple conversion between `time.Time` objects and JSON timestamps in *milliseconds from epoch*.

