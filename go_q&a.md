# Go (Golang) Core Concepts — Clean Notes

## The Blank Identifier (`_`)

The blank identifier (`_`) is used to ignore values in Go.

It prevents compiler errors like "declared and not used" when you receive a value you don't care about.

Example:

```go
value, _ := someFunction()
```

Here, the second return value is discarded.

It's commonly used in:

- Ignoring errors (when explicitly intentional)
- Loop indices or values you don't need
- Importing packages for side effects

## Declare vs Assign vs Initialize

### Declaration

Declares that a variable exists.

```go
var x int
```

- `x` exists
- It has a zero value (`0` for `int`)
- No explicit value assigned yet

### Assignment

Gives a value to an existing variable.

```go
x = 10
```

- Variable must already exist
- Updates its value

### Short Variable Declaration (Initialization)

Declares + assigns in one step.

```go
x := 10
```

- Creates the variable
- Infers its type
- Initializes it with a value

## Allocation vs Initialization

### Allocation

Allocation means reserving memory.

```go
var p *int // pointer declared, but nil (no int allocated yet)
p = new(int)
```

- `new(int)` allocates memory for an `int`
- Returns a pointer to zero-initialized memory (`0`)

So after:

```go
p = new(int)
```

- Memory exists
- Value is `0`

### Initialization

Initialization means setting a meaningful value in allocated memory.

```go
p := new(int)
*p = 42
```

- `new(int)` allocates memory
- `*p = 42` assigns a real value into that memory

## Slices: `make`

```go
s1 := make([]int, 10)
```

- length = 10
- capacity = 10
- initialized with zeros

```go
s2 := make([]int, 0, 10)
```

- length = 0
- capacity = 10
- empty slice, but preallocated space exists

```go
s3 := make([]int, 10, 10)
```

- length = 10
- capacity = 10
- same as `s1` in practice

## Map Key Types in Go

Map keys must be comparable types.

That means the type must support:

- `==`
- `!=`

Valid key types:

- strings
- integers
- booleans
- arrays (if elements are comparable)
- structs (if all fields are comparable)
- pointers
- interfaces (if dynamic value is comparable)

Invalid key types:

- slices
- maps
- functions

Because they are not comparable.

## What is Idiomatic Go?

Idiomatic Go means writing code the "Go way":

- Simple and readable over clever
- Explicit over implicit
- Composition over inheritance
- Concurrency as a first-class concept (goroutines, channels)
- No premature optimization

A common mindset:

> "Make it work, make it right, make it fast — in that order."

## What Does "Go is Strongly Typed" Mean?

Go is strongly typed because:

- Every variable has a fixed type at compile time
- Types are enforced strictly
- You cannot freely mix types without explicit conversion

Example:

```go
var x int = 10
var y float64 = float64(x)
```

No automatic implicit conversions between incompatible types.

## `var` Keyword

`var` declares variables explicitly.

```go
var x int
var y int = 10
```

Use `var` when:

- Declaring package-level variables
- You want explicit type clarity
- You are not initializing immediately

At function scope, `:=` is often preferred.

## What is a Compiler?

A compiler translates source code into machine code.

In Go:

- Your `.go` files → compiled into a single binary
- That binary runs directly on the OS

## What is Garbage Collection?

Garbage collection is automatic memory management.

The Go runtime:

- Tracks memory usage
- Finds objects no longer in use
- Frees them automatically

This avoids manual memory freeing like in C/C++.

## Go Runtime in the Binary

When Go compiles a program, it includes the runtime inside the binary.

This includes:

- Garbage collector
- Memory allocator
- Goroutine scheduler
- Channel implementation
- Core concurrency primitives

So when you run:

```sh
./myprogram
```

You are running:

- Your code
- The Go runtime

## What is a Go Package?

A package is a directory of Go files that are compiled together and share functionality.

Every Go file starts with a package declaration:

```go
package math
```

Example:

```go
// math/add.go
package math

func Add(a, b int) int {
    return a + b
}
```

Key points:

- A package = a namespace + grouping unit
- Code inside the same package can access unexported identifiers
- Exported identifiers start with a capital letter
- Used to structure Go projects cleanly

## Can I declare and assign in Golang without using the := operator.

Yes, you can declare a variable and assign a value in one step without using := by using the var keyword:

```go
var name string = "Alice"
// However, Go also allows type inference with var, so you can omit the type:

var name = "Alice"  // Go infers the type as string
```

## Stack vs Heap

An analogy: think of a hotel.

**Stack** — you always get the room at the top of the current floor.

```
Checkout
   ⬆
Room 103
Room 102
Room 101
```

When you leave, the next guest gets Room 103. Very simple: allocation and
deallocation just move a pointer up or down.

**Heap** — you can book any available room in the entire hotel.

```
101 Occupied
102 Free
103 Occupied
104 Free
105 Occupied
...
```

The hotel has to search for an available room, mark it occupied, and later
clean it after you check out. More flexible, but more work.

Function-scoped variables are on the stack (unless they escape — e.g. a
pointer to them is returned or captured — in which case Go's escape analysis
moves them to the heap).

## Handling Long-Running Operations with Riverqueue

**Q: You have a system where users can trigger long-running operations (for
example, generating reports, processing files, or running data exports).
These operations can take several minutes, and many users may trigger them
at the same time. How would you use Riverqueue to handle this concurrently
while allowing users to track their requests?**

I'd use Riverqueue as an asynchronous job processor. When a request comes
in, I create an operation ID, enqueue a River job containing that ID, and
return the ID immediately. Multiple River workers process jobs concurrently.
The worker updates the operation status in the database, and the client
polls or subscribes using the operation ID to track progress and
completion.