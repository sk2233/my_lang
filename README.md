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

func max(n1,n2){
    if(n1>n2){
        return n1;
    }else{
        return n2;
    }
}
print(max(22,33));

class A{
    test(){
        print("AAA");
    }
}
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
0087 0008 OP_JUMP 166
0096 0009 OP_LGET 0
0105 0009 OP_LGET 1
0114 0009 OP_GT
0115 0009 OP_FJUMP 143
0124 0010 OP_LGET 0
0133 0000 OP_RETURN
0134 0000 OP_JUMP 153
0143 0012 OP_LGET 1
0152 0000 OP_RETURN
0153 0008 OP_CONSTANT 6 <num 0.000000>
0162 0008 OP_FIX_RETURN
0163 0000 OP_POP
0164 0000 OP_POP
0165 0008 OP_END_FUNC
0166 0008 OP_CONSTANT 7 <func max 96 162 2>
0175 0008 OP_GDEFINE 8 <str "max">
0184 0015 OP_GGET 4 <str "print">
0193 0015 OP_GGET 8 <str "max">
0202 0015 OP_CONSTANT 0 <num 22.000000>
0211 0015 OP_CONSTANT 2 <num 33.000000>
0220 0000 OP_CALL 2
0222 0000 OP_CALL 1
0224 0015 OP_POP
0225 0018 OP_JUMP 258
0234 0019 OP_GGET 4 <str "print">
0243 0019 OP_CONSTANT 9 <str "AAA">
0252 0000 OP_CALL 1
0254 0019 OP_POP
0255 0018 OP_THIS
0256 0018 OP_FIX_RETURN
0257 0018 OP_END_FUNC
0258 0017 OP_CONSTANT 10 <class A>
0267 0017 OP_GDEFINE 11 <str "A">
0276 0023 OP_JUMP 344
0285 0024 OP_SUPER
0286 0024 OP_FGET 12 <str "test">
0295 0000 OP_CALL 0
0297 0024 OP_POP
0298 0025 OP_GGET 4 <str "print">
0307 0025 OP_CONSTANT 13 <str "BBB">
0316 0000 OP_CALL 1
0318 0025 OP_POP
0319 0026 OP_GGET 4 <str "print">
0328 0026 OP_THIS
0329 0026 OP_FGET 14 <str "name">
0338 0000 OP_CALL 1
0340 0026 OP_POP
0341 0023 OP_THIS
0342 0023 OP_FIX_RETURN
0343 0023 OP_END_FUNC
0344 0028 OP_JUMP 399
0353 0029 OP_THIS
0354 0029 OP_FGET 14 <str "name">
0363 0029 OP_LGET 0
0372 0029 OP_SET
0373 0029 OP_POP
0374 0030 OP_GGET 4 <str "print">
0383 0030 OP_CONSTANT 15 <str "init">
0392 0000 OP_CALL 1
0394 0030 OP_POP
0395 0028 OP_THIS
0396 0028 OP_FIX_RETURN
0397 0000 OP_POP
0398 0028 OP_END_FUNC
0399 0022 OP_CONSTANT 16 <class B>
0408 0022 OP_INHERIT
0409 0022 OP_GDEFINE 17 <str "B">
0418 0033 OP_GGET 17 <str "B">
0427 0033 OP_CONSTANT 18 <str "LLLL">
0436 0000 OP_CALL 1
0438 0033 OP_GDEFINE 19 <str "b">
0447 0034 OP_GGET 19 <str "b">
0456 0034 OP_FGET 12 <str "test">
0465 0000 OP_CALL 0
0467 0034 OP_POP
0468 0036 OP_CONSTANT 0 <num 22.000000>
0477 0036 OP_GDEFINE 20 <str "num">
0486 0037 OP_GGET 20 <str "num">
0495 0037 OP_CONSTANT 2 <num 33.000000>
0504 0037 OP_GT
0505 0037 OP_FJUMP 544
0514 0038 OP_GGET 4 <str "print">
0523 0038 OP_CONSTANT 0 <num 22.000000>
0532 0000 OP_CALL 1
0534 0038 OP_POP
0535 0000 OP_JUMP 565
0544 0040 OP_GGET 4 <str "print">
0553 0040 OP_CONSTANT 2 <num 33.000000>
0562 0000 OP_CALL 1
0564 0040 OP_POP
0565 0043 OP_CONSTANT 6 <num 0.000000>
0574 0043 OP_GDEFINE 21 <str "i">
0583 0043 OP_GGET 21 <str "i">
0592 0043 OP_CONSTANT 22 <num 10.000000>
0601 0043 OP_LT
0602 0000 OP_FJUMP 689
0611 0000 OP_JUMP 659
0620 0043 OP_GGET 21 <str "i">
0629 0043 OP_GGET 21 <str "i">
0638 0043 OP_CONSTANT 23 <num 1.000000>
0647 0043 OP_ADD
0648 0043 OP_SET
0649 0000 OP_POP
0650 0000 OP_JUMP 583
0659 0044 OP_GGET 4 <str "print">
0668 0044 OP_GGET 21 <str "i">
0677 0000 OP_CALL 1
0679 0044 OP_POP
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