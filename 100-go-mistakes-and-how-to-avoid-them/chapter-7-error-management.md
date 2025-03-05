# Chapter 7: Error management

## #48: Panicking

- Example:
    ```go
    func main() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("recover", r)
            }
        }()

        f()
    }

    func f() {
        fmt.Println("a")
        panic("foo")
        fmt.Println("b")
    }
    ```
- Once a panic is triggered, it continues up the call stack until either the current goroutine has returned or panic is caught with `recover`.
- Panicking in Go should be used **sparingly**. We have seen two prominent cases:
  - 👍 One to signal a **programmer error**:
    - Invalid HTTP status code: `code < 100 || code > 999 `
    - SQL driver is `nil` (`driver.Driver` is an interface) or has already been registered: `driver == nil`
  - 👍 And another where our app fails to create a **mandatory dependency**. Hence, there are exceptional conditions that lead us to stop the app.
      - We depend on a service that needs to validate the provided email address with `MustCompile`.
  -  In most other cases, error management should be done with a function that returns a **proper error** type as the last return argument.

## #49: Ignoring when to wrap an error

- Error wrapping is about wrapping or packing an error inside a wrapper container that also makes the source error available.
- In general, the two main use cases for error wrapping are the following:
    - Adding additional context to an error
    - Marking an error as a specific error
- Before Go 1.13, to wrap an error, the only option without using an external library was to create a custom error type:
    ```go
    type BarError struct {
        Err error
    }
    func (b BarError) Error() string {
        return "bar failed:" + b.Err.Error()
    }
    ```
- To overcome this situation, Go 1.13 introduced the `%w` directive:
    ```go
    if err != nil {
        return fmt.Errorf("bar failed: %w", err)
    }
    ```
- The last option we will discuss is to use the `%v` directive, instead:
    ```go
    if err != nil {
        return fmt.Errorf("bar failed: %v", err)
    }
    ```
- The difference is that the error itself isn’t wrapped. We transform it into another error to add context, and the source error is no longer available.
- Let’s review all the different options we tackled:
    | Option                   | Extra Context                                                     | Marking an error | Source error available                                                |
    | ------------------------ | ----------------------------------------------------------------- | ---------------- | --------------------------------------------------------------------- |
    | Returning error directly | No                                                                | No               | Yes                                                                   |
    | Custom error type        | Possible (if the error type contains a string field, for example) | Yes              | Possible (if the source error is exported or accessible via a method) |
    | fmt.Errorf with %w       | Yes                                                               | No               | Yes                                                                   |
    | fmt.Errorf with %v       | Yes                                                               | No               | No                                                                    |
- To summarize, when handling an error, we can decide to wrap it. Wrapping is about adding additional context to an error and/or marking an error as a specific type.
  - If we need to mark an error, we should create a custom error type.
  - However, if we just want to add extra context, we should use `fmt.Errorf` with the `%w` directive as it doesn’t require creating a new error type.
- Yet, error wrapping creates potential **coupling** as it makes the source error available for the caller.
  - If we want to prevent it, we shouldn’t use error wrapping but error transformation, for example, using `fmt.Errorf` with the `%v` directive.

## #50: Checking an error type inaccurately

<p align="center"><img src="./assets/wrap-errors.png" width="500px" height="auto"></p>

- Go 1.13 came with a directive to wrap an error and a way to check whether the **wrapped error** is of a certain type with `errors.As`.
- This function **recursively** unwraps an error and returns true if an error in the chain matches the expected type.
    ```go
    // Get transaction ID
    amount, err := getTransactionAmount(transactionID)
    if err != nil {
        if errors.As(err, &transientError{}) {
            http.Error(w, err.Error(), http.StatusServiceUnavailable)
        } else {
            http.Error(w, err.Error(), http.StatusBadRequest)
        }
        return
    }
    ```
- ▶️ Regardless of whether the error is returned directly by the function we call or wrapped inside an error, `errors.As` will be able to recursively unwrap our main error and see if one of the errors is a specific type.

## #51: Checking an error value inaccurately

- A **sentinel error** is an error defined as a global variable:
    ```go
    import "errors"
    var ErrFoo = errors.New("foo") // the convention is to start with Err followed by the error type
    ```
- The general principle behind sentinel errors is to convey **expected** error that clients will expect to check. Therefore, as general guidelines:
  - 👍 **Expected** errors should be designed as error **values** (sentinel errors): `var ErrFoo = errors.New("foo")`.
  - 👍 **Unexpected** errors should be designed as error **types**: `type BarError struct { … }`, with `BarError` implementing the error interface.
