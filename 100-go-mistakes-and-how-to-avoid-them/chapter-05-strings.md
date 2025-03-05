# Chapter 5: Strings

## #36: Not understanding the concept of a rune

- 🎗️ In Go, a `rune` is a **Unicode code point**.
- `UTF-8` encodes characters into 1 to 4 bytes, hence, up to 32 bits. This is why in Go, a `rune` is an **alias** of `int32`: `type rune = int32`.
- 📌 In Go, a source code is encoded in `UTF-8`. So, all string literals are encoded into a sequence of bytes using `UTF-8`. However, a string is a **sequence of arbitrary bytes**; it’s not necessarily based on UTF-8.
- A character isn’t always encoded into a **single byte**:
    ```go
    s := "汉"
    fmt.Println(len(s)) // 3 - len built-in function applied on a string doesn’t return the number of characters; it returns the number of bytes.
    ```
- Conversely, we can create a string from a list of bytes. We mentioned that the `汉` character was encoded using three bytes, `0xE6`,` 0xB1`, and `0x89`:
    ```go
    s := string([]byte{0xE6, 0xB1, 0x89})
    fmt.Printf("%s\n", s // 汉
    ```

## #37: Inaccurate string iteration

- Let’s look at a concrete example. Here, we want to print the different `runes` in a string and their corresponding positions:
    ```go
    s := "hêllo"
    for i := range s {
        fmt.Printf("position %d: %c\n", i, s[i])
    }
    fmt.Printf("len=%d\n", len(s))
    ```
- We have to recognize that in this example, we don’t iterate over each `rune`; instead, we iterate over each **starting index** of a `rune`.
- Printing `s[i]` doesn’t print the *ith* `rune`; it prints the `UTF-8` representation of the byte at index `i`. To fix this, we have to use the value element of the range operator:
    ```go
    s := "hêllo"
    for i, r := range s {
        fmt.Printf("position %d: %c\n", i, r)
    }
    ```
- The other approach is to convert the string into a slice of `runes` and iterate over it:
    ```go
    s := "hêllo"
    runes := []rune(s)
    for i, r := range runes {
        fmt.Printf("position %d: %c\n", i, r)
    }
    ```

## #38: Misusing trim functions

- One common mistake made by Go developers when using the `strings` package is to **mix** `TrimRight` and `TrimSuffix`.
- `TrimRight` iterates backward over each `rune`. If a rune is part of the provided set, the function removes it. If not, the function stops its iteration and returns the remaining string.
- On the other hand, `TrimSuffix` returns a string without a provided trailing suffix.
  - Also, removing the trailing suffix **isn’t a repeating** operation, so `TrimSuffix("123xoxo", "xo")` returns `123xo`.
  - The principle is the same for the left-hand side of a string with `TrimLeft` and `TrimPrefix`.


## #39: Under-optimized string concatenation

- Concatenating strings using `+=` does not perform well when we need to concatenate many strings. 🎯 Don't forget one of the core characteristics of a string: its **immutability**. Therefore, each iteration doesn’t update the string; it reallocates a new string in memory, which significantly impacts performance.
- Solution is to use `strings.Builder`. Using this struct, we can also append:
    - A byte slice using `Write`.
    - A single byte using `WriteByte`.
    - A single rune using `WriteRune`.
- **Internally**, `strings.Builder` holds a **byte slice**. Each call to `WriteString` results in a call to `append` on this slice.
- There are two impacts:
  - First, this struct shouldn’t be used **concurrently**, as the calls to `append` would lead to **race conditions**.
  - The second impact is something that we saw in mistake #21, “Inefficient slice initialization”: if the future length of a slice is already known, we should **preallocate** it. For that purpose, `strings.Builder` exposes a method `Grow(n int)` to guarantee space for another `n` bytes.
```go
func concat(values []string) string {
    total := 0
    for i := 0; i < len(values); i++ {
        total += len(values[i])
    }
    sb := strings.Builder{}
    sb.Grow(total)
    for _, value := range values {
        _, _ = sb.WriteString(value)
    }
    return sb.String()
}
```
- 👍 `strings.Builder` is the recommended solution to concatenate a list of strings. Usually, this solution should be used within a **loop**.

## #40: Useless string conversions

- When choosing to work with a `string` or a `[]byte`, most programmers tend to favor strings for convenience. But most I/O is actually done with `[]byte`.
- There is a price to pay when converting a `[]byte` into a `string` and then converting a `string` into a `[]byte`. Memory-wise, each of these conversions requires an extra **allocation**. Indeed, even though a string is backed by a `[]byte`, converting a `[]byte` into a `string` requires a **copy** of the byte slice. It means a new memory allocation and a copy of all the bytes.
- Indeed, all the **exported functions** of the `strings` package also have alternatives in the `bytes` package: `Split`, `Count`, `Contains`, `Index`, and so on. Hence, whether we’re doing I/O or not, we should first check whether we could implement a whole workflow using bytes instead of strings and avoid the price of additional conversions.

## #41: Substrings and memory leaks

- To extract a subset of a string, we can use the following syntax:
    ```go
    s1 := "Hello, World!"
    s2 := s1[:5] // Hello
    ```
- `s2` is constructed as a substring of `s1`. This example creates a string from the **first five bytes**, not the **first five runes**. Hence, we shouldn’t use this syntax in the case of runes encoded with multiple bytes. Instead, we should convert the input string into a `[]rune` type first:
    ```go
    s1 := "Hêllo, World!"
    s2 := string([]rune(s1)[:5]) // Hêllo
    ```
- When doing a substring operation, the Go specification doesn’t specify whether the resulting string and the one involved in the substring operation should share the
same data. However, the standard Go compiler does let them **share the same backing array**, which is probably the best solution **memory-wise** and **performance-wise** as it prevents a new allocation and a copy.
- We mentioned that log messages can be quite heavy. `log[:36] `will create a new string referencing the same backing array. Therefore, each uuid string that we store in
memory will contain not just 36 bytes but the number of bytes in the initial log string: potentially, thousands of bytes.
- How can we fix this? By making a **deep copy** of the substring so that the internal byte slice of uuid references a new backing array of only 36 bytes:
    ```go
    func (s store) handleLog(log string) error {
        if len(log) < 36 {
            return errors.New("log is not correctly formatted")
        }
        uuid := string([]byte(log[:36])) // The copy is performed by converting the substring into a []byte first and then into a string again.
        s.store(uuid)
        // Do something
    }
    ```
- As of Go 1.18, the standard library also includes a solution with `strings.Clone` that returns a fresh copy of a string: `uuid := strings.Clone(log[:36])`.
