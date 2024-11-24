/*
@author: sk
@date: 2024/11/16
*/
package main

import (
	"fmt"
	"time"
)

type Stack struct {
	Data  []*Value
	Count int
}

func NewStack() *Stack {
	return &Stack{Data: make([]*Value, 0), Count: 0}
}

func (s *Stack) Push(val *Value) {
	s.Data = append(s.Data, val)
	s.Count++
}

func (s *Stack) Pop() *Value {
	s.Count--
	res := s.Data[s.Count]
	s.Data = s.Data[:s.Count]
	return res
}

func (s *Stack) Peek() *Value {
	return s.Data[s.Count-1]
}

func (s *Stack) Get(index uint64) *Value {
	return s.Data[index]
}

func (s *Stack) Set(index uint64, val *Value) {
	s.Data[index] = val
}

func (s *Stack) Swap(i1 int, i2 int) {
	s.Data[i1], s.Data[i2] = s.Data[i2], s.Data[i1]
}

type CallFrame struct {
	Func       *Func  // 对应的函数
	Ip         int    // 返回地址
	SlotOffset uint64 // 开始的栈局部变量偏移 最外层偏移为0
}

type FrameStack struct {
	Frames []*CallFrame
	Count  int
}

func (s *FrameStack) Push(frame *CallFrame) {
	s.Frames = append(s.Frames, frame)
	s.Count++
}

func (s *FrameStack) Peek() *CallFrame {
	return s.Frames[s.Count-1]
}

func (s *FrameStack) Pop() *CallFrame {
	s.Count--
	res := s.Frames[s.Count]
	s.Frames = s.Frames[:s.Count]
	return res
}

func NewFrameStack() *FrameStack {
	return &FrameStack{Frames: make([]*CallFrame, 0), Count: 0}
}

func RegisterNative(global map[string]*Value) {
	global["print"] = NewNative(&Native{
		Name:   "print",
		ArgCnt: 1,
		Func: func(values []*Value) *Value {
			fmt.Println("sys-call", values[0].String())
			return NewNum(0)
		},
	})
	global["unix"] = NewNative(&Native{
		Name:   "unix",
		ArgCnt: 0,
		Func: func(values []*Value) *Value {
			return NewNum(float64(time.Now().Unix()))
		},
	})
}

