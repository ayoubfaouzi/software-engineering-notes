# Templates

- A template is a class or a function that we parameterize with a set of types or values.

## Parameterized Types

- Using `class` to introduce a type parameter is equivalent to using `typename`, and in older code we often see `template<class T>` as the prefix.
    ```cpp
    template<typename T>
    class Vector {
    private:
        T∗ elem; // elem points to an array of sz elements of type T
            int sz;
    public:
        explicit Vector(int s); // constructor: establish invariant, acquire resources
        ˜Vector() { delete[] elem; } // destructor: release resources
        // ... copy and move operations ...
        T& operator[](int i); // for non-const Vectors
        const T& operator[](int i) const; // for const Vectors
        int size() const { return sz; }
    }
    ```
- To support the `range-for` **loop** for our `Vector`, we must define suitable` begin()` and `end()` functions:
    ```cpp
    template<typename T>
    T∗ begin(Vector<T>& x) {
        return &x[0]; // pointer to first element or to one-past-the-last element
    }

    template<typename T>
    T∗ end(Vector<T>& x) {
        return &x[0]+x.size(); // pointer to one-past-the-last element
    }

    // Given those, we can write:
    void write2(Vector<string>& vs) {// Vector of some strings
        for (auto& s : vs)
            cout << s << '\n';
    }
    ```
- A template plus a set of template arguments is called an **instantiation** or a **specialization**. Late in the compilation process, at instantiation time, code is generated for each instantiation used in a program.
- This `template<Element T>` prefix is C++’s version of mathematic’s *for all T such that Element(T)*; that is, `Element` is a **predicate** that checks whether `T` has all the properties that a `Vector` requires.
  - Such a predicate is called a **concept**.
  - A template argument for which a concept is specified is called a **constrained argument** and a template for which an argument is constrained is called a **constrained template**.
  - It is a **compile-time error** to try to use a template with a type that does not meet its requirements !

## Value Template Arguments

- In addition to **type arguments**, a template can take **value arguments**. For example:
    ```cpp
    template<typename T, int N>
    struct Buffer {
        constexpr int size() { return N; }
        T elem[N];
    // ...
    }

    Buffer<char,1024> glob; // global buffer of characters (statically allocated)
    ```

## Template Argument Deduction

- When defining a type as an instantiation of a template we must specify its template arguments.
- Consider using the standard-library template pair: `pair<int,double> p = {1, 5.2};`
- Having to specify the template argument types can be **tedious**. Fortunately, in many contexts, we can simply let pair’s constructor **deduce** the template arguments from an **initializer**:
    ```cpp
    pair p = {1, 5.2}; // p is a pair<int,double>
    ```
- Containers provide another example:
    ```cpp
    template<typename T>
    class Vector {
    public:
        Vector(int);
        Vector(initializer_list<T>); // initializer-list constructor
        // ...
    };
    Vector v1 {1, 2, 3}; // deduce v1’s element type from the initializer element type: int
    Vector v2 = v1; // deduce v2’s element type from v1’s element type: int
    auto p = new Vector{1, 2, 3}; // p is a Vector<int>*
    Vector<int> v3(1); // here we need to be explicit about the element type (no element type is mentioned).
    ```
- This simplifies notation and can eliminate annoyances caused by mistyping redundant template argument types. However, it is not a panacea. Deduction can cause surprises. Consider:
    ```cpp
    Vector<string> vs {"Hello", "World"}; // OK: Vector<string>
    Vector vs1 {"Hello", "World"}; // OK: deduces to Vector<const char*> (Surprise?)
    Vector vs2 {"Hello"s, "World"s}; // OK: deduces to Vector<string>
    Vector vs3 {"Hello"s, "World"}; // error: the initializer list is not homogenous
    Vector<string> vs4 {"Hello"s, "World"}; // OK: the element type is explicit
    ```
