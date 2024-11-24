# 使用 go 语言编写的编译型语言
## 输入 case
```javascript
var a=22;
var c=33;
{
    var a=33;
    print(a*c);
}

class A{
    test(){
        print("AAA");
    }
}

class B < A{
    test(){
        super.test();
        print("BBB");
    }
    init(){
        print("init");
        print(this);
    }
}

var b=B();
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
0072 0005 OP_MULTIPLY
0073 0000 OP_CALL 1
0075 0005 OP_POP
0076 0000 OP_POP
0077 0009 OP_JUMP 110
0086 0010 OP_GGET 4 <str "print">
0095 0010 OP_CONSTANT 5 <str "AAA">
0104 0000 OP_CALL 1
0106 0010 OP_POP
0107 0009 OP_THIS
0108 0009 OP_FIX_RETURN
0109 0009 OP_END_FUNC
0110 0008 OP_CONSTANT 6 <class A>
0119 0008 OP_GDEFINE 7 <str "A">
0128 0015 OP_JUMP 174
0137 0016 OP_SUPER
0138 0016 OP_FGET 8 <str "test">
0147 0000 OP_CALL 0
0149 0016 OP_POP
0150 0017 OP_GGET 4 <str "print">
0159 0017 OP_CONSTANT 9 <str "BBB">
0168 0000 OP_CALL 1
0170 0017 OP_POP
0171 0015 OP_THIS
0172 0015 OP_FIX_RETURN
0173 0015 OP_END_FUNC
0174 0019 OP_JUMP 220
0183 0020 OP_GGET 4 <str "print">
0192 0020 OP_CONSTANT 10 <str "init">
0201 0000 OP_CALL 1
0203 0020 OP_POP
0204 0021 OP_GGET 4 <str "print">
0213 0021 OP_THIS
0214 0000 OP_CALL 1
0216 0021 OP_POP
0217 0019 OP_THIS
0218 0019 OP_FIX_RETURN
0219 0019 OP_END_FUNC
0220 0014 OP_CONSTANT 11 <class B>
0229 0014 OP_INHERIT
0230 0014 OP_GDEFINE 12 <str "B">
0239 0025 OP_GGET 12 <str "B">
0248 0000 OP_CALL 0
0250 0025 OP_GDEFINE 13 <str "b">
0259 0026 OP_GGET 13 <str "b">
0268 0026 OP_FGET 8 <str "test">
0277 0000 OP_CALL 0
0279 0026 OP_POP
0280 0028 OP_CONSTANT 0 <num 22.000000>
0289 0028 OP_GDEFINE 14 <str "num">
0298 0029 OP_GGET 14 <str "num">
0307 0029 OP_CONSTANT 2 <num 33.000000>
0316 0029 OP_GT
0317 0029 OP_FJUMP 356
0326 0030 OP_GGET 4 <str "print">
0335 0030 OP_CONSTANT 0 <num 22.000000>
0344 0000 OP_CALL 1
0346 0030 OP_POP
0347 0000 OP_JUMP 377
0356 0032 OP_GGET 4 <str "print">
0365 0032 OP_CONSTANT 2 <num 33.000000>
0374 0000 OP_CALL 1
0376 0032 OP_POP
0377 0035 OP_CONSTANT 15 <num 0.000000>
0386 0035 OP_GDEFINE 16 <str "i">
0395 0035 OP_GGET 16 <str "i">
0404 0035 OP_CONSTANT 17 <num 10.000000>
0413 0035 OP_LT
0414 0000 OP_FJUMP 501
0423 0000 OP_JUMP 471
0432 0035 OP_GGET 16 <str "i">
0441 0035 OP_GGET 16 <str "i">
0450 0035 OP_CONSTANT 18 <num 1.000000>
0459 0035 OP_ADD
0460 0035 OP_SET
0461 0000 OP_POP
0462 0000 OP_JUMP 395
0471 0036 OP_GGET 4 <str "print">
0480 0036 OP_GGET 16 <str "i">
0489 0000 OP_CALL 1
0491 0036 OP_POP
0492 0000 OP_JUMP 432
=====================Interpret=====================
sys-call <num 1089.000000>
sys-call <str "init">
sys-call <inst [class B]>
sys-call <str "AAA">
sys-call <str "BBB">
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