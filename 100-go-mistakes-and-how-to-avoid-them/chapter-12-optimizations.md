# Chapter 12 : Optimizations

## #91: Not understanding CPU caches

- When a specific memory location is accessed (for example, by reading a variable), one of the following is likely to happen in the near future:
    - The same location will be referenced again.
    - Nearby memory locations will be referenced.
- The former refers to **temporal locality**, and the latter refers to **spatial locality**. Both are part of a principle called *locality of reference* ⭐.
- Because of spatial locality, the CPU copies what we call a **cache line** instead of copying a single variable from the main memory to a cache.
- If we benchmark these two functions:
  - `sumFoo` which receives a **slice of struct**, and sums the first field of the struct.
  - `sumBar` also computes a sum. But this time, the argument is a **struct containing slices**. <p align="center"><img src="./assets/slice-of-structs-vs-struct-of-slices.png" width="300px" height="auto"></p>
- ➡️ `sumBar` is faster (about 20% on my machine). The main reason is a **better spatial locality** (all the elements of a
are allocated contiguously ) that makes the CPU **fetch fewer cache lines** from memory.
- This example demonstrates how spatial locality can have a substantial impact on performance. To optimize an application, we should **organize data** to get the most value out of each individual cache line.
- **Predictability** refers to the ability of a CPU to anticipate what the application will do to speed up its execution. Let’s see a concrete example where a lack of predictability negatively impacts application performance.
  - Again, let’s look at two functions that sum a list of elements:
    - `linkedList` iterates over a linked list (allocated contiguously) and sums all the values.
    - `sum2` iterates over a slice, one element out of two. <p align="center"><img src="./assets/linked-list-vs-slice.png" width="300px" height="auto"></p>
  - The two data structures have the **same spatial locality**, so we may expect a similar execution time for these two functions. But the function iterating on the slice is significantly faster (about 70% on my machine). What’s the reason? 🤔:
- **Striding** relates to how CPUs work through data. There are three different types of strides:
    - **Unit stride** — All the values we want to access are allocated **contiguously**: for example, a slice of `int64` elements. This stride is **predictable** for a CPU and the most efficient because it requires a minimum number of cache lines to walk through the elements.
    - **Constant stride** — Still predictable for the CPU: for example, a slice that iterates over every two elements. This stride requires more cache lines to walk through data, so it’s less efficient than a unit stride.
    - **Non-unit stride** — A stride the CPU **can’t predict**: for example, a linked list or a slice of pointers. Because the CPU **doesn’t know** whether data is allocated **contiguously**, it won’t fetch any cache lines. <p align="center"><img src="./assets/cpu-striding.png" width="300px" height="auto"></p>
