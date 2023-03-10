# go-pingen-sdk
![Build Status](https://github.com/goneup/go-pingen-sdk/workflows/release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/goneup/go-pingen-sdk)](https://goreportcard.com/report/github.com/goneup/go-pingen-sdk)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/goneup/go-pingen-sdk?style=plastic)


A personal hobby project that implements the API of postal provider [pingen.com](https://www.pingen.com) in Go.

**Not affiliated in any way with pingen.com**

## Usage
> **Warning**
> This API sends real letters and costs money. After you send a letter it is validated and printed.   There is only a very short timeframe before printing where the sending process can be canceled.

Import it with
```
    go get github.com/goneup/go-pingen-sdk
```

Simple usage with. A full example can be found at [full_example.go](examples/full_example.go).
```go
    clientID := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    org := os.Getenv("PINGEN_ORG")
    useProd := true

    c := pingen.NewClient(clientID, clientSecret, useProd, org, context.Background())
    c.ListLetters()
	
```

For a complete overview of methods, refer to the [go package documentation](https://pkg.go.dev/github.com/goneup/go-pingen-sdk) 

## API Information

[Documentation](https://api.v2.pingen.com/documentation)

You need a API key, type Client Credentials. Instructions can be found [here](https://api.v2.pingen.com/documentation#section/Authentication/How-to-obtain-a-Client-ID). The organisation UUID can be found in any URL on the website.

Pingen has a seperate staging environment if you want to test the complete flow, see [here](https://api.v2.pingen.com/documentation#section/Basics)

