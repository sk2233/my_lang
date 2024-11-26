/*
@author: sk
@date: 2024/11/16
*/
package main

import (
	"fmt"
)

func Disassemble(chunk *Chunk) {
	offset := 0
	for offset < len(chunk.Data) {
		offset = DisassembleInstruction(chunk, offset)
	}
}

func DisassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d %04d ", offset, chunk.Line[offset])

	opCode := OpCode(chunk.Data[offset])
	switch opCode {
	case OpReturn:
		return SimpleInstruction("OP_RETURN", offset)
	case OpConstant:
		return ConstantInstruction("OP_CONSTANT", chunk, offset)
	case OpGDefine:
		return ConstantInstruction("OP_GDEFINE", chunk, offset)
	case OpGGet:
		return ConstantInstruction("OP_GGET", chunk, offset)
	case OpGSet:
		return ConstantInstruction("OP_GSET", chunk, offset)
	case OpLGet:
		return IndexInstruction("OP_LGET", chunk, offset)
	case OpLSet:
		return IndexInstruction("OP_LSET", chunk, offset)
	case OpFJump:
		return IndexInstruction("OP_FJUMP", chunk, offset)
	case OpMarkBC:
		return MarkBCInstruction("OP_MARKBC", chunk, offset)
	case OpJump:
		return IndexInstruction("OP_JUMP", chunk, offset)
	case OpCall:
		return CallInstruction("OP_CALL", chunk, offset)
	case OpNegate:
		return SimpleInstruction("OP_NEGATE", offset)
	case OpNot:
		return SimpleInstruction("OP_NOT", offset)
	case OpAdd:
		return SimpleInstruction("OP_ADD", offset)
	case OpSubtract:
		return SimpleInstruction("OP_SUBTRACT", offset)
	case OpMultiply:
		return SimpleInstruction("OP_MULTIPLY", offset)
	case OpDivide:
		return SimpleInstruction("OP_DIVIDE", offset)
	case OpGT:
		return SimpleInstruction("OP_GT", offset)
	case OpLT:
		return SimpleInstruction("OP_LT", offset)
	case OpEQ:
		return SimpleInstruction("OP_EQ", offset)
	//case OpPrint:
	//	return SimpleInstruction("OP_PRINT", offset)
	case OpPop:
		return SimpleInstruction("OP_POP", offset)
	case OpPush:
		return SimpleInstruction("OP_PUSH", offset)
	case OpFixReturn:
		return SimpleInstruction("OP_FIX_RETURN", offset)
	case OpSet:
		return SimpleInstruction("OP_SET", offset)
	case OpThis:
		return SimpleInstruction("OP_THIS", offset)
	case OpSuper:
		return SimpleInstruction("OP_SUPER", offset)
	case OpContinue:
		return SimpleInstruction("OP_CONTINUE", offset)
	case OpBreak:
		return SimpleInstruction("OP_BREAK", offset)
	case OpInherit:
		return SimpleInstruction("OP_INHERIT", offset)
	case OpEndFunc:
		return SimpleInstruction("OP_END_FUNC", offset)
	case OpFGet:
		return ConstantInstruction("OP_FGET", chunk, offset)
	case OpFSet:
		return ConstantInstruction("OP_FSET", chunk, offset)
	default:
		fmt.Printf("unkonw opcode %v\n", opCode)
		return offset + 1
	}
}

func FieldInstruction(name string, chunk *Chunk, offset int) int {
	instIndex := chunk.ReadIndex(offset + 1)
	inst := chunk.Values[instIndex]
	fieldIndex := chunk.ReadIndex(offset + 1 + 8)
	field := chunk.Values[fieldIndex]
	fmt.Printf("%s %d %s %d %s\n", name, instIndex, inst, fieldIndex, field)
	return offset + 1 + 8 + 8
}

func CallInstruction(name string, chunk *Chunk, offset int) int {
	count := chunk.ReadU8(offset + 1)
	fmt.Printf("%s %d\n", name, count)
	return offset + 1 + 1
}

func ConstantInstruction(name string, chunk *Chunk, offset int) int {
	index := chunk.ReadIndex(offset + 1)
	value := chunk.Values[index]
	fmt.Printf("%s %d %s\n", name, index, value)
	return offset + 1 + 8
}

func IndexInstruction(name string, chunk *Chunk, offset int) int {
	index := chunk.ReadIndex(offset + 1)
	fmt.Printf("%s %d\n", name, index) // 没有名称只能打出索引了 例如局部变量编译是不存储名称的
	return offset + 1 + 8
}

func MarkBCInstruction(name string, chunk *Chunk, offset int) int {
	cIndex := chunk.ReadIndex(offset + 1)
	bIndex := chunk.ReadIndex(offset + 1 + 8)
	fmt.Printf("%s %d %d\n", name, cIndex, bIndex)
	return offset + 1 + 8 + 8
}

func SimpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
