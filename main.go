package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	checkErr(err)

	var request pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(input, &request)
	checkErr(err)
	response, err := generate(&request)
	checkErr(err)

	out, err := proto.Marshal(response)
	checkErr(err)
	fmt.Fprint(os.Stdout, string(out))
}
