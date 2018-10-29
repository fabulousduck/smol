# Smol
Smol interpreted language.

# Table of contents

* [table of contents](#Table-of-contents)
* [Installing](#Installing)
* [Running a file](#Running-a-file)
* [Documentation](#Documentation)
    * [General](#General)
    * [Syntax](#Syntax)
    * [Keywords](#Keywords)
        * [MEM](#MEM-K-V)
        * [PRI](#PRI-V)
        * [PRU](#PRU-V)
        * [INC](#INC-V)
        * [BRK](#BRK)
    * [Functions](#DEF)
    * [switch](#SWT)
        * [CAS](#CAS-A)
        * [EOS](#EOS)
    * [a not b](#ANB)
    * [logical operators](#Logical-operators)
        * [EQ](#EQ[A,B])
        * [NEQ](#NEQ[A,B])
        * [GT](#GT[A,B])
        * [LT](#LT[A,B])
    * [Math](#Math)
        * [ADD](#ADD-A-B)
        * [SUB](#SUB-A-B)
        * [MUL](#MUL-A-B)
        * [DIV](#DIV-A-B)
        * [POW](#SRQ-A-B)
    * [Comments](#Comments)




# Installing

1. install the golang language. How to do this can be found [here](https://golang.org/doc/install)
2. execute the following commands
```bash
$ go get github.com/fabulousduck/smol
$ cd $GOPATH/src/github.com/fabulousduck/smol/cmd
$ go build main.go
```

# Running a file

Once you have built the main.go file, you can execute it with any .lo file like so:
```bash
    ./main ../examples/example.lo
```

# Documentation

## General

* The file extention for smol is `.lo`

## Syntax

The syntax of smol is extremely small and aims to resemble something between lolcode and x86_64 ASM. Although it does not implement any of the x86_64 keywords. It does however keep to the style of 3 letter keywords.

All statements must end with a `;`. Not doing so will result in syntax errors.

## Keywords


### `MEM K V`
In smol `MEM` is used to declare a variable on the stack. Losp only supports whole integers as variable types. This is done on purpose to make the programmer use arithmatic to accomplish tasks like you would in assembly.


Example:
```asm
MEM A 20;

PRI A;
```

outputs:

```
20
```

### `PRI V`

`PRI` stands for "PRint Integer, so as the name suggests, the second parameter must be a number litteral or a variable name.

Example:
```asm
MEM A 20;

PRI A;
```

outputs:

```
20
```

### `PRU V`

Pru is the same as `PRI` but instead of simply printing to the screen what it is given. I will use the given variable and lookup its associated value on the unicode table. So this can be used to print a character instead of an integer.

Example

Example:
```asm
MEM A 72;

PRU A;
```

outputs:

```
H
```


### `INC V`

`INC` increments value `V` where `V` must be a variable. trying to call `INC` on a number litteral will result in a syntax error.

Example:

Example:
```asm
MEM A 20;
INC A;
PRI A;
```

outputs:

```
21
```

### `BRK`
Simply prints a `\n` character to the terminal

## `DEF`

Smol has support for simple functions. The can be defined like so:

```
FUNCTION_NAME[1,2];

DEF FUNCTION_NAME<PARAM_A,PARAM_B>:
    PRI A;
    PRI B;
END
```

Functions do not support return values yet. 


## `SWT`

`SWT` is a the switch equivelant of smol. It supports cases using either number litterals or variables. it also supports default cases. it can be used like so:

Example:

```asm
MEM A 30;
MEM B 10;

SWT[B]: #SWiTch
    CAS 10: #case
        PRI 700;
        BRK;
    END
    CAS 20:
        PRI 20;
        BRK; 
    END
    CAS A:
        PRI A;
        BRK;
    END
    EOS: #End Of Switch
        PRI 30;
        BRK;
    END
END

```

outputs

```asm
700
```

### `CAS A`

`CAS` defines a case within a switch.

Example: 

```asm
MEM A 30;
MEM B 10;

SWT[B]:
    CAS 10: #case
        PRI 700;
        BRK;
    END
    CAS A:
        PRI A;
        BRK;
    END
END

```

outputs

```asm
700
```

### `EOS`

`EOS` can be used to declare a default case in a `SWT` statement

Example

```asm
MEM A 100;
MEM B 44;

SWT[B]: #SWiTch
    CAS 10: #case
        PRI 700;
        BRK;
    END
    CAS 20:
        PRI 20;
        BRK; 
    END
    CAS A:
        PRI A;
        BRK;
    END
    EOS: #End Of Switch
        PRI 30;
        BRK;
    END
END

```

outputs

```asm
30
```

## `ANB`
`ANB` is the while loop of smol. It will run its body untill `A == B`. So it can be seen as a simple `while a != b {}` loop.

Example:

```
MEM A 0;
MEM B 10;

ANB[A,B]:
    PRI A;
    BRK;
    INC A;
END

PRU 0;
```

outputs

```
0
1
2
3
4
5
6
7
8
9
```

## Logical operators

### `EQ[A,B]`

EQ stands for "equals" and checks if `A == B`.

Example

```asm
MEM A 10;

EQ[A, 10]:
    PRI A;
END

```

outputs:

```
10
```

### `NEQ[A,B]`

NEQ stands for "not equals" and checks if `A == B`.

Example

```asm
MEM A 11;

NEQ[A, 10]:
    PRI A;
END

```

outputs:

```
10
```

### `GT[A,B]`

GT stands for "greater than" and checks is `A < B`

Example

```asm
MEM A 10;

GT[A, 9]:
    PRI A;
END

```

outputs:

```
10
```

### `LT[A,B]`

LT stands for "less than" and checks if `A < B`

Example

```asm
MEM A 10;

LT[A, 11]:
    PRI A;
END

```

outputs:

```
10
```

## Math

Smol supports the basic mathematical operators and all work the same way.
When called, like in assebly, the result of the calculation will be stored in the left hand variable given.
this means it is not possible for the left hand side to be a number litteral.

### ADD A B

Adds A and B

```asm

MEM A 10;
MEM B 20;

ADD A B;

```


outputs
```
30
```

### SUB A B

Subtracts B from A

```asm

MEM A 20;
MEM B 10;

SUB A B;

```


outputs
```
10
```

### MUL A B

multiplies A and B

```asm

MEM A 10;
MEM B 20;

MUL A B;

```

outputs
```
200
```

### DIV A B

Devides A by B

```asm

MEM A 20;
MEM B 10;

DIV A B;

```

outputs
```
2
```

## Comments

Smol has support for code comments using the `#` symbol.

Example
```asm
MEM A 10; #side comment 

#top comment
MEM B 20;
```