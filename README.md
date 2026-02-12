# gologger-nazalog

A [nazalog](https://github.com/q191201771/naza/tree/master/pkg/nazalog) adapter for [gologger](https://github.com/kordar/gologger).

## Installation

```bash
go get github.com/kordar/gologger_nazalog
```

## Usage

```go
package main

import (
	"github.com/kordar/gologger"
	"github.com/kordar/gologger_nazalog"
)

func main() {
	// Optional: Initialize nazalog with a filename (or empty for default ./logs/naza.log)
	gologger_nazalog.Init001("app.log")
	
	// Create the adapter
	adapter := gologger_nazalog.NewNazalogAdapt()
	
	// Initialize gologger with the adapter
	logger.InitGlobal(adapter)
	
	// Use gologger
	logger.Info("This message is logged via nazalog")
}
```