func Interpret(chunk *Chunk) {
	stack := NewStack()
	frames := NewFrameStack()
	frames.Push(&CallFrame{})         // 压入默认函数
	global := make(map[string]*Value) // 全局参数
	RegisterNative(global)
	var this *Value // bad impl
	ip := 0

	for ip < len(chunk.Data) {
		opCode := chunk.ReadOpCode(ip)
		ip++

		switch opCode {
		case OpReturn: // 结束函数跳转到收尾位置
			frame := frames.Peek()
			ip = frame.Func.EndIp
		case OpConstant:
			val := chunk.ReadConstant(ip)
			ip += 8
			stack.Push(val)
		case OpNegate:
			val := stack.Pop()
			stack.Push(val.Negate())
		case OpNot:
			val := stack.Pop()
			stack.Push(val.Not())
		case OpAdd:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.Add(val1))
		case OpSubtract:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.Subtract(val1))
		case OpMultiply:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.Multiply(val1))
		case OpDivide:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.Divide(val1))
		case OpGT:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.GT(val1))
		case OpLT:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.LT(val1))
		case OpEQ:
			val1 := stack.Pop()
			val2 := stack.Pop()
			stack.Push(val2.EQ(val1))
		//case OpPrint:
		//	val := stack.Pop()
		//	fmt.Printf("%s\n", val)
		case OpPop:
			stack.Pop()
		case OpPush:
			stack.Push(&Value{}) // 只是占位，没有意义
		case OpFixReturn:
			frame := frames.Peek()
			stack.Swap(int(frame.SlotOffset), stack.Count-1)
		case OpEndFunc:
			frame := frames.Pop() // 弹出调用栈返回到之前的位置
			ip = frame.Ip
		case OpCall:
			count := chunk.ReadU8(ip)
			ip++
			args := make([]*Value, count)
			for i := uint8(0); i < count; i++ {
				args[i] = stack.Pop()
			}
			val := stack.Pop() // 要被调用到对象
			for i := uint8(0); i < count; i++ {
				stack.Push(args[i]) // 前面弹出是为了获取被调用到对象 获取后再写回去
			}
			if val.Type == ValueFunc { // 大部分校验可以放在编译时期
				if count != val.Func.ArgCnt {
					panic(fmt.Sprintf("arg cnt not match %d != %d", count, val.Func.ArgCnt))
				}
				// 调用函数 准备栈帧，跳转到函数执行
				frames.Push(&CallFrame{Func: val.Func, Ip: ip, SlotOffset: uint64(stack.Count - int(count))})
				ip = val.Func.Ip
			} else if val.Type == ValueNative {
				if count != val.Native.ArgCnt {
					panic(fmt.Sprintf("arg cnt not match %d != %d", count, val.Native.ArgCnt))
				} // 准备参数
				args = make([]*Value, count)
				for i := uint8(0); i < count; i-- {
					args[count-i-1] = stack.Pop()
				} // 调用写回结果
				res := val.Native.Func(args)
				stack.Push(res)
			} else if val.Type == ValueClass { // 构造其实例
				inst := NewInst(&Inst{Class: val.Class, Fields: make(map[string]*Value)})
				if method, ok := val.Class.Methods[InitMethod]; ok { // 初始化函数必须返回 this
					// 有初始化方法为初始化方法做准备
					if count != method.Func.ArgCnt {
						panic(fmt.Sprintf("arg cnt not match %d != %d", count, method.Func.ArgCnt))
					}
					// 调整参数与返回值位置
					args = make([]*Value, count)
					for i := uint8(0); i < count; i++ {
						args[i] = stack.Pop()
					}
					for i := uint8(0); i < count; i++ {
						stack.Push(args[i]) // 前面弹出是为了获取被调用到对象 获取后再写回去
					}
					// 调用参数进行执行
					this = inst
					frames.Push(&CallFrame{Func: method.Func, Ip: ip, SlotOffset: uint64(stack.Count - int(count))})
					ip = method.Func.Ip
				} else {
					stack.Push(inst)
				}
			} else {
				panic(fmt.Sprintf("val type %d not func , class or native", val.Type))
			}
		case OpGDefine:
			key := chunk.ReadConstant(ip)
			ip += 8
			global[key.Str] = stack.Pop()
		case OpGGet:
			key := chunk.ReadConstant(ip)
			ip += 8
			if val, ok := global[key.Str]; ok {
				stack.Push(val)
			} else {
				panic(fmt.Sprintf("global key not found: %s", key.Str))
			}
		case OpGSet:
			key := chunk.ReadConstant(ip)
			ip += 8
			if _, ok := global[key.Str]; !ok {
				panic(fmt.Sprintf("global key not found: %s", key.Str))
			}
			global[key.Str] = stack.Peek() // 这里不能 pop 赋值语句视为表达式
		case OpFGet:
			key := chunk.ReadConstant(ip)
			ip += 8
			if val := stack.Pop(); val.Type == ValueInst {
				inst := val.Inst
				if method, ok := inst.Class.Methods[key.Str]; ok {
					stack.Push(method) // 先尝试获取定义的方法
					this = val         // 记录 this 值
				} else { // 再去获取定义的属性
					if _, ok = inst.Fields[key.Str]; !ok { // 没有就去创建
						inst.Fields[key.Str] = &Value{}
					}
					stack.Push(inst.Fields[key.Str])
				}
			} else {
				panic(fmt.Sprintf("val %v not inst", val))
			}
		case OpFSet:
			key := chunk.ReadConstant(ip)
			ip += 8
			res := stack.Pop()
			if val := stack.Pop(); val.Type == ValueInst {
				val.Inst.Fields[key.Str] = res
				stack.Push(res) // 设置也是有返回值的
			} else {
				panic(fmt.Sprintf("val %v not inst", val))
			}
		case OpSet:
			src := stack.Pop() // 先弹出的是最终值
			tar := stack.Peek()
			tar.Type = src.Type
			tar.Num = src.Num
			tar.Str = src.Str
			tar.Bool = src.Bool
			tar.Func = src.Func
			tar.Native = src.Native
			tar.Class = src.Class
			tar.Inst = src.Inst
		case OpThis: // this 是关键字，实际就是一种特殊的字段获取
			stack.Push(this)
		case OpSuper:
			if this == nil {
				panic(fmt.Sprintf("invalid use super"))
			}
			super := this.Inst.Class.Super
			if len(super) == 0 {
				panic(fmt.Sprintf("class %s no supper class", this.Inst.Class.Name))
			}
			if val, ok := global[super]; ok && val.Type == ValueClass {
				inst := &Inst{ // super 只是使用父方法  属性还是子属性
					Class:  val.Class,
					Fields: this.Inst.Fields,
				}
				stack.Push(NewInst(inst))
			} else {
				panic(fmt.Sprintf("supper class %s not found or type err", super))
			}
		case OpInherit: // 继承操作
			temp := stack.Pop()
			class := temp.Class
			if val, ok := global[class.Super]; ok {
				// 继承方法，注意覆盖顺序 不采用链式查找的方式节约效率
				methods := make(map[string]*Value)
				for key, value := range val.Class.Methods {
					methods[key] = value
				}
				for key, value := range class.Methods {
					methods[key] = value
				}
				class.Methods = methods
			} else {
				panic(fmt.Sprintf("super class %s not found", class.Super))
			}
			stack.Push(temp)
		case OpLSet:
			frame := frames.Peek()       // 局部变量要受栈帧影响
			index := chunk.ReadIndex(ip) // 对局部变量的操作目标是一定存在的，否者就是全局变量了
			ip++
			stack.Set(index+frame.SlotOffset, stack.Peek())
		case OpLGet:
			frame := frames.Peek()
			index := chunk.ReadIndex(ip) // 对局部变量的操作目标是一定存在的，否者就是全局变量了
			ip += 8
			stack.Push(stack.Get(index + frame.SlotOffset))
		case OpFJump:
			index := chunk.ReadIndex(ip)
			val := stack.Pop()
			if val.Type != ValueBool {
				panic(fmt.Sprintf("jump to non-boolean value: %d", val.Type))
			}
			if val.Bool {
				ip += 8
			} else {
				ip = int(index)
			}
		case OpJump:
			index := chunk.ReadIndex(ip)
			ip = int(index)
		default:
			panic(fmt.Sprintf("unknown opcode %d", opCode))
		}
	}
	fmt.Printf("stack size %d\n", stack.Count) // 一般要求为 0
}
