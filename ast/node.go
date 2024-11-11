package ast

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/file"
)

var (
	anyType = reflect.TypeOf(new(any)).Elem()
)

// Node represents items of abstract syntax tree.
type Node interface {
	Location() file.Location
	SetLocation(file.Location)
	Nature() nature.Nature
	SetNature(nature.Nature)
	Type() reflect.Type
	SetType(reflect.Type)
	String() string
}

type AnyNodeType string

const (
	AnyNodeTypeNil                AnyNodeType = "nil"
	AnyNodeTypeIdentifier         AnyNodeType = "identifier"
	AnyNodeTypeInteger            AnyNodeType = "integer"
	AnyNodeTypeFloat              AnyNodeType = "float"
	AnyNodeTypeBool               AnyNodeType = "bool"
	AnyNodeTypeString             AnyNodeType = "string"
	AnyNodeTypeConstant           AnyNodeType = "constant"
	AnyNodeTypeUnary              AnyNodeType = "unary"
	AnyNodeTypeBinary             AnyNodeType = "binary"
	AnyNodeTypeChain              AnyNodeType = "chain"
	AnyNodeTypeMember             AnyNodeType = "member"
	AnyNodeTypeSlice              AnyNodeType = "slice"
	AnyNodeTypeCall               AnyNodeType = "call"
	AnyNodeTypeBuiltin            AnyNodeType = "builtin"
	AnyNodeTypePredicate          AnyNodeType = "predicate"
	AnyNodeTypePointer            AnyNodeType = "pointer"
	AnyNodeTypeConditional        AnyNodeType = "conditional"
	AnyNodeTypeVariableDeclarator AnyNodeType = "variable_declarator"
	AnyNodeTypeArray              AnyNodeType = "array"
	AnyNodeTypeMap                AnyNodeType = "map"
	AnyNodeTypePair               AnyNodeType = "pair"
)

type AnyNode struct {
	Base
	NodeType AnyNodeType `json:"type"`

	Operator  *string     `json:"operator,omitempty"`
	Value     any         `json:"value,omitempty"`
	Node      *AnyNode    `json:"node,omitempty"`
	Left      *AnyNode    `json:"left,omitempty"`
	Right     *AnyNode    `json:"right,omitempty"`
	From      *AnyNode    `json:"from,omitempty"`
	To        *AnyNode    `json:"to,omitempty"`
	Property  *AnyNode    `json:"property,omitempty"`
	Optional  *bool       `json:"optional,omitempty"`
	Method    *bool       `json:"method,omitempty"`
	Callee    *AnyNode    `json:"callee,omitempty"`
	Arguments *[]*AnyNode `json:"arguments,omitempty"`
	Name      *string     `json:"name,omitempty"`
	Throws    *bool       `json:"throws,omitempty"`
	Map       *AnyNode    `json:"map,omitempty"`
	Cond      *AnyNode    `json:"cond,omitempty"`
	Exp1      *AnyNode    `json:"exp1,omitempty"`
	Exp2      *AnyNode    `json:"exp2,omitempty"`
	ValueNode *AnyNode    `json:"value_node,omitempty"`
	Expr      *AnyNode    `json:"expr,omitempty"`
	Nodes     *[]*AnyNode `json:"nodes,omitempty"`
	Pairs     *[]*AnyNode `json:"pairs,omitempty"`
	Key       *AnyNode    `json:"key,omitempty"`
}

