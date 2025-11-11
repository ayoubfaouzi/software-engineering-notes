## Strings

- The standard library provides a string type to complement the **string literals**.
- You can concatenate a string, a string literal, a C-style string, or a character to a string.
  - String literals are fixed sequences of characters enclosed in double quotes, like "Hello, World!".
  - They represent constant arrays of characters, terminated by a null character ('\0').
- The standard string has a `move` constructor, so returning even **long strings** by **value** is efficient.
- A string is mutable.
- A string literal is by definition a `const char∗`. To get a literal of type `std::string` use an `s` suffix. For example:
    ```cpp
    auto cat = "Cat"s; // a std::string
    auto dog = "Dog"; // a C-style string: a const char*
    To use the s suffix, you need to use the namespace std::literals::string_literals
    ```

### string Implementation

- These days, `string` is usually implemented using the **short-string optimization**.
  - That is, short string (about 15 characters) values are kept in the `string` object itself and only **longer** strings are placed on **free** store.
  - improve performance by reducing heap allocations and improving cache locality.
- To handle multiple character sets, `string` is really an **alias** for a general template `basic_string` with the character type `char`:
    ```cpp
    template<typename Char>
    class basic_string {
        // ... string of Char ...
    };
    using string = basic_string<char>;
    // A user can define strings of arbitrary character types. For example, assuming we have a Japanese character type Jchar, we can write:
    using Jstring = basic_string<Jchar>;
    ```

### String Views

- A `string_view` (basically a {pointer,length} pair) denoting a contiguous sequence (read only ⚠️) of characters.
- The characters can be stored in many possible ways, including in a `string` and in a C-style string.
- `string_view` is like a pointer or a reference in that it does not own the characters it points to.
- Lifecycle management is up to the caller.
- `string_view` is semantically a **string&** but conceptually a **value type**.
- `string_view` has the same interface as an `std::string` minus the mutable operations.
- There is no `c_str()`, so if you want a null terminated string, you need to do some work.
- Why we are trying to avoid temporary `std:string` objects.
  - They are not small ( 3 pointers)
  - They can cause memory allocation.
  - SSO mitigates this.. somewhat
  - You have to copy data around unnecessary - CoW
- When to use it instead of `std::string`:
  - Passing as a parameter to a pure function.
  - Returning from a function.
  - A reference to a part of a long-lived structure.
- Compared to `std::string`:
  - `string_views` do not extend the lifetime of temporaries.
  - Not null terminated.
- ⚠️ DO NOT:
  - storing `string_views` in a container is potentially risky. You can end up holding onto freed memory.
  - Use a `string_view` to initialize am `std::string` member.
  - Don't assign to `auto`.
  - Constructing a `string_view` from a`nullptr` is UB.
- `std:string_view` is a *borrow type*"
  - Borrow types are essentially "borrowed" references to existing objects.
  - they lack ownership
  - they are short-lived
  - they generally can do without an assignment operator
  - they generally appear only in function parameter lists
  - they generally cannot be stored in data structures or returned safely from functions (no ownership semantics)
- There are extra complexities when we want to pass a **substring**. To address this, the standard library offers `string_view`.
- Consider a simple function concatenating two strings:
    ```cpp
    string cat(string_view sv1, string_view sv2)
    {
        string res {sv1}; // initialize from sv1
        return res += sv2; // append from sv2 and return
    }

    // We can call this cat():
    string king = "Harold";
    auto s1 = cat(king,"William"); // HaroldWilliam: string and const char*
    auto s2 = cat(king,king); // HaroldHarold: string and string
    auto s3 = cat("Edward","Stephen"sv); // EdwardStephen: const char * and string_view
    auto s4 = cat("Canute"sv,king); // CanuteHarold
    auto s5 = cat({&king[0],2},"Henry"sv); // HaHenry
    auto s6 = cat({&king[0],2},{&king[2],4}); // Harold
    ```
- This `cat()` has three advantages over the `compose()` that takes const string& arguments:
  - It can be used for character sequences managed in many different ways.
  - We can easily pass a substring.
  - We don’t have to create a string to pass a C-style string argument

## Regular expressions

- let us define and print a pattern:
    ```cpp
    regex pat {R"(\w{2}\s∗\d{5}(−\d{4})?)"}; // U.S. postal code pattern: XXddddd-dddd and variant
    ```
- To express the pattern, I use a **raw string literal** : `R"( pattern )"`.
  -  This allows **backslashes** and **quotes** to be used directly in the string.
  -  Raw strings are particularly suitable for regular expressions because they tend to contain a lot of backslashes.
  -  Had I used a conventional string, the pattern definition would have been: `regex pat {"\\w{2}\\s∗\\d{5}(−\\d{4})?"};`
