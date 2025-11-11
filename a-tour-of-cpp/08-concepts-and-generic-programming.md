# Generic Programming

- Templates offer:
  - The ability to pass types (as well as values and templates) as arguments without **loss of information**. This implies great flexibility in what can be expressed and excellent opportunities for inlining, of which current implementations take great advantage.
  - Opportunities to weave together information from different contexts at **instantiation** time. This implies **optimization opportunities**.
  - The ability to pass values as template arguments. This implies opportunities for **compile-time computation**.
- In other words, templates provide a powerful mechanism for compile-time computation and type manipulation that can lead to **very compact** and **efficient code**.
- Template is C++’s main support for **generic programming**. Templates provide (compile-time) **parametric polymorphism**.

## Concepts:

- Consider the example below:
  ```cpp
  template<typename Seq, typename Value>
  Value sum(Seq s, Value v) {
    for (const auto& x : s)
      v+=x;
    return v;
  }
  ```
- This `sum()` requires that:
  - its first template argument is some kind of sequence of elements, and
  - its second template argument is some kind of number.
- To be more specific, `sum()` can be invoked for a pair of arguments:
  - A sequence, `Seq`, that supports` begin()` and `end()` so that the `range-for` will work.
  - An arithmetic type, `Value`, that supports `+=` so that elements of the sequence can be added.
- We call such requirements **concepts**.
- Usually, we can do better than that. Consider that `sum()` again:
  ```cpp
    template<Sequence Seq, Number Num>
    requires Arithmetic<range_value_t<Seq>,Num>
    Num sum(Seq s, Num n);
    for (const auto& x : s)
      v+=x;
    return v;
  }
  ```
- The `range_value_t` of a sequence is the type of the elements in that sequence; it comes from the standard library where it names the type of the elements of a range.
- `Arithmetic<X,Y>` is a concept specifying that we can do arithmetic with numbers of types `X` and `Y`.
- Unsurprisingly, `requires Arithmetic<range_value_t<Seq>,Num>` is called a **requirements-clause**.
- The `template<Sequence Seq>` notation is simply a **shorthand** for an explicit use of `requires Sequence<Seq>`. If I liked verbosity, I could equivalently have written:
  ```cpp
  template<typename Seq, typename Num>
  requires Sequence<Seq> && Number<Num> && Arithmetic<range_value_t<Seq>,Num>
  Num sum(Seq s, Num n);
  ```

### Concept-based Overloading

- Consider a slightly simplified standard-library function `advance()` that advances an iterator:
  ```cpp
  template<forward_iterator Iter>
  void advance(Iter p, int n) // move p n elements forward
  {
    while (n−−)
      ++p; // a forward iterator has ++, but not + or +=
  }

  template<random_access_iterator Iter>
  void advance(Iter p, int n) // move p n elements forward
  {
    p+=n; // a random-access iterator has +=
  }
  ```
- The compiler will select the template with the **strongest requirements** met by the arguments. In this case, a `list` only supplies forward iterators, but a `vector` offers random-access iterators, so we get:
  ```cpp
  void user(vector<int>::iterator vip, list<string>::iterator lsp) {
    advance(vip,10); // uses the fast advance()
    advance(lsp,10); // uses the slow advance()
  }
  ```
- The rules for concept-based **overloading** are far simpler than the rules for general overloading. Consider first a single argument for several alternative functions:
  - If the argument doesn’t match the concept, that alternative cannot be chosen.
  - If the argument matches the concept for just one alternative, that alternative is chosen.
  - If arguments from **two alternatives** match a concept and one is stricter than the other (match all the requirements of the other and more), that **alternative** is chosen.
  - If arguments from two alternatives are **equally** good matches for a concept, we have an **ambiguity**.
- For an alternative to be chosen it must be
  - a match for all of its arguments, and
  - at least an equally good match for all arguments as other alternatives, and
  - a **better match** for at least one argument.

### Valid Code

