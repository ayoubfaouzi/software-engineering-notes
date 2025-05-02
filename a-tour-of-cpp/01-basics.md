# Basics

## Variables

- The `=` form is traditional and dates back to C, but if in doubt, use the general `{}` list form.
- If nothing else, it saves you from conversions that **lose information**:
    ```cpp
    int i1 = 7.8; // i1 becomes 7 (surprise?)
    int i2 {7.8}; // error: floating-point to integer conversion
    ```

- We use `auto` where we don’t have a specific reason to mention the type explicitly. *Specific reasons* include:
  - The definition is in a large scope where we want to make the type **clearly visible to readers** of our code.
  - The type of the initializer **isn’t obvious**.
  - We want to be **explicit** about a variable’s **range** or **precision** (e.g., `double` rather than `float`).

## Scope

A declaration introduces its name into a scope:
- **Local scope**: A name declared in a **function** or **lambda** is called a local name. Its scope extends from its point of declaration to the end of the block in which its declaration occurs. A block is delimited by a `{}` pair. Function argument names are considered local names.
- **Class scope**: A name is called a **member name** (or a class member name) if it is defined in a class, outside any function, lambda, or enum class. Its scope extends from the opening `{` of its enclosing declaration to the matching `}`.
- **Namespace scope**: A name is called a namespace member name if it is defined in a namespace outside any function, lambda, class, or enum class. Its scope extends from the point of declaration to the end of its namespace.
- A name not declared inside any other construct is called a global name and is said to be in the **global namespace**.

## Constants

- **const**: primarily to specify interfaces so that data can be passed to functions using pointers and references without fear of it being modified.
- The value of a `const` may be calculated at run time!
    ```cpp
    const double s1 = sum(v); // OK: sum(v) is evaluated at run time
    ```
- **constexpr**: meaning roughly *to be evaluated at compile time.*.
  - This is used primarily to specify constants, to allow placement of data in read-only memory (where it is unlikely to be corrupted), and for performance.
  - The value of a `constexpr` must be calculated by the compiler.
    ```cpp
    constexpr double s2 = sum(v); // error: sum(v) is not a constant expression
    ```
  - For a function to be usable in a **constant expression**, that is, in an expression that will be evaluated by the compiler, it must be defined `constexpr` or `consteval`. For example:
    ```cpp
    constexpr double square(double x) { return x∗x; }
    constexpr double max1 = 1.4∗square(17); // OK: 1.4*square(17) is a constant expression
    constexpr double max2 = 1.4∗square(var); // error: var is not a constant, so square(var) is not a constant
    const double max3 = 1.4∗square(var); // OK: may be evaluated at run time
    ```
- **consteval**: When we want a function to be used only for evaluation at compile time, we declare it `consteval` rather than `constexpr`.
  - It fails if used in a runtime context.
  - Only for functions (not variables).
  - For example:
    ```cpp
    consteval double square2(double x) { return x∗x; }

    constexpr double max1 = 1.4∗square2(17); // OK: 1.4*square(17) is a constant expression
    const double max3 = 1.4∗square2(var); // error: var is not a constant
    ```
- `const` before and after a method:
  - Using `const` before means it will return a `const` reference to T (here data_)
    ```c
    Class c;
    T& t = c.get_data()             // Not allowed.
    const T& tc = c.get_data()      // OK.
    ```
  - Using `const` after means the method will not modify any member variables of the class (unless the member is **mutable**).
    - A const member function can be invoked for **both const and non-const objects**, but a non-const member function can only be invoked for non-const objects ⚠️.
    ```c
    const T& get_data() const { return data_; }
    ```
- Pointer to const V const Pointer:
  ```c
  int i = 10;
  int* p = &i;  // pointer
  const int* p = &i; // pointer to const
  int* const p = &i // const pointer
  const int* const p = &i // const pointer to const value
  ```

## Pointers, arrays and References

- Passing a value by reference is just **syntax sugar**.
- There is nothing you can do in reference that you can't do without a pointer 😏.
- References are **cleaner** and simpler to read.
- You cannot set the value of a **references multiple times**, once you ref a variable, you can't assign it to another variable ⚠️.
- References has to be **initialized** when declared.
- References is not **re-bindable** ⚠️.
- Pointers are **addressable**, you can get their `@`, ref you can't, it's hidden.
- Pointers are **nullable**.
- using **nullptr** eliminates potential confusion between integers (such as 0 or NULL) and pointers.



A name declared in a condition is in scope on both branches of the if-statement:
```cpp
void do_something(vector<int>& v)
{
  if (auto n = v.size(); n!=0) {
    // ... we get here if n!=0 ...
  }
    // ...
}
```

## Reference lifetime extension

- Prevents temporary objects from being **destroyed prematurely** when bound to const **lvalue** references or **rvalue** references.
- It enhances safety and efficiency by avoiding dangling references and unnecessary copies.
- It does not apply to **non-const references**, references stored in **containers**, or references returned from **function**.
