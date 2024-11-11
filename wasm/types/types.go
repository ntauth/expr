package types

import (
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

type CompileResult struct {
	Program *vm.Program     `json:"program"`
	AnyTree *parser.AnyTree `json:"tree"`
}

type PatchRequest struct {
	AnyTree      *parser.AnyTree `json:"tree"`
	Loc          file.Location   `json:"loc"`
	PatchAnyTree *parser.AnyTree `json:"patch_tree"`
}

type PatchResult struct {
	Program *vm.Program     `json:"program"`
	AnyTree *parser.AnyTree `json:"tree"`
}
