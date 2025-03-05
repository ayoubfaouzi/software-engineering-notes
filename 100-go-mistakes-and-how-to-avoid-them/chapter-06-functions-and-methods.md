# Chapter 6: Functions and methods

## #42: Not knowing which type of receiver to use

- Choosing between value and pointer receivers isn’t always straightforward. Let’s discuss some of the conditions to help us choose.
- **A receiver must be a pointer:**
    - If the method needs to **mutate** the receiver. This rule is also valid if the receiver is a **slice** and a method needs to append elements:
        ```go
        type slice []int
        func (s *slice) add(element int) {
            *s = append(*s, element)
        }
        ```

    - If the method receiver contains a field that cannot be copied: for example, a type part of the `sync` package.
- **A receiver should be a pointer**:
    - If the receiver is a **large** object. Using a pointer can make the call more efficient, as doing so prevents making an **extensive copy**. When in doubt about how
large is large, benchmarking can be the solution; it’s pretty much impossible to state a specific size, because it depends on many factors.
- **A receiver must be a value**:
    - If we have to enforce a receiver’s **immutability**.
    - If the receiver is a `map`, `function`, or `channel`. Otherwise, a compilation error occurs.
- **A receiver should be a value**:
    - If the receiver is a slice that doesn’t have to be mutated.
    - If the receiver is a **small array** or **struct** that is naturally a value type without mutable fields, such as `time.Time`.
    - If the receiver is a basic type such as `int`, `float64`, or `string`.

> ⚠️ Mixing receiver types should be avoided in general but is not forbidden in 100% of cases.

## #43: Never using named result parameters

- What are the rules regarding named result parameters?
  - In most cases, using named result parameters in the context of an **interface definition** can increase **readability** without leading to any side effects. But there’s no strict rule in the context of a **method implementation**.
- In some cases, named result parameters can also increase readability: for example, if two parameters have the same type:
    ```go
    type locator interface {
        // Just by reading this code, can you guess what these two float32 results are?
        // Perhaps they are a latitude and a longitude, but in which order? using named
        // result parameters makes it clear.
        getCoordinates(address string) (lat, lng float32, err error)
    }
    ```
- In other cases, they can also be used for **convenience**:
    ```go
    func ReadFull(r io.Reader, buf []byte) (n int, err error) {
        // Because both n and err are initialized to their zero value, the implementation is shorter.
        for len(buf) > 0 && err == nil {
            var nr int
            nr, err = r.Read(buf)
            n += nr
            buf = buf[nr:]
        }
        return
    }
    ```
- 👍 Therefore, we should use named result parameters sparingly when there’s a **clear benefit**.

> 🎯 One note regarding naked returns (returns without arguments): they are considered acceptable in **short functions**; otherwise, they can harm readability because the reader must remember the outputs throughout the entire function. We should also be consistent within the scope of a function, using either only naked returns or only returns with arguments.

## #44: Unintended side effects with named result parameters

- Here’s the new implementation of the `getCoordinates` method. Can you spot what’s wrong with this code?
    ```go
    func (l loc) getCoordinates(ctx context.Context, address string) (
        lat, lng float32, err error) {
        isValid := l.validateAddress(address)
        if !isValid {
            return 0, 0, errors.New("invalid address")
        }
        if ctx.Err() != nil {
            return 0, 0, err
        }
        // Get and return coordinates
    }
    ```
- The error might not be obvious at first glance. Here, the error returned in the if `ctx.Err() != nil` scope is `err`. But we haven’t assigned any value to the `err` variable. It’s still assigned to the zero value of an error type: `nil`. Hence, this code will always return a nil error ‼️
- ⚠️ Remain cautious when using named result parameters, to avoid potential side effects.

## #45: Returning a nil receiver

- Consider the example below:
    ```go
    func (c Customer) Validate() error {
        var m *MultiError
        if c.Age < 0 {
            m = &MultiError{}
            m.Add(errors.New("age is negative"))
        }
        if c.Name == "" {
            if m == nil {
                m = &MultiError{}
            }
            m.Add(errors.New("name is nil"))
        }
        return m
    }
    ```
