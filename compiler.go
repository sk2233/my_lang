/*
@author: sk
@date: 2024/11/16
*/
package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	Source string
	Index  int
	Line   int
}

func (s *Scanner) Scan() []*Token {
	res := make([]*Token, 0)
	for s.HasMore() {
		if token := s.ScanToken(); token != nil {
			res = append(res, token)
		}
	}

	res = append(res, &Token{
		Type: TokenEOF,
		Line: s.Line,
	})
	return res
}

var (
	keywords = map[string]TokenType{
		"and":   TokenAnd,
		"class": TokenClass,
		"else":  TokenElse,
		"false": TokenFalse,
		"for":   TokenFor,
		"func":  TokenFunc,
		"if":    TokenIf,
		"or":    TokenOr,
		//"print":  TokenPrint,
		"return": TokenReturn,
		"super":  TokenSuper,
		"this":   TokenThis,
		"true":   TokenTrue,
		"var":    TokenVar,
		"while":  TokenWhile,
	}
)

func (s *Scanner) ScanToken() *Token {
	ch := s.Read()
	switch ch {
	// 单字符
	case '(':
		return &Token{Value: "(", Type: TokenLParen, Line: s.Line}
	case ')':
		return &Token{Value: ")", Type: TokenRParen, Line: s.Line}
	case '{':
		return &Token{Value: "{", Type: TokenLBrace, Line: s.Line}
	case '}':
		return &Token{Value: "}", Type: TokenRBrace, Line: s.Line}
	case ';':
		return &Token{Value: ";", Type: TokenSemi, Line: s.Line}
	case ',':
		return &Token{Value: ",", Type: TokenComma, Line: s.Line}
	case '.':
		return &Token{Value: ".", Type: TokenDot, Line: s.Line}
	case '-':
		return &Token{Value: "-", Type: TokenSub, Line: s.Line}
	case '+':
		return &Token{Value: "+", Type: TokenAdd, Line: s.Line}
	case '*':
		return &Token{Value: "*", Type: TokenMul, Line: s.Line}
	// 单 or 双字符
	case '!':
		if s.Match('=') {
			return &Token{Value: "!=", Type: TokenNE, Line: s.Line}
		}
		return &Token{Value: "!", Type: TokenNot, Line: s.Line}
	case '=':
		if s.Match('=') {
			return &Token{Value: "==", Type: TokenEQ, Line: s.Line}
		}
		return &Token{Value: "=", Type: TokenAssign, Line: s.Line}
	case '>':
		if s.Match('=') {
			return &Token{Value: ">=", Type: TokenGE, Line: s.Line}
		}
		return &Token{Value: ">", Type: TokenGT, Line: s.Line}
	case '<':
		if s.Match('=') {
			return &Token{Value: "<=", Type: TokenLE, Line: s.Line}
		}
		return &Token{Value: "<", Type: TokenLT, Line: s.Line}
	case '/':
		if s.Match('/') { // 单行注释支持
			for s.HasMore() && s.Source[s.Index] != '\n' { // 忽略一行
				s.Index++
			}
			return nil
		}
		return &Token{Value: "/", Type: TokenDiv, Line: s.Line}
	// 忽略空字符
	case ' ', '\r', '\t':
		return nil
	case '\n': // 比较特别的空字符
		s.Line++
		return nil
	// 字面量
	case '"': // 字符串必须是单行的
		start := s.Index
		for s.HasMore() && s.Source[s.Index] != '"' {
			s.Index++
		}
		if !s.HasMore() {
			panic(fmt.Sprintf("no end string err line %d", s.Line))
		}
		s.Index++
		return &Token{Value: s.Source[start : s.Index-1], Type: TokenStr, Line: s.Line}
	default:
		if IsDigit(ch) { // 数字解析
			start := s.Index - 1
			for s.HasMore() && IsDigit(s.Source[s.Index]) {
				s.Index++
			} // 注意短路与运算顺序
			if s.Match('.') && s.HasMore() && IsDigit(s.Source[s.Index]) {
				for s.HasMore() && IsDigit(s.Source[s.Index]) { // 处理小数点后面的数字
					s.Index++
				}
			}
			return &Token{Value: s.Source[start:s.Index], Type: TokenNum, Line: s.Line}
		} else if IsAlpha(ch) { // 标识符&关键字 处理
			start := s.Index - 1
			for s.HasMore() && (IsAlpha(s.Source[s.Index]) || IsDigit(s.Source[s.Index]) || s.Source[s.Index] == '_') {
				s.Index++
			}
			str := s.Source[start:s.Index]
			if tokenType, ok := keywords[str]; ok { // 关键字
				return &Token{Value: str, Type: tokenType, Line: s.Line}
			} // 自定义标识符
			return &Token{Value: str, Type: TokenId, Line: s.Line}
		} else {
			panic(fmt.Sprintf("unknown ch %v line %d", ch, s.Line))
		}
	}
}

func IsAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func IsDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (s *Scanner) HasMore() bool {
	return s.Index < len(s.Source)
}

func (s *Scanner) Read() byte {
	s.Index++
	return s.Source[s.Index-1]
}

func (s *Scanner) Match(ch byte) bool {
	if !s.HasMore() {
		return false
	}
	if s.Source[s.Index] == ch {
		s.Index++
		return true
	}
	return false
}

func NewScanner(source string) *Scanner {
	return &Scanner{Source: source, Index: 0, Line: 1}
}

type Token struct {
	Value string
	Type  TokenType
	Line  int
}

type Local struct { // 局部变量
	Name  string
	Depth int
}

type Parser struct {
	Tokens []*Token
	Index  int
	Chunk  *Chunk
	Depth  int      // 深度
	Locals []*Local // 局部变量
}

func (p *Parser) Parse() *Chunk {
	for !p.Match(TokenEOF) {
		p.Declaration()
	}
	return p.Chunk
}

func (p *Parser) Declaration() {
	if p.Match(TokenVar) {
		p.VarDeclaration()
	} else if p.Match(TokenClass) {
		p.ClassDeclaration()
	} else if p.Match(TokenFunc) {
		p.FuncDeclaration()
	} else {
		p.Statement()
	}
}

func (p *Parser) ClassDeclaration() {
	name := p.Must(TokenId) // 类名称
	super := ""
	if p.Match(TokenLT) { // class < sup_class
		super = p.Must(TokenId).Value
	}
	if name.Value == super { // 不能继承自己
		panic(fmt.Sprintf("class %v super can't be self line %d", name.Value, name.Line))
	}
	p.Must(TokenLBrace)
	methods := make(map[string]*Value)
	for !p.Match(TokenRBrace) {
		// 解析方法，只可能是方法
		method := p.Method()
		methods[method.Func.Name] = method
	}
	class := &Class{Name: name.Value, Super: super, Methods: methods}
	p.MakeConstant(NewClass(class), name.Line)
	if len(super) > 0 { // 存在继承处理继承
		p.Chunk.WriteOpCode(OpInherit, name.Line)
	}
	p.MakeDeclaration(NewStr(name.Value), name.Line)
}

func (p *Parser) Method() *Value {
	name := p.Must(TokenId)
	skipOffset := p.MakeJump(OpJump, name.Line) // 执行时跳过整个函数定义
	funcIp := len(p.Chunk.Data)
	p.BeginScope() // 入参也在范围内
	p.Must(TokenLParen)
	count := uint8(0)
	if !p.Match(TokenRParen) { // 有参数准备参数
		param := p.Must(TokenId) // 参数会在函数调用时埋进来 这里准备一下局部变量表就行
		p.Locals = append(p.Locals, &Local{Name: param.Value, Depth: p.Depth})
		count++
		for p.Match(TokenComma) {
			param = p.Must(TokenId)
			p.Locals = append(p.Locals, &Local{Name: param.Value, Depth: p.Depth})
			count++
		}
		p.Must(TokenRParen)
	}
	p.Must(TokenLBrace)
	p.Block() // 这里使用的局部变量偏移都是从 0 开始的
	// 一般前面应该有 return 返回值的，若到了这里还没有返回值就塞一个 这里默认返回 this 所有实例方法默认支持链式调用 init 默认返回实例
	p.Chunk.WriteOpCode(OpThis, name.Line)
	endIp := len(p.Chunk.Data)                  // 函数中的 return 会返回到这里，因为还有参数栈要处理
	p.Chunk.WriteOpCode(OpFixReturn, name.Line) // 直接弹出参数会先把返回值弹出，需要交换一下返回值的位置
	p.EndScope()
	p.Chunk.WriteOpCode(OpEndFunc, name.Line) // 结束函数，返回之前的位置
	// 返回地址
	p.Chunk.SetIndex(skipOffset, uint64(len(p.Chunk.Data)))
	func0 := &Func{Ip: funcIp, EndIp: endIp, ArgCnt: count, Name: name.Value} // 先只考虑全局函数  局部函数存储在栈中且需要额外记录参数偏移受局部函数位置的影响
	return NewFunc(func0)
}

