package ast

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expression
}

func (n ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	VariableName  string
	IsConstant    bool
	AssignedValue Expression
	ExplicitType  Type
}

func (n VarDeclStmt) stmt() {}

type StructProperty struct {
	IsStatic bool
	// Property string
	Type Type
}

type StructMethod struct {
	IsStatic bool
	// Property string
	// Type     Type
}

type StructDeclStmt struct {
	// Public     bool
	StructName string
	Properties map[string]StructProperty
	Methods    map[string]StructMethod
}

func (n StructDeclStmt) stmt() {}
