# go-pingen-sdk

A personal hobby project that implements the API of postal provider [pingen.com](pingen.com).

**Not affiliated in any way with pingen.com**

## Usage

Import it with
```
    go get github.com/GoneUp/go-pingen-sdk
```

Simple usage with (or see [main.go](main.go))
```
    clientID := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    org := os.Getenv("PINGEN_ORG")

    c := pingen.NewClient(clientID, clientSecret, true, org, context.Background())
    c.ListLetters()
	
```