func (p *Parser) FuncDeclaration() { // 暂时只支持函数定义在最外层 方便局部变量偏移计算
	name := p.Must(TokenId)
	skipOffset := p.MakeJump(OpJump, name.Line) // 执行时跳过整个函数定义
	funcIp := len(p.Chunk.Data)
	p.BeginScope() // 入参也在范围内
	p.Must(TokenLParen)
	count := uint8(0)
	if !p.Match(TokenRParen) { // 有参数准备参数
		param := p.Must(TokenId) // 参数会在函数调用时埋进来 这里准备一下局部变量表就行
		p.Locals = append(p.Locals, &Local{Name: param.Value, Depth: p.Depth})
		count++
		for p.Match(TokenComma) {
			param = p.Must(TokenId)
			p.Locals = append(p.Locals, &Local{Name: param.Value, Depth: p.Depth})
			count++
		}
		p.Must(TokenRParen)
	}
	p.Must(TokenLBrace)
	p.Block() // 这里使用的局部变量偏移都是从 0 开始的
	// 一般前面应该有 return 返回值的，若到了这里还没有返回值就塞一个
	p.MakeConstant(NewNum(0), name.Line)
	endIp := len(p.Chunk.Data)                  // 函数中的 return 会返回到这里，因为还有参数栈要处理
	p.Chunk.WriteOpCode(OpFixReturn, name.Line) // 直接弹出参数会先把返回值弹出，需要交换一下返回值的位置
	p.EndScope()
	p.Chunk.WriteOpCode(OpEndFunc, name.Line) // 结束函数，返回之前的位置
	// 返回地址
	p.Chunk.SetIndex(skipOffset, uint64(len(p.Chunk.Data)))
	func0 := &Func{Ip: funcIp, EndIp: endIp, ArgCnt: count, Name: name.Value} // 先只考虑全局函数  局部函数存储在栈中且需要额外记录参数偏移受局部函数位置的影响
	p.MakeConstant(NewFunc(func0), name.Line)
	p.MakeDeclaration(NewStr(name.Value), name.Line)
}

func (p *Parser) VarDeclaration() {
	token := p.Must(TokenId)
	val := NewStr(token.Value)
	p.Must(TokenAssign)
	p.Expression()
	token = p.Must(TokenSemi)
	if p.Depth > 0 { // 局部变量 下标正好与栈内元素下标对其 Expression 结果不弹出了，就留在栈内
		for i := len(p.Locals) - 1; i >= 0; i-- {
			if p.Locals[i].Depth < p.Depth { // 深度在数组中是从小到大排列的 可以提前结束
				break
			}
			if p.Locals[i].Name == val.Str { // 防止局部变量重复定义
				panic(fmt.Sprintf("repeated declaration local %v line %d", val.Str, token.Line))
			}
		}
		p.Locals = append(p.Locals, &Local{Name: val.Str, Depth: p.Depth})
	} else { // 全局变量
		p.MakeDeclaration(val, token.Line)
	}
}

func (p *Parser) Statement() {
	//if p.Match(TokenPrint) {
	//	p.PrintStatement()
	//} else
	if p.Match(TokenIf) {
		p.IfStatement()
	} else if p.Match(TokenWhile) {
		p.WhileStatement()
	} else if p.Match(TokenReturn) {
		p.ReturnStatement()
	} else if p.Match(TokenFor) {
		p.ForStatement()
	} else if p.Match(TokenLBrace) {
		p.BeginScope()
		p.Block()
		p.EndScope()
	} else {
		p.ExpressionStatement()
	}
}

func (p *Parser) ReturnStatement() {
	if p.Match(TokenSemi) {
		p.MakeConstant(NewNum(0), 0) // 默认返回值
	} else {
		p.Expression() // 有返回值
		p.Must(TokenSemi)
	}
	p.Chunk.WriteOpCode(OpReturn, 0)
}

