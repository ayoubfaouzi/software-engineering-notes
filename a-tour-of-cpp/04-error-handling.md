# Error Handling

## Exceptions

- use the standard-library exception `length_error` to report a negative number of elements because some standard-library operations use that exception to report problems of this kind.

```cpp
void f(Vector& v) {
    try { // out_of_range exceptions thrown in this block are handled by the handler defined below
        compute1(v); // might try to access beyond the end of v
    }
    catch (const out_of_range& err) {   // caught the exception by reference to avoid copying
        cerr << err.what() << '\n';     // print the error message put into it at the throw-point.
    }
}
```

- Often, a function has no way of completing its assigned task after an exception is thrown. Then, handling an exception means doing some m**inimal local cleanup** and **rethrowing** the
exception. For example:
```cpp
void test(int n) {
    try {
        Vector v(n);
    }
    catch (std::length_error&) { // do something and rethrow
        cerr << "test failed: length error\n";
        throw; // rethrow
    }
    catch (std::bad_alloc&) { // ouch! this program is not designed to handle memory exhaustion
        std::terminate(); // terminate the program
    }
}
```
- There are languages where exceptions are designed simply to provide an alternate mechanism for returning values. C++ is not such a language: exceptions are designed to be used to report failure to complete a given task.
- Compilers are optimized to make returning a value much cheaper than throwing the same value as an exception.

## Return error V Throw an Exception V Terminate the Program.

- The logic behind this in my opinion, is that with return codes, the onus is up to the caller of your code to check return codes to handle errors. This then, allows users of your code to potentially get into a bad situation. Exceptions, on the other hand, must be dealt with, or they bubble up the call stack until it terminates the program.
- If you are writing C++ (and wish to adhere to C++ idioms), and want to use error handling, exceptions are the way to go.
- The guideline to follow is that exceptions are called exceptions because they are exceptional. If a condition can be reasonably expected then it should not be signaled with an exception.

- We **return an error code** when:
  - A failure is normal and expected. For example, it is quite normal for a request to open a file to fail (maybe there is no file of that name or maybe the file cannot be opened with the permissions requested).
  - An immediate caller can reasonably be expected to handle the failure.
  - An error happens in one of a set of parallel tasks and we need to know which task failed.
  - A system has so little memory that the run-time support for exceptions would crowd out essential functionality.
- We **throw an exception** when:
    - An error is so rare that a programmer is likely to forget to check for it. For example, when did you last check the return value of printf()?
    - An error cannot be handled by an immediate caller. Instead, the error has to percolate back up the call chain to an *ultimate caller.* For example, it is infeasible to have every function in an application reliably handle every allocation failure and network outage. Repeatedly checking an error-code would be tedious, expensive, and error-prone. The tests for errors and passing error-codes as return values can easily obscure the main logic of a function.
    - New kinds of errors can be added in lower-modules of an application so that higher-level modules are not written to cope with such errors. For example, when a previously single-threaded application is modified to use multiple threads or resources are placed remotely to be accessed over a network.
    - No suitable return path for errors codes is available. For example, a constructor does not have a return value for a ‘‘caller’’ to check. In particular, constructors may be invoked for several local variables or in a partially constructed complex object so that clean-up based on error codes would be quite complicated. Similarly, an operators don’t usually have an obvious return path for error codes. For example,` a∗b+c/d`.
    - The return path of a function is made more complicated or more expensive by a need to pass both a value and an error indicator back (a pair), possibly leading to the use of out-parameters, non-local error-status indicators, or other workarounds.
    - The recovery from errors depends on the results of several function calls, leading to the need to maintain local state between calls and complicated control structures.
    - The function that found the error was a callback (a function argument), so the immediate caller may not even know what function was called.
    - An error implies that some ‘‘undo action’’ is needed.
- We **terminate** when:
  - An error is of a kind from which we cannot recover. For example, for many – but not all – systems there is no reasonable way to recover from memory exhaustion.
  - The system is one where error-handling is based on restarting a thread, process, or computer whenever a non-trivial error is detected.
- *RAII* is essential for simple and efficient error-handling using exceptions. Code littered with try-blocks often simply reflects the worst aspects of error-handling strategies conceived for error codes.

## Assertions

- The standard library offers the **debug macro**, `assert()`, to assert that a condition must hold at run time.
  - If the condition of an `assert()` fails in **debug mode**, the program terminates.
  - If not in debug mode, the `assert()` is not checked.

## Static assertions

- We can perform simple checks on most properties that are known at **compile time** and report failures to meet our expectations as compiler error messages. For example:
    ```cpp
    `static_assert(4<=sizeof(int), "integers are too small"); // check integer size`.
    ```
- The `static_assert` mechanism can be used for anything that can be expressed in terms of constant expressions

## noexcept keyword

- By declaring a function, a method, or a lambda-function as `noexcept`, you specify that these does **not throw an exception** and if they throw, `std::terminate()` is called to immediately terminate the program.
  - `noexcept` specification is equivalent to the `noexcept(true)` specification.
  - `throw()` is equivalent to `noexcept(true)` but was deprecated with C++11 and will be removed with C++20.
  - In contrast, `noexcept(false)` means that the function may throw an exception. The `noexcept` specification is part of the function type but can not be used for function **overloading**.
- why:
  - semantics, if a function is specified as `noexcept`, it can be safely used in a non-throwing function.
  - an **optimization** opportunity for the compiler. `noexcept` may not call `std::unexpected` and may not unwind the stack.
- `noexcept` on functions is **hazardous** 🤷:
  - If a noexcept function calls a function that throws an exception expecting it to be caught and handled, the `noexcept` turns that into a fatal error.
  - Also, noexcept forces the writer to handle errors through some form of error codes that can be complex, error-prone, and expensive
