# Chapter 8 Concurrency: Foundations

## #55: Mixing up concurrency and parallelism

- In a parallel implementation of a coffee shop, every part of the system is **independent**. The coffee shop should serve consumers twice as fast. <p align="center"><img src="./assets/parallelism.png" width="400px" height="auto"></p>
- With this new design, we don’t make things **parallel**. But the overall structure is affected: we split a given role into two roles, and we introduce another queue. Unlike parallelism, which is about **doing the same thing multiple times at once**, concurrency is about **structure**. <p align="center"><img src="./assets/concurrency.png" width="400px" height="auto"></p>
- We have increased the level of parallelism by introducing more machines. Again, the structure hasn’t changed; it remains a three-step design. But **throughput** should increase because the level of **contention** for the coffee-grinding threads should decrease.
- ▶️ With this design, we can notice something important: **concurrency enables parallelism**. Indeed, concurrency provides a **structure** to solve a problem with parts that may be **parallelized**.
- 🧠 In summary, concurrency and parallelism are different. Concurrency is about structure, and we can change a sequential implementation into a concurrent one by introducing different steps that separate **concurrent threads** can tackle. Meanwhile, parallelism is about execution, and we can use it at the step level by adding more parallel threads.

## #56: Thinking concurrency is always faster

- As Go developers, we can’t create threads directly, but we can create **goroutines**, which can be thought of as application-level threads.
- However, whereas an OS thread is **context-switched** on and off a CPU core by the **OS**, a goroutine is context-switched on and off an OS thread by the **Go runtime**.
- Goroutines start with a small stack size of **2 KB** (as of Go 1.4 and later), which can dynamically grow and shrink as needed.
- Context switching a goroutine versus a thread is about **80% to 90% faster**, depending on the architecture.
- The Go [scheduler](http://mng.bz/N611) uses the following terminology:
  - G—Goroutine
  - M—OS thread (stands for machine)
  - P—CPU core (stands for processor)
- Each OS thread (`M`) is assigned to a CPU core (`P`) by the OS scheduler. Then, each goroutine (`G`) runs on an M.
- The `GOMAXPROCS` variable defines the limit of `M`s in charge of executing user-level code simultaneously.
- A goroutine has a simpler lifecycle than an OS thread. It can be doing one of the following:
  - **Executing**  - The goroutine is scheduled on an M and executing its instructions.
  - **Runnable** - The goroutine is waiting to be in an executing state.
  - **Waiting** - The goroutine is stopped and pending something completing, such as a system call or a synchronization operation (such as acquiring a mutex).
- The Go runtime handles two kinds of **queues**: one **local** queue per `P` and a **global** queue shared among all the `Ps`.
<p align="center"><img src="./assets/go-scheduler.png" width="600px" height="auto"></p>

- Every sixty-first execution, the Go scheduler will check whether goroutines from the global queue are available. If not, it will check its local queue. Meanwhile, if both the global and local queues are empty, the Go scheduler can pick up goroutines from other local queues.
  - This principle in scheduling is called **work stealing**, and it allows an **underutilized** processor to actively look for another processor’s goroutines and steal some.
-  Since Go 1.14, the Go scheduler is now **preemptive**: when a goroutine is running for a specific amount of time *(10 ms)*, it will be marked preemptible and can be context-switched off to be replaced by another goroutine. This allows a **long-running job** to be forced to **share CPU** time.
- 🧠 If the workload that we want to parallelize is too **small**, meaning we’re going to compute it too **fast**, the benefit of distributing a job across cores is destroyed 🤒:
  -  The time it takes to create a goroutine and have the scheduler execute it is much too high compared to directly merging a tiny number of items in the current goroutine.
  -  Although goroutines are lightweight and faster to start than threads, we can still face cases where a workload is too small.
- 📑 So, where should we go from here? We must keep in mind that **concurrency isn’t always faster** and shouldn’t be considered the default way to go for all problems.
  - First, it makes things more complex. Also, modern CPUs have become incredibly efficient at executing **sequential** code and **predictable** code.
  - For example, a **superscalar processor** can parallelize instruction execution over a single core with high efficiency.

## #57: Being puzzled about when to use channels or mutexes

<p align="center"><img src="./assets/mutex-vs-channels.png" width="400px" height="auto"></p>

- **Synchronization** is enforced with **mutexes** but not with any channel types (not with buffered channels). Hence, in general, synchronization between parallel goroutines should be achieved via mutexes.
- Conversely, in general, concurrent goroutines have to **coordinate and orchestrate**. For example, if `G3` needs to aggregate results from both `G1` and `G2`, `G1` and `G2` need to signal to `G3` that a new intermediate result is available. This coordination falls under the scope of **communication** — therefore, **channels**.
- Regarding concurrent goroutines, there’s also the case where we want to transfer the ownership of a resource from one step (`G1` and `G2`) to another (`G3`); for example, if `G1` and `G2` are enriching a shared resource and at some point, we consider this job as complete. Here, we should use **channels** to **signal** that a specific resource is ready and handle the ownership transfer.
- Mutexes and channels have different semantics. Whenever we want to **share** a state or **access a shared resource**, mutexes ensure exclusive access to this resource.
- Conversely, channels are a mechanic for **signaling** with or without data (chan struct{} or not).
- **Coordination** or **ownership** transfer should be achieved via **channels**.
- It’s important to know whether goroutines are parallel or concurrent because, in general, we need **mutexes** for **parallel goroutines** and **channels** for **concurrent** ones.

## #58: Not understanding race problems

### Data races vs. race conditions

- Data race occurs when two or more goroutines **simultaneously** access the **same memory** location and **at least one is writing**.
- Here is an example where two goroutines increment a shared variable:
    ```go
    i := 0
    go func() {
        i++
    }()
    go func() {
        i++
    }()
    ```
- The first option is to make the increment operation **atomic**, meaning it’s done in a single operation. This prevents entangled running operations:
    ```go
    var i int64
    go func() {
        atomic.AddInt64(&i, 1)
    }()
    go func() {
        atomic.AddInt64(&i, 1)
    }()
    ```
- An atomic operation **can’t be interrupted**, thus preventing two accesses at the same time.
- Or we can use a **mutex** to ensures that at most one goroutine accesses a so-called **critical section**.
    ```go
    i := 0
    mutex := sync.Mutex{}
    go func() {
        mutex.Lock()
        i++
        mutex.Unlock()
    }()
    go func() {
        mutex.Lock()
        i++
        mutex.Unlock()
    }()
    ```
- Another possible option is to prevent sharing the same memory location and instead favor communication across the goroutines. For example, we can create a **channel** that each goroutine uses to produce the value of the increment:
    ```go
    i := 0
    ch := make(chan int)
    go func() {
        ch <- 1
    }()
    go func() {
        ch <- 1
    }()
    i += <-ch
    i += <-ch
    ```
- Each goroutine sends a notification via the channel that we should increment i by 1. The parent goroutine collects the notifications and increments `i`. Because it’s the only goroutine writing to `i`, this solution is also **free of data races**.
- Does a **data-race-free** application necessarily mean a deterministic result? Let’s explore this question with another example:
    ```go
    i := 0
    mutex := sync.Mutex{}
    go func() {
        mutex.Lock()
        defer mutex.Unlock()
        i = 1
    }()
    go func() {
        mutex.Lock()
        defer mutex.Unlock()
        i = 2
    }()
    ```
- Depending on the execution order, `i` will eventually equal either 1 or 2. This example doesn’t lead to a **data race**. But it has a **race condition** ‼️
- A **race condition** occurs when the behavior depends on the **sequence** or the **timing** of events that can’t be controlled. Here, the timing of events is the goroutines’ execution order.

### The Go memory model

- The [Go memory model](https://golang.org/ref/mem) is a specification that defines the conditions under which a read from a variable in one goroutine can be guaranteed to happen after a write to the same variable in a different goroutine.
- Let’s examine these guarantees:
  1. Creating a goroutine happens before the goroutine’s execution begins. Therefore, reading a variable and then spinning up a new goroutine that writes to this variable doesn’t lead to a data race:
    ```go
    i := 0
    go func() {
        i++
    }()
    ```
  2. Conversely, the exit of a goroutine isn’t guaranteed to happen before any event. Thus, the following example has a data race:
    ```go
    i := 0
    go func() {
        i++
    }()
    fmt.Println(i)
    ```
  3. A send on a channel happens before the corresponding receive from that channel completes. In the next example, a parent goroutine increments a variable before a send, while another goroutine reads it after a channel read:
    ```go
    i := 0
    ch := make(chan struct{})
    go func() {
        <-ch
        fmt.Println(i)
    }()
    i++
    ch <- struct{}{}
    ```
  4. Closing a channel happens before a receive of this closure. The next example is similar to the previous one, except that instead of sending a message, we close the channel:
    ```go
    i := 0
    ch := make(chan struct{})
    go func() {
        <-ch
        fmt.Println(i)
    }()
    i++
    close(ch)
    ```
  5. A receive from an **unbuffered** channel happens before the send on that channel completes.
    ```go
    i := 0
    ch := make(chan struct{})
    go func() {
        i = 1
        <-ch
    }()
    ch <- struct{}{}
    fmt.Println(i)
    ```

## #59: Not understanding the concurrency impacts of a workload type

- It it important to classify a **workload** in the context of a concurrent application. Let’s illustrate this alongside one concurrency pattern: **worker pooling**.
- Doing so involves creating workers (goroutines) of a **fixed size** that poll tasks from a **common channel**:
    ```go
    func read(r io.Reader) (int, error) {
        var count int64
        wg := sync.WaitGroup{}
        var n = 10 // define the pool size.
        ch := make(chan []byte, n) // create a channel with the same capacity as the pool
        wg.Add(n)
        for i := 0; i < n; i++ {
            go func() {
                defer wg.Done()
                for b := range ch {
                    v := task(b)
                    atomic.AddInt64(&count, int64(v))
                }
            }()
        }
        for {
            b := make([]byte, 1024)
            // Read from r to b
            ch <- b
        }
        close(ch)
        wg.Wait()
        return int(count), nil
    }
    ```
- If the workload is **I/O-bound**, the answer mainly depends on the external system. How many concurrent accesses can the system cope with if we want to maximize **throughput**?
- If the workload is **CPU-bound**, a best practice is to rely on `GOMAXPROCS`. `GOMAXPROCS` is a variable that sets the number of OS threads allocated to running goroutines. By default, this value is set to the number of logical CPUs.
- Let’s take a concrete example and say that we will run our application on a four-core machine: thus Go will instantiate four OS threads where goroutines will be executed. At first, things may not be ideal: we may face a scenario with four CPU cores and four goroutines but only one goroutine being executed 😐:
<p align="center"><img src="./assets/one-goroutine-running.png" width="500px" height="auto"></p>

- Eventually, given the work-stealing concept we already described, `P1` may steal goroutines from the local `P0` queue. However, since one of the main goals of the Go scheduler is to optimize resources (here, the distribution of the goroutines), we should end up in such a scenario given the nature of the workloads.
<p align="center"><img src="./assets/at-most-two-goroutines-running.png" width="500px" height="auto"></p>

- This scenario is still not optimal, because at most two goroutines are running. If there are enough resources in the machine, eventually, the OS should move M2 and M3 as shown: <p align="center"><img src="./assets/optimal-goroutines-gomaxproc.png" width="500px" height="auto"></p>

- ⚠️ There is no guarantee about when this situation will happen. This global picture cannot be designed and requested by us (Go developers), However, as we have seen, we can enable it with favorable conditions in the case of CPUbound workloads: having a worker pool based on `GOMAXPROCS`.
- Last but not least, let’s bear in mind that we should validate our assumptions via **benchmarks** in most cases. **Concurrency isn’t straightforward**, and it can be pretty easy to make hasty assumptions that turn out to be invalid 🙃.

## #60: Misunderstanding Go contexts

- A Context carries a **deadline**, a **cancellation signal**, and **other value**s across API boundaries.

### Deadline

```go
func (h publishHandler) publishPosition(position flight.Position) error {
    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
    defer cancel() //
    return h.pub.Publish(ctx, position)
}
```
What’s the rationale for calling the `cancel` function as a `defer` function?
    - Internally, `context.WithTimeout` creates a **goroutine** that will be retained in memory for 4 seconds or until `cancel` is called.
    - Therefore, calling `cancel` as a `defer` function means that when we exit the parent function, the context will be canceled, and the goroutine created will be stopped.
    - It’s a **safeguard** so that when we return, we don’t leave retained objects in memory.

### Cancellation signals

- A possible approach is to use `context.WithCancel`, which returns a context (first variable returned) that will cancel once the `cancel` function (second variable returned) is called:
    ```go
    func main() {
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()
        go func() {
            CreateFileWatcher(ctx, "foo.txt")
        }()
        // ...
    }
    ```
- When `main` returns, it calls the `cancel` function to cancel the context passed to `CreateFileWatcher` so that the file descriptor is **closed gracefully**.

### Context values

- A context conveying values can be created this way:
    ```go
    ctx := context.WithValue(parentCtx, "key", "value")
    ```
- Just like `context.WithTimeout`, `context.WithDeadline`, and `context.WithCancel`, `context.WithValue` is created from a **parent** context.
- We can access the value using the `Value` method:
    ```go
    ctx := context.WithValue(context.Background(), "key", "value")
    fmt.Println(ctx.Value("key"))
    ```
- The key and values provided are `any` types. Indeed, for the value, we want to pass `any `types . But why should the key be an empty interface as well and not a **string**, for example ❓
  - That could lead to **collisions**: two functions from different packages could use the **same string** value as a key. Hence, the latter would **override** the former value.
  - Consequently, a best practice while handling context keys is to create an **unexported** custom type:
    ```go
    package provider
    type key string
    const myCustomKey key = "key"
    func f(ctx context.Context) {
        ctx = context.WithValue(ctx, myCustomKey, "foo")
    // ...
    }
    ```
- Use cases:
  - For example, if we use tracing, we may want different subfunctions to share the **same correlation ID**.
  - Another example is if we want to implement an HTTP middleware.

### Catching a context cancellation

- `context.Context` type exports a `Done` method that returns a receive-only notification channel: `<-chan struct{}`. This channel is closed when the work associated with the context should be canceled.
- `context.Context` exports an `Err` method that returns nil if the `Done` channel isn’t yet closed. Otherwise, it returns a **non-nil** error explaining why the Done channel was closed: for example:
    - A `context.Canceled` error if the channel was canceled
    - A `context.DeadlineExceeded` error if the context’s deadline passed.
- Let’s see a concrete example in which we want to keep receiving messages from a channel. Meanwhile, our implementation should be context aware and return if the provided context is done:
    ```go
    func handler(ctx context.Context, ch chan Message) error {
        for {
            select {
            case msg := <-ch:
                // Do something with msg
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
    ```
- Within a function that receives a context conveying a possible cancellation or timeout, the action of receiving or sending a message to a channel **shouldn’t** be done in a **blocking** way. For example, in the following function, we send a message to a channel and receive one from another channel:
    ```go
    func f(ctx context.Context) error {
        // ...
        ch1 <- struct{}{}
        v := <-ch2
        // ...
    }
- The problem with this function is that if the context is canceled or times out, we may have to wait until a message is sent or received, without benefit. Instead, we should use `select` to either wait for the channel actions to complete or wait for the context cancellation:
    ```go
    func f(ctx context.Context) error {
        // ...
        select {
            case <-ctx.Done():
                return ctx.Err()
            case ch1 <- struct{}{}:
        }
        select {
            case <-ctx.Done():
                return ctx.Err()
            case v := <-ch2:
        // ...
        }
    }
    ```
- With this new version, if `ctx` is canceled or times out, we return immediately, without blocking the channel send or receive.
- 👍 In general, a function that users wait for should take a context, as doing so allows upstream callers to decide when calling this function should be aborted.
- 👍 When in doubt about which context to use, we should use `context.TODO()` instead of passing an empty context with `context.Background`. `context.TODO()` returns an empty context, but semantically, it conveys that the context to be used is either unclear or not yet available.