- If elements of an initializer list have differing types, we cannot deduce a unique element type, so we get an **ambiguity error**.
- Sometimes, we need to resolve an ambiguity. For example, the standard-library vector has a constructor that takes a **pair of iterators** delimiting a sequence and also an initializer constructor that can take a pair of values. Consider:
    ```cpp
    template<typename T>
    class Vector {
    public:
        Vector(initializer_list<T>); // initializer-list constructor
        template<typename Iter>
        Vector(Iter b, Iter e); // [b:e) iterator-pair constructor
        struct iterator { using value_type = T; /* ... */ };
        iterator begin();
        // ...
    };

    Vector v1 {1, 2, 3, 4, 5}; // element type is int
    Vector v2(v1.begin(),v1.begin()+2); // a pair of iterators or a pair of values (of type iterator)?
    Vector v3(9,17); // error: ambiguous
    ```
- For those, we need a way of saying *a pair of values of the same type should be considered iterators.*. Adding a **deduction guide** after the declaration of `Vector` does exactly that:
    ```cpp
    template<typename Iter>
    Vector(Iter,Iter) −> Vector<typename Iter::value_type>;
    ```
- Now we have:
    ```cpp
    Vector v1 {1, 2, 3, 4, 5}; // element type is int
    Vector v2(v1.begin(),v1.begin()+2); // pair-of-iterators: element type is int
    Vector v3 {v1.begin(),v1.begin()+2}; // element type is Vector2::iterator
    ```
- The `{}` **initializer** syntax always prefers the `initializer_list` constructor (if present), so v3 is a vector of iterators: `Vector<Vector<int>::iterator>`.
- The `()` **initialization** syntax is conventional for when we don’t want an `initializer_list`.
- The effects of deduction guides are often subtle, so it is best to design class templates so that deduction guides are not needed.
- People who like acronyms refer to "*class template argument deduction*’ as CTAD.

## Parameterized Operations

### Function Templates

- We can write a function that calculates the sum of the element values of any sequence that a rangefor can traverse (e.g., a container) like this:
    ```cpp
    template<typename Sequence, typename Value>
    Value sum(const Sequence& s, Value v) {
        for (auto x : s)
            v+=x;
        return v;
    }
    ```
- A function template can be a **member function**, but **not a virtual member**. The compiler would not know all instantiations of such a template in a program, so it could not generate a **vtbl**.

### Function Objects

- One particularly useful kind of template is the function object (sometimes called a **functor**), which is used to define objects that can be **called like functions**.
- For example:
    ```cpp
    template<typename T>
    class Less_than {
        const T val; // value to compare against
    public:
        Less_than(const T& v) :val{v} { }
        bool operator()(const T& x) const { return x<val; } // call operator
    }
    ```
- The beauty of function objects is that they **carry the value** to be compared against with them. We don’t have to write a separate function for each value (and each type), and we don’t have to
introduce **nasty global variables** to hold values.
- Also, for a simple function object like `Less_than`, **inlining** is simple, so a call of `Less_than` is far more efficient than an indirect function call.
- The ability to carry data plus their efficiency makes function objects particularly useful as arguments to **algorithms**.

### Lambda Expressions

- In the previous section, we defined `Less_than` separately from its use. That can be inconvenient. Consequently, there is a notation for implicitly generating function objects:
    ```cpp
        void f(const Vector<int>& vec, const list<string>& lst, int x, const string& s) {
        cout << "number of values less than " << x
        << ": " << count(vec,[&](int a){ return a<x; })
        << '\n';
        cout << "number of values less than " << s
        << ": " << count(lst,[&](const string& a){ return a<s; })
        << '\n';
    }
    ```
- The notation `[&](int a){ return a<x; }` is called a **lambda expression**. It generates a function object similar to `Less_than<int>{x}`.
- The `[&]` is a capture list specifying that all local names used in the lambda body (such as `x`) will be accessed through **references**.
- Had we wanted to capture only `x`, we could have said so: `[&x]`.
- Had we wanted to give the generated object a **copy** of `x`, we could have said so: `[x]`.
- Capture nothing is `[ ]`, capture all local names used by **reference** is `[&]`, and capture all local names used by **value** is `[=]`.
- For a lambda defined within a member function:
  - `[this]` captures the **current object** by **reference** so that we can refer to class members.
  - If we want a **copy** of the current object, we say `[∗this]`.
- If we want to capture several specific objects, we can list them. The use of `[i,this]` in the use of` expect()` is an example.

### Lamdas as function arguments

