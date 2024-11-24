/*
@author: sk
@date: 2024/11/16
*/
package main

import (
	"encoding/binary"
	"fmt"
)

type Value struct {
	Type   ValueType // 也可以使用 any 的方式存储
	Num    float64
	Str    string
	Bool   bool
	Func   *Func
	Native *Native // 本地方法
	Class  *Class
	Inst   *Inst
}

func NewInst(inst *Inst) *Value {
	return &Value{Inst: inst, Type: ValueInst}
}

type Inst struct {
	Class  *Class
	Fields map[string]*Value
}

func (i *Inst) String() string {
	return fmt.Sprintf("[class %s]", i.Class.Name)
}

func NewClass(class *Class) *Value {
	return &Value{Class: class, Type: ValueClass}
}

type Class struct {
	Name    string
	Super   string
	Methods map[string]*Value
}

func (c *Class) String() string {
	return fmt.Sprintf("%s", c.Name)
}

func NewNative(native *Native) *Value {
	return &Value{Native: native, Type: ValueNative}
}

type Func struct {
	Name   string
	Ip     int
	EndIp  int
	ArgCnt uint8
}

func (f *Func) String() string {
	return fmt.Sprintf("%s %d %d %d", f.Name, f.Ip, f.EndIp, f.ArgCnt)
}

type Native struct {
	Name   string
	ArgCnt uint8
	Func   func([]*Value) *Value
}

func (n *Native) String() string {
	return fmt.Sprintf("%s %d", n.Name, n.ArgCnt)
}

func NewNum(num float64) *Value {
	return &Value{Num: num, Type: ValueNum}
}

func NewStr(str string) *Value {
	return &Value{Str: str, Type: ValueStr}
}

func NewBool(bool0 bool) *Value {
	return &Value{Bool: bool0, Type: ValueBool}
}

func NewFunc(func0 *Func) *Value {
	return &Value{Func: func0, Type: ValueFunc}
}

func (v *Value) String() string {
	switch v.Type {
	case ValueNum:
		return fmt.Sprintf("<num %f>", v.Num)
	case ValueBool:
		return fmt.Sprintf("<bool %t>", v.Bool)
	case ValueStr:
		return fmt.Sprintf("<str %q>", v.Str)
	case ValueFunc:
		return fmt.Sprintf("<func %s>", v.Func.String())
	case ValueNative:
		return fmt.Sprintf("<native %s>", v.Native.String())
	case ValueClass:
		return fmt.Sprintf("<class %s>", v.Class.String())
	case ValueInst:
		return fmt.Sprintf("<inst %s>", v.Inst.String())
	default:
		return fmt.Sprintf("unkonw type %v", v.Type)
	}
}

func (v *Value) Negate() *Value {
	switch v.Type {
	case ValueNum:
		return &Value{Type: ValueNum, Num: -v.Num}
	default:
		panic(fmt.Sprintf("unknown type %v for negate", v.Type))
	}
}

func (v *Value) Add(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueNum, Num: v.Num + val.Num}
	} else if v.Type == ValueStr && val.Type == ValueStr {
		return &Value{Type: ValueStr, Str: v.Str + val.Str}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for add", v.Type, val.Type))
	}
}

func (v *Value) Subtract(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueNum, Num: v.Num - val.Num}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for sub", v.Type, val.Type))
	}
}

func (v *Value) Multiply(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueNum, Num: v.Num * val.Num}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for mul", v.Type, val.Type))
	}
}

func (v *Value) Divide(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueNum, Num: v.Num / val.Num}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for div", v.Type, val.Type))
	}
}

func (v *Value) Not() *Value {
	switch v.Type {
	case ValueBool:
		return &Value{Type: ValueBool, Bool: !v.Bool}
	default:
		panic(fmt.Sprintf("unknown type %v not", v.Type))
	}
}

func (v *Value) GT(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueBool, Bool: v.Num > val.Num}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for gt", v.Type, val.Type))
	}
}

func (v *Value) LT(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueBool, Bool: v.Num < val.Num}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for lt", v.Type, val.Type))
	}
}

func (v *Value) EQ(val *Value) *Value {
	if v.Type == ValueNum && val.Type == ValueNum {
		return &Value{Type: ValueBool, Bool: v.Num == val.Num}
	} else if v.Type == ValueBool && val.Type == ValueBool {
		return &Value{Type: ValueBool, Bool: v.Bool == val.Bool}
	} else if v.Type == ValueStr && val.Type == ValueStr {
		return &Value{Type: ValueBool, Bool: v.Str == val.Str}
	} else {
		panic(fmt.Sprintf("unsupport type %v , %v for eq", v.Type, val.Type))
	}
}

type Chunk struct {
	Data   []byte
	Line   map[int]int // 对应偏移的指令对应的行号信息
	Values []*Value
}

func NewChunk() *Chunk {
	return &Chunk{Data: make([]byte, 0), Line: make(map[int]int), Values: make([]*Value, 0)}
}

func (c *Chunk) WriteOpCode(code OpCode, line int) {
	c.Line[len(c.Data)] = line
	c.Data = append(c.Data, byte(code))
}

func (c *Chunk) AddValue(value *Value) uint64 {
	if value.Type == ValueNum || value.Type == ValueStr || value.Type == ValueBool {
		for i, temp := range c.Values { // 对常量进行复用 只对 num str bool 进行复用
			if temp.Type == value.Type && temp.Bool == value.Bool &&
				temp.Str == value.Str && temp.Num == value.Num {
				return uint64(i)
			}
		}
	}
	c.Values = append(c.Values, value)
	return uint64(len(c.Values) - 1)
}

func (c *Chunk) WriteIndex(index uint64) {
	c.Data = binary.LittleEndian.AppendUint64(c.Data, index)
}

func (c *Chunk) ReadIndex(offset int) uint64 {
	return binary.LittleEndian.Uint64(c.Data[offset:])
}

func (c *Chunk) ReadOpCode(offset int) OpCode {
	return OpCode(c.Data[offset])
}

func (c *Chunk) WriteConstant(val *Value) {
	index := c.AddValue(val)
	c.WriteIndex(index)
}

func (c *Chunk) ReadConstant(offset int) *Value {
	index := c.ReadIndex(offset)
	return c.Values[index]
}

func (c *Chunk) SetIndex(offset int, index uint64) {
	binary.LittleEndian.PutUint64(c.Data[offset:], index)
}

func (c *Chunk) WriteU8(val uint8) {
	c.Data = append(c.Data, val)
}

func (c *Chunk) ReadU8(ip int) uint8 {
	return c.Data[ip]
}
