package jit

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func Parse(expr string) (root expr, e error) {
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %v", expr, err)
	}
	defer func() {
		if err := recover(); err != nil {
			root = nil
			e = fmt.Errorf("parse %q: %v", expr, err)
		}
	}()
	return parseExpr(node), nil
}

func parseExpr(node ast.Expr) expr {
	switch node := node.(type) {
	default:
		panic(fmt.Sprintf("syntax error: %T", node))
	case *ast.BasicLit:
		return parseBasicLit(node)
	case *ast.BinaryExpr:
		return parseBinaryExpr(node)
	case *ast.CallExpr:
		return parseCallExpr(node)
	case *ast.Ident:
		return parseIdent(node)
	case *ast.ParenExpr:
		return parseExpr(node.X)
	case *ast.UnaryExpr:
		return parseUnaryExpr(node)
	}
}

func parseBasicLit(node *ast.BasicLit) expr {
	switch node.Kind {
	default:
		panic(fmt.Sprintf("syntax error: %v (%T)", node.Value, node))
	case token.FLOAT, token.INT:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			panic(err)
		}
		return constant{value: v}
	}
}

func parseBinaryExpr(node *ast.BinaryExpr) expr {
	x := parseExpr(node.X)
	y := parseExpr(node.Y)
	switch node.Op {
	default:
		panic(fmt.Sprintf("syntax error:", node.Op))
	case token.ADD, token.SUB, token.MUL, token.QUO:
		return binexpr{node.Op.String(), x, y}
	}
}

func parseCallExpr(node *ast.CallExpr) expr {
	//args := make([]expr, 0, len(node.Args))
	//for _, a := range node.Args {
	//	args = append(args, parseExpr(a))
	//}
	fun := node.Fun.(*ast.Ident).Name
	if len(node.Args) != 1 {
		panic(fmt.Sprintf("%v needs 1 argument, have %v", fun, len(node.Args)))
	}
	arg := parseExpr(node.Args[0])
	if funcs[fun] == 0 {
		panic(fmt.Sprintf("undefined: %q", fun))
	}
	return callexpr{fun, arg}
}

func parseIdent(node *ast.Ident) expr {
	switch node.Name {
	default:
		panic(fmt.Sprintf("undefined: %v", node.Name))
	case "x", "y":
		return variable{name: node.Name}
	}
}

func parseUnaryExpr(node *ast.UnaryExpr) expr {
	switch node.Op {
	default:
		panic(fmt.Sprintf("syntax error: %v", node.Op))
	case token.ADD:
		return parseExpr(node.X)
	case token.SUB:
		return binexpr{node.Op.String(), constant{value: 0}, parseExpr(node.X)}
	}
}
