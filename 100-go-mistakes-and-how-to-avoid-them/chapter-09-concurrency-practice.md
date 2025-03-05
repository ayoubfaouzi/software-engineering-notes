# Chapter 9: Concurrency: Practice

## #61: Propagating an inappropriate context

- Context propagation can sometimes lead to subtle bugs, preventing subfunctions from being correctly executed.
- Consider the example below:
    ```go
    func handler(w http.ResponseWriter, r *http.Request) {
        response, err := doSomeTask(r.Context(), r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        go func() {
            err := publish(r.Context(), response)
            // Do something with err
        }()
        writeResponse(response)
    }
    ```
- We have to know that the context attached to an HTTP request can cancel in different conditions:
    - When the client’s connection closes
    - In the case of an HTTP/2 request, when the request is canceled
    - ⚠️ When the response has been written back to the client
- When the response has been written to the client, the context associated with the request will be **canceled**. Therefore, we are facing a **race condition**.
- If the response is written before or during the Kafka publication, the message shouldn’t be published. Calling `publish` will return an error because we returned the HTTP response quickly.
- Ideally, we would like to have a new context that is **detached** from the potential parent cancellation but still conveys the **values**.
- The context’s deadline is managed by the `Deadline` method and the cancellation signal is managed via the `Done` and `Err` methods. When a deadline has passed or the context has been canceled, `Done` should return a **closed channel**, whereas `Err` should return an error. Finally, the values are carried via the `Value` method.
- Let’s create a custom context that detaches the cancellation signal from a parent context:
    ```go
    type detach struct {
        ctx context.Context
    }
    func (d detach) Deadline() (time.Time, bool) {
        return time.Time{}, false
    }
    func (d detach) Done() <-chan struct{} {
        return nil
    }
    func (d detach) Err() error {
        return nil
    }
    func (d detach) Value(key any) any {
        return d.ctx.Value(key)
    }
    ```
- Thanks to our custom context, we can now call publish and detach the cancellation signal:
    ```go
    err := publish(detach{ctx: r.Context()}, response)
    ```

## #62: Starting a goroutine without knowing when to stop it

- A goroutine is a resource like any other that must eventually be closed to free memory or other resources.
- Starting a goroutine without knowing when to stop it is a design issue. Whenever a goroutine is started, we should have a clear plan about when it will stop, if we don't, it can lead to **leaks** ⚠️.
-  In terms of memory, a goroutine starts with a **minimum stack size** of `2 KB`, which can grow and shrink as needed (the **maximum stack** size is `1 GB` on 64-bit and `250 MB` on 32-bit). Memory-wise, a goroutine can also hold variable references allocated to the heap. Meanwhile, a goroutine can hold resources such as HTTP or database connections, open files, and network sockets that should eventually be closed gracefully. If a goroutine is leaked, these kinds of resources will also be leaked.
-  Here’s a first implementation:
    ```go
    func main() {
        newWatcher()
        // Run the application
    }

    type watcher struct { /* Some resources */ }

    func newWatcher() {
        w := watcher{}
        go w.watch()
    }
    ```
- The problem with this code is that when the main goroutine exits, the application is stopped. Hence, the resources created by watcher aren’t closed **gracefully**.
- One option could be to pass to `newWatcher` a context that will be **canceled** when `main` returns:
    ```go
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    newWatcher(ctx)
    ```
  - However, can we guarantee that watch will have time to do so? Absolutely not — and that’s a design flaw !
- The problem is that we used signaling to convey that a goroutine had to be stopped. We didn’t **block the parent** goroutine until the resources had been closed. Let’s make sure we do:
```go
func main() {
    w := newWatcher()
    defer w.close() // Defers the call to the close method
    // Run the application
}
func newWatcher() watcher {
    w := watcher{}
    go w.watch()
    return w
}
func (w watcher) close() {
    // Close the resources
}
```
- `watcher` has a new method: `close`. Instead of signaling watcher that it’s time to close its resources, we now call this `close` method, using `defer` to guarantee that the resources are closed before the application exits.

## #63: Not being careful with goroutines and loop variables

- In the following example, we initialize a slice. Then, within a closure executed as a new goroutine, we access this element:
    ```go
    s := []int{1, 2, 3}
    for _, i := range s {
        go func() {
            fmt.Print(i)
        }()
    }
    ```
