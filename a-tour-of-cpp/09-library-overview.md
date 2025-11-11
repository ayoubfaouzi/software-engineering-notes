# Standard Library

## Namespaces

- The facilities of the standard library are placed in namespace `std` and made available to users through modules or headed files.
- It is generally in poor taste to dump every name from a namespace into the global namespace.
- To use a suffix from a **sub-namespace**, we have to introduce it into the namespace in which we want to use it. For example:
    ```cpp
    // no mention of complex_literals
    auto z1 = 2+3i; // error: no suffix ’i’
    using namespace literals::complex_literals; // make the complex literals visible
    auto z2 = 2+3i; // ok: z2 is a complex<double>
    ```
- There is no coherent philosophy for what should be in a sub-namespace. However, suffixes **cannot be explicitly qualified** so we can only bring in a single set of suffixes into a scope without risking ambiguities. Therefore suffixes for a library meant to work with other libraries (that might define their own suffixes) are placed in sub-namespaces.
- The standard-library offers algorithms, such as `sort()` and `copy()`, in two versions:
  - A **traditional** sequence version taking a pair of iterators; e.g., `sort(begin(v),v.end())`
  - A **range** version taking a single range; e.g., `sort(v)`.
    ```cpp
    using namespace std;
    using namespace ranges;
    void f(vector<int>& v)
    {
        sort(v.begin(),v.end()); // error:ambiguous
        sort(v); // error: ambiguous
    }
    ```
- To protect against ambiguities when using traditional unconstrained templates, the standard requires that we **explicitly** introduce the range version of a standard-library algorithm into a scope:
    ```cpp
    using namespace std;
    void g(vector<int>& v)
    {
        sort(v.begin(),v.end()); // OK
        sort(v); // error: no matching function (in std)
        ranges::sort(v); // OK
        using ranges::sort; // sort(v) OK from here on
        sort(v); // OK
    }
    ```