- First, we need a function that applies an operation to each object pointed to by the elements of a container of pointers:
    ```cpp
    template<typename C, typename Oper>
    void for_each(C& c, Oper op) // assume that C is a container of pointers
    {
        for (auto& x : c)
        op(x); // pass op() a reference to each element pointed to
    }
    ```
- This is a simplified version of the standard-library `for_each` algorithm. Now, we can write a version of `user()` without writing a set of `_all` functions:
    ```cpp
    void user() {
        vector<unique_ptr<Shape>> v;
        while (cin)
        v.push_back(read_shape(cin));
        for_each(v,[](unique_ptr<Shape>& ps){ ps−>draw(); }); // draw_all()
        for_each(v,[](unique_ptr<Shape>& ps){ ps−>rotate(45); }); // rotate_all(45)
    }
    ```
- Like a function, a lambda can be **generic**. For example:
    ```cpp
    template<class S>
    void rotate_and_draw(vector<S>& v, int r) {
        for_each(v,[](auto& s){ s−>rotate(r); s−>draw(); });
    }
    ```
- Here, like in variable declarations, `auto` means that a value of any type is accepted as an initializer. This makes a lambda with an `auto` parameter a template, a **generic lambda**.
- We can define a function, `finally()` that takes an action to be executed on the exit from the scope:
    ```cpp
    void old_style(int n) {
    void∗ p = malloc(n∗sizeof(int)); // C-style
    auto act = finally([&]{free(p);}); // call the lambda upon scope exit
    // ...
    } // p is implicitly freed upon scope exit
- This is ad hoc, but far better than trying to correctly and consistently call` free(p)` on all exits from the function. The `finally()` function is **trivial**:
    ```cpp
    template <class F>
    [[nodiscard]] auto finally(F f)
    {
        return Final_action{f};
    }
    ```
- The class `Final_action` that supplies the necessary destructor can look like this:
    ```cpp
    template <class F>
    struct Final_action {
    explicit Final_action(F f) :act(f) {}
        ˜Final_action() { act(); }
        F act;
    }
    ```

## Template Mechanisms

### Variable Templates

- Introduced in C++14 that allows you to define a template for variables, similar to how you can define templates for functions and classes.
- The standard library uses variable templates to provide mathematical constants, such as `pi` and `log2e`:
    ```cpp
    template <class T>
        constexpr T viscosity = 0.4;

    template <class T>
        constexpr space_vector<T> external_acceleration = { T{}, T{−9.8}, T{} };

    auto vis2 = 2∗viscosity<double>;
    auto acc = external_acceleration<float>;
    ```
- Naturally, we can use arbitrary expressions of suitable types as initializers. Consider:
    ```cpp
    template<typename T, typename T2>
    constexpr bool Assignable = is_assignable<T&,T2>::value; // is_assignable is a type trait
    template<typename T>
    void testing()
    {
    static_assert(Assignable<T&,double>, "can't assign a double to a T");
    static_assert(Assignable<T&,string>, "can't assign a string to a T");
    }
    ```
- Variable templates are often used with type traits to simplify the syntax and make the code more readable.

### Aliases

- Surprisingly often, it is useful to introduce a synonym for a type or a template: `using size_t = unsigned int;`.
- It is very common for a parameterized type to provide an alias for types related to their template arguments. For example:
    ```cpp
    template<typename T>
    class Vector {
    public:
        using value_type = T;
        // ...
    }
    ```
- In fact, every standard-library container provides `value_type` as the name for the type of its elements. This allows us to write code that will work for every container that follows this convention. For example:
    ```cpp
    template<typename C>
    using Value_type = C::value_type; // the type of C’s elements

    template<typename Container>
    void algo(Container& c) {
        Vector<Value_type<Container>> vec; // keep results here
    // ...
    }
    ```

### Compile-Time if

- We can use a compile-time if as:
    ```cpp
    template<typename T>
    void update(T& target) {
        if constexpr(is_trivially_copyable_v<T>)
            simple_and_fast(target); // for "plain old data"
        else
            slow_and_safe(target); // for more complex types
    }
    ```
- ⚠️ `if constexpr` is not a text-manipulation mechanism and cannot be used to break the usual rules of grammar, type, and scope.