func NodeToAnyNode(node Node) *AnyNode {
	anyNode := &AnyNode{}

	switch node := node.(type) {
	case *NilNode:
		anyNode.NodeType = AnyNodeTypeNil
		anyNode.Base = node.Base
		anyNode.Value = nil
	case *IdentifierNode:
		anyNode.NodeType = AnyNodeTypeIdentifier
		anyNode.Base = node.Base
		anyNode.Value = node.Value
	case *IntegerNode:
		anyNode.NodeType = AnyNodeTypeInteger
		anyNode.Base = node.Base
		anyNode.Value = node.Value
	case *FloatNode:
		anyNode.NodeType = AnyNodeTypeFloat
		anyNode.Base = node.Base
		anyNode.Value = node.Value
	case *BoolNode:
		anyNode.NodeType = AnyNodeTypeBool
		anyNode.Base = node.Base
		anyNode.Value = node.Value
	case *StringNode:
		anyNode.NodeType = AnyNodeTypeString
		anyNode.Base = node.Base
		anyNode.Value = node.Value
	case *ConstantNode:
		anyNode.NodeType = AnyNodeTypeConstant
		anyNode.Base = node.Base
		anyNode.Value = node.Value
	case *UnaryNode:
		anyNode.NodeType = AnyNodeTypeUnary
		anyNode.Base = node.Base
		anyNode.Operator = &node.Operator
		anyNode.Node = NodeToAnyNode(node.Node)
	case *BinaryNode:
		anyNode.NodeType = AnyNodeTypeBinary
		anyNode.Base = node.Base
		anyNode.Operator = &node.Operator
		anyNode.Left = NodeToAnyNode(node.Left)
		anyNode.Right = NodeToAnyNode(node.Right)
	case *ChainNode:
		anyNode.NodeType = AnyNodeTypeChain
		anyNode.Base = node.Base
		anyNode.Node = NodeToAnyNode(node.Node)
	case *MemberNode:
		anyNode.NodeType = AnyNodeTypeMember
		anyNode.Base = node.Base
		anyNode.Node = NodeToAnyNode(node.Node)
		anyNode.Property = NodeToAnyNode(node.Property)
		anyNode.Optional = &node.Optional
		anyNode.Method = &node.Method
	case *SliceNode:
		anyNode.NodeType = AnyNodeTypeSlice
		anyNode.Base = node.Base
		anyNode.Node = NodeToAnyNode(node.Node)
		anyNode.From = NodeToAnyNode(node.From)
		anyNode.To = NodeToAnyNode(node.To)
	case *CallNode:
		anyNode.NodeType = AnyNodeTypeCall
		anyNode.Base = node.Base
		anyNode.Callee = NodeToAnyNode(node.Callee)
		arguments := make([]*AnyNode, len(node.Arguments))
		anyNode.Arguments = &arguments
		for i, arg := range node.Arguments {
			arguments[i] = NodeToAnyNode(arg)
		}
	case *BuiltinNode:
		anyNode.NodeType = AnyNodeTypeBuiltin
		anyNode.Base = node.Base
		anyNode.Name = &node.Name
		arguments := make([]*AnyNode, len(node.Arguments))
		anyNode.Arguments = &arguments
		for i, arg := range node.Arguments {
			arguments[i] = NodeToAnyNode(arg)
		}
		anyNode.Map = NodeToAnyNode(node.Map)
	case *PredicateNode:
		anyNode.NodeType = AnyNodeTypePredicate
		anyNode.Base = node.Base
		anyNode.Node = NodeToAnyNode(node.Node)
	case *PointerNode:
		anyNode.NodeType = AnyNodeTypePointer
		anyNode.Base = node.Base
		anyNode.Name = &node.Name
	case *ConditionalNode:
		anyNode.NodeType = AnyNodeTypeConditional
		anyNode.Base = node.Base
		anyNode.Cond = NodeToAnyNode(node.Cond)
		anyNode.Exp1 = NodeToAnyNode(node.Exp1)
		anyNode.Exp2 = NodeToAnyNode(node.Exp2)
	case *VariableDeclaratorNode:
		anyNode.NodeType = AnyNodeTypeVariableDeclarator
		anyNode.Base = node.Base
		anyNode.Name = &node.Name
		anyNode.Value = NodeToAnyNode(node.Value)
		anyNode.Expr = NodeToAnyNode(node.Expr)
	case *ArrayNode:
		anyNode.NodeType = AnyNodeTypeArray
		anyNode.Base = node.Base
		nodes := make([]*AnyNode, len(node.Nodes))
		anyNode.Nodes = &nodes
		for i, node := range node.Nodes {
			nodes[i] = NodeToAnyNode(node)
		}
	case *MapNode:
		anyNode.NodeType = AnyNodeTypeMap
		anyNode.Base = node.Base
		pairs := make([]*AnyNode, len(node.Pairs))
		anyNode.Pairs = &pairs
		for i, pair := range node.Pairs {
			pairs[i] = NodeToAnyNode(pair)
		}
	case *PairNode:
		anyNode.NodeType = AnyNodeTypePair
		anyNode.Base = node.Base
		anyNode.Key = NodeToAnyNode(node.Key)
		anyNode.Value = NodeToAnyNode(node.Value)
	}

	return anyNode
}

