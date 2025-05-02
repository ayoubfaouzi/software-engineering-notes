# User Defined Types

## Classes

- The `public` and `private` parts of a class declaration can appear in any order, but conventionally we place the public declarations first and the `private` declarations later.
    ```cpp
    class Vector:
        public:
            Vector(int s) : elem{new double[s]}, sz{s} {} // initializes the Vector members using a member initializer list
            double& operator[](int i) {
                return elem[i];
            }
            int size () { return sz };


        private:
            double *elem;
            int sz;
    ```

- What is the difference between a `struct` and a `class` ❓
  - From the programmer's point-of-view, there is a very minor difference. Members of a `struct` have **public** visibility by **default**, whereas members of a class have **private** visibility by **default**.
  - `struct` **inherits publicly** by default while `class` **inherits privately** by default.
  - `struct` is typically used for plain old data (POD) types (simple data containers without complex logic). `class` is used for **encapsulated objects** with **behavior** (methods, private data, etc.).
  - The general rule to follow is that `structs` should be **small**, **simple** (one-level) collections of related properties, that are **immutable** once created; for anything else, use a `class`.
  - Structs were left in C++ for compatibility reasons with C.

## Enumerations

- Used to represent small sets of integer values.
- They are used to make code more readable and less error-prone than it would have been had the symbolic (and mnemonic) enumerator names not been used.
    ```cpp
    enum class Color { red, blue, green };
    Color c = 2; // initialization error: 2 is not a Color
    Color x = Color{5}; // OK, but verbose
    Color y {6}; // also OK
    int x = int(Color::red); // explicitly convert an enum value to its underlying type
    ```
- If you don’t ever want to explicitly qualify enumerator names and want enumerator values to be `ints` (without the need for an explicit conversion), you can remove the `class` from `enum` class to get a **plain** enum.
- Advantages of using **class enums** compared to **traditional C** enums:
  - The enumerators are **scoped** inside the enum (avoid namespace pollution).
  - **No implicit conversion** to int (strongly typed), must use `static_cast` for explicit conversion.
    - C enums can lead to unintended comparisons between unrelated enums.
  - Allows explicit specification of the underlying type (e.g., int, char, short).
    - Useful for **memory optimization** or **serialization**.
        ```cpp
        - enum Color { Red, Green, Blue }; // Could be `int`, `short`, etc.
        - enum class Color : char { Red, Green, Blue }; // Stored as `char`
        - enum class BigEnum : uint64_t { Value1, Value2 }; // Guaranteed 64-bit
        ```
  - Can always be forward-declared (since the underlying type defaults to int unless specified).
    - Could not be forward-declared (before C++11).
    - In C++11, can be forward-declared only if the underlying type is specified.
        ```cpp
        enum Color : int; // Forward declaration (C++11+)
        enum Color : int { Red, Green, Blue };
        ```

## Unions

- `std::variant` (C++17) eliminates the need for manual tagging:
    ```cpp
        std::variant<int, float, std::string> data;
        data = 42; // Stores int

        // Safe access
        if (std::holds_alternative<int>(data)) {
            std::cout << std::get<int>(data);
        }
    ```
 - For many uses, a **variant** is simpler and safer to use than a `union`.
