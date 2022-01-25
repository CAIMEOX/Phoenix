# Fast Builder Lambda
> Fast BuilderLambda is a special version of Fast Builder that based on gopher tunnel.
> This edition is much more powerful and professional than Project Alpha.

## Installation / Building

### Linux

- First of all clone the source code

  ```shell
  git clone https://github.com/CAIMEOX/Phoenix.git
  ```

- cd into the project library

  ```shell
  cd Phoenix
  ```

- Build Fast Builder lambda

  ```shell
  go build
  ```

- Run Fast Builder Lambda

  ```shell
  ./phoenix
  ```

## Usage

### Prepare

- Create a new file named  `config.toml` in project root directory with content below.

  ```shell
  [Connection]
  # RemoteAddress is the address of the server you want to entry.
  RemoteAddress = "127.0.0.1:19132"
  [User]
  # Enable auth if you need to authenticate to the server.
  Auth = true
  # Set to the operator
  User = "CAIMEOX"
  Bot = "CAIMEO_Bot"
  [Debug]
  Enabled = false 
  [Lib]
  # Load Standard library, includes basic arithmetic functions like +, - and array operations.
  Std = true
  # Load your own Fast Builder Script.
  # Check scripts/example.scm for more information.
  Script = ["path_to_script1.scm","path_to_script2.scm"]
  ```

- Launch Fast Builder, soon the Bot will join the targeted server.

- Entry the world.

### REPL

​	Bot will start a interactive interpreter session on your chat screen. You can send text chat to interact with the Fast Builder System. The "/" is needless and you can execute it by simply send it as a chat message.

### Language reference

> There are two major reasons why use scheme :
>
> - Scheme places an emphasis on **Functional programming** and **Recursive Implementation.** 
> - Scheme is quite **expressive** and **extensible**. (which makes it much powerful than shell)

​	Fast Builder use scheme, a dialect of lisp, as the command language. Similar to lisp, a simple command or function call is done by surrounding the function with their parameters in a parenthesis.

```lisp
(function param1 param2)
; Anything following a ';' is a comment
```

​	A single function can return a single value which can be a integer, floating point, string, a Go native struct, or even a nil value. Unlike in C/C++, Go, Java, or Rust, lisp doesn't have to start from a main function. Like python, the interpreter just reads the given file or any kind of input, and evaluates them on the fly.

​	By default the interpreter includes many functions built that allows to do some basic arithmetic. 

### Language Inbuilt

​	Defining variables, assignment, loops, etc., are built into the interpreter. In fact those are keywords. There are a handful of keywords in scheme.

- var 

  - defining a new variable, if the passed name is already defined this throws an error.

  - syntax : `(var VARIABLE_NAME INIT_VALUE)`

  - examples:

    ```lisp
    (var name "CAIMEO")
    (var age 15)
    ```

- set

  - setting value to the variable. if the variable name passed is not defined, this throws an error.

  - syntax: `(set VARIABLE_NAME VALUE)`

  - examples:

    ```lisp
    ; Set default block
    (set block "tnt")
    ```

- fn

  - function declaration

- return

  - return a value

  - syntax : `(return VALUE|VARIABLE_NAME)`

  - examples:

    ```lisp
    (return "Charon is lovely!")
    (return "Lisp is awesome")
    (return 114514)
    ```

- progn

  - run a list of lisp expressions, to be discussed later.

- loop

  - `loop` is a type of loop construct.
  - This loop is like `while` loop in C.

- in

  - `in` is another kind of loop construct.
  - similar to the one in python. `for i in list: ...`

- if

  - condition construct of scheme.

- match

  - `switch...case` condition construct of scheme

- eval

  - `eval` keyword is used to evaluate a string as a scheme expression and return the evaluated value.
  - **syntax** : `(eval LISP_EXPRESSION)`
  - **example** : `(eval "(+ 9 7)")` => 16

- delete

  - `delete` is used to delete a variable from the interpreter's memory.
  - **syntax** : `(delete VARIABLE_NAME)`
  - **example** : `(delete age)`

### Basic Types

​	Like any other language, there are some inbuilt types like int, string, float etc., Defining them is very simple. `1` is a simple integer. `3.14` is a simple floating point decimal. `"simple string"` is a simple string. `true` or `false` can be used for denoting the Boolean.

```lisp
; Interpreter automatically recognizes the number and assigns the variable age with an integer value of 16.
(var age 16)
; Get the type of age
; => Eval : int
(type age)
; Similar to this all other variables can be set up.
; Arrays can be set up like following :
(var fruits ["apple" "orange" "banana" "Papaya"])
; The arrays defined can contain any types of variables.
(var any ["caimeo" 11 "torrekie" 233])
```