func AnyNodeToNode(anyNode *AnyNode) (Node, error) {
	var node Node

	switch anyNode.NodeType {
	case AnyNodeTypeNil:
		node = &NilNode{Base: anyNode.Base}
	case AnyNodeTypeIdentifier:
		value, ok := anyNode.Value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid identifier value: %v", anyNode.Value)
		}
		node = &IdentifierNode{Base: anyNode.Base, Value: value}
	case AnyNodeTypeInteger:
		value, ok := anyNode.Value.(int)
		if !ok {
			return nil, fmt.Errorf("invalid integer value: %v", anyNode.Value)
		}
		node = &IntegerNode{Base: anyNode.Base, Value: value}
	case AnyNodeTypeFloat:
		value, ok := anyNode.Value.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid float value: %v", anyNode.Value)
		}
		node = &FloatNode{Base: anyNode.Base, Value: value}
	case AnyNodeTypeBool:
		value, ok := anyNode.Value.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid boolean value: %v", anyNode.Value)
		}
		node = &BoolNode{Base: anyNode.Base, Value: value}
	case AnyNodeTypeString:
		value, ok := anyNode.Value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid string value: %v", anyNode.Value)
		}
		node = &StringNode{Base: anyNode.Base, Value: value}
	case AnyNodeTypeConstant:
		node = &ConstantNode{Base: anyNode.Base, Value: anyNode.Value}
	case AnyNodeTypeUnary:
		subNode, err := AnyNodeToNode(anyNode.Node)
		if err != nil {
			return nil, err
		}
		node = &UnaryNode{Base: anyNode.Base, Operator: *anyNode.Operator, Node: subNode}
	case AnyNodeTypeBinary:
		leftNode, err := AnyNodeToNode(anyNode.Left)
		if err != nil {
			return nil, err
		}
		rightNode, err := AnyNodeToNode(anyNode.Right)
		if err != nil {
			return nil, err
		}
		node = &BinaryNode{Base: anyNode.Base, Operator: *anyNode.Operator, Left: leftNode, Right: rightNode}
	case AnyNodeTypeChain:
		subNode, err := AnyNodeToNode(anyNode.Node)
		if err != nil {
			return nil, err
		}
		node = &ChainNode{Base: anyNode.Base, Node: subNode}
	case AnyNodeTypeMember:
		subNode, err := AnyNodeToNode(anyNode.Node)
		if err != nil {
			return nil, err
		}
		propertyNode, err := AnyNodeToNode(anyNode.Property)
		if err != nil {
			return nil, err
		}
		node = &MemberNode{Base: anyNode.Base, Node: subNode, Property: propertyNode, Optional: *anyNode.Optional, Method: *anyNode.Method}
	case AnyNodeTypeSlice:
		subNode, err := AnyNodeToNode(anyNode.Node)
		if err != nil {
			return nil, err
		}
		fromNode, err := AnyNodeToNode(anyNode.From)
		if err != nil {
			return nil, err
		}
		toNode, err := AnyNodeToNode(anyNode.To)
		if err != nil {
			return nil, err
		}
		node = &SliceNode{Base: anyNode.Base, Node: subNode, From: fromNode, To: toNode}
	case AnyNodeTypeCall:
		calleeNode, err := AnyNodeToNode(anyNode.Callee)
		if err != nil {
			return nil, err
		}
		callNode := &CallNode{Base: anyNode.Base, Callee: calleeNode}
		for _, argNode := range *anyNode.Arguments {
			arg, err := AnyNodeToNode(argNode)
			if err != nil {
				return nil, err
			}
			callNode.Arguments = append(callNode.Arguments, arg)
		}
		node = callNode
	case AnyNodeTypeBuiltin:
		builtinNode := &BuiltinNode{Base: anyNode.Base, Name: *anyNode.Name}
		for _, arg := range *anyNode.Arguments {
			argNode, err := AnyNodeToNode(arg)
			if err != nil {
				return nil, err
			}
			builtinNode.Arguments = append(builtinNode.Arguments, argNode)
		}
		node = builtinNode
	case AnyNodeTypePredicate:
		subNode, err := AnyNodeToNode(anyNode.Node)
		if err != nil {
			return nil, err
		}
		node = &PredicateNode{Base: anyNode.Base, Node: subNode}
	case AnyNodeTypePointer:
		if anyNode.Name == nil {
			return nil, fmt.Errorf("pointer name is nil")
		}
		node = &PointerNode{Base: anyNode.Base, Name: *anyNode.Name}
	case AnyNodeTypeConditional:
		condNode, err := AnyNodeToNode(anyNode.Cond)
		if err != nil {
			return nil, err
		}
		exp1Node, err := AnyNodeToNode(anyNode.Exp1)
		if err != nil {
			return nil, err
		}
		exp2Node, err := AnyNodeToNode(anyNode.Exp2)
		if err != nil {
			return nil, err
		}
		node = &ConditionalNode{Base: anyNode.Base, Cond: condNode, Exp1: exp1Node, Exp2: exp2Node}
	case AnyNodeTypeVariableDeclarator:
		valueNode, err := AnyNodeToNode(anyNode.ValueNode)
		if err != nil {
			return nil, err
		}
		exprNode, err := AnyNodeToNode(anyNode.Expr)
		if err != nil {
			return nil, err
		}
		node = &VariableDeclaratorNode{Base: anyNode.Base, Name: *anyNode.Name, Value: valueNode, Expr: exprNode}
	case AnyNodeTypeArray:
		nodes := make([]Node, len(*anyNode.Nodes))
		for i, elemNode := range *anyNode.Nodes {
			elem, err := AnyNodeToNode(elemNode)
			if err != nil {
				return nil, err
			}
			nodes[i] = elem
		}
		node = &ArrayNode{Base: anyNode.Base, Nodes: nodes}
	case AnyNodeTypeMap:
		pairs := make([]Node, len(*anyNode.Pairs))
		for i, pairNode := range *anyNode.Pairs {
			pair, err := AnyNodeToNode(pairNode)
			if err != nil {
				return nil, err
			}
			pairs[i] = pair
		}
		node = &MapNode{Base: anyNode.Base, Pairs: pairs}
	case AnyNodeTypePair:
		keyNode, err := AnyNodeToNode(anyNode.Key)
		if err != nil {
			return nil, err
		}
		valueNode, err := AnyNodeToNode(anyNode.ValueNode)
		if err != nil {
			return nil, err
		}
		node = &PairNode{Base: anyNode.Base, Key: keyNode, Value: valueNode}
	default:
		return nil, fmt.Errorf("unknown node type: %s", anyNode.NodeType)
	}

	return node, nil
}

