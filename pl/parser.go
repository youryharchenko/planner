package pl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

const SliceDefaultCap = 10

type Parser struct {
	lexer *Lexer
	env   *Env

	Done         chan bool
	reqStop      chan bool
	ReqReset     chan io.RuneScanner
	AddInput     chan io.RuneScanner
	ParsedOutput chan []ParserReply

	sendMe            []ParserReply
	FlagSendNeedInput bool
}

type ParserReply struct {
	Expr []Expression
	Err  error
}

func (env *Env) NewParser() *Parser {
	p := &Parser{
		env:          env,
		Done:         make(chan bool),
		reqStop:      make(chan bool),
		ReqReset:     make(chan io.RuneScanner),
		AddInput:     make(chan io.RuneScanner),
		ParsedOutput: make(chan []ParserReply),
		sendMe:       make([]ParserReply, 0, 1),
	}

	p.lexer = NewLexer(p)
	return p
}

func (p *Parser) Start() {
	go func() {
		//log.Println("Parser: started")
		defer close(p.Done)
		expressions := make([]Expression, 0, SliceDefaultCap)

		// maybe we already have input, be optimistic!
		// no need to call p.GetMoreInput() before staring
		// our loop.

		for {
			expr, err := p.ParseExpression(0)
			//log.Println("Parser: expression:", expr.String())
			//if err != nil || expr == SexpEnd {
			if err != nil || expr == ExprEnd {
				if err == ParserHaltRequested {
					log.Println("Parser: ParserHaltRequested")
					return
				}
				err = p.GetMoreInput(expressions, err)
				if err == ParserHaltRequested {
					log.Println("Parser: ParserHaltRequested")
					return
				}
				// GetMoreInput will have delivered what we gave them. Reset since we
				// don't own that memory any more.
				expressions = make([]Expression, 0, SliceDefaultCap)
			} else {
				// INVAR: err == nil && expr is not SexpEnd
				expressions = append(expressions, expr)
				//log.Println("Parser: expressions:", expressions)
			}
		}
	}()
}

var ParserHaltRequested = fmt.Errorf("parser halt requested")
var ResetRequested = fmt.Errorf("parser reset requested")

var ErrMoreInputNeeded = fmt.Errorf("parser needs more input")

func (p *Parser) GetMoreInput(deliverThese []Expression, errorToReport error) error {

	if len(deliverThese) == 0 && errorToReport == nil {
		p.FlagSendNeedInput = true
	} else {
		p.sendMe = append(p.sendMe,
			ParserReply{
				Expr: deliverThese,
				Err:  errorToReport,
			})
	}

	for {
		select {
		case <-p.reqStop:
			return ParserHaltRequested
		case input := <-p.AddInput:
			p.lexer.AddNextStream(input)
			p.FlagSendNeedInput = false
			return nil
		case input := <-p.ReqReset:
			p.lexer.Reset()
			p.lexer.AddNextStream(input)
			p.FlagSendNeedInput = false
			return ResetRequested
		case p.HaveStuffToSend() <- p.sendMe:
			p.sendMe = make([]ParserReply, 0, 1)
			p.FlagSendNeedInput = false
		}
	}
}

func (p *Parser) HaveStuffToSend() chan []ParserReply {
	if len(p.sendMe) > 0 || p.FlagSendNeedInput {
		return p.ParsedOutput
	}
	return nil
}

func (p *Parser) ResetAddNewInput(s io.RuneScanner) {
	select {
	case p.ReqReset <- s:
	case <-p.reqStop:
	}
}

func (p *Parser) ParseTokens() ([]Expression, error) {
	select {
	case out := <-p.ParsedOutput:
		Q("ParseTokens got p.ParsedOutput out: '%#v'", out)
		r := make([]Expression, 0)
		for _, k := range out {
			r = append(r, k.Expr...)
			//Q("\n ParseTokens k.Expr = '%v'\n\n", (&SexpArray{Val: k.Expr, Env: p.env}).SexpString(nil))
			if k.Err != nil {
				return r, k.Err
			}
		}
		return r, nil
	case <-p.reqStop:
		return nil, ErrShuttingDown
	}
}

var ErrShuttingDown error = fmt.Errorf("lexer shutting down")