- The question of whether a set of template arguments offers what a template requires of its template parameters ultimately boils down to whether some **expressions are valid**.
- Using a `requires-expression`, we can check if a set of expressions is valid. For example, we might try to write `advance()` **without** the use of the standard-library concept `random_access_iterator`:
  ```cpp
  template<forward_iterator Iter>
  requires requires(Iter p, int i) { p[i]; p+i; } // Iter has subscripting and integer addition
  void advance(Iter p, int n) // move p n elements forward
  {
    p+=n;
  }
  ```
- A `requires−expression` is a predicate that is true if the statements in it are valid code and false if not.
- `requires-expressions` are extremely flexible and impose no programming discipline 🤷
- The use of **requires requires** in `advance()` is deliberately **inelegant** and **hackish**. Note that I ‘forgot’ to specify `+=` and the required return types for the operations. Therefore, some uses of the version of `advance()` will pass concept checking and still **not compile**.

### Definition of Concepts

- A concept is a **compile-time predicate** specifying how one or more types can be used. Consider first one of the simplest examples:
  ```cpp
  template<typename T>
  concept Equality_comparable =
    requires (T a, T b) {
      { a == b } −> Boolean; // compare Ts with ==
      { a != b } −> Boolean; // compare Ts with !=
  };
  ```
- `Equality_comparable` is the concept we use to ensure that we can compare values of a type equal and non-equal. We simply say that, given two values of the type, they must be comparable using `==` and `!=` and the result of those operations must be Boolean. For example:
  ```cpp
  static_assert(Equality_comparable<int>); // succeeds
  struct S { int a; };
  static_assert(Equality_comparable<S>); // fails because structs don’t automatically get == and !=
  ```
- The result of an `{...}` specified after a `−>` must be a **concept**. Unfortunately, there isn’t a standard-library **boolean** concept, so I defined one (§14.5). Boolean simply means a type that can be used as a condition.
- Defining `Equality_comparable` to handle nonhomogeneous comparisons is almost as easy:
  ```cpp
  template<typename T, typename T2 =T>
  concept Equality_comparable =
    requires (T a, T2 b) {
      { a == b } −> Boolean; // compare a T to a T2 with ==
      { a != b } −> Boolean; // compare a T to a T2 with !=
      { b == a } −> Boolean; // compare a T2 to a T with ==
      { b != a } −> Boolean; // compare a T2 to a T with !=
  }

  static_assert(Equality_comparable<int,double>); // succeeds
  static_assert(Equality_comparable<int>); // succeeds (T2 is defaulted to int)
  static_assert(Equality_comparable<int,string>); // fails
  ```
- The typename `T2 =T` says that if we don’t specify a second template argument, `T2` will be the same as `T`; `T` is a **default template argument**.
- This `Equally_comparable` is almost identical with the standard-library `equality_comparable`.

#### Concepts and auto

- The keyword `auto` denotes the least constrained concept for a value: it simply requires that it must be a value of some type. Taking an `auto` parameter makes a function into a **function template**.
- Given concepts, we can strengthen requirements of all such initializations by preceding `auto` by a concept. For example:
  ```cpp
  auto twice(Arithmetic auto x) { return x+x; } // just for numbers
  auto thrice(auto x) { return x+x+x; } // for anything with a +
  ```
- In addition to their use for constraining function arguments, concepts can constrain the initialization of variables:
  ```cpp
  auto ch1 = open_channel("foo"); // works with whatever open_channel() returns
  Arithmetic auto ch2 = open_channel("foo"); // error: a channel is not Arithmetic
  Channel auto ch3 = open_channel("foo"); // OK: assuming Channel is an appropriate concept and that open_channel() returns one
  ```

## Abstraction Using Templates

- Good abstractions are carefully grown from **concrete examples**. It is not a good idea to try to **abstract** by trying to prepare for every conceivable need and technique; in that direction lies **inelegance** and **code bloat**. Instead, start with one – and preferably more – concrete examples from real use and try to eliminate inessential details. Consider:
  ```cpp
  double sum(const vector<int>& v) c{
    double res = 0;
    for (auto x : v)
      res += x;
    return res;
  }
  ```