func (p *Parser) ForStatement() {
	p.Must(TokenLParen)
	if p.Match(TokenVar) { // 初始化语句
		p.VarDeclaration()
	} else {
		p.Must(TokenSemi)
	}
	startIp := len(p.Chunk.Data)
	condOffset := -1
	if !p.Match(TokenSemi) { // 有条件
		p.Expression()
		condOffset = p.MakeJump(OpFJump, 0)
		p.Must(TokenSemi)
	}
	bodyOffset := p.MakeJump(OpJump, 0)
	changeIp := len(p.Chunk.Data)
	if !p.Match(TokenRParen) { // 有变化
		p.Expression()
		p.Chunk.WriteOpCode(OpPop, 0) // 抛弃其返回值 只是需要Expression的副作用
		p.Must(TokenRParen)
	}
	p.MakeDirectJump(OpJump, startIp, 0)
	p.Chunk.SetIndex(bodyOffset, uint64(len(p.Chunk.Data)))
	p.Statement()
	p.MakeDirectJump(OpJump, changeIp, 0)
	if condOffset >= 0 { // 有跳出条件
		p.Chunk.SetIndex(condOffset, uint64(len(p.Chunk.Data)))
	}
}

func (p *Parser) WhileStatement() {
	p.Must(TokenLParen)
	startIp := len(p.Chunk.Data)
	p.Expression()
	token := p.Must(TokenRParen)
	exitOffset := p.MakeJump(OpFJump, token.Line)
	p.Statement()
	p.MakeDirectJump(OpJump, startIp, 0) // 不停跳回判断条件
	p.Chunk.SetIndex(exitOffset, uint64(len(p.Chunk.Data)))
}

func (p *Parser) MakeDirectJump(code OpCode, ip int, line int) {
	p.Chunk.WriteOpCode(code, line)
	p.Chunk.WriteIndex(uint64(ip))
}

func (p *Parser) IfStatement() {
	p.Must(TokenLParen)
	p.Expression()
	token := p.Must(TokenRParen)
	elseOffset := p.MakeJump(OpFJump, token.Line)
	p.Statement()
	ifOffset := p.MakeJump(OpJump, 0) // 应该使用上一个的 line
	p.Chunk.SetIndex(elseOffset, uint64(len(p.Chunk.Data)))
	if p.Match(TokenElse) {
		p.Statement()
	}
	p.Chunk.SetIndex(ifOffset, uint64(len(p.Chunk.Data)))
}

func (p *Parser) MakeJump(code OpCode, line int) int {
	p.Chunk.WriteOpCode(code, line)
	p.Chunk.WriteIndex(0) // 预先使用 0 占位 并返回其地址
	return len(p.Chunk.Data) - 8
}

func (p *Parser) BeginScope() {
	p.Depth++
}

func (p *Parser) EndScope() {
	p.Depth--
	if len(p.Locals) == 0 {
		return // 没有局部变量需要处理
	}
	for i := len(p.Locals) - 1; i >= 0; i-- {
		if p.Locals[i].Depth < p.Depth {
			p.Locals = p.Locals[:i+1] // 该结束了 局部变量列表也要裁剪
			return
		} // 移除该范围的所有局部变量
		p.Chunk.WriteOpCode(OpPop, 0)
	}
	p.Locals = make([]*Local, 0) // 没有中断全部局部变量都没了
}

func (p *Parser) Block() {
	for !p.Match(TokenRBrace) { // 块声明无需 ;
		p.Declaration()
	}
}

func (p *Parser) ExpressionStatement() {
	p.Expression()
	token := p.Must(TokenSemi)
	p.Chunk.WriteOpCode(OpPop, token.Line) // 表达式的返回值直接丢弃
}

//func (p *Parser) PrintStatement() {
//	p.Expression()
//	token := p.Must(TokenSemi)
//	p.Chunk.WriteOpCode(OpPrint, token.Line)
//}

func (p *Parser) Match(tokenType TokenType) bool {
	if !p.HasMore() {
		return false
	}
	if p.Tokens[p.Index].Type == tokenType {
		p.Index++
		return true
	}
	return false
}