- The output of this code **isn’t deterministic**. For example, sometimes it prints `233` and other times `333`.
- All the goroutines refer to the exact same variable! When a goroutine runs, it prints the value of `i` at the time `fmt.Print` is executed. Hence, i may have been modified since the goroutine was launched.
- The first solution, if we want to keep using a closure, involves **creating a new variable**:
    ```go
    for _, i := range s {
        val := i
        go func() {
            fmt.Print(val)
        }()
    }
    ```
- The second option no longer relies on a closure and instead uses an **actual function**:
    ```go
    for _, i := range s {
        go func(val int) {
            fmt.Print(val)
        }(i)
    }
    ```

## #64: Expecting deterministic behavior using select and channels

Imagine an example where we receive from two channels, but we want to prioritize `messageCh`. For example, if a disconnection occurs, we want to ensure that we have received all the messages before returning. We may decide to handle the prioritization like so:
```go
for {
    select {
        case v := <-messageCh: fmt.Println(v)
    case <-disconnectCh:
        fmt.Println("disconnection, return")
        return
    }
}

// Dummy code to produce messages.
for i := 0; i < 10; i++ {
    messageCh <- i
}
disconnectCh <- struct{}{}
```

- If we run this example, here is a possible output if `messageCh` is buffered:
```
0 1 2 3 4
disconnection, return
```
- Why 🤔? If one or more of the communications can proceed, a single one that can proceed is chosen via a uniform **pseudo-random selection**.
- Unlike a `switch` statement, where the first case with a match wins, the `select` statement selects **randomly** if multiple options are possible.
- There are different possibilities if we want to receive all the messages before returning in case of a disconnection.
    - If there’s a **single producer** goroutine, we have two options:
        - Make `messageCh` an **unbuffered** channel instead of a buffered channel.
        - Use a **single channel** instead of two channels. For example, we can define a struct that conveys either both messages.
    - If we fall into the case where we have multiple producer goroutines, it may be impossible to guarantee which one writes first. Hence, whether we have an **unbuffered** `messageCh` channel or a **single** channel, it will lead to a **race condition** among the producer goroutines.
      - In that case, we can implement the following solution: Receive from either `messageCh` or `disconnectCh` and if a disconnection is received, read all the existing messages in `messageCh`, if any, then return.
```go
for {
    select {
    case v := <-messageCh: fmt.Println(v)
    case <-disconnectCh:
        for {
            select {
            case v := <-messageCh: fmt.Println(v)
            default:
                fmt.Println("disconnection, return")
                return
            }
        }
    }
}
```

## #65: Not using notification channels

- If we don’t need a specific value to convey some information, we need a channel **without data**. The idiomatic way to handle it is a channel of **empty structs**: `chan struct{}`.
- In Go, an empty struct is a struct without any fields. Regardless of the architecture, it occupies **zero bytes** of storage. But why not use an **empty interface** `(var i interface{})`? Because an empty interface isn’t free; it occupies 8 bytes on 32-bit architecture and 16 bytes on 64-bit architecture.
- For example, if we need a **hash set** structure (a collection of unique elements), we should use an empty struct as a value: `map[K]struct{}`.
- An empty struct clarifies for receivers that they shouldn’t expect any meaning from a message’s content—only the fact that they have received a message. In Go, such channels are called **notification channels**.

## #66: Not using nil channels

- In this example of merging two channels, we can use **nil channels** to implement an elegant state machine that will remove one case from a `select` statement, it doesn’t require a **busy loop** that will waste CPU cycles.
- If `ch1` is closed, we assign `ch1` to `nil`. Hence, during the next loop iteration, the `select` statement will only wait for two conditions:
    - `ch2` has a new message.
    - `ch2` is closed
```go
func merge(ch1, ch2 <-chan int) <-chan int {
    ch := make(chan int, 1)
    go func() {
        for ch1 != nil || ch2 != nil {
            select {
                case v, open := <-ch1:
                    if !open {
                        ch1 = nil
                        break
                    }
                    ch <- v
                case v, open := <-ch2:
                    if !open {
                        ch2 = nil
                        break
                    }
                    ch <- v
            }
        }
        close(ch)
    }()
    return ch
}
```

## #67: Being puzzled about channel size

- A **buffered channel** doesn’t provide any **strong synchronization**. Buffered channels can lead to obscure **deadlocks** that would be immediately apparent with unbuffered channels.
- While using a **worker pooling-like** pattern, meaning spinning a fixed number of goroutines that need to send data to a shared channel. In that case, we can tie
the channel size to the number of goroutines created 👍.
- When using channels for **rate-limiting** problems. For example, if we need to enforce resource utilization by bounding the number of requests, we should set
up the channel size according to the limit 👍.
- It’s pretty common to see a codebase using magic numbers for setting a channel size: `ch := make(chan int, 40)`.
    - Why 40? What’s the rationale? Why not `50` or even `1000`? Setting such a value should be done for a good reason.
    - Perhaps it was decided following a **benchmark** or performance tests.
    - In many cases, it’s probably a good idea to **comment** on the rationale for such a value.
