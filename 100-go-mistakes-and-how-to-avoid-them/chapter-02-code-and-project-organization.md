# Chapter 2: Code and project organization

## #1: Unintended variable shadowing

- In Go, a variable name declared in a block can be redeclared in an inner block.
- **Variable shadowing** occurs when a variable name is redeclared in an inner block, but we saw that this practice is prone to mistakes:
```go
var client *http.Client
if tracing {
    client, err := createClientWithTracing()
    if err != nil {
        return err
    }
    log.Println(client)
} else {
    client, err := createDefaultClient()
    if err != nil {
        return err
    }
    log.Println(client)
}
// Use client
```

## #2: Unnecessary nested code

- A critical aspect of **readability** is the number of **nested** levels.
- In general, the more nested levels a function requires, the more complex it is to read and understand:
    ```go
    func join(s1, s2 string, max int) (string, error) {
        if s1 == "" {
            return "", errors.New("s1 is empty")
        } else {
            if s2 == "" {
                return "", errors.New("s2 is empty")
            } else {
                concat, err := concatenate(s1, s2)
                if err != nil {
                    return "", err
                } else {
                    if len(concat) > max {
                        return concat[:max], nil
                    } else {
                        return concat, nil
                    }
                }
            }
        }
    }
    ```
- *Align the happy path to the left; you should quickly be able to scan down one column to see the expected execution flow*.
- When an `if` block returns, we should omit the `else` block in all cases.
- If we encounter a **non happy-path**, we should flip the condition like so:
    ```go
    if s == "" {
        return errors.New("empty string")
    }
    // ...
    ```

## #3: Misusing init functions

- Consider the example below:
    ```go
    var db *sql.DB
    func init() {
        dataSourceName := os.Getenv("MYSQL_DATA_SOURCE_NAME")
        d, err := sql.Open("mysql", dataSourceName)
        if err != nil {
            log.Panic(err)
        }
        err = d.Ping()
        if err != nil {
            log.Panic(err)
        }
        db = d
    }
    ```
-  Let’s describe three main downsides of the code above:
   - It shouldn’t necessarily be **up to the package** itself to decide whether to **stop** the application. Perhaps a caller might have preferred implementing a retry or using a fallback mechanism. In this case, opening the database within an `init` function prevents client packages from implementing their error-handling logic.
   - If we add tests to this file, the `init` function will be executed before running the test cases, which isn’t necessarily what we want (for example, if we add unit tests on a utility function that doesn’t require this connection to be created). Therefore, the `init` function in this example **complicates writing unit tests**.
   - The last downside is that the example requires assigning the database connection pool to a **global variable**. Global variables have some severe drawbacks; for example:
       - Any functions can alter global variables within the package.
       - Unit tests can be more complicated because a function that depends on a global variable won’t be **isolated anymore**.
- We should be cautious with `init` functions. They can be helpful in some situations, however, such as defining static configuration. Otherwise, and in most cases, we should handle initializations through ad-hoc functions.

## #4: Overusing getters and setters

- Using **getters** and **setters** presents some advantages, including these:
    - They **encapsulate** a **behavior** associated with getting or setting a field, allowing new functionality to be added later (for example, validating a field, returning a computed value, or wrapping the access to a field around a mutex).
    - They **hide** the **internal representation**, giving us more flexibility in what we expose.
    - They provide a **debugging** interception point for when the property changes at run time, making debugging easier.
- If we fall into these cases or foresee a possible use case while guaranteeing **forward compatibility**, using getters and setters can bring some value. For example, if we use them with a field called `balance`, we should follow these naming conventions:
    - The getter method should be named `Balance` (not `GetBalance`).
    - The setter method should be named `SetBalance`.
- We shouldn’t **overwhelm** our code with getters and setters on structs if they don’t bring any value. We should be pragmatic and strive to find the right **balance** between efficiency and following idioms that are sometimes considered indisputable in other programming paradigms.

## #5: Interface pollution