func (p *Parser) HasMore() bool {
	return p.Index < len(p.Tokens)
}

// Expression
// Binary
func (p *Parser) Expression() {
	p.Binary(len(binaryTokenOrder) - 1)
}

var (
	binaryTokenOrder = [][]TokenType{ // 优先级高的在前
		{TokenMul, TokenDiv},
		{TokenAdd, TokenSub},
		{TokenGT, TokenLT, TokenGE, TokenLE},
		{TokenEQ, TokenNE},
		{TokenAnd}, // 暂时没有考虑短路
		{TokenOr},
	}
)

// Binary
// OR
// AND
// != ==
// > < >= <=
// + -
// * /
// Unary
func (p *Parser) Binary(order int) {
	if order > 0 { // 先读取一项
		p.Binary(order - 1)
	} else {
		p.Unary()
	}

	for { // 循环处理
		token := p.Peek()
		if !TypeMatch(token, binaryTokenOrder[order]...) { // 处理结束
			return
		}
		p.Read()       // 消耗了
		if order > 0 { // 再读取一项
			p.Binary(order - 1)
		} else {
			p.Unary()
		}
		switch token.Type {
		case TokenAdd:
			p.Chunk.WriteOpCode(OpAdd, token.Line)
		case TokenSub:
			p.Chunk.WriteOpCode(OpSubtract, token.Line)
		case TokenMul:
			p.Chunk.WriteOpCode(OpMultiply, token.Line)
		case TokenDiv:
			p.Chunk.WriteOpCode(OpDivide, token.Line)
		case TokenGT:
			p.Chunk.WriteOpCode(OpGT, token.Line)
		case TokenGE:
			p.Chunk.WriteOpCode(OpLT, token.Line)
			p.Chunk.WriteOpCode(OpNot, token.Line)
		case TokenLT:
			p.Chunk.WriteOpCode(OpLT, token.Line)
		case TokenLE:
			p.Chunk.WriteOpCode(OpGT, token.Line)
			p.Chunk.WriteOpCode(OpNot, token.Line)
		case TokenEQ:
			p.Chunk.WriteOpCode(OpEQ, token.Line)
		case TokenNE:
			p.Chunk.WriteOpCode(OpEQ, token.Line)
			p.Chunk.WriteOpCode(OpNot, token.Line)
		default:
			panic(fmt.Sprintf("unknown token %v line %v", token.Type, token.Line))
		}
	}
}

func TypeMatch(token *Token, types ...TokenType) bool {
	for _, tokenType := range types {
		if tokenType == token.Type {
			return true
		}
	}
	return false
}

// Primary
// 12 33
// Group
func (p *Parser) Primary() { // 单个基础元素
	if p.Match(TokenLParen) {
		p.Expression()
		p.Must(TokenRParen)
		return
	}
	token := p.Read()
	if token.Type == TokenNum { // 数字
		val, err := strconv.ParseFloat(token.Value, 64)
		HandleErr(err)
		p.MakeConstant(NewNum(val), token.Line)
		return
	}
	if token.Type == TokenTrue || token.Type == TokenFalse { // true or false
		p.MakeConstant(NewBool(token.Type == TokenTrue), token.Line)
		return
	}
	if token.Type == TokenStr {
		p.MakeConstant(NewStr(token.Value), token.Line)
		return
	} // test.age()   test.age=""
	if token.Type == TokenId || token.Type == TokenThis || token.Type == TokenSuper { // 使用变量
		// 先把当前变量放到栈顶
		if token.Type == TokenThis {
			p.Chunk.WriteOpCode(OpThis, token.Line)
		} else if token.Type == TokenSuper {
			p.Chunk.WriteOpCode(OpSuper, token.Line)
		} else {
			index := p.ParseLocal(token.Value)
			if index >= 0 { // 局部变量
				p.Chunk.WriteOpCode(OpLGet, token.Line)
				p.Chunk.WriteIndex(uint64(index)) // 直接写入栈下标
			} else { // 全局变量
				p.MakeGGet(NewStr(token.Value), token.Line)
			}
		}
		// 再不停的对栈顶元素变换
		for {
			if p.Match(TokenAssign) { // 赋值操作，赋值会导致直接结束
				p.Expression()
				p.Chunk.WriteOpCode(OpSet, token.Line)
				return
			}
			if p.Match(TokenLParen) {
				p.CallFunc()
				continue
			}
			if p.Match(TokenDot) {
				p.FieldGet()
				continue
			}
			return
		}
	}
	panic(fmt.Sprintf("unknown token %v line %v", token.Type, token.Line))
}