- Now, let’s test this implementation by running a case with a valid `Customer`:
    ```go
    customer := Customer{Age: 33, Name: "John"}
    if err := customer.Validate(); err != nil {
        log.Fatalf("customer is invalid: %v", err)
    }
    // Output:
    // > 2021/05/08 13:47:28 customer is invalid: <nil>
    ```
- 🎯 In Go, we have to know that a pointer receiver can be `nil`. In Go, a method is just **syntactic sugar** for a function whose **first parameter** is the receiver.
- `m` is initialized to the zero value of a pointer: `nil`. Then, if all the checks are valid, the argument provided to the return statement isn’t `nil` **directly** but a **nil pointer** ⚠️.
- Because a `nil` pointer is a **valid receiver**, **converting the result into an interface** won’t **yield** a `nil` value. In other words, the caller of `Validate` will always get a **non-nil** error.
- To make this point clear, let’s remember that in Go, an interface is a **dispatch wrapper**. Here, the *wrappee* is `nil` (the `MultiError` pointer), whereas the *wrapper* isn’t (the error interface).
- Therefore, regardless of the `Customer` provided, the caller of this function will always receive a non-nil error.
<p align="center"><img src="./assets/error-wrapper-is-not-nil.png" width="300px" height="auto"></p>

- Remember: An interface converted from a `nil` pointer isn’t a `nil` interface ‼️ For that reason, when we have to return an **interface**, we should return not a `nil` **pointer** but a `nil` **value** directly.

## #46: Using a filename as a function input

- When creating a new function that needs to read a file, passing a filename isn’t considered a best practice and can have negative effects, such as making **unit tests harder to write**.
- In Go, the idiomatic way is to start from the **reader’s abstraction** 👍.
- What are the benefits of this approach? First, this function **abstracts the data source**. Is it a file? An HTTP request? A socket input? It’s not important for the function. Because `*os.File` and the `Body` field of `http.Request` implement `io.Reader`, we can reuse the same function regardless of the input type.
- Another benefit is related to **testing**. We mentioned that creating one file per test case could quickly become **cumbersome** 👎. Now that `countEmptyLines` accepts an `io.Reader`, we can implement unit tests by creating an `io.Reader` from a string:
    ```go
    func TestCountEmptyLines(t *testing.T) {
        emptyLines, err := countEmptyLines(strings.NewReader(
            `foo
                bar
                baz
                `))
        // Test logic
    }
    ```
- In this test, we create an `io.Reader` using `strings.NewReader` from a **string literal directly**. Therefore, we don’t have to create one file per test case.
- Each test case can be **self-contained**, improving the test **readability** and **maintainability** as we don’t have to open another file to see the content.

## #47: Ignoring how defer arguments and receivers are evaluated

### Argument evaluation

- Consider the example below:
    ```go
    func f() error {
        var status string
        defer notify(status)
        defer incrementCounter(status)

        if err := foo(); err != nil {
            status = StatusErrorFoo
            return err
        }

        if err := bar(); err != nil {
            status = StatusErrorBar
            return err
        }
        status = StatusSuccess
        return nil
    }
    ```
- We need to understand something crucial about argument evaluation in a `defer` function:
  - The arguments are **evaluated right away**, not once the surrounding function returns ⚠️.
- How can we solve this problem if we want to keep using `defer`? There are two leading solutions.
    - The first solution is to pass a string **pointer** to the `defer` functions: `defer notify(&status)`.
    - There’s another solution: calling a **closure** as a `defer` statement:
        ```go
        func f() error {
            var status string
            defer func() {
                notify(status)
                incrementCounter(status)
            }()

            // The rest of the function is unchanged
        }
        ```
      - `status` is evaluated once the **closure is executed**, not when we call `defer`.
      - This solution also works and doesn’t require `notify` and `incrementCounter` to change their **signature**.


### Pointer and value receivers

- Consider the example below:
    ```go
    func main() {
        s := Struct{id: "foo"}
        defer s.print()
        s.id = "bar"
    }
    type Struct struct {
        id string
    }
    func (s Struct) print() {
        fmt.Println(s.id)
    }
    ```
- As with arguments, calling `defer` makes the receiver be evaluated **immediately**. Hence, `defer` delays the method’s execution with a struct that contains an `id` field equal to `foo`.
- Conversely, if the pointer is a receiver, the potential changes to the receiver **after the call** to `defer` are visible.