- Interface pollution is about **overwhelming** our code with **unnecessary abstractions**, making it harder to understand.
- What makes Go interfaces so different is that they are satisfied **implicitly**. There is no explicit keyword like *implements* to mark that an object `X` implements interface `Y`.
- While designing interfaces, the **granularity** (how many methods the interface contains) is also something to keep in mind:
  - 🌟 *The bigger the interface, the weaker the abstraction*.
  - Examples: `io.Reader`.
- When to use interfaces ?
  - When multiple types implement a **common behavior**. For example, This `Interface` has a strong potential for reusability because it encompasses the common behavior to **sort** any collection that is index-based:
  ```go
    type Interface interface {
        Len() int
        Less(i, j int) bool
        oSwap(i, j int)
    }
  ```
  - **Decoupling** our code from an implementation:
    - **Liskov Substitution Principle** (the L in *Robert C. Martin’s* SOLID design principles).
    - If we rely on an **abstraction** instead of a **concrete** implementation, the implementation itself can be replaced with another without even having to change our code. Example:
    ```go
    type customerStorer interface {
        StoreCustomer(Customer) error
    }
    type CustomerService struct {
        storer customerStorer
    }
    func (cs CustomerService) CreateNewCustomer(id string) error {
        customer := Customer{id: id}
        return cs.storer.StoreCustomer(customer)
    }
    ```
    - This gives us more **flexibility** in how we want to test the method:
      - Use the concrete implementation via integration tests.
      - Use a mock (or any kind of test double) via unit tests.
  - **Restrict** a type to a specific behavior for various reasons, such as semantics enforcement:
    ```go
    type Foo struct {
        threshold intConfigGetter
    }
    func NewFoo(threshold intConfigGetter) Foo {
        return Foo{threshold: threshold}
    }
    func (f Foo) Bar() {
        threshold := f.threshold.Get()
        // ...
    }
    ```
    - The configuration getter is **injected** into the `NewFoo` factory method. It doesn’t impact a client of this function because it can still pass an `IntConfig` struct as
it implements `intConfigGetter`. Then, we can only read the configuration in the `Bar` method, not modify it.
- The main caveat when programming meets abstractions is remembering that abstractions should be **discovered**, not **created**.
  - 🌟 *Don’t design with interfaces, discover them*.

## #6: Interface on the producer side

- **Producer** side — An interface defined in the same package as the concrete implementation.
- **Consumer** side — An interface defined in an external package where it’s used. <p align="center"><img src="./assets/iface-producer-consumer.png" width="500px" height="auto"></p>
- It’s common to see developers creating interfaces on the producer side, alongside the concrete implementation. This design is perhaps a habit from developers having a *C#* or a *Java* background. But in Go, in most cases this is not what we should do 🤨.
- Let’s discuss the following example:
    ```go
    package store
    // CustomerStorage is a good way 🤷 to decouple the client code from the actual implementation.
    // Or, perhaps we can foresee that it will help clients in creating test doubles.
    type CustomerStorage interface {
        StoreCustomer(customer Customer) error
        GetCustomer(id string) (Customer, error)
        UpdateCustomer(customer Customer) error
        GetAllCustomers() ([]Customer, error)
        GetCustomersWithoutContract() ([]Customer, error)
        GetCustomersWithNegativeBalance() ([]Customer, error)
    }
    ```
- This isn’t a best practice in Go. Why ? It’s not up to the producer to **force a given abstraction** for all the clients. Instead, it’s up to the client to decide whether it needs some form of abstraction and then determine the best abstraction level for its needs. For example:
    ```go
    package client

    type customersGetter interface {
        GetAllCustomers() ([]store.Customer, error)
    }
    ```
- The main point is that the `client` package can now define the most **accurate** abstraction for its need (here, only one method). It relates to the concept of the *Interface Segregation Principle* (the `I` in *SOLID*) ▶️ No client should be forced to depend on methods it doesn’t use.
- An interface should live on the **consumer** side in most cases. However, in particular contexts (for example - in the standard library `encoding, encoding/json, encoding/binary`, when we know — not foresee—that an abstraction will be helpful for consumers), we may want to have it on the **producer** side. If we do, we should strive to keep it as **minimal** as possible, increasing its **reusability** potential and making it more easily **composable**.

