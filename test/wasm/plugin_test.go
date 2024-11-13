package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
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

	var programTree parser.AnyTree
	err = json.Unmarshal([]byte(programTreeJson), &programTree)
	require.NoError(t, err)

	// expr := []byte("test || test")
	// _, out, err := plugin.Call("compile", expr)
	// require.NoError(t, err)

	// var compileExprRes wasm_types.CompileResult
	// err = json.Unmarshal(out, &compileExprRes)
	// require.NoError(t, err)

	// patchSubExpr := []byte("false")
	// _, out, err = plugin.Call("compile", patchSubExpr)
	// require.NoError(t, err)

	patchProgram := []byte("x==2.3")
	_, out, err := plugin.Call("compile", patchProgram)
	require.NoError(t, err)

	// var compileSubExprRes wasm_types.CompileResult
	// err = json.Unmarshal(out, &compileSubExprRes)
	// require.NoError(t, err)

	var compilePatchProgramRes wasm_types.CompileResult
	err = json.Unmarshal(out, &compilePatchProgramRes)
	require.NoError(t, err)

	patchReq := wasm_types.PatchRequest{
		ProgramTree:      &programTree,
		PatchNodeID:      "5942c4bb-f4af-41a5-8380-8fe3ab84d787",
		PatchProgramTree: compilePatchProgramRes.ProgramTree,
	}
	patchReqJson, err := json.Marshal(patchReq)
	require.NoError(t, err)

	_, out, err = plugin.Call("patch", patchReqJson)
	require.NoError(t, err)

	var patchResult wasm_types.PatchResult
	err = json.Unmarshal(out, &patchResult)
	require.NoError(t, err)

	response := string(out)
	fmt.Println(response)
}

func TestPlugin_Run(t *testing.T) {
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

	program, err := expr.Compile("x || y", expr.Optimize(false))
	require.NoError(t, err)

	env := map[string]any{
		"x": true,
		"y": false,
	}

	tree := program.Tree.ToAnyTree()
	req := wasm_types.RunRequest{
		ProgramTree: &tree,
		Env:         env,
	}

	reqJson, err := json.Marshal(req)
	require.NoError(t, err)

	data := string(reqJson)
	_, out, err := plugin.Call("run", []byte(data))
	require.NoError(t, err)

	response := string(out)
	fmt.Println(response)
}

type patcher struct {
	patchNodeID string
	patchNode   ast.Node
}

func (v *patcher) Visit(node *ast.Node) {
	if (*node).ID() == v.patchNodeID {
		ast.Patch(node, v.patchNode)
	}
}

func TestPatcher(t *testing.T) {
	// program, err := expr.Compile("true || (x==1 && y==2) || (z in [1,2,3] && w==4 || (a!=5 && b==6))", expr.Optimize(false))
	// require.NoError(t, err)

	patchProgram, err := expr.Compile("x==2.3", expr.Optimize(false))
	require.NoError(t, err)

	// tree := program.Tree.ToAnyTree()
	// treeJson, err := json.MarshalIndent(tree, "", "  ")
	// require.NoError(t, err)
	// fmt.Println(string(treeJson))

	var tree parser.AnyTree
	err = json.Unmarshal([]byte(programTreeJson), &tree)
	require.NoError(t, err)

	programTree, err := parser.FromAnyTree(tree)
	require.NoError(t, err)

	// Patch
	patchedProgram, err := expr.CompileTree(
		programTree,
		expr.Patch(&patcher{
			patchNodeID: "5942c4bb-f4af-41a5-8380-8fe3ab84d787",
			patchNode:   patchProgram.Node(),
		}),
		expr.Optimize(false),
	)
	require.NoError(t, err)

	_ = patchedProgram
}

