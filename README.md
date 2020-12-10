# go-base2048

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
