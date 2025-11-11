# Input and Output

- An `ostream` converts typed objects to a stream of characters (bytes).
- An `istream` converts a stream of characters (bytes) to typed objects.
- The operations on `istreams` and `ostreams` are **type-safe**, **type-sensitive**, and extensible to handle **user-defined types**.
- The streams can be used for input into and output from strings, for formatting into string buffers, into areas of memory, and for file I/O.
- The I/O stream classes all have destructors that free all resources owned (such as buffers and file handles). That is, they are examples of "Resource Acquisition Is Initialization"

## Output

- The operator `<<` (*put to*) is used as an output operator on objects of type `ostream`.
- `cout` is the standard output stream and `cerr` is the standard stream for reporting errors.

## Input

- The operator `>>` (*get from*) is used as an input operator; `cin` is the standard input stream.
    ```cpp
    int i;
    double d;
    cin >> i >> d; // read into i and d
    ```
- In both cases, the read of the integer is terminated by **any character that is not a digit**. By default, `>>` skips initial whitespace, so a suitable complete input sequence would be:
    ```
    1234
    12.34e5
    ```
- By default, a **whitespace** character, such as a space or a newline, terminates the read. You can read a whole line using the `getline()` function. For example:
    ```cpp
    cout << "Please enter your name\n";
    string str;
    getline(cin,str);
    cout << "Hello, " << str << "!\n
    ```
- The newline that terminated the line is discarded, so cin is ready for the next input line.

## I/O State:

- An `iostream` has a **state** that we can examine to determine whether an operation succeeded. The most common use is to read a sequence of values:
    ```cpp
    vector<int> read_ints(istream& is) {
        vector<int> res;
        for (int i; is>>i; )
            res.push_back(i);
        return res;
    }
    ```
- This reads from `is` until something that is not an integer is encountered. That something will typically be the end of input. What is happening here is that the operation `is>>i` returns a reference to
`is`, and testing an `iostream` yields *true* if the stream is ready for another operation.


## I/O of User-Defined Types

- In addition to the I/O of built-in types and standard strings, the iostream library allows us to define
I/O for our own types:
    ```cpp
    struct Entry {
        string name;
        int number;
    }

    ostream& operator<<(ostream& os, const Entry& e) {
        return os << "{\"" << e.name << "\", " << e.number << "}";
    }
    ```
- The corresponding input operator is more complicated because it has to check for correct formatting and deal with errors:
```cpp
// read { "name" , number } pair. Note: formatted with { " " , and }
istream& operator>>(istream& is, Entry& e) {
    char c, c2;
    if (is>>c && c=='{' && is>>c2 && c2=='"') { // start with a { followed by a "
        string name; // the default value of a string is the empty string: ""
        while (is.get(c) && c!='"') // anything before a " is part of the name
            name+=c;
        if (is>>c && c==',') {
            int number = 0;
            if (is>>number>>c && c=='}') { // read the number and a }
                e = {name,number}; // assign to the entry
                return is;
            }
        }
    }
    is.setstate(ios_base::failbit); // register the failure in the stream
    return is;
}
```

## Output Formatting

