# Casting

- **"unchecked casts"** typically refer to type conversions that are performed without any runtime checks to ensure the validity of the conversion.
  - This means that the compiler does not verify whether the cast is safe, potentially leading to undefined behavior if the cast is invalid.
  - Unchecked casts are often performed using **C-style casts**, `reinterpret_cast`, or `const_cast`.
  - To avoid the pitfalls of unchecked casts, consider using:
      - `static_cast`, `dynamic_cast` or avoid Casting all together: whenever possible, design your code to minimize the need for casts, which can often be a sign of poor design or code that needs refactoring.

## const_cast

- is typically used to cast away the *constness* of objects. It is the only C++-style cast that can do this.
- It’s commonly used when interfacing with APIs that require a non-const argument, but you have a `const` object.
- `const_cast` does not change the actual object’s **constness**. If you use `const_cast` to remove the const qualifier and modify a `const` object, the behavior is **undefined**.

## dynamic_cast

- is used mainly for safely **downcasting** pointers or references within an **inheritance hierarchy**.
- It is particularly useful in the context of **polymorphism**, where you need to ensure that a pointer or reference of a **base** class type can be safely converted to a **derived** class type.
- `dynamic_cast` performs a runtime check to ensure that the cast is valid. If the cast fails (e.g., if you attempt to cast to a type that is not actually the type of the object), it returns `nullptr` for pointers or throws a `bad_cast` exception for **references**.
- It is the only cast that cannot be performed using the old-style syntax.
- It is also the only cast that may have a significant runtime cost.

## reinterpret_cast

- is used for low-level casting, such as converting one pointer type to another, or casting an integer to a pointer, and vice versa.
- It can be used to cast any pointer type to any other pointer type, even if they are unrelated.
- `reinterpret_cast` is the most dangerous and least safe of the C++ casts because it does not perform any type checks and can easily lead to undefined behavior if misused. It’s typically used in systems programming or when interfacing with hardware (i.e., non-portable code).

## static_cast

- is used for conversions between **related types**, such as between **base** and **derived** classes (when the relationship is known at compile time), or between numeric types (like `int` to `float`).
- It is also used for **explicitly** performing standard conversions like **implicit** conversions, e.g., `int` to `double`, `pointer` to `void pointer` or vice versa, non-const object to const object.
- Though it cannot cast from **const to non-const** objects. (Only const_cast can do that.)

## bit_cast

- is a function template introduced in C++20 that allows for safe bitwise casting between **objects of different types**.
- It performs a reinterpretation of the bits of the object as if they were of another type, without invoking any undefined behavior, as long as certain conditions are met.
  - The source and destination types must be of the same size. If they are not, the code will not compile, ensuring that you don't accidentally reinterpret memory of the wrong size.
- Unlike `reinterpret_cast`, which can result in **undefined** behavior if used improperly (e.g., casting between unrelated types or violating alignment requirements), `std::bit_cast` is designed to be a **safer** alternative that doesn't result in undefined behavior as long as the size requirement is met.