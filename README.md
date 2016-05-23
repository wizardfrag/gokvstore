# `gokvstore`

## Running

1. `go get -u github.com/wizardfrag/gokvstore/...`
2. `kvstore`
3. Done!

## Usage

### HTTP

There is one HTTP endpoint, accessible via `GET` and `PUSH`. This is the following, where `:key` is replaced by the key you wish to store to/read from:

    /:key

#### POST

| Argument | Value                                                                  | Required |
|----------|------------------------------------------------------------------------|----------|
| key      | the key to store to, included in HTTP URI                              | true     |
| val      | the value to store                                                     | true     |
| type     | the type of the val, valid types are: `int`, `bool`, `float`, `string` | true     |

#### GET

All that is required for GET requests is the key, as part of the HTTP URI.

### TCP / UDP

TCP and UDP requests are similar, but slightly different. UDP requests must be made with one command per packet, otherwise an error will be returned.

#### PUT

Example:

```json
{"cmd":"put","item":{"key":"test","value":1}}
```

#### GET

Example:

```json
{"cmd":"get","item":{"key":"test"}}
```