## #7: Returning interfaces

- We will consider two packages: `client`, which contains a `Store` interface and `store`, which contains an implementation of `Store`. <p align="center"><img src="./assets/store-client-dependency.png" width="500px" height="auto"></p>
- The `client` package can’t call the `NewInMemoryStore` function anymore; otherwise, there would be a **cyclic dependency**.
- In general, returning an interface restricts **flexibility** because we force all the clients to use one particular type of abstraction.
- *Be conservative in what you do, be liberal in what you accept from others.* If we apply this idiom to Go, it means
    - 👍 Returning structs instead of interfaces.
    - 👍 Accepting interfaces if possible.
- We shouldn’t return interfaces but concrete implementations. Otherwise, it can make our design more complex due to package dependencies and can restrict flexibility because all the clients would have to rely on the same
abstraction.
- Also, we will only be able to use the methods **defined** in the interface, and not the methods defined in the **concrete type**.

## #8: any says nothing

- With Go 1.18, the predeclared type `any` became an alias for an **empty interface {}**.
- In assigning a value to an `any` type, we **lose** all **type information**, which requires a type assertion to get anything useful out of the `i` variable.
    ```go
    func main() {
        var i any
        i = 42
        i = "foo"
        i = struct {
            s string
        }{
            s: "bar",
        }
        i = f
        _ = i
    }
    func f() {}
    ```
- In methods, accepting or returning an `any` type doesn’t **convey meaningful information**:
  - Because there is no **safeguard** at compile time, nothing prevents a caller from calling these methods with whatever data type.
  - Also, the methods lack **expressiveness**. If future developers need to use the parameters of type `any`, they will probably have to dig into the documentation or read the code to understand how to use these methods.
- What are the cases when `any` is helpful?:
  - In the `encoding/json` package. Because we can marshal any type, the `Marshal` function accepts an `any` argument.
  - Another example is in the `database/sql` package. If the query is parameterized (for example, `SELECT * FROM FOO WHERE id = ?`), the parameters could be any kind.

## #9: Being confused about when to use generics

- Go 1.18 adds generics to the language 🥳.
-  👍 Few common uses where generics are recommended:
    - **Data structures** : We can use generics to **factor out** the element type if we implement a binary tree, a linked list, or a heap, for example.
    - **Functions working with slices, maps, and channels of any type** : A function to merge two channels would work with any `channel` type, for example
      - Hence, we could use type parameters to factor out the channel type:
        ```go
        func merge[T any](ch1, ch2 <-chan T) <-chan T {
        // ...
        }
        ```
    - **Factoring out behaviors instead of types**:  The `sort` package, for example, contains a `sort.Interface` interface with three methods:
        ```go
        type Interface interface {
            Len() int
            Less(i, j int) bool
            Swap(i, j int)
        }
        ```
      - This interface is used by different functions such as `sort.Ints` or `sort.Float64s`. Using type parameters, we could factor out the sorting behavior(for example, by defining a `struct` holding a slice and a comparison function):
        ```go
        type SliceFn[T any] struct {
            S []T
            Compare func(T, T) bool
        }
        func (s SliceFn[T]) Len() int { return len(s.S) }
        func (s SliceFn[T]) Less(i, j int) bool { return s.Compare(s.S[i], s.S[j]) }
        func (s SliceFn[T]) Swap(i, j int) { s.S[i], s.S[j] = s.S[j], s.S[i] }
        ```
      - Then, because the `SliceFn` struct implements `sort.Interface`, we can sort the provided slice using the `sort.Sort(sort.Interface)` function:
        ```go
        s := SliceFn[int]{
            S: []int{3, 2, 1},
            Compare: func(a, b int) bool {
                return a < b
            },
        }
        sort.Sort(s)
        fmt.Println(s.S)
        ```