### Conditions

​	There are 2 conditional constructs in scheme language.

- if
- match

#### if

- syntax

  ```lisp
  (if (CONDITION)
    (SUCCESS_CLAUSE) ;; single expression
    (FAILURE_CLAUSE) ;; single clause, optional.
    )
  ```

- example

  ```lisp
  (var age 22)
  (if (> age 18)
    (println "You are an adult")
    (println "You are not an adult"))
  ; Tips: If you want to do more than a single expression in any of the clauses, throw a progn.
  ; If the else clause is not required just leave the FAILURE_CLAUSE empty or () .
  (if (> age 18)
      (progn
        (println "You are an adult")
        (set r18_access true)))
  ```

#### Return

​	In scheme `if` can also return a value based on a condition. So consider the previous example. The variable `r18_access` is set to true. Here is another way to do it.

```lisp
(var r18_access (if (> age 18) true false))
; This can be written as:
(if (> age 18)
  (var r18_access true)
  (var r18_access false))
```

#### match

​	`match` conditional is similar to `switch...case` in `C/C++`. But there are no other keywords used unlike in `C/C++` (like `case`, `break`). You can match any kind of variable like strings, floats, etc.,

- syntax

  ```lisp
  (match MATCH_VARIABLE
      CASE_1_VARIABLE (CASE_1_BODY)
      CASE_2_VARIABLE (CASE_2_BODY)
      CASE_3_VARIABLE (CASE_3_BODY)
      CASE_4_VARIABLE (CASE_4_BODY)
      _               (DEFAULT_CASE_BODY)  ;; `default:` is here denoted as `_`
      )
  ```

- example

  ```lisp
  (var number 1)
  (match number
      0    (println "zero")
      1    (println "one")           ;; In this case this will succeed
      2    (println "two")
      3    (println "three")
      4    (println "four")
      5    (println "five")
      6    (println "six")
      7    (println "seven")
      8    (println "eight")
      9    (println "nine")
      _    (println "greater than 10")
  )
  ; Similar to if, match can return values
  ; numberString = "one"
  (var number 1)
  (var numberString (match number
                      0    "zero"
                      1    "one"
                      2    "two"
                      3    "three"
                      4    "four"
                      5    "five"
                      6    "six"
                      7    "seven"
                      8    "eight"
                      9    "nine"
                      _    "greater than 10"
                  ))
  
  ```

### Loops

There are two loop constructs

- loop
- in

#### loop

​	`loop` is a simple C/C++ while loop kind of construct. It will run the loop body unless the condition is false.

- syntax

  ```lisp
  (loop CONDITION
  	(LOOP_BODY))
  ```

- example

  ```lisp
  (var numbers [1 2 3 4 5 6 7])
  (var sum 0)
  
  (var i 0)
  (var sum2 0)
  
  (loop (< i (len numbers))
        (progn
          (set sum2 (+ sum2 (array-index numbers i)))
          (set i (+ i 1))
          ))
  ```

#### in

​	`in` loop is used for traversing through an array easily. It is like python's `for i in range(5)` construct.

- syntax

  ```lisp
  (in ARRAY ARRAY_VARIABLE LOOP_BODY)
  ```

- example

  ```lisp
  (var numbers [1 2 3 4 5 6 7])
  (var sum 0)
  
  (in numbers number
      (progn
        (set sum (+ sum number))))
  ```

### Basic Fast Builder Functions

​	After learning the basic elements of the Scheme, it's time to learn about Fast Builder functions. Fast Builder will create a new Space that based on Overworld where all the operations will occur here by default. Certainly you can create your own Space that used for Linear transformation or applying other advanced space mapping.  

#### Fetcher

- get

  ```lisp
  (get) 
  ; Set the default position to the current position of you. This function always returns Nil
  ```

#### Builder

​	Fast Builder offers many simple geometry structures generator. They will return an Array contains Vectors.

- circle

  ```lisp
  ; facing, available values: x,y,z
  (circle radius inner-radius height facing)
  ```

- sphere

  ```lisp
  (sphere raidus inner-radius)
  ```

- ellipse

  ```lisp
  (ellipse length width height facing)
  ```
  
You have to apply `plot` function to generate structures in Minecraft world.
```lisp
(plot Vectors)
; Examples:
(plot (round 5 3 4 "y"))
(plot (sphere 10 9))
```