func (parser *Parser) ParseListOld(depth int) (sx Expression, err error) {
	lexer := parser.lexer
	var tok Token

tokFilled:
	for {
		tok, err = lexer.PeekNextToken()
		//Q("\n ParseList(depth=%d) got lexer.PeekNextToken() -> tok='%v' err='%v'\n", depth, tok, err)
		if err != nil {
			return ExprNull, err
		}
		if tok.typ != TokenEnd {
			break tokFilled
		}
		// instead of returning UnexpectedEnd, we:
		err = parser.GetMoreInput(nil, ErrMoreInputNeeded)
		//Q("\n ParseList(depth=%d) got back from parser.GetMoreInput(): '%v'\n", depth, err)
		switch err {
		case ParserHaltRequested:
			return ExprNull, err
		case ResetRequested:
			return ExprEnd, err
		}
		// have to still fill tok, so
		// loop to the top to PeekNextToken
	}

	if tok.typ == TokenRParen {
		_, _ = lexer.GetNextToken()
		return ExprNull, nil
	}

	//var start = &SexpPair{}

	expr, err := parser.ParseExpression(depth + 1)
	if err != nil {
		return ExprNull, err
	}

	//start.Head = expr
	var start = NewPair(expr, ExprNull)

	tok, err = lexer.PeekNextToken()
	if err != nil {
		return ExprNull, err
	}

	// backslash '\' replaces dot '.' in zygomys
	if tok.typ == TokenBackslash {
		// eat up the backslash
		_, _ = lexer.GetNextToken()
		expr, err = parser.ParseExpression(depth + 1)
		if err != nil {
			return ExprNull, err
		}

		// eat up the end paren
		tok, err = lexer.GetNextToken()
		if err != nil {
			return ExprNull, err
		}
		// make sure it was actually an end paren
		if tok.typ != TokenRParen {
			return ExprNull, errors.New("extra value in dotted pair")
		}
		start.tail = expr
		return start, nil
	}

	expr, err = parser.ParseList(depth + 1)
	if err != nil {
		return start, err
	}
	start.tail = expr

	return start, nil
}

func (parser *Parser) ParseList(depth int) (Expression, error) {
	lexer := parser.lexer
	arr := make([]Expression, 0, SliceDefaultCap)

	var tok Token
	var err error
	for {
	getTok:
		for {
			tok, err = lexer.PeekNextToken()
			if err != nil {
				return ExprEnd, err
			}

			if tok.typ == TokenComma {
				// pop off the ,
				_, _ = lexer.GetNextToken()
				continue getTok
			}

			if tok.typ != TokenEnd {
				break getTok
			} else {
				//instead of return SexpEnd, UnexpectedEnd
				// we ask for more, and then loop
				err = parser.GetMoreInput(nil, ErrMoreInputNeeded)
				switch err {
				case ParserHaltRequested:
					return ExprNull, err
				case ResetRequested:
					return ExprEnd, err
				}
			}
		}

		if tok.typ == TokenRParen {
			// pop off the ]
			_, _ = lexer.GetNextToken()
			break
		}

		expr, err := parser.ParseExpression(depth + 1)
		if err != nil {
			return ExprNull, err
		}
		arr = append(arr, expr)
	}

	//return &SexpArray{Val: arr, Env: parser.env}, nil
	return NewLlist(arr...), nil
}

func (parser *Parser) ParseArray(depth int) (Expression, error) {
	lexer := parser.lexer
	arr := make([]Expression, 0, SliceDefaultCap)

	var tok Token
	var err error
	for {
	getTok:
		for {
			tok, err = lexer.PeekNextToken()
			if err != nil {
				return ExprEnd, err
			}

			if tok.typ == TokenComma {
				// pop off the ,
				_, _ = lexer.GetNextToken()
				continue getTok
			}

			if tok.typ != TokenEnd {
				break getTok
			} else {
				//instead of return SexpEnd, UnexpectedEnd
				// we ask for more, and then loop
				err = parser.GetMoreInput(nil, ErrMoreInputNeeded)
				switch err {
				case ParserHaltRequested:
					return ExprNull, err
				case ResetRequested:
					return ExprEnd, err
				}
			}
		}

		if tok.typ == TokenRSquare {
			// pop off the ]
			_, _ = lexer.GetNextToken()
			break
		}

		expr, err := parser.ParseExpression(depth + 1)
		if err != nil {
			return ExprNull, err
		}
		arr = append(arr, expr)
	}

	//return &SexpArray{Val: arr, Env: parser.env}, nil
	return NewPlist(arr...), nil
}

