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