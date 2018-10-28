# Smol
Smol interpreted language.

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

`MEM` does not yet support variable resolution, so creating a variable with the second parameter being a reference to a variable does not work. This will throw a syntax error.

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

## Functions

Smol has support for simple functions. The can be defined like so:

```
FUNCTION_NAME[1,2];

DEF FUNCTION_NAME<PARAM_A,PARAM_B>:
    PRI A;
    PRI B;
END
```

Functions do not support return values yet. 


## `ANB`
`ANB` is the while loop of smol. It will run its body untill `A == B`. So it can be seen as a simple `while a < b {}` loop.

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