// Patch replaces the node with a new one.
// Location information is preserved.
// Type information is lost.
func Patch(node *Node, newNode Node) {
	newNode.SetLocation((*node).Location())
	*node = newNode
}

// Base is a Base struct for all nodes.
type Base struct {
	Loc file.Location
	Nat *nature.NatureBase
	nat nature.Nature
}

// Location returns the location of the node in the source code.
func (n *Base) Location() file.Location {
	return n.Loc
}

// SetLocation sets the location of the node in the source code.
func (n *Base) SetLocation(loc file.Location) {
	n.Loc = loc
}

// Nature returns the nature of the node.
func (n *Base) Nature() nature.Nature {
	return n.nat
}

// SetNature sets the nature of the node.
func (n *Base) SetNature(nature nature.Nature) {
	n.nat = nature
	n.Nat = &nature.NatureBase
}

// Type returns the type of the node.
func (n *Base) Type() reflect.Type {
	if n.nat.Type == nil {
		return anyType
	}
	return n.nat.Type
}

// SetType sets the type of the node.
func (n *Base) SetType(t reflect.Type) {
	n.nat.Type = t
}

// NilNode represents nil.
type NilNode struct {
	Base
}

// IdentifierNode represents an identifier.
type IdentifierNode struct {
	Base
	Value string // Name of the identifier. Like "foo" in "foo.bar".
}

// IntegerNode represents an integer.
type IntegerNode struct {
	Base
	Value int // Value of the integer.
}

// FloatNode represents a float.
type FloatNode struct {
	Base
	Value float64 // Value of the float.
}