func (p *Parser) FieldGet() {
	// 写入操作与属性
	field := p.Must(TokenId) // token.field
	p.Chunk.WriteOpCode(OpFGet, field.Line)
	p.Chunk.WriteConstant(NewStr(field.Value))
}

func (p *Parser) CallFunc() {
	count := uint8(0)
	if !p.Match(TokenRParen) {
		p.Expression() // 所有入参放入参数列表 局部变量表
		count++
		for p.Match(TokenComma) {
			p.Expression() // 所有入参放入参数列表 局部变量表
			count++
		}
		p.Must(TokenRParen)
	} // 函数参数数量
	p.MakeCall(count)
}

func (p *Parser) VarGetSet(token *Token) {
	index := p.ParseLocal(token.Value)
	if index >= 0 { // 局部变量
		if p.Match(TokenAssign) {
			p.Expression()
			p.Chunk.WriteOpCode(OpLSet, token.Line)
			p.Chunk.WriteIndex(uint64(index)) // 直接写入栈下标
		} else {
			p.Chunk.WriteOpCode(OpLGet, token.Line)
			p.Chunk.WriteIndex(uint64(index)) // 直接写入栈下标
		}
	} else { // 全局变量
		if p.Match(TokenAssign) {
			p.Expression()
			p.MakeGSet(NewStr(token.Value), token.Line)
		} else {
			p.MakeGGet(NewStr(token.Value), token.Line)
		}
	}
}

// Unary
// - !
// Primary
func (p *Parser) Unary() {
	token := p.Peek()
	if !TypeMatch(token, TokenSub, TokenNot) {
		p.Primary()
		return
	}
	p.Read()
	p.Primary()
	switch token.Type {
	case TokenSub:
		p.Chunk.WriteOpCode(OpNegate, token.Line)
	case TokenNot:
		p.Chunk.WriteOpCode(OpNot, token.Line)
	default:
		panic(fmt.Sprintf("unknown token %v line %v", token.Type, token.Line))
	}
}

func (p *Parser) Read() *Token {
	p.Index++
	return p.Tokens[p.Index-1]
}

func (p *Parser) Peek() *Token {
	return p.Tokens[p.Index]
}

func (p *Parser) MakeConstant(val *Value, line int) {
	p.Chunk.WriteOpCode(OpConstant, line)
	p.Chunk.WriteConstant(val)
}

func (p *Parser) MakeDeclaration(val *Value, line int) {
	p.Chunk.WriteOpCode(OpGDefine, line)
	p.Chunk.WriteConstant(val)
}

func (p *Parser) MakeGGet(val *Value, line int) {
	p.Chunk.WriteOpCode(OpGGet, line)
	p.Chunk.WriteConstant(val)
}

func (p *Parser) MakeGSet(val *Value, line int) {
	p.Chunk.WriteOpCode(OpGSet, line)
	p.Chunk.WriteConstant(val)
}

func (p *Parser) MakeCall(count uint8) {
	p.Chunk.WriteOpCode(OpCall, 0)
	p.Chunk.WriteU8(count)
}

func (p *Parser) ParseLocal(name string) int {
	for i := len(p.Locals) - 1; i >= 0; i-- { // 按层次来，自然支持覆盖
		if p.Locals[i].Name == name {
			return i // 既是局部变量下标也是栈下标
		}
	}
	return -1
}

func (p *Parser) Must(tokenTypes ...TokenType) *Token {
	token := p.Read()
	for _, tokenType := range tokenTypes {
		if tokenType == token.Type {
			return token
		}
	}
	panic(fmt.Sprintf("line %d err need %v but has %v", token.Line, tokenTypes, token.Type))
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{Tokens: tokens, Index: 0, Chunk: NewChunk(), Depth: 0, Locals: make([]*Local, 0)}
}

func Compile(source string) *Chunk {
	scanner := NewScanner(source)
	tokens := scanner.Scan()
	parser := NewParser(tokens)
	return parser.Parse()
}