- 👎 when is it recommended that we not use generics:
    - When **calling a method of the type argument**: Consider a function that receives an `io.Writer` and calls the `Write` method, for example:
        ```go
        func foo[T io.Writer](w T) {
            b := getBytes()
            _, _ = w.Write(b)
        }
        ```
        - In this case, using generics won’t bring any value to our code whatsoever. We should make the `w` argument an `io.Writer` directly.
    - When it makes our code **more complex**: Generics are never mandatory, and as Go developers, we have lived without them for more than a decade. If we’re writing generic functions or structures and we figure out that it doesn’t make our code clearer, we should probably reconsider our decision for that particular use case.

## #10: Not being aware of the possible problems with type embedding

- 👎 The wolloing example is a wrong usage of type embedding. Since `sync.Mutex` is an embedded type, the `Lock` and `Unlock` methods will be **promoted**. Therefore, both methods become **visible** to external clients using `InMem`:
    ```go
    type InMem struct {
        sync.Mutex
        m map[string]int
    }
- We want to write a custom logger that contains an `io.WriteCloser` and exposes two methods, `Write` and `Close`. If `io.WriteCloser` wasn’t **embedded**, we would need to write it like so:
    ```go
    type Logger struct {
        writeCloser io.WriteCloser
    }
    func (l Logger) Write(p []byte) (int, error) {
        return l.writeCloser.Write(p) // Forwards the call to writeCloser
    }
    func (l Logger) Close() error {
        return l.writeCloser.Close() // Forwards the call to writeCloser
    }
    func main() {
        l := Logger{writeCloser: os.Stdout}
        _, _ = l.Write([]byte("foo"))
        _ = l.Close()
    }
    ```
- 👍 Logger `would` have to provide both a `Write` and a `Close` method that would only **forward** the call to `io.WriteCloser`. However, if the field now becomes **embedded**, we can remove these forwarding methods:
    ```go
    type Logger struct {
        io.WriteCloser
    }
    func main() {
        l := Logger{WriteCloser: os.Stdout}
        _, _ = l.Write([]byte("foo"))
        _ = l.Close()
    }
    ```
- If we decide to use type embedding, we need to keep two main constraints in mind:
    - It shouldn’t be used solely as some **syntactic sugar** to simplify accessing a field (such as `Foo.Baz()` instead of `Foo.Bar.Baz()`). If this is the only rationale, let’s not embed the inner type and use a field instead.
    - It shouldn’t promote data (fields) or a behavior (methods) we want to **hide** from the outside: for example, if it allows clients to access a locking behavior (`sync.Mutex`) that should remain **private** to the struct.

## #11 Not using the functional options pattern

- How can we implement passing an configuration option to a function in an API-friendly way? Let’s look at the different options.
- **Config struct**:
  - The **mandatory** parameters could live as function parameters, whereas the **optional** parameters could be handled in the `Config` struct:
    ```go
    type Config struct {
        Port int
    }
    func NewServer(addr string, cfg Config) { }
    ```
    - 👍 This solution fixes the **compatibility** issue. Indeed, if we add new options, it will not break on the client side.
    - 👎 However, this approach does not distinguish between a field purposely set to 0 and a missing field:
        - 0 for an integer, 0.0 for a floating-point type
        - "" for a string
        - Nil for slices, maps, channels, pointers, interfaces.
        - ▶️ One option might be to handle all the parameters of the configuration struct as **pointers**, however, it’s not handy for clients to work with pointers as they have to create a variable and then pass a pointer 🤷, also client using our library with the default configuration will need to pass an **empty struct**.
- **Builder pattern**:
  - The construction of `Config` is separated from the struct itself. It requires an extra struct, `ConfigBuilder`, which receives methods to configure and build a `Config`:
    ```go
    type Config struct {
        Port int
    }
    type ConfigBuilder struct {
        port *int
    }
    func (b *ConfigBuilder) Port(
        port int) *ConfigBuilder {
        b.port = &port
        return b
    }
    func (b *ConfigBuilder) Build() (Config, error) {
        cfg := Config{}
        if b.port == nil {
            cfg.Port = defaultHTTPPort
        } else {
            if *b.port == 0 {
                cfg.Port = randomPort()
            } else if *b.port < 0 {
                return Config{}, errors.New("port should be positive")
            } else {
                cfg.Port = *b.port
            }
        }
        return cfg, nil
    }
    ```
  - The `ConfigBuilder` struct holds the client configuration. It exposes a `Port` method to set up the port. Usually, such a configuration method returns the **builder itself** so that we can use method **chaining** (for example, `builder.Foo("foo").Bar("bar")`). It also exposes a `Build` method that holds the logic on initializing the port value (whether the pointer was nil, etc.) and returns a `Config` struct once created.
    - 👍 This approach makes port management handier. It’s **not required** to pass an integer pointer, as the `Port` method accepts an integer. However, we still need to pass a config struct that can be empty if a client wants to use the default configuration 🤷.
    - 👎 In programming languages where exceptions are thrown, builder methods such as `Port` can raise **exceptions** if the input is invalid. If we want to keep the ability to **chain** the calls, the
function **can’t return an error**. Therefore, we have to delay the validation in the `Build()`.
- **Functional options pattern**:
- The main idea is as follows:
  - An **unexported** struct holds the configuration: options.
  - Each option is a function that returns the **same type**: `type Option func(options *options) error`. For example, `WithPort` accepts an `int` argument that represents the port and returns an `Option` type that represents how to update the options struct.
    ```go
    type options struct {
        port *int
    }
    type Option func(options *options) error

    func WithPort(port int) Option {
        return func(options *options) error {
            if port < 0 {
                return errors.New("port should be positive")
            }
            options.port = &port
            return nil
        }
    }
    ```
  - Each config field requires creating a public function (that starts with the `With` prefix by convention) containing similar logic: validating inputs if needed and updating the config struct.
    ```go
    func NewServer(addr string, opts ...Option) (*http.Server, error) {
        var options options
        for _, opt := range opts {
            err := opt(&options)
            if err != nil {
                return nil, err
            }
        }
        // At this stage, the options struct is built and contains the config
        // Therefore, we can implement our logic related to port configuration
        var port int
        if options.port == nil {
            port = defaultHTTPPort
        } else {
            if *options.port == 0 {
                port = randomPort()
            } else {
                port = *options.port
            }
        }
        // ...
    }
  - Because `NewServer` accepts **variadic** Op`tion arguments, a client can now call this API by passing multiple options following the mandatory address argument. For example:
    ```go
    server, err := httplib.NewServer("localhost",
        httplib.WithPort(8080),
        httplib.WithTimeout(time.Second)
    )
    ```
  - 👍 Provides a handy and API-friendly way to handle options and represent the most idiomatic way. If the client needs the default configuration, it doesn’t have to provide an argument.

## #12: Project misorganization

- **Project structure**:
  - Go language maintainer has no strong convention about structuring a project in Go. However, one layout has emerged over the years: [project-layout](https://github.com/golang-standards/project-layout).
  - There’s no **mandatory** convention. This layout may be helpful for you or not, but what’s important here is that **indecision** is the only **wrong decision**.
- **Package organization**:
  - In Go, there is no concept of subpackages. However, we can decide to organize packages within subdirectories.
  - The main benefit of subdirectories is to keep packages in a place where they live with **high cohesion**.
  - Avoid **premature packaging** because it might cause us to over-complicate a project. Sometimes, it’s better to use a simple organization and have our project evolve when we understand what it contains rather than **forcing** ourselves to make the perfect structure up front.
  - Avoid having **dozens** of **nano** packages containing only one or two files. Conversely, we should also avoid **huge** packages that dilute the meaning of a package name.
  - We should name our packages after what they **provide**, not what they **contain**. Also, naming should be **meaningful**. Therefore, a package name should be **short**, **concise**, **expressive**, and, by convention, a **single lowercase word**.
  - We should **minimize** what should be **exported** as much as possible to **reduce the coupling** between packages and keep unnecessary exported elements hidden. In doubt, default to not exporting it!

## #13: Creating utility packages

- As a rule of thumb, creating **shared** packages without meaningful names isn’t a good idea; this includes utility packages such as `utils`, `common`, or `base`. Also, bear in mind that naming a package after what it provides and not what it contains can be an efficient way to increase its **expressiveness**.
- Consider the following example:
    ```go
    package util

    func NewStringSet(...string) map[string]struct{} {
        // ...
    }
    func SortStringSet(map[string]struct{}) []string {
        // ...
    }

    // A client will use this package like this:
    set := util.NewStringSet("c", "a", "b")
    fmt.Println(util.SortStringSet(set))
    ```
- Instead of a utility package, we should create an expressive package name such as `stringset`. For example:
    ```go
    package stringset

    func New(...string) map[string]struct{} { ... }
    func Sort(map[string]struct{}) []string { ... }

    // In this example, we removed the suffixes for NewStringSet and SortStringSet,
    // which respectively became New and Sort. On the client side, it now looks like this:
    set := stringset.New("c", "a", "b")
    fmt.Println(stringset.Sort(set))
    ```
- We could even go a step further. Instead of exposing utility functions, we could create a specific type and expose `Sort` as a method this way:
    ```go
    type Set map[string]struct{}

    func New(...string) Set { ... }
    func (s Set) Sort() []string { ... }

    // This change makes the client even simpler. There would only be one reference to the
    set := stringset.New("c", "a", "b")
    fmt.Println(set.Sort())
    ```

## #14 Ignoring package name collisions

- Package collisions occur when a **variable name** collides with an existing **package name**, preventing the package from being reused:
    ```go
    redis := redis.NewClient() // Calls NewClient from the redis package
    v, err := redis.Get("foo") // Uses the redis variable
    ```
- Two solutions:
  - Change the variable name to `redisClient` for example.
  - Create an alias for the `redis` package: `import redisapi "mylib/redis"`
- One option could also be to use **dot imports** to access all the public elements of a package without the package qualifier ▶️ increase confusion !
- We should avoid naming collisions between a **variable** and a **built-in** function. For example, we could do something like this:
    ```go
    copy := copyFile(src, dst) // The copy variable collides with the copy built-in function.
    ```
- 👍 In summary, we should prevent variable name collisions to avoid **ambiguity**.

## #15: Missing code documentation

- First, every **exported** element must be **documented**.
- The convention is to add comments, starting with the **name** of the exported element.
- As a convention, each comment should be a **complete sentence** that ends with **punctuation**.
- When we document a function, we should highlight what the function **intends to do**, not how it does it ⚠️; this belongs to the **core** of a function and comments, not documentation.
- 💡 We should ideally provide enough information that the consumer does not have to look at our code to understand how to use an exported element.
- When it comes to documenting a **variable** or a **constant**, we might be interested in conveying two aspects: its **purpose** and its **content**.
  - The **former** should live as **code documentation** to be useful for external clients.
  - The latter, though, **shouldn’t** necessarily be **public**. For example:
    ```go
    // DefaultPermission is the default permission used by the store engine.
    const DefaultPermission = 0o644 // Need read and write accesses.
    ```
- This constant represents the default permission. The code documentation conveys its purpose, whereas the comment alongside the constant describes its actual content.
- To help clients and maintainers understand a package’s scope, we should also document each **package**. The convention is to start the comment with *// Package* followed by the package name:
    ```go
    // Package math provides basic constants and mathematical functions.
    //
    // This package does not guarantee bit-identical results
    // across architectures.
    package math
    ```
- The first line of a package comment should be **concise**. That’s because it will appear in the package. Then, we can provide all the information we need in the following lines.

## #16: Not using linters

- A linter is an automatic tool to analyze code and catch errors.
- `Golint` (deprecated now) is the linter that is developer by Google. There's no drop-in replacement for it, but tools such as `Staticcheck` and `go vet` should be used instead.
- `golangci-lint` integrate almost all community driven linters.
