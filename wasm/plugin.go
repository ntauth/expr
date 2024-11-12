package main

import (
	"encoding/json"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/wasm/types"
	pdk "github.com/extism/go-pdk"
)

//export compile
func compile() int32 {
	exprInput := pdk.InputString()
	program, err := expr.Compile(exprInput, expr.Optimize(false))
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	anyTree := program.Tree.ToAnyTree()
	result := types.CompileResult{
		Program:     program,
		ProgramTree: &anyTree,
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	mem := pdk.AllocateString(string(resultJson))

	// zero-copy output to host
	pdk.OutputMemory(mem)

	return 0
}

//export compileTree
func compileTree() int32 {
	anyTreeInput := pdk.InputString()

	var anyTree parser.AnyTree
	err := json.Unmarshal([]byte(anyTreeInput), &anyTree)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	tree, err := parser.FromAnyTree(anyTree)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	program, err := expr.CompileTree(tree, expr.Optimize(false))
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	anyTree = program.Tree.ToAnyTree()
	result := types.CompileResult{
		Program:     program,
		ProgramTree: &anyTree,
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	mem := pdk.AllocateString(string(resultJson))

	// zero-copy output to host
	pdk.OutputMemory(mem)

	return 0
}

type patcher struct {
	patchLoc  file.Location
	patchNode ast.Node
}

func (v *patcher) Visit(node *ast.Node) {
	if (*node).Location().From == v.patchLoc.From || (*node).Location().To == v.patchLoc.To {
		ast.Patch(node, v.patchNode)
	}
}

//export patch
func patch() int32 {
	patchRequestInput := pdk.InputString()
	var patchRequest types.PatchRequest
	err := json.Unmarshal([]byte(patchRequestInput), &patchRequest)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	tree, err := parser.FromAnyTree(*patchRequest.ProgramTree)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	patchTree, err := parser.FromAnyTree(*patchRequest.PatchProgramTree)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	// Patch
	patchedProgram, err := expr.CompileTree(
		tree,
		expr.Patch(&patcher{
			patchLoc:  patchRequest.Loc,
			patchNode: patchTree.Node,
		}),
		expr.Optimize(false),
	)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	anyTree := patchedProgram.Tree.ToAnyTree()
	result := types.PatchResult{
		Program:     patchedProgram,
		ProgramTree: &anyTree,
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	mem := pdk.AllocateString(string(resultJson))
	// zero-copy output to host
	pdk.OutputMemory(mem)

	return 0
}

//export run
func run() int32 {
	runRequestInput := pdk.InputString()
	var runRequest types.RunRequest
	err := json.Unmarshal([]byte(runRequestInput), &runRequest)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	tree, err := parser.FromAnyTree(*runRequest.ProgramTree)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	program, err := expr.CompileTree(tree, expr.Optimize(false))
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	result, err := expr.Run(program, runRequest.Env)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return -1
	}

	mem := pdk.AllocateString(string(resultJson))
	// zero-copy output to host
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
