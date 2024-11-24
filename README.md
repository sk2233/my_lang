# 使用 go 语言编写的编译型语言
## 支持功能
1. 字符串，数字，布尔值常量池复用
2. 反汇编支持
3. 常见表达式解析
4. 全局变量，局部变量栈上分配，变量覆盖支持
5. 常见流程结构支持
6. 函数调用，本地方法接入
7. 类与实例的字段，方法，初始化方法，this支持
8. 继承与 super支持
## 输入 case
```javascript
var a=22;
var c=33;
{
    var a=33;
    print(a*(c+2));
}



class A{
    test(){
        print("AAA");
    }
}

func max(n1,n2){
    if(n1>n2){
        return n1;
    }else{
        return n2;
    }
}
print(max(22,33));

class B < A{
    test(){
        super.test();
        print("BBB");
        print(this.name);
    }
    init(name){
        this.name=name;
        print("init");
    }
}

var b=B("LLLL");
b.test();

var num=22;
if(num>33){
    print(22);
}else{
    print(33);
}

for(var i=0;i<10;i=i+1){
    print(i);
}
```
## 输出 case
支持对生成的字节码反编译出 ip地址，行号，指令与其参数
```shell
=====================Disassemble=====================
0000 0001 OP_CONSTANT 0 <num 22.000000>
0009 0001 OP_GDEFINE 1 <str "a">
0018 0002 OP_CONSTANT 2 <num 33.000000>
0027 0002 OP_GDEFINE 3 <str "c">
0036 0004 OP_CONSTANT 2 <num 33.000000>
0045 0005 OP_GGET 4 <str "print">
0054 0005 OP_LGET 0
0063 0005 OP_GGET 3 <str "c">
0072 0005 OP_CONSTANT 5 <num 2.000000>
0081 0005 OP_ADD
0082 0005 OP_MULTIPLY
0083 0000 OP_CALL 1
0085 0005 OP_POP
0086 0000 OP_POP
0087 0011 OP_JUMP 120
0096 0012 OP_GGET 4 <str "print">
0105 0012 OP_CONSTANT 6 <str "AAA">
0114 0000 OP_CALL 1
0116 0012 OP_POP
0117 0011 OP_THIS
0118 0011 OP_FIX_RETURN
0119 0011 OP_END_FUNC
0120 0010 OP_CONSTANT 7 <class A>
0129 0010 OP_GDEFINE 8 <str "A">
0138 0016 OP_JUMP 217
0147 0017 OP_LGET 0
0156 0017 OP_LGET 1
0165 0017 OP_GT
0166 0017 OP_FJUMP 194
0175 0018 OP_LGET 0
0184 0000 OP_RETURN
0185 0000 OP_JUMP 204
0194 0020 OP_LGET 1
0203 0000 OP_RETURN
0204 0016 OP_CONSTANT 9 <num 0.000000>
0213 0016 OP_FIX_RETURN
0214 0000 OP_POP
0215 0000 OP_POP
0216 0016 OP_END_FUNC
0217 0016 OP_CONSTANT 10 <func max 147 213 2>
0226 0016 OP_GDEFINE 11 <str "max">
0235 0023 OP_GGET 4 <str "print">
0244 0023 OP_GGET 11 <str "max">
0253 0023 OP_CONSTANT 0 <num 22.000000>
0262 0023 OP_CONSTANT 2 <num 33.000000>
0271 0000 OP_CALL 2
0273 0000 OP_CALL 1
0275 0023 OP_POP
0276 0026 OP_JUMP 344
0285 0027 OP_SUPER
0286 0027 OP_FGET 12 <str "test">
0295 0000 OP_CALL 0
0297 0027 OP_POP
0298 0028 OP_GGET 4 <str "print">
0307 0028 OP_CONSTANT 13 <str "BBB">
0316 0000 OP_CALL 1
0318 0028 OP_POP
0319 0029 OP_GGET 4 <str "print">
0328 0029 OP_THIS
0329 0029 OP_FGET 14 <str "name">
0338 0000 OP_CALL 1
0340 0029 OP_POP
0341 0026 OP_THIS
0342 0026 OP_FIX_RETURN
0343 0026 OP_END_FUNC
0344 0031 OP_JUMP 399
0353 0032 OP_THIS
0354 0032 OP_FGET 14 <str "name">
0363 0032 OP_LGET 0
0372 0032 OP_SET
0373 0032 OP_POP
0374 0033 OP_GGET 4 <str "print">
0383 0033 OP_CONSTANT 15 <str "init">
0392 0000 OP_CALL 1
0394 0033 OP_POP
0395 0031 OP_THIS
0396 0031 OP_FIX_RETURN
0397 0000 OP_POP
0398 0031 OP_END_FUNC
0399 0025 OP_CONSTANT 16 <class B>
0408 0025 OP_INHERIT
0409 0025 OP_GDEFINE 17 <str "B">
0418 0037 OP_GGET 17 <str "B">
0427 0037 OP_CONSTANT 18 <str "LLLL">
0436 0000 OP_CALL 1
0438 0037 OP_GDEFINE 19 <str "b">
0447 0038 OP_GGET 19 <str "b">
0456 0038 OP_FGET 12 <str "test">
0465 0000 OP_CALL 0
0467 0038 OP_POP
0468 0040 OP_CONSTANT 0 <num 22.000000>
0477 0040 OP_GDEFINE 20 <str "num">
0486 0041 OP_GGET 20 <str "num">
0495 0041 OP_CONSTANT 2 <num 33.000000>
0504 0041 OP_GT
0505 0041 OP_FJUMP 544
0514 0042 OP_GGET 4 <str "print">
0523 0042 OP_CONSTANT 0 <num 22.000000>
0532 0000 OP_CALL 1
0534 0042 OP_POP
0535 0000 OP_JUMP 565
0544 0044 OP_GGET 4 <str "print">
0553 0044 OP_CONSTANT 2 <num 33.000000>
0562 0000 OP_CALL 1
0564 0044 OP_POP
0565 0047 OP_CONSTANT 9 <num 0.000000>
0574 0047 OP_GDEFINE 21 <str "i">
0583 0047 OP_GGET 21 <str "i">
0592 0047 OP_CONSTANT 22 <num 10.000000>
0601 0047 OP_LT
0602 0000 OP_FJUMP 689
0611 0000 OP_JUMP 659
0620 0047 OP_GGET 21 <str "i">
0629 0047 OP_GGET 21 <str "i">
0638 0047 OP_CONSTANT 23 <num 1.000000>
0647 0047 OP_ADD
0648 0047 OP_SET
0649 0000 OP_POP
0650 0000 OP_JUMP 583
0659 0048 OP_GGET 4 <str "print">
0668 0048 OP_GGET 21 <str "i">
0677 0000 OP_CALL 1
0679 0048 OP_POP
0680 0000 OP_JUMP 620
=====================Interpret=====================
sys-call <num 1155.000000>
sys-call <num 33.000000>
sys-call <str "init">
sys-call <str "AAA">
sys-call <str "BBB">
sys-call <str "LLLL">
sys-call <num 33.000000>
sys-call <num 0.000000>
sys-call <num 1.000000>
sys-call <num 2.000000>
sys-call <num 3.000000>
sys-call <num 4.000000>
sys-call <num 5.000000>
sys-call <num 6.000000>
sys-call <num 7.000000>
sys-call <num 8.000000>
sys-call <num 9.000000>
```
# 参考
https://readonly.link/books/https://raw.githubusercontent.com/GuoYaxiang/craftinginterpreters_zh/main/book.json?front-matter=contents