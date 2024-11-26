/*
@author: sk
@date: 2024/11/16
*/
package main

type OpCode byte

const (
	OpReturn OpCode = iota + 1
	OpConstant
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNot
	OpEQ // != 使用 EQ 与 NOT 组合
	OpGT
	OpLT
	//OpPrint
	OpPop
	OpGDefine
	OpGGet
	OpGSet
	OpLGet
	OpLSet
	OpFJump // 为假跳转
	OpJump  // 无条件跳转
	OpPush  // push 一些占位
	OpCall
	OpFixReturn // 修复返回值的位置
	OpEndFunc
	OpFGet
	OpFSet
	OpSet
	OpThis
	OpInherit
	OpSuper
	OpMarkBC // 标记 break 与 continue 的位置
	OpBreak
	OpContinue
)

type ValueType uint8

const (
	ValueNum ValueType = iota + 1
	ValueStr
	ValueBool
	ValueFunc
	ValueNative
	ValueClass
	ValueInst
)

type TokenType uint8

const (
	// 单字符词法
	TokenLParen TokenType = iota + 1 // (
	TokenRParen                      // )
	TokenLBrace                      // {
	TokenRBrace                      // }
	TokenComma                       // ,
	TokenDot                         // .
	TokenSub                         // -
	TokenAdd                         // +
	TokenSemi                        // ;
	TokenDiv                         // /
	TokenMul                         // *
	// 一或两字符词法
	TokenNot    // !
	TokenNE     // !=
	TokenAssign // =
	TokenEQ     // ==
	TokenGT     // >
	TokenGE     // >=
	TokenLT     // <
	TokenLE     // <=
	// 字面量
	TokenId  // abc
	TokenStr // "sdas"
	TokenNum // 22.33
	// 关键字
	TokenAnd      // and
	TokenOr       // or
	TokenTrue     // true
	TokenFalse    // false
	TokenIf       // if
	TokenElse     // else
	TokenFor      // for
	TokenWhile    // while
	TokenBreak    // break
	TokenContinue // continue
	TokenThis     // this
	TokenSuper    // super
	TokenClass    // class
	TokenFunc     // func
	TokenVar      // var
	TokenReturn   // return
	//TokenPrint  // print
	// 其他
	TokenEOF // 文件结尾
)

const (
	InitMethod = "init"
)
