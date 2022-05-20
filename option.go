package main

import (
	"io"
	"log"
	"os"
	"strings"
)

var (
	logger *log.Logger
	opt    *option
)

type option struct {
	replaceKeyword bool
	outdir         string
}

func initLogger(enable bool) {
	logger = log.Default()
	if enable {
		logger.SetOutput(os.Stderr)
	} else {
		logger.SetOutput(io.Discard)
	}
}

func applyOptions(optString string) {
	opt = &option{}
	for _, o := range strings.Split(optString, ",") {
		tmp := strings.SplitN(o, "=", 2)
		if len(tmp) == 2 {
			k, v := strings.TrimSpace(tmp[0]), strings.TrimSpace(tmp[1])
			switch k {
			case "outdir":
				opt.outdir = v
			case "replace_keyword":
				if v == "true" || v == "1" {
					opt.replaceKeyword = true
				}
			default:
			}
		}
	}
}