const programTreeJson = `
{
  "node": {
    "id": "65ee2de0-80e0-43f3-912c-598312cfa0f0",
    "loc": {
      "from": 23,
      "to": 25
    },
    "nat": {
      "type_name": "bool",
      "func_name": "",
      "array_of": null,
      "predicate_out": null,
      "fields": null,
      "default_map_value": null,
      "strict": false,
      "nil": false,
      "method": false,
      "method_index": 0,
      "field_index": null
    },
    "type": "binary",
    "operator": "||",
    "left": {
      "id": "7fe3b7b9-ac16-4702-9b66-57433cc09e81",
      "loc": {
        "from": 5,
        "to": 7
      },
      "nat": {
        "type_name": "bool",
        "func_name": "",
        "array_of": null,
        "predicate_out": null,
        "fields": null,
        "default_map_value": null,
        "strict": false,
        "nil": false,
        "method": false,
        "method_index": 0,
        "field_index": null
      },
      "type": "binary",
      "operator": "||",
      "left": {
        "id": "fbcf5acd-37f6-47f8-a7e6-3ddbe387ba65",
        "loc": {
          "from": 0,
          "to": 4
        },
        "nat": {
          "type_name": "bool",
          "func_name": "",
          "array_of": null,
          "predicate_out": null,
          "fields": null,
          "default_map_value": null,
          "strict": false,
          "nil": false,
          "method": false,
          "method_index": 0,
          "field_index": null
        },
        "type": "bool",
        "value": true
      },
      "right": {
        "id": "64e05ed2-9e86-40e1-9f16-32d34e3840fe",
        "loc": {
          "from": 14,
          "to": 16
        },
        "nat": {
          "type_name": "bool",
          "func_name": "",
          "array_of": null,
          "predicate_out": null,
          "fields": null,
          "default_map_value": null,
          "strict": false,
          "nil": false,
          "method": false,
          "method_index": 0,
          "field_index": null
        },
        "type": "binary",
        "operator": "&&",
        "left": {
          "id": "5942c4bb-f4af-41a5-8380-8fe3ab84d787",
          "loc": {
            "from": 10,
            "to": 12
          },
          "nat": {
            "type_name": "bool",
            "func_name": "",
            "array_of": null,
            "predicate_out": null,
            "fields": null,
            "default_map_value": null,
            "strict": false,
            "nil": false,
            "method": false,
            "method_index": 0,
            "field_index": null
          },
          "type": "binary",
          "operator": "==",
          "left": {
            "id": "f6972e5b-e92c-47e6-88ec-84a3d9535590",
            "loc": {
              "from": 9,
              "to": 10
            },
            "nat": {
              "type_name": "",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "identifier",
            "value": "x"
          },
          "right": {
            "id": "231f9bde-33ae-4551-80d5-457e38b09c6f",
            "loc": {
              "from": 12,
              "to": 13
            },
            "nat": {
              "type_name": "int",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "integer",
            "value": 1
          }
        },
        "right": {
          "id": "f94310f9-5a19-4886-8242-19ed42d60e8a",
          "loc": {
            "from": 18,
            "to": 20
          },
          "nat": {
            "type_name": "bool",
            "func_name": "",
            "array_of": null,
            "predicate_out": null,
            "fields": null,
            "default_map_value": null,
            "strict": false,
            "nil": false,
            "method": false,
            "method_index": 0,
            "field_index": null
          },
          "type": "binary",
          "operator": "==",
          "left": {
            "id": "27a70038-ba7f-49c1-9091-94b770ab346e",
            "loc": {
              "from": 17,
              "to": 18
            },
            "nat": {
              "type_name": "",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "identifier",
            "value": "y"
          },
          "right": {
            "id": "af1783fb-32b2-4dcc-b67a-b60edf3ba201",
            "loc": {
              "from": 20,
              "to": 21
            },
            "nat": {
              "type_name": "int",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "integer",
            "value": 2
          }
        }
      }
    },
    "right": {
      "id": "d3146965-1a0e-43ee-ad83-f354ab39e3ac",
      "loc": {
        "from": 48,
        "to": 50
      },
      "nat": {
        "type_name": "bool",
        "func_name": "",
        "array_of": null,
        "predicate_out": null,
        "fields": null,
        "default_map_value": null,
        "strict": false,
        "nil": false,
        "method": false,
        "method_index": 0,
        "field_index": null
      },
      "type": "binary",
      "operator": "||",
      "left": {
        "id": "8c07b73d-9ba8-49ef-b7a2-38c9012bcfe4",
        "loc": {
          "from": 40,
          "to": 42
        },
        "nat": {
          "type_name": "bool",
          "func_name": "",
          "array_of": null,
          "predicate_out": null,
          "fields": null,
          "default_map_value": null,
          "strict": false,
          "nil": false,
          "method": false,
          "method_index": 0,
          "field_index": null
        },
        "type": "binary",
        "operator": "\u0026\u0026",
        "left": {
          "id": "7617a039-9d16-42a8-9385-6577a24bd7bd",
          "loc": {
            "from": 29,
            "to": 31
          },
          "nat": {
            "type_name": "bool",
            "func_name": "",
            "array_of": null,
            "predicate_out": null,
            "fields": null,
            "default_map_value": null,
            "strict": false,
            "nil": false,
            "method": false,
            "method_index": 0,
            "field_index": null
          },
          "type": "binary",
          "operator": "in",
          "left": {
            "id": "834874a6-49be-4659-9dcb-ed0b51817537",
            "loc": {
              "from": 27,
              "to": 28
            },
            "nat": {
              "type_name": "",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "identifier",
            "value": "z"
          },
          "right": {
            "id": "b0fcdf5f-9a24-40b8-aa9d-a60e9471e550",
            "loc": {
              "from": 32,
              "to": 33
            },
            "nat": {
              "type_name": "[]interface {}",
              "func_name": "",
              "array_of": {
                "type_name": "int",
                "func_name": "",
                "array_of": null,
                "predicate_out": null,
                "fields": null,
                "default_map_value": null,
                "strict": false,
                "nil": false,
                "method": false,
                "method_index": 0,
                "field_index": null
              },
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "array",
            "nodes": [
              {
                "id": "0bb852f9-16c8-4690-91db-286f549ab0b6",
                "loc": {
                  "from": 33,
                  "to": 34
                },
                "nat": {
                  "type_name": "int",
                  "func_name": "",
                  "array_of": null,
                  "predicate_out": null,
                  "fields": null,
                  "default_map_value": null,
                  "strict": false,
                  "nil": false,
                  "method": false,
                  "method_index": 0,
                  "field_index": null
                },
                "type": "integer",
                "value": 1
              },
              {
                "id": "90060fcb-68c7-49e0-8518-00e996d9ffef",
                "loc": {
                  "from": 35,
                  "to": 36
                },
                "nat": {
                  "type_name": "int",
                  "func_name": "",
                  "array_of": null,
                  "predicate_out": null,
                  "fields": null,
                  "default_map_value": null,
                  "strict": false,
                  "nil": false,
                  "method": false,
                  "method_index": 0,
                  "field_index": null
                },
                "type": "integer",
                "value": 2
              },
              {
                "id": "9ce5d1b4-ebc8-45f6-ad5c-f96092a55433",
                "loc": {
                  "from": 37,
                  "to": 38
                },
                "nat": {
                  "type_name": "int",
                  "func_name": "",
                  "array_of": null,
                  "predicate_out": null,
                  "fields": null,
                  "default_map_value": null,
                  "strict": false,
                  "nil": false,
                  "method": false,
                  "method_index": 0,
                  "field_index": null
                },
                "type": "integer",
                "value": 3
              }
            ]
          }
        },
        "right": {
          "id": "03a0e849-1536-42ef-9335-d8db90c07c30",
          "loc": {
            "from": 44,
            "to": 46
          },
          "nat": {
            "type_name": "bool",
            "func_name": "",
            "array_of": null,
            "predicate_out": null,
            "fields": null,
            "default_map_value": null,
            "strict": false,
            "nil": false,
            "method": false,
            "method_index": 0,
            "field_index": null
          },
          "type": "binary",
          "operator": "==",
          "left": {
            "id": "577fa4d5-2d9c-4aa0-8bbb-634124f8fe2a",
            "loc": {
              "from": 43,
              "to": 44
            },
            "nat": {
              "type_name": "",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "identifier",
            "value": "w"
          },
          "right": {
            "id": "cefb0dd0-bbbe-4d94-afcf-fc54a57776e1",
            "loc": {
              "from": 46,
              "to": 47
            },
            "nat": {
              "type_name": "int",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "integer",
            "value": 4
          }
        }
      },
      "right": {
        "id": "a69a4a94-b0a1-448f-9f37-a85803d6693d",
        "loc": {
          "from": 57,
          "to": 59
        },
        "nat": {
          "type_name": "bool",
          "func_name": "",
          "array_of": null,
          "predicate_out": null,
          "fields": null,
          "default_map_value": null,
          "strict": false,
          "nil": false,
          "method": false,
          "method_index": 0,
          "field_index": null
        },
        "type": "binary",
        "operator": "\u0026\u0026",
        "left": {
          "id": "4f41e4cd-84b6-426a-94e9-715484a67deb",
          "loc": {
            "from": 53,
            "to": 55
          },
          "nat": {
            "type_name": "bool",
            "func_name": "",
            "array_of": null,
            "predicate_out": null,
            "fields": null,
            "default_map_value": null,
            "strict": false,
            "nil": false,
            "method": false,
            "method_index": 0,
            "field_index": null
          },
          "type": "binary",
          "operator": "!=",
          "left": {
            "id": "a0bc4b17-7fab-4920-8d8f-d21916329e40",
            "loc": {
              "from": 52,
              "to": 53
            },
            "nat": {
              "type_name": "",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "identifier",
            "value": "a"
          },
          "right": {
            "id": "a44c54fd-895f-476d-aca8-7989d61f61ac",
            "loc": {
              "from": 55,
              "to": 56
            },
            "nat": {
              "type_name": "int",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "integer",
            "value": 5
          }
        },
        "right": {
          "id": "c7f14785-39c0-47cc-adb4-16f2c2bae218",
          "loc": {
            "from": 61,
            "to": 63
          },
          "nat": {
            "type_name": "bool",
            "func_name": "",
            "array_of": null,
            "predicate_out": null,
            "fields": null,
            "default_map_value": null,
            "strict": false,
            "nil": false,
            "method": false,
            "method_index": 0,
            "field_index": null
          },
          "type": "binary",
          "operator": "==",
          "left": {
            "id": "a72b6b65-07e3-4a93-895b-f0709926641c",
            "loc": {
              "from": 60,
              "to": 61
            },
            "nat": {
              "type_name": "",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "identifier",
            "value": "b"
          },
          "right": {
            "id": "ac11ba1f-ba2e-43eb-ad6f-b1ca26cf782a",
            "loc": {
              "from": 63,
              "to": 64
            },
            "nat": {
              "type_name": "int",
              "func_name": "",
              "array_of": null,
              "predicate_out": null,
              "fields": null,
              "default_map_value": null,
              "strict": false,
              "nil": false,
              "method": false,
              "method_index": 0,
              "field_index": null
            },
            "type": "integer",
            "value": 6
          }
        }
      }
    }
  },
  "source": [
    116,
    114,
    117,
    101,
    32,
    124,
    124,
    32,
    40,
    120,
    61,
    61,
    49,
    32,
    38,
    38,
    32,
    121,
    61,
    61,
    50,
    41,
    32,
    124,
    124,
    32,
    40,
    122,
    32,
    105,
    110,
    32,
    91,
    49,
    44,
    50,
    44,
    51,
    93,
    32,
    38,
    38,
    32,
    119,
    61,
    61,
    52,
    32,
    124,
    124,
    32,
    40,
    97,
    33,
    61,
    53,
    32,
    38,
    38,
    32,
    98,
    61,
    61,
    54,
    41,
    41
  ]
}
`