- Let’s bear in mind that deciding about an accurate **queue size** isn’t an easy problem. First, it’s a balance between CPU and memory. The smaller the value, the more **CPU contention** we can face. But the bigger the value, the more **memory** will need to be allocated.
- Another point to consider is the one mentioned in a 2011 white paper about [LMAX Disruptor](https://lmax-exchange.github.io/disruptor/files/Disruptor-1.0.pdf):
    > Queues are typically always close to full or close to empty due to the differences in pace between consumers and producers. They very rarely operate in a balanced middle ground where the rate of production and consumption is evenly matched
- So, it’s rare to find a channel size that will be steadily accurate, meaning an accurate value that won’t lead to too much contention or a waste of memory allocation.

## #68: Forgetting about possible side effects with string formatting

- Let's demonstrate with this first example how one formatting a key from a context can lead to a **data race**:
    ```go
    func (w *watcher) Watch(ctx context.Context, key string, opts ...OpOption) WatchChan {
        ctxKey := fmt.Sprintf("%v", ctx) // Formats the map key depending on the provided context
        wgs := w.streams[ctxKey]
    ```
- When formatting a string from a context created with values (`context.WithValue`), Go will read **all the values** in this context.
  - In this case, the *etcd* developers found that the context provided to `Watch` was a context containing mutable values (for example, a pointer to a struct) in some conditions.
  - They found a case where one goroutine was updating one of the context values, whereas another was executing `Watch`, hence reading all the values in this context.
  - ⚠️ This led to a data race.
- The [fix](https://github.com/etcd-io/etcd/pull/7816) was to not rely on `fmt.Sprintf` to format the map’s key to prevent traversing and reading the chain of wrapped values in the context. Instead, the solution was to implement a custom `streamKeyFromCtx` function to extract the key from a specific context value that **wasn’t mutable**.
- The second example deals with a `Customer` struct that can be accessed **concurrently**:
    ```go
    type Customer struct {
        mutex sync.RWMutex
        id string
        age int
    }
    func (c *Customer) UpdateAge(age int) error {
        c.mutex.Lock()
        defer c.mutex.Unlock()
        if age < 0 {
            return fmt.Errorf("age should be positive for customer %v", c)
        }
        c.age = age
        return nil
    }

    func (c *Customer) String() string {
        c.mutex.RLock()
        defer c.mutex.RUnlock()
        return fmt.Sprintf("id %s, age %d", c.id, c.age)
    }
    ```
- If the provided `age` is negative, we return an error. Because the error is formatted, using the `%s` directive on the receiver, it will call the `String` method
to format `Customer`. But because `UpdateAge` already acquires the mutex lock, the `String` method won’t be able to acquire it.
- ⚠️ Hence, this leads to a deadlock situation. (That's why we should create unit tests also for edge cases 🧠).
- In our case, locking the mutex only after the age has been checked avoids the deadlock situation or we change the way we format the error so that it doesn’t call the `String` method.

## #69: Creating data races with append

- In the following example, we will initialize a slice and create two goroutines that will use `append` to create a new slice with an additional element:
    ```go
    s := make([]int, 1)
    go func() {
        s1 := append(s, 1)
        fmt.Println(s1)
    }()
    go func() {
        s2 := append(s, 1)
        fmt.Println(s2)
    }()
    ```
  - In this example, we create a slice with `make([]int, 1)`. The code creates a one length, one-capacity slice. Thus, because the slice is **full**, using `append` in each goroutine returns a slice backed by a new array. It **doesn’t mutate** the existing array; hence, it doesn’t lead to a data race.
- Instead of creating a slice with a length of 1, we create it with a length of 0 but a capacity of 1: `s := make([]int, 0, 1)`:
    ```
    WARNING: DATA RACE
    Write at 0x00c00009e080 by goroutine 10:
    ...
    Previous write at 0x00c00009e080 by goroutine 9:
    ```
  - The array isn’t full. Both goroutines attempt to update the **same index** of the backing array (index 1), which is a data race ⚠️.
- To solve this issue, we can make a copy and uses `append` on the copied slice. This prevents a data race because both goroutines work on isolated data.

🎯 Data races with slices and maps

- How much do data races impact slices and maps? When we have multiple goroutines the following is true:
    - Accessing the same slice index with at least one goroutine updating the value is a data race. The goroutines access the same memory location.
    - Accessing different slice indices regardless of the operation isn’t a data race; different indices mean different memory locations.
    - Accessing the same map (regardless of whether it’s the same or a different key) with at least one goroutine updating it is a data race. Why is this different from a slice data structure? As we mentioned in chapter 3, a map is an array of buckets, and each bucket is a pointer to an array of key-value pairs. A hashing algorithm is used to determine the array index of the bucket. Because this algorithm contains some randomness during the map initialization, one execution may lead to the same array index, whereas another execution may not. The race detector handles this case by raising a warning regardless of whether an actual data race occurs.

## #70: Using mutexes inaccurately with slices and maps

- Consider the example below:
    ```go
    func (c *Cache) AverageBalance() float64 {
        c.mu.RLock()
        balances := c.balances
        c.mu.RUnlock()

        sum := 0.
        for _, balance := range balances {
            sum += balance
        }
        return sum / float64(len(balances))
    }
    ```
- If we run a test using the `-race` flag with two concurrent goroutines, one calling `AddBalance` and another calling `AverageBalance`, a data race occurs.
The reason is that `balances := c.balances` (same for a slice) creates a new slice that has the same length and the same capacity and is backed by the same array as `c.balances` ⚠️.
- There are two leading solutions to prevent this: **protect the whole function**, or work on a **deep copy** of the actual data.

## #71: Misusing sync.WaitGroup

- Consider the example below:
    ```go
    wg := sync.WaitGroup{}
    var v uint64
    for i := 0; i < 3; i++ {
        go func() {
            wg.Add(1)
            atomic.AddUint64(&v, 1)
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Println(v)
    ```
- If we run this example, we get a non-deterministic value: the code can print any value from 0 to 3. Also, if we enable the `-race `flag, Go will even catch a data race. How is this possible, given that we are using the `sync/atomic` package to update v? What’s wrong with this code?
- The problem is that `wg.Add(1)` is called within the newly created goroutine, not in the parent goroutine. Hence, there is no guarantee that we have indicated to the wait group that we want to wait for three goroutines before calling `wg.Wait()` 🤷.
- When dealing with goroutines, it’s crucial to remember that the execution **isn’t deterministic without synchronization**. For example, the following code could print either `ab or ba`:
    ```go
    go func() {
        fmt.Print("a")
    }()
    go func() {
        fmt.Print("b")
    }()
    ```
- Both goroutines can be assigned to different threads, and there’s no guarantee which thread will be executed first.
- The CPU has to use a **memory fence** (also called a memory barrier) to ensure order. Go provides different synchronization techniques for implementing memory fences: for example, `sync.WaitGroup` enables a **happens-before** relationship between `wg.Add` and `wg.Wait`.
- There are two options to fix our issue:
  - First, we can call `wg.Add` before the loop with 3: `wg.Add(3)`.
  - Or, second, we can call `wg.Add` during each loop iteration **before spinning** up the child goroutines.

## #72: Forgetting about sync.Cond

- The example in this section implements a donation goal mechanism: an application that raises alerts whenever specific goals are reached. We will have one goroutine in charge of incrementing a balance (an updater goroutine). In contrast, other goroutines will receive updates and print a message whenever a specific goal is reached (listener goroutines).
    ```go
    // Listener goroutines
    f := func(goal int) {
        donation.mu.RLock()
        for donation.balance < goal {
            donation.mu.RUnlock()
            donation.mu.RLock()
        }
        fmt.Printf("$%d goal reached\n", donation.balance)
        donation.mu.RUnlock()
    }
    ```
- The main issue—and what makes this a terrible implementation—is the **busy loop**. Each listener goroutine keeps looping until its donation goal is met, which wastes a lot of CPU cycles and makes the CPU usage gigantic 🤒.
- If we think about **signaling in Go**, we should consider **channels**. So, let’s try another version using the channel primitive:
    ```go
    // Listener goroutines
    f := func(goal int) {
        for balance := range donation.ch {
            if balance >= goal {
                fmt.Printf("$%d goal reached\n", balance)
                return
            }
        }
    }
    ```
    - A message sent to a channel is received by only one goroutine.
    - The default distribution mode with multiple goroutines receiving from a **shared channel** is **round-robin**. It can change if one goroutine isn’t **ready** to receive messages; in that case, Go distributes the message to the next available goroutine.
    - Only a **channel closure** event can be **broadcast** to multiple goroutines. But here we don’t want to close the channel, because then the updater goroutine couldn’t send messages.
    - Another issue is that the listener goroutines return whenever their donation goal is met. Hence, the updater goroutine has to know when all the listeners stop receiving messages to the channel. Otherwise, the channel will eventually become **full** and **block** the sender.
- Ideally, we need to find a way to repeatedly broadcast notifications whenever the balance is updated to multiple goroutines. Fortunately, Go has a solution: `sync.Cond`:
    ```go
    // Listener goroutines
    f := func(goal int) {
        donation.cond.L.Lock()
        for donation.balance < goal {
            donation.cond.Wait()
        }
        fmt.Printf("%d$ goal reached\n", donation.balance)
        donation.cond.L.Unlock()
    }
    // Updater goroutine
    for {
        time.Sleep(time.Second)
        donation.cond.L.Lock()
        donation.balance++
        donation.cond.L.Unlock()
        donation.cond.Broadcast() // wakes all the goroutines waiting on the condition.
    }
    ```
- The call to `Wait` must happen within a critical section, which may sound odd 😶‍🌫️ Won’t the lock prevent other goroutines from waiting for the same condition? Actually, the implementation of `Wait` is the following:
    1. Unlock the mutex 🌝.
    2. Suspend the goroutine, and wait for a notification.
    3. Lock the mutex when the notification arrives.

> Let’s also note one possible drawback when using sync.Cond. When we send a notification—for example, to a chan struct—even if there’s no active receiver, the message is buffered, which guarantees that this notification will be received eventually. Using `sync.Cond` with the Broadcast method wakes all goroutines currently waiting on the condition; if there are none, the notification will be **missed**.

🧠 `Signal()` vs. `Broadcast()`:

> We can wake a single goroutine using Signal() instead of Broadcast(). In terms of semantics, it is the same as sending a message in a chan struct in a non-blocking fashion:
```go
ch := make(chan struct{})
select {
    case ch <- struct{}{}:
    default:
}
```

## #73: Not using errgroup

- `golang.org/x` is a repository providing extensions to the standard library. The *sync* sub-repository contains a handy package: `errgroup`.
- It exports a single `WithContext` function that returns a `*Group` struct given a context. This struct provides **synchronization**, **error propagation**, and *context cancellation* for a group of goroutines and exports only two methods:
    - `Go` to trigger a call in a new goroutine.
    - `Wait` to block until all the goroutines have completed. It returns the first non-nil error, if any.
```go
func handler(ctx context.Context, circles []Circle) ([]Result, error) {
    results := make([]Result, len(circles))
    g, ctx := errgroup.WithContext(ctx)
    for i, circle := range circles {
        i := i
        circle := circle
        g.Go(func() error {
            result, err := foo(ctx, circle)
            if err != nil {
                return err
            }
            results[i] = result
            return nil
        })

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```
- This solution is inherently more straightforward as we don’t have to rely on extra concurrency primitives, and the `errgroup.Group` is sufficient to tackle our use case.
- When want to return an error, if any. Hence, there’s no point in waiting until the second and third calls are complete.
  - Using `errgroup.WithContext` creates a shared context used in all the parallel calls. Because the first call returns an error in 1ms, it will cancel the context and thus the other goroutines. So, we won’t have to wait 5 seconds to return an error 👍. This is another benefit when using `errgroup`.

## #74: Copying a sync type

- The `sync` package provides basic synchronization primitives such as mutexes, condition variables, and wait groups. For all these types, there’s a hard rule to follow: they should never be **copied** ⚠️.
    ```go
    type Counter struct {
        mu sync.Mutex
        counters map[string]int
    }

    func NewCounter() Counter {
        return Counter{counters: map[string]int{}}
    }

    func (c Counter) Increment(name string) {
        c.mu.Lock()
        defer c.mu.Unlock()
        c.counters[name]++
    }
    ```
- If we run this example, it raises a data race !
- The problem in our `Counter` implementation is that the **mutex is copied**. Because the receiver of `Increment` is a value, whenever we call `Increment`, it performs a copy of the `Counter` struct, which also copies the mutex. Therefore, the increment isn't done in a shared critical section.
- `sync` types shouldn’t be copied. This rule applies to the following types:
  - sync.Cond
  - sync.Map
  - sync.Mutex
  - sync.RWMutex
  - sync.Once
  - sync.Pool
  - sync.WaitGroup
- We may face the issue of unintentionally copying a `sync` field in the following conditions:
    - Calling a method with a value receiver (as we have seen)
    - Calling a function with a `sync` argument
    - Calling a function with an argument that contains a `sync` field.
