# Compound Go Monitor
Interact with [Compound](https://app.compound.finance/) DeFi protocol. 
The compound [API](https://compound.finance/docs/api#compound-api) supports pagination for querying all records by specifying a current page number. 
This program fetches all accounts utilizing golang's channel primitive for maximum concurrency.  There is some internal rate-limiting to prevent being throttled by Compound / DynamoDB servers and/or opening too many file descriptors locally.

## Install
``` go get github.com/davidiola/compound_go_monitor```

## Usage
```go run main.go```