- In <regex>, the standard library provides support for regular expressions:
  - `regex_match():` Match a regular expression against a string (of known size)
  - `regex_search()`: Search for a string that matches a regular expression in an (arbitrarily long) stream of data
  - `regex_replace()`: Search for strings that match a regular expression in an (arbitrarily long) stream of data and replace them.
  - `regex_iterator`: Iterate over matches and submatches
  - `regex_token_iterator`: Iterate over non-matches.

### Searching

- Example how to use `regex_search`:
  ```cpp
  regex pat {R"(\w{2}\s∗\d{5}(−\d{4})?)"}; // U.S. postal code pattern
  int lineno = 0;
  for (string line; getline(in,line); ) {
    ++lineno;
    smatch matches; // matched strings go here
    if (regex_search(line, matches, pat)) {
      cout << lineno << ": " << matches[0] << '\n'; // the complete match
      if (1<matches.size() && matches[1].matched) // if there is a sub-pattern and if it is matched
        cout << "\t: " << matches[1] << '\n'; // submatch
    }
  }
  ```

### Regular Expression Notation

- The regex library can recognize several variants of the notation for regular expressions. Here, I use the default notation, a variant of the ECMA standard used for ECMAScript (more commonly known as JavaScript).
  | Regular Expression Special Characters   |                                             |
  | --------------------------------------- | ------------------------------------------- |
  | . Any single character (a ‘‘wildcard’’) | \ Next character has a special meaning      |
  | [ Begin character class                 | ∗ Zero or more (suffix operation)           |
  | ] End character class                   | + One or more (suffix operation)            |
  | { Begin count                           | ? Optional (zero or one) (suffix operation) |
  | } End count                             | Alternative (or)                            |
  | ( Begin grouping                        | ˆ Start of line; negation                   |
  | ) End grouping                          | $ End of line                               |
- A pattern can be optional or repeated (the default is **exactly once**) by adding a suffix:
  | Repetition                              |
  | --------------------------------------- |
  | { n } Exactly n times                   |
  | { n, } n or more times                  |
  | {n,m} At least n and at most m times    |
  | ∗ Zero or more, that is, {0,}           |
  | + One or more, that is, {1,}            |
  | ? Optional (zero or one), that is {0,1} |
- A suffix `?` after any of the repetition notations **(?, ∗, +, and { })** makes the pattern matcher **lazy** or **non-greedy**.
- That is, when looking for a pattern, it will look for the **shortest** match rather than the longest. By **default**, the pattern matcher always looks for the **longest** match; this is known as the **Max Munch rule**:
- The pattern `(ab)+ `matches all of `ababab`. However, `(ab)+?` matches only the first `ab`.
- Several character classes are supported by shorthand notation:
  | Character Class | Abbreviations                                                |
  | --------------- | ------------------------------------------------------------ |
  | \d              | A decimal digit [[:digit:]]                                  |
  | \s              | A space (space, tab, etc.) [[:space:]]                       |
  | \w              | A letter (a-z) or digit (0-9) or underscore (_) [_[:alnum:]] |
  | \D              | Not \d [ˆ[:digit:]]                                          |
  | \S              | Not \s [ˆ[:space:]]                                          |
  | \W              | Not \w [ˆ_[:alnum:]]                                         |
- A **group** (a subpattern) potentially to be represented by a **sub_match** is delimited by parentheses.
- If you need parentheses that should not define a subpattern, use `(?:` rather than plain `(`. For example:
  - `(\s|:|,)∗(\d∗)` // optional spaces, colons, and/or commas followed by an optional number
- Assuming that we were not interested in the characters before the number (presumably separators), we could write:
  - `(?:\s|:|,)∗(\d∗)` // optional spaces, colons, and/or commas followed by an optional number
- This would save the regular expression engine from having to store the first characters: the `(?:` variant has only one subpattern.
  | Regular Expression | Grouping Examples                                                                      |
  | ------------------ | -------------------------------------------------------------------------------------- |
  | \d∗\s\w+           | No groups (subpatterns)                                                                |
  | (\d∗)\s(\w+)       | Two groups                                                                             |
  | (\d∗)(\s(\w+))+    | Two groups (groups do not nest)                                                        |
  | (\s∗\w∗)+          | One group; one or more subpatterns; only the last subpattern is saved as a `sub_match` |
  | <(.∗?)>(.∗?)</\1>  | Three groups; the \1 means "same as group 1"                                           |


### Iterators

- We can define a `regex_iterator` for iterating over a sequence of characters finding matches for a pattern
- For example, we can use a `sregex_iterator` (a regex_iterator<string>) to output all whitespace separated words in a string:
  ```cpp
  void test()
  {
    string input = "aa as; asd ++eˆasdf asdfg";
    regex pat {R"(\s+(\w+))"};
    for (sregex_iterator p(input.begin(),input.end(),pat); p!=sregex_iterator{}; ++p)
      cout << (∗p)[1] << '\n';
  }

  // This outputs:
  as
  asd
  asdfg
  ```
