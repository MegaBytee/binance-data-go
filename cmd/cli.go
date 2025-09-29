package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	binance_data "github.com/MegaBytee/binance-data-go"
	"github.com/MegaBytee/binance-data-go/config"
	"github.com/MegaBytee/binance-data-go/types"
)

func main() {

	cfg := config.Config{
		WithProxy: false,
		WithCache: false,
	}

	downloader := binance_data.NewDownloader(&cfg)
	if downloader == nil {
		panic("stop here")
	}

	from := flag.String("from", "s", "-from s or f (s:spot , f:futures) (required)")
	dt := flag.String("dt", "m", "-dt d or m (d:daily , m:monthly) (required)")
	symbol := flag.String("symbol", "BTCUSDT", "-symbol BTCUSDT (required)")
	tf := flag.String("tf", "1m", fmt.Sprintf("timeframes (choices: %s)", strings.Join(types.TimeFrames, ", ")))

	help := flag.Bool("help", false, "show help")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		//fmt.Fprintf(flag.CommandLine.Output(), "\timeframes choices:\n  %s\n", strings.Join(types.TimeFrames, ", "))
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	// Required flag check
	if *from == "" || *dt == "" || *symbol == "" || *tf == "" {
		fmt.Fprintln(os.Stderr, "Error:  is required")
		flag.Usage()
		os.Exit(2)
	}
	// Validate tf choice
	if !types.IsTimeFrameValidChoice(*tf) {
		fmt.Fprintf(os.Stderr, "Error: invalid --mode %q. Valid choices: %s\n", *tf, strings.Join(types.TimeFrames, ", "))
		flag.Usage()
		os.Exit(2)
	}
	params := types.NewDataParamsFromCli(*from, *dt, *symbol, *tf)

	fmt.Println(params)

	downloader.Run(params)
	downloader.Close()
}