// BoolNode represents a boolean.
type BoolNode struct {
	Base
	Value bool // Value of the boolean.
}

// StringNode represents a string.
type StringNode struct {
	Base
	Value string // Value of the string.
}

// ConstantNode represents a constant.
// Constants are predefined values like nil, true, false, array, map, etc.
// The parser.Parse will never generate ConstantNode, it is only generated
// by the optimizer.
type ConstantNode struct {
	Base
	Value any // Value of the constant.
}

// UnaryNode represents a unary operator.
type UnaryNode struct {
	Base
	Operator string // Operator of the unary operator. Like "!" in "!foo" or "not" in "not foo".
	Node     Node   // Node of the unary operator. Like "foo" in "!foo".
}

// BinaryNode represents a binary operator.
type BinaryNode struct {
	Base
	Operator string // Operator of the binary operator. Like "+" in "foo + bar" or "matches" in "foo matches bar".
	Left     Node   // Left node of the binary operator.
	Right    Node   // Right node of the binary operator.
}

// ChainNode represents an optional chaining group.
// A few MemberNode nodes can be chained together,
// and will be wrapped in a ChainNode. Example:
//
//	foo.bar?.baz?.qux
//
// The whole chain will be wrapped in a ChainNode.
type ChainNode struct {
	Base
	Node Node // Node of the chain.
}

// MemberNode represents a member access.
// It can be a field access, a method call,
// or an array element access.
// Example:
//
//	foo.bar or foo["bar"]
//	foo.bar()
//	array[0]
type MemberNode struct {
	Base
	Node     Node // Node of the member access. Like "foo" in "foo.bar".
	Property Node // Property of the member access. For property access it is a StringNode.
	Optional bool // If true then the member access is optional. Like "foo?.bar".
	Method   bool
}

// SliceNode represents access to a slice of an array.
// Example:
//
//	array[1:4]
type SliceNode struct {
	Base
	Node Node // Node of the slice. Like "array" in "array[1:4]".
	From Node // From an index of the array. Like "1" in "array[1:4]".
	To   Node // To an index of the array. Like "4" in "array[1:4]".
}

// CallNode represents a function or a method call.
type CallNode struct {
	Base
	Callee    Node   // Node of the call. Like "foo" in "foo()".
	Arguments []Node // Arguments of the call.
}

// BuiltinNode represents a builtin function call.
type BuiltinNode struct {
	Base
	Name      string // Name of the builtin function. Like "len" in "len(foo)".
	Arguments []Node // Arguments of the builtin function.
	Throws    bool   // If true then accessing a field or array index can throw an error. Used by optimizer.
	Map       Node   // Used by optimizer to fold filter() and map() builtins.
}

// PredicateNode represents a predicate.
// Example:
//
//	filter(foo, .bar == 1)
//
// The predicate is ".bar == 1".
type PredicateNode struct {
	Base
	Node Node // Node of the predicate body.
}

// PointerNode represents a pointer to a current value in predicate.
type PointerNode struct {
	Base
	Name string // Name of the pointer. Like "index" in "#index".
}

// ConditionalNode represents a ternary operator.
type ConditionalNode struct {
	Base
	Cond Node // Condition of the ternary operator. Like "foo" in "foo ? bar : baz".
	Exp1 Node // Expression 1 of the ternary operator. Like "bar" in "foo ? bar : baz".
	Exp2 Node // Expression 2 of the ternary operator. Like "baz" in "foo ? bar : baz".
}

// VariableDeclaratorNode represents a variable declaration.
type VariableDeclaratorNode struct {
	Base
	Name  string // Name of the variable. Like "foo" in "let foo = 1; foo + 1".
	Value Node   // Value of the variable. Like "1" in "let foo = 1; foo + 1".
	Expr  Node   // Expression of the variable. Like "foo + 1" in "let foo = 1; foo + 1".
}

// ArrayNode represents an array.
type ArrayNode struct {
	Base
	Nodes []Node // Nodes of the array.
}

// MapNode represents a map.
type MapNode struct {
	Base
	Pairs []Node // PairNode nodes.
}

// PairNode represents a key-value pair of a map.
type PairNode struct {
	Base
	Key   Node // Key of the pair.
	Value Node // Value of the pair.
}