- For `sum2`, we face a constant stride. However, for the linked list, we face a non-unit stride. Even though we know the data is allocated contiguously, the CPU doesn’t know that. Therefore, it can’t predict how to walk through the linked list 🤷.
- Because of the different stride and similar spatial locality, iterating over a linked list is **significantly slower** than a slice of values. We should generally favor unit strides over constant strides because of the better spatial locality. But a non-unit stride cannot be predicted by the CPU regardless of how the data is allocated, leading to negative performance impacts.
- why changing the overall number of columns (seen in *Mistake#89*) impacted the benchmark results. It might sound counterintuitive: because we need to read only the first eight columns, why does changing the total number of columns affect the execution time ❓
  - When these two functions (`calculateSum512` and `513`) are benchmarked each time with a new matrix, we don’t observe **any difference**. However, if we keep reusing the same matrix, `calculateSum513` is about **50% faster** on my machine. The reason lies in CPU caches and how a memory block is copied to a cache line.
  - Now, let’s say the benchmark executes the function with a slice pointing to the same matrix starting at address `0000000000000`. When the function reads `s[0][0]`, the address isn’t in the cache. This block was **already replaced**.
  - Instead of using **CPU caches from one execution** to **another**, the benchmark will lead to more **cache misses**.
  - This type of cache miss is called a **conflict miss**: a miss that wouldn’t occur if the cache wasn’t **partitioned**. All the variables we iterate belong to a memory block whose set index is `00`. Therefore, we use only **one cache set** instead of having a distribution across the entire cache.
  - In this example, this stride is called a **critical stride**: it leads to accessing memory addresses with the same set index that are hence stored to the same cache set.
  - Let’s come back to our real-world example with the two functions calculateSum512 and calculateSum513. The benchmark was executed on a `32 KB` eight-way set-associative L1D cache: 64 sets total. Because a cache line is `64 bytes`, the critical stride equals 64 × 64 bytes = `4 KB`. Four KB of `int64` types represent `512` elements.
  - ➡️ Therefore, we reach a critical stride with a matrix of **512 columns**, so we have a **poor caching distribution**. Meanwhile, if the matrix contains **513 columns**, it doesn’t lead to a critical stride. This is why we observed such a massive difference between the two benchmarks 😵‍💫.

## #92: Writing concurrent code that leads to false sharing

- To illustrate the concept of false sharing, we use two structs, `Input` and `Result`:
    ```go
    type Input struct {
        a int64
        b int64
    }
    type Result struct {
        sumA int64
        sumB int64
    }
    ```
- We spin up two goroutines: one that iterates over each `a` field and another that iterates over each `b` field:
    ```go
    go func() {
        for i := 0; i < len(inputs); i++ {
            result.sumA += inputs[i].a
        }
        wg.Done()
    }()
    go func() {
        for i := 0; i < len(inputs); i++ {
            result.sumB += inputs[i].b
        }
        wg.Done()
    }()
    ```
- Because `sumA` and `sumB` are allocated contiguously, in most cases (seven out of eight), both variables are allocated to the **same memory block**. <p align="center"><img src="./assets/false-sharing-same-block.png" width="300px" height="auto"></p>
- Now, let’s assume that the machine contains two cores. In most cases, we should eventually have two threads scheduled on different cores. So if the CPU decides to copy this memory block to a cache line, it is copied twice: <p align="center"><img src="./assets/cache-line-copy-multi-core.png" width="400px" height="auto"></p>
- Both cache lines are replicated because L1D is per core. Recall that in our example, each goroutine updates its own variable: `sumA` on one side, and `sumB` on the other side.
- Because these cache lines are **replicated**, one of the goals of the CPU is to **guarantee cache coherency**. For example, if one goroutine updates `sumA` and another reads `sumA` (after some synchronization), we expect our application to get the latest value.
- However, our example doesn’t do exactly this. Both goroutines access their own variables, not a shared one. We might expect the CPU to know about this and understand that it **isn’t a conflict**, but this isn’t the case 🤷‍♂️.
- When we write a variable that’s in a cache, the **granularity** tracked by the CPU isn’t the variable: it’s the **cache line**.
- When a cache line is **shared** across **multiple cores** and at least one goroutine is a **writer**, the **entire cache line** is **invalidated**. - This happens even if the updates are logically independent (for example, `sumA` and `sumB`). This is the problem of false sharing, and it degrades performance ⚠️.
- So how do we solve false sharing? There are two main solutions.
    - The first solution is to use the same approach we’ve shown but ensure that `sumA` and `sumB` aren’t part of the same cache line. For example, we can update the `Result` struct to add **padding** between the fields:
        ```go
        type Result struct {
            sumA int64
            _ [56]byte
            sumB int64
        }
        ```
        - Using padding, `sumA` and `sumB` will always be part of different memory blocks and hence **different cache lines**.
        - If we benchmark both solutions (with and without padding), we see that the padding solution is **significantly faster** (about 40% on my machine) 😮‍💨.
    - The second solution is to **rework the structure** of the algorithm. For example, instead of having both goroutines share the same struct, we can make them communicate their local result via channels.

## #93: Not taking into account instruction-level parallelism

- CPU designers stopped focusing solely on **clock speed** to improve CPU performance. They developed multiple optimizations, including **ILP (Instruction-Level Parallelism)**.
- If we have a sequence of 3 instructions:
  - If executed sequentially, this would have taken the following time: `total time = t(I1) + t(I2) + t(I3)`.
  - Thanks to ILP, the total time is the following: `total time = max(t(I1), t(I2), t(I3))`.
- ILP looks 🤹‍♂️, theoretically. But it leads to a few challenges called **hazards**:
  - For example, what if `I3` sets a variable to 42 but `I2` is a conditional instruction. In theory, this scenario should prevent executing `I2` and `I3` in parallel. This is called a **control hazard** or **branching hazard**. In practice, CPU designers solved control hazards using **branch prediction**.
  - For example, if `I1` adds the numbers in registers A and B to C and `I2` adds the numbers in registers C and D to D. Because `I2` depends on the outcome of `I1` concerning the value of register C, the two instructions cannot be executed simultaneously ➡️ **data hazard**.
    - CPU designers have come up with a trick called **forwarding** that basically bypasses writing to a register. This technique doesn’t solve the problem but rather tries to alleviate the effects 🤷.
  - There are also **structural hazards**, when at least two instructions in the pipeline need the **same resource**. As Go developers, we can’t really impact these kinds of hazards.
- Let’s get back to our initial problem and focus on the content of the loop:
    ```go
    s[0]++
    if s[0]%2 == 0 {
        s[1]++
    }
    ```
- If we highlight the hazards between the instructions, we get: <p align="center"><img src="./assets/hazards-between-instructions.png" width="400px" height="auto"></p>
    - The only independent instructions are the `s[0]` check and the `s[1]` increment, so these two instruction sets can be executed in parallel thanks to branch prediction. <p align="center"><img src="./assets/ilp-v1.png" width="300px" height="auto"></p>

- Can we improve our code to minimize the number of data hazards ❓ Let’s write another version that introduces a **temporary variable**:
    ```go
    v := s[0]
    s[0] = v + 1
    if v%2 != 0 {
        s[1]++
    }
    ```
<p align="center"><img src="./assets/hazards-between-instructions_improved.png" width="400px" height="auto"></p>

- The significant difference is regarding the data hazards: the `s[0]` increment step and the check `v` step now depend on the **same instruction** (`read s[0] into v`).
- Why does this matter? Because it allows the CPU to increase the level of parallelism:
<p align="center"><img src="./assets/ilp-v2.png" width="300px" height="auto"></p>

- Despite having the same number of steps, the second version increases how many steps can be executed in parallel: three parallel routes instead of two. Meanwhile, the execution time should be optimized because the longest path has been reduced. If we benchmark these two functions, we see a significant **speed improvement** for the second version (about 20% on my machine), mainly because of ILP.

## #94: Not being aware of data alignment

- Suppose we allocate two variables, an `int32` (32 bytes) and an `int64` (64 bytes):
    ```go
    var i int32
    var j int64
    ```
- Without data alignment, on a 64-bit architecture, these two variables could be allocated as below:
    <p align="center"><img src="./assets/bad-data-alignment.png" width="200px" height="auto"></p>

- The `j` variable allocation could be spread over two words. If the CPU wanted to read `j`, it would require **two memory accesses** instead of one.
- To prevent such a case, a variable’s memory address should be a **multiple** of its own **size**. This is the concept of **data alignment**. In Go, the alignment is guaranteed for common variable types such as: `byte`, `float64`, `complex128`, ..
- During the compilation, the Go compiler adds **padding** to guarantee data alignment:
    ```go
    type Foo struct {
        b1 byte
        _ [7]byte // Added by the compiler
        i int64
        b2 byte
        _ [7]byte // Added by the compiler, because a struct’s size must be a multiple of the word size (8 bytes)
    }
    ```
- Every time a `Foo` struct is created, it requires **24 bytes** in memory, but only **10 bytes** contain data—the remaining 14 bytes are padding ❗
- Because a struct is an **atomic unit**, it will never be **reorganized**, even after a GC; it will always occupy 24 bytes in memory.
- Note that the compiler **doesn’t rearrange** the fields; it only adds padding to guarantee data alignment.
- How can we reduce the amount of memory allocated? The rule of thumb is to reorganize a struct so that its fields are **sorted by type size** in descending order 👍.
- Besides the size overhead of padding, if we created `Foo` variables frequently and they were allocated to the heap, the result would be more frequent GCs, impacting overall application performance ⚠️.
- Speaking of performance, there’s another effect on **spatial locality**:
    ```go
    // Consider the example of iterating over the slice and sums all the i fields.
    for i := 0; i < len(foos); i++ {
        s += foos[i].i
    }
    ```
  - Each gray bar represents 8 bytes of data, and the darker bars are the i variables.
    <p align="center"><img src="./assets/padding-and-spacial-locality.png" width="500px" height="auto"></p>
  - Because each cache line contains more `i` variables, iterating over a slice of `Foo` requires fewer cache lines total.
  - ➡️  Each cache line is more useful because it contains on average **33%** more `i` variables. Therefore, iterating over a `Foo` slice to sum all the int64 elements is more efficient 👍.
