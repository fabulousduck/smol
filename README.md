
<h1 align="center">Smol</h1>
<p align="center">Smol interpreted language.</p>

<p align="center">This is a small language i made to have something other than raw assembly to write when i want to make a piece of software for an older system. Right now, to build the system and test it. It only compiles down to Chip-8 opcodes. It has a interpreter too and a REPL if you so wish. In the future i plan to support older chips like the 6502 and the Z80 flavours. 


# Table of contents

* [table of contents](#Table-of-contents)
* [Installing](#Installing)
* [Running a file](#Running-a-file)
* [Documentation](#Documentation)
    * [General](#General)
    * [Syntax](#Syntax)
    * [Variables](#Variables)
    * [Operators](#Operators)
        * [++](#a++)
        * [--](#a\-\-)
    * [Inbuilt Functions](#Inbuilt-functions)
        * [print](#print\(v\))
    * [Functions](#def)
    * [switch](#switch)
        * [case](#case-a)
        * [default](#default)
    * [whileNot(a,b)](#whileNot\(a,b\))
    * [logical operators](#Logical-operators)
        * [eq](#eq\(a,b\))
        * [neq](#neq\(a,b\))
        * [gt](#gt\(a,b\))
        * [lt](#lt\(a,b\))
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

## Variables


In smol it is possible to define variables of a number of different types. The following types are supported as of now.

* `Uint32`
* `Uint64`
* `Bool`
* `String`

The syntax for declaring a variable is as follows:
```
<type> <name> = <value>
```

Example:
```asm
Uint32 myUint32 = 100
Uint64 myUint64 = 200
Bool myBool = True
String myString = "hello world!"
```

outputs:

```
20
```

## Operators

### `a++`

`++` is a direct operator on variables that increments the value by one.
Example:

Example:
```asm
mem a 20;
a++
print(a)
```

outputs:

```
21
```

### `a--`

`--` is a direct operator on variables that decrements the value by one.
Example:

Example:
```asm
mem a = 20;
a--
print(a)
```

outputs:

```
19
```


## Inbuilt functions

### `print(v)`

`print` is a general printing function that prints to STDOUT. This function does not get embedded into bytecode unless the target machine has a form of STDOUT

Example:
```asm
mem a = 20;

print(a)
```

outputs:

```
20
```


## `def`

Smol has support for simple functions. The can be defined like so:

```asm
functionName(1,2)

def function_name(a,b):
    print(a)
    print(b)
end
```

Functions do not support return values yet. 


## `switch`

`switch` is a basic implementation of a switch. It supports cases using either number litterals or variables. it also supports default cases. it can be used like so:

Example:

```asm
mem a = 30;
mem b = 10;

switch(b):
    case 10: #case
        print(700)
    end
    case 20:
        print(20) 
    end
    case a:
        print(a)
    end
    default:
        print(30)
    end
end

```

outputs

```asm
700
```

### `case a`

`case` defines a case within a switch.

Example: 

```asm
mem a = 30;
mem b = 10;

switch(b):
    case 10: #case
        print(700)
    end
    case a:
        print(a)
    end
END

```

outputs

```asm
700
```

### `default`

`default` can be used to declare a default case in a `switch` statement

Example

```asm
mem a = 100;
mem b = 44;

switch(b): #SWiTch
    case 10:
        print(700)
    end
    case 20:
        print(20) 
    end
    case a:
        print(A)
    end
    default:
        print(30)
    end
end

```

outputs

```asm
30
```

## `whileNot(a,b)`
`whileNot` is the while loop of smol. It will run its body untill `A == B`. So it can be seen as a simple `while a != b {}` loop.

Example:

```
mem a = 0;
mem b = 10;

whileNot(a,b):
    print(a)
    a++
end
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

### `eq(a,b)`

eq stands for "equals" and checks if `A == B`.

Example

```asm
mem a = 10;

eq(a, 10):
    print(a)
end

```

outputs:

```
10
```

### `neq(a,b)`

neq stands for "not equals" and checks if `A != B`.

Example

```asm
mem a = 11;

neq(a, 10):
    print(a)
end

```

outputs:

```
10
```

### `gt(a,b)`

gt stands for "greater than" and checks is `A < B`

Example

```asm
mem a = 10;

gt(a, 9):
    print(a)
end
```

outputs:

```
10
```

### `lt(a,b)`

lt stands for "less than" and checks if `A < B`

Example

```asm
mem a = 10;

lt(a, 11):
    print(a)
end
```

outputs:

```
10
```

## Comments

Smol has support for code comments using the `#` symbol.

Example
```asm
mem a = 10; #side comment 

#top comment
mem b 20;
```
