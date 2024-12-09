package ast

type BlockStatement struct {
	Body []Stmt
}

func (b BlockStatement) stmt() {}

type VarDeclarationStatement struct {
	Identifier    string
	Constant      bool
	AssignedValue Expr
	ExplicitType  Type
}

func (n VarDeclarationStatement) stmt() {}

type ExpressionStatement struct {
	Expression Expr
}

func (n ExpressionStatement) stmt() {}

type Parameter struct {
	Name string
	Type Type
}

type FunctionDeclarationStatement struct {
	Parameters []Parameter
	Name       string
	Body       []Stmt
	ReturnType Type
}

func (n FunctionDeclarationStatement) stmt() {}

type IfStatement struct {
	Condition  Expr
	Consequent Stmt
	Alternate  Stmt
}

func (n IfStatement) stmt() {}

type ImportStatement struct {
	Name string
	From string
}

func (n ImportStatement) stmt() {}

type ForeachStatement struct {
	Value    string
	Index    bool
	Iterable Expr
	Body     []Stmt
}

func (n ForeachStatement) stmt() {}

type ClassDeclarationStatement struct {
	Name string
	Body []Stmt
}

func (n ClassDeclarationStatement) stmt() {}

type PublicDeclarationStatement struct {
	Value    Stmt
	Function Expr
}

func (p PublicDeclarationStatement) stmt() {}

type PrivateDeclarationStatement struct {
	Value    Stmt
	Function Expr
}

func (p PrivateDeclarationStatement) stmt() {}

type ProtectedDeclarationStatement struct {
	Value    Stmt
	Function Expr
}

func (p ProtectedDeclarationStatement) stmt() {}

type StaticDeclarationStatement struct {
	Value Stmt
}

func (p StaticDeclarationStatement) stmt() {}

type ScopeStatement struct {
	Path   string
	Parent string
}

func (n ScopeStatement) stmt() {}
