# go-pingen-sdk
![Build Status](https://github.com/goneup/go-pingen-sdk/workflows/release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/goneup/go-pingen-sdk)](https://goreportcard.com/report/github.com/goneup/go-pingen-sdk
![GitHub release (latest by date)](https://img.shields.io/github/v/release/goneup/go-pingen-sdk?style=plastic)


A personal hobby project that implements the API of postal provider [pingen.com](https://www.pingen.com).

**Not affiliated in any way with pingen.com**

## Usage

Import it with
```
    go get github.com/GoneUp/go-pingen-sdk
```

Simple usage with (or see [main.go](main.go))
```go
    clientID := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    org := os.Getenv("PINGEN_ORG")
    useProd := true

    c := pingen.NewClient(clientID, clientSecret, useProd, org, context.Background())
    c.ListLetters()
	
```


## API Information

[Documentation](https://api.v2.pingen.com/documentation)

You need a API key, type Client Credentials. Instructions can be found [here](https://api.v2.pingen.com/documentation#section/Authentication/How-to-obtain-a-Client-ID). The organisation UUID can be found in any URL on the website.

Pingen has a seperate staging environment if you want to test the complete flow, see [here](https://api.v2.pingen.com/documentation#section/Basics)

