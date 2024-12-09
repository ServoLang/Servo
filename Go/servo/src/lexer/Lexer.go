package lexer

func New(input string) *Lexer {
	l := &Lexer{Input: input}
	return l
}
