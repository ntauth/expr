package types

import (
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

type CompileResult struct {
	Program     *vm.Program     `json:"program"`
	ProgramTree *parser.AnyTree `json:"tree"`
}

type PatchRequest struct {
	ProgramTree      *parser.AnyTree `json:"tree"`
	PatchNodeID      string          `json:"patch_node_id"`
	PatchProgramTree *parser.AnyTree `json:"patch_tree"`
}

type PatchResult struct {
	Program     *vm.Program     `json:"program"`
	ProgramTree *parser.AnyTree `json:"tree"`
}

type RunRequest struct {
	ProgramTree *parser.AnyTree `json:"tree"`
	Env         map[string]any  `json:"env"`
}
