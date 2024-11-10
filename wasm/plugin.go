package main

import (
	"encoding/json"

	"github.com/expr-lang/expr"
	pdk "github.com/extism/go-pdk"
)

//export test_run
func test_run() int32 {
	pdk.Log(pdk.LogInfo, "test_run")
	return 1
}

//export compile
func compile() string {
	exprInput := pdk.InputString()
	program, err := expr.Compile(exprInput)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return ""
	}

	programJson, err := json.Marshal(program)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return ""
	}

	return string(programJson)
}

func main() {}
