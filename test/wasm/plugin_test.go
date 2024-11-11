package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/parser"
	wasm_types "github.com/expr-lang/expr/wasm/types"
	extism "github.com/extism/go-sdk"
)

func TestPlugin_Compile(t *testing.T) {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: "../../expr.wasm",
			},
		},
	}

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
	require.NoError(t, err)

	data := []byte("test || test")
	_, out, err := plugin.Call("compile", data)
	require.NoError(t, err)

	response := string(out)
	fmt.Println(response)
}

func TestPlugin_CompileTree(t *testing.T) {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: "../../expr.wasm",
			},
		},
	}

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
	require.NoError(t, err)

	program, err := expr.Compile("true || true", expr.Optimize(false))
	require.NoError(t, err)

	treeJson, err := json.Marshal(program.Tree.ToAnyTree())
	require.NoError(t, err)

	data := string(treeJson)
	_, out, err := plugin.Call("compileTree", []byte(data))
	require.NoError(t, err)

	response := string(out)
	fmt.Println(response)
}

func TestPlugin_Patch(t *testing.T) {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: "../../expr.wasm",
			},
		},
	}

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
	require.NoError(t, err)

	expr := []byte("test || test")
	_, out, err := plugin.Call("compile", expr)
	require.NoError(t, err)

	var compileExprRes wasm_types.CompileResult
	err = json.Unmarshal(out, &compileExprRes)
	require.NoError(t, err)

	patchSubExpr := []byte("false")
	_, out, err = plugin.Call("compile", patchSubExpr)
	require.NoError(t, err)

	var compileSubExprRes wasm_types.CompileResult
	err = json.Unmarshal(out, &compileSubExprRes)
	require.NoError(t, err)

	patchReq := wasm_types.PatchRequest{
		AnyTree:      compileExprRes.AnyTree,
		Loc:          file.Location{From: 0, To: 4},
		PatchAnyTree: compileSubExprRes.AnyTree,
	}
	patchReqJson, err := json.Marshal(patchReq)
	require.NoError(t, err)

	_, out, err = plugin.Call("patch", patchReqJson)
	require.NoError(t, err)

	response := string(out)
	fmt.Println(response)
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

func TestPatcher(t *testing.T) {
	program, err := expr.Compile("true || true", expr.Optimize(false))
	require.NoError(t, err)

	// var treeGob bytes.Buffer
	// gob.Register(&ast.NilNode{})
	// gob.Register(&ast.IdentifierNode{})
	// gob.Register(&ast.IntegerNode{})
	// gob.Register(&ast.FloatNode{})
	// gob.Register(&ast.BoolNode{})
	// gob.Register(&ast.StringNode{})
	// gob.Register(&ast.ConstantNode{})
	// gob.Register(&ast.UnaryNode{})
	// gob.Register(&ast.BinaryNode{})
	// gob.Register(&ast.ChainNode{})
	// gob.Register(&ast.MemberNode{})
	// gob.Register(&ast.SliceNode{})
	// gob.Register(&ast.CallNode{})
	// gob.Register(&ast.BuiltinNode{})
	// gob.Register(&ast.PredicateNode{})
	// gob.Register(&ast.PointerNode{})
	// gob.Register(&ast.ConditionalNode{})
	// gob.Register(&ast.VariableDeclaratorNode{})
	// gob.Register(&ast.ArrayNode{})
	// gob.Register(&ast.MapNode{})
	// gob.Register(&ast.PairNode{})

	// enc := gob.NewEncoder(&treeGob)
	// err = enc.Encode(program.Tree)
	// require.NoError(t, err)

	// dec := gob.NewDecoder(&treeGob)
	// var programTree parser.Tree
	// err = dec.Decode(&programTree)
	// require.NoError(t, err)

	patchProgram, err := expr.Compile("false", expr.Optimize(false))
	require.NoError(t, err)

	tree := program.Tree.ToAnyTree()

	programTree, err := parser.FromAnyTree(tree)
	require.NoError(t, err)

	// Patch
	patchedProgram, err := expr.CompileTree(
		programTree,
		expr.Patch(&patcher{
			patchLoc:  file.Location{From: 0, To: 4},
			patchNode: patchProgram.Node(),
		}),
		expr.Optimize(false),
	)
	require.NoError(t, err)

	_ = patchedProgram
}
