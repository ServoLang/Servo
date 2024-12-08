package ast

type Stmt interface {
	stmt()
}

type Expression interface {
	expr()
}

type Type interface {
	_type()
}
