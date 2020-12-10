# go-base2048

[![Go](https://github.com/Milly/go-base2048/workflows/Go/badge.svg)](https://github.com/Milly/go-base2048/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/Milly/go-base2048/branch/master/graph/badge.svg)](https://codecov.io/gh/Milly/go-base2048)
[![Go Report Card](https://goreportcard.com/badge/github.com/Milly/go-base2048)](https://goreportcard.com/report/github.com/Milly/go-base2048)

A base2048 encoding of binary data

## Usage

```go
input := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x0}
enc := base2048.DefaultEncoding.EncodeToString(input)
// out should be "นנநØ"
```

```go
input := "นנநØ"
enc := base2048.DefaultEncoding.DecodeString(input)
// out should be []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x0}
```

# Thanks

This is based on [rust-base2048](https://github.com/llfourn/rust-base2048).

# License

under BSD Zero-Clause License