- We have seen how `errors.As` is used to check an error against a **type**. With error **values**, we can use its counterpart: `errors.Is`:
    ```go
    err := query()
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
        // ...
        } else {
        // ...
        }
    }
    ```
- ▶️ if we use error wrapping in our app with the `%w` directive and `fmt.Errorf`, checking an error against a specific value should be done using` errors.Is` instead of `==`. Thus, even if the sentinel error is **wrapped**, `errors.Is` can recursively unwrap it and compare each error in the chain against the provided value.

## #52: Handling an error twice

- Consider the log below:
    ```
    2021/06/01 20:35:12 invalid latitude: 200.000000
    2021/06/01 20:35:12 failed to validate source coordinates
    ```
- Having **two log lines** for a **single error** is a problem. Why?
  - Because it makes **debugging harder**. For example, if this function is called multiple times concurrently, the two messages may not be one after the other in the logs, making the debugging process more complex.
- As a rule of thumb, an error should be handled **only once**. Logging an error is handling an error, and so is returning an error. Hence, we should **either log or return** an error, **never both** ❗.
-  Let’s rewrite our implementation to handle errors only once:
    ```go
    func GetRoute(srcLat, srcLng, dstLat, dstLng float32) (Route, error) {
        err := validateCoordinates(srcLat, srcLng)
        if err != nil {
            return Route{}, err
        }
        err = validateCoordinates(dstLat, dstLng)
        if err != nil {
            return Route{}, err
        }
        return getRoute(srcLat, srcLng, dstLat, dstLng)
    }
    ```
- The issue with this implementation is that we lost the origin of the error, so we need to **add additional context**:
- Let’s rewrite the latest version of our code using Go 1.13** error wrapping**:
    ```go
    func GetRoute(srcLat, srcLng, dstLat, dstLng float32) (Route, error) {
        err := validateCoordinates(srcLat, srcLng)
        if err != nil {
            return Route{}, fmt.Errorf("failed to validate source coordinates: %w", err)
        }
        err = validateCoordinates(dstLat, dstLng)
        if err != nil {
            return Route{}, fmt.Errorf("failed to validate target coordinates: %w", err)
        }
        return getRoute(srcLat, srcLng, dstLat, dstLng)
    }
    ```

## #53: Not handling an error

- When we want to ignore an error in Go, there’s only one way to write it:
    ```go
    _ = notify() // good
    notify()     // bad
    ```
- 👍 It may be a good idea to write a comment that indicates the **rationale** for **why** the error is **ignored**.
- Even if we are sure that an error can and should be ignored, we must do so **explicitly** by assigning it to the blank identifier. This way, a future reader will understand that we ignored the error intentionally.

## #54: Not handling defer errors

- As discussed in the previous section, if we don’t want to handle the error, we should ignore it explicitly using the blank identifier:
    ```go
    defer func() {
        _ = rows.Close()
    }()
    ```
-  In this case, calling `Close()` returns an error when it fails to free a DB connection from the pool. Hence, ignoring this error is probably not what we want to do.
- Most likely, a better option would be to log a message, or propagate it to the caller of `getBalance` so that they can decide how to handle it?
    ```go
    defer func() {
        err := rows.Close()
        if err != nil {
            return err
        }
    }()
    ```
- This implementation doesn’t compile. Indeed, the `return` statement is associated with the **anonymous** `func()` function, not `getBalance`. If we want to tie the error returned by `getBalance` to the error caught in the `defer` call, we must use **named result parameters**. Let’s write the first version:
    ```go
    func getBalance(db *sql.DB, clientID string) (balance float32, err error) {
        rows, err := db.Query(query, clientID)
        if err != nil {
            return 0, err
        }
        defer func() {
         err = rows.Close()
        }()
        if rows.Next() {
            err := rows.Scan(&balance)
            if err != nil {
                return 0, err
            }
            return balance, nil
        }
    }
    ```
- This code may look okay, but there’s a problem with it. If `rows.Scan` returns an error, `rows.Close` is executed anyway; but because this call overrides the error returned by `getBalance`, instead of returning an error, we may return a `nil` error if `rows.Close` returns successfully.
- Here’s our final implementation of the anonymous function:
    ```go
    defer func() {
        closeErr := rows.Close()
        if err != nil {
            if closeErr != nil {
                log.Printf("failed to close rows: %v", err)
            }
            return
        }
        err = closeErr
    }()
    ```