func (parser *Parser) ParseExpression(depth int) (res Expression, err error) {
	defer func() {
		if res != nil {
			//Q("returning from ParseExpression at depth=%v with res='%s'\n", depth, res.SexpString())
		} else {
			//Q("returning from ParseExpression at depth=%v, res = nil", depth)
		}
	}()

	lexer := parser.lexer
	//env := parser.env

	//getAnother:
	tok, err := lexer.GetNextToken()
	//log.Println("ParseExpression: next token:", tok, err)
	if err != nil {
		return ExprEnd, err
	}

	switch tok.typ {
	case TokenLParen:
		exp, err := parser.ParseList(depth + 1)
		//log.Println("ParseExpression: parsed list:", exp, err)
		return exp, err
	case TokenLSquare:
		exp, err := parser.ParseArray(depth + 1)
		//log.Println("ParseExpression: parsed array:", exp, err)
		return exp, err
	//case TokenLCurly:
	//	exp, err := parser.ParseInfix(depth + 1)
	//	return exp, err
	//case TokenQuote:
	//	expr, err := parser.ParseExpression(depth + 1)
	//	if err != nil {
	//		return SexpNull, err
	//	}
	//	return MakeList([]Sexp{env.MakeSymbol("quote"), expr}), nil
	//case TokenCaret:
	//	// '^' is now our syntax-quote symbol, not TokenBacktick, to allow go-style `string literals`.
	//	expr, err := parser.ParseExpression(depth + 1)
	//	if err != nil {
	//		return SexpNull, err
	//	}
	//	return MakeList([]Sexp{env.MakeSymbol("syntaxQuote"), expr}), nil
	//case TokenTilde:
	//	expr, err := parser.ParseExpression(depth + 1)
	//	if err != nil {
	//		return SexpNull, err
	//	}
	//	return MakeList([]Sexp{env.MakeSymbol("unquote"), expr}), nil
	//case TokenTildeAt:
	//	expr, err := parser.ParseExpression(depth + 1)
	//	if err != nil {
	//		return SexpNull, err
	//	}
	//	return MakeList([]Sexp{env.MakeSymbol("unquote-splicing"), expr}), nil
	//case TokenFreshAssign:
	//	return env.MakeSymbol(tok.str), nil
	//case TokenColonOperator:
	//	return env.MakeSymbol(tok.str), nil
	//case TokenDollar:
	//	return env.MakeSymbol(tok.str), nil
	//case TokenBool:
	//	return &SexpBool{Val: tok.str == "true"}, nil
	case TokenDecimal:
		i, err := strconv.ParseInt(tok.str, 10, ExprIntSize)
		if err != nil {
			return ExprNull, err
		}
		//return &SexpInt{Val: i}, nil
		return NewInt(i), nil
	case TokenHex:
		i, err := strconv.ParseInt(tok.str, 16, ExprIntSize)
		if err != nil {
			return ExprNull, err
		}
		//return &SexpInt{Val: i}, nil
		return NewInt(i), nil
	case TokenOct:
		i, err := strconv.ParseInt(tok.str, 8, ExprIntSize)
		if err != nil {
			return ExprNull, err
		}
		//return &SexpInt{Val: i}, nil
		return NewInt(i), nil
	case TokenBinary:
		i, err := strconv.ParseInt(tok.str, 2, ExprIntSize)
		if err != nil {
			return ExprNull, err
		}
		//return &SexpInt{Val: i}, nil
		return NewInt(i), nil
	//case TokenChar:
	//	return &SexpChar{Val: rune(tok.str[0])}, nil
	case TokenString:
		//return &SexpStr{S: tok.str}, nil
		return NewWord(tok.str), nil
	//case TokenBacktickString:
	//	return &SexpStr{S: tok.str, backtick: true}, nil
	case TokenFloat:
		f, err := strconv.ParseFloat(tok.str, ExprFloatSize)
		if err != nil {
			return ExprNull, err
		}
		//return &SexpFloat{Val: f}, nil
		return NewFloat(f), nil
	case TokenEnd:
		return ExprEnd, nil
	case TokenSymbol:
		//	return env.MakeSymbol(tok.str), nil
		return NewWord(tok.str), nil
	//case TokenSymbolColon:
	//	sym := env.MakeSymbol(tok.str)
	//	sym.colonTail = true
	//	return sym, nil
	//case TokenDot:
	//	sym := env.MakeSymbol(tok.str)
	//	sym.isDot = true
	//	return sym, nil
	case TokenDotSymbol:
		//	sym := env.MakeSymbol(tok.str)
		//	sym.isDot = true
		//	return sym, nil
		return NewRef(LocalValue, NewWord(tok.str[1:])), nil
	case TokenComment:
		//Q("parser making SexpComment from '%s'", tok.str)
		//return &SexpComment{Comment: tok.str}, nil
		return NewComment(tok.str), nil
		// parser skips comments
		//goto getAnother
		//case TokenBeginBlockComment:
		// parser skips comments
		//	return parser.ParseBlockComment(&tok)
		//parser.ParseBlockComment(&tok)
		//goto getAnother
		//case TokenComma:
		//	return &SexpComma{}, nil
		//case TokenSemicolon:
		//	return &SexpSemicolon{}, nil
	}
	return ExprNull, fmt.Errorf("Invalid syntax, don't know what to do with %v '%v'", tok.typ, tok)
}