- This is obviously one of many ways to compute the sum of a sequence of numbers. Consider what makes this code less general than it needs to be:
  - Why just ints?
  - Why just vectors?
  - Why accumulate in a double?
  - Why start at 0?
  - Why add?
- Answering the first four questions by making the concrete types into template arguments, we get the simplest form of the standard-library accumulate algorithm:
  ```cpp
  template<forward_iterator Iter, Arithmetic<iter_value_t<Iter>> Val>
  Val accumulate(Iter first, Iter last, Val res) {
    for (auto p = first; p!=last; ++p)
      res += ∗p;
    return res;
  }
  ```
- Here, we have:
  - The data structure to be traversed has been abstracted into a pair of iterators representing a sequence.
  - The type of the accumulator has been made into a parameter.
  - The type of the accumulator must be arithmetic .
  - The type of the accumulator must work with the iterator’s value type (the element type of the sequence)
  - The initial value is now an input; the type of the accumulator is the type of this initial value.
- Conversely, the best way to develop a template is often to:
  - first, write a concrete version
  - then, debug, test, and measure it
  - finally, replace the concrete types with template arguments.
- Naturally, the repetition of `begin()` and `end()` is tedious, so we can simplify the user interface a bit:
  ```cpp
  template<forward_range R, Arithmetic<value_type_t<R>> Val>
  Val accumulate(const R& r, Val res = 0) {
    for (auto x : r)
      res += x;
    return res;
  }
  ```

## Variadic Templates

- Traditionally, implementing a variadic template has been to separate the first argument from the rest and then **recursively** call the variadic template for the tail of the arguments:
  ```cpp
  template<typename T>
  concept Printable = requires(T t) { std::cout << t; } // just one operation!
    void print() {
    // what we do for no arguments: nothing
    }

  template<Printable T, Printable... Tail>
  void print(T head, Tail... tail) {
    cout << head << ' '; // first, what we do for the head
    print(tail...); // then, what we do for the tail
  }
  ```
- A parameter declared with a ... is called a **parameter pack**.
- A call of `print()` separates the arguments into a **head** (the first) and a **tail** (the rest).
- The head is printed and then `print()` is called for the **tail**. Eventually, of course, tail will become empty, so we need the no-argument version of `print()` to deal with that. If we don’t want to allow the zero-argument case, we can eliminate that `print()` using a **compile-time if**:
  ```cpp
  template<Printable T, Printable... Tail>
  void print(T head, Tail... tail) {
    cout << head << ' ';
    if constexpr(sizeof...(tail)> 0)
      print(tail...);
  }
  ```
- To simplify the implementation of simple variadic templates, C++ offers a limited form of iteration over elements of a parameter pack. For example:
  ```cpp
  template<Number... T>
  int sum(T... v) { // The body of sum uses a (right) fold expression, alternatively we can do a left fold (0 + ... + v);
    return (v + ... + 0); // add all elements of v starting with 0
  }
  // This sum() can take any number of arguments of any types:
  int x = sum(1, 2, 3, 4, 5); // x becomes 15
  int y = sum('a', 2.4, x); // y becomes 114 (2.4 is truncated and the value of ’a’ is 97)
  ```

### Forwarding Arguments

- Passing arguments unchanged through an interface is an important use of variadic templates.
- Consider a notion of a network input channel for which the actual method of moving values is a parameter. Different transport mechanisms have different sets of constructor parameters:
  ```cpp
  template<concepts::InputTransport Transport>
  class InputChannel {
    public:
      // ...
      InputChannel(Transport::Args&&... transportArgs) : _transport(std::forward<TransportArgs>(transportArgs)...) {}
      // ...
      Transport _transport;
  }
  ```
- The standard-library function `forward()` is used to move the arguments unchanged from the `InputChannel` constructor to the `Transport` constructor.
- The point here is that the writer of `InputChannel` can construct an object of type `Transport` without having to know what arguments are required to construct a particular `Transport`. The implementer of `InputChannel` needs only to know the common user interface for all `Transport` objects.
