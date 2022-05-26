package main

import (
	"os"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

var (
	debugEnvName = "GORM_TAG_DEBUG"
	enableLogger bool
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	enableLogger = os.Getenv(debugEnvName) != ""
	opt := pgs.DebugEnv(debugEnvName)
	mod := newMod()
	pgs.Init(opt).
		RegisterModule(mod).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
