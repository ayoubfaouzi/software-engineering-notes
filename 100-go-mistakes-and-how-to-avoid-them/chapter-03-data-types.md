# Chapter 3: Data types

## #17: Creating confusion with octal literals

- In Go, an integer literal starting with `0` is considered an **octal integer** (base 8).
- Octal integers are useful in different scenarios. For instance, suppose we want to open a file using `os.OpenFile`. This function requires passing a permission as a `uint32`. If we want to match a Linux permission, we can pass an octal number for readability instead of a base 10 number:
    ```go
    file, err := os.OpenFile("foo", os.O_RDONLY, 0644)
    ```
- Using `0o` as a prefix instead of only `0` means the same thing. However, it can help make the code clearer.
- Finally, we can also use an underscore character (`_`) as a separator for **readability**. For
example, we can write 1 billion this way: `1_000_000_000`. We can also use the underscore character with other representations (for example, `0b00_00_01`).
- In summary, Go handles **binary**, **hexadecimal**, **imaginary**, and **octal** numbers.
  - Octal numbers start with a 0. However, to improve readability and **avoid potential mistakes** for future code readers, make octal numbers explicit using a `0o` prefix.

## #18: Neglecting integer overflows

- Suppose we want to initialize an `int32` to its maximum value and then increment it. What should be the behavior of this code?
    ```go
    var counter int32 = math.MaxInt32
    counter++
    fmt.Printf("counter=%d\n", counter)
    ```
- This code compiles and doesn’t panic at run time. However, the `counter++` statement generates an integer overflow: `counter=-2147483648` ‼️
- An **integer overflow** occurs when an arithmetic operation creates a value outside the range that can be represented with a given number of bytes.
- Because an `int32` is a **signed** integer, the bit on the left represents the integer’s sign: `0` for positive, `1` for negative. If we increment this integer, there is no space left to represent the new value. Hence, this leads to an integer overflow.
- In Go, an integer overflow that can be detected at **compile** time generates a compilation error. However, at **run time**, an integer **overflow** or **underflow** is **silent**; this does not lead to an application panic ⚠️.
- How can we detect an integer overflow during an **addition**? The answer is to reuse `math.MaxInt`:
    ```go
    func AddInt(a, b int) int {
        if a > math.MaxInt-b {
            panic("int overflow")
        }
        return a + b
    }
    ```
- **Multiplication** is a bit more complex to handle. We have to perform checks against the minimal integer, `math.MinInt`:
    ```go
    func MultiplyInt(a, b int) int {
        if a == 0 || b == 0 {
            return 0
        }
        result := a * b
        if a == 1 || b == 1 {
            return result
        }
        if a == math.MinInt || b == math.MinInt {
            panic("integer overflow")
        }
        if result/b != a {
            panic("integer overflow")
        }
        return result
    }
    ```

## #19: Not understanding floating points

- To avoid bad surprises, we need to know that floating-point arithmetic is an approximation of real arithmetic.
  - 💡 Let’s take the `float64` type as an example. Note that there’s an **infinite** number of real values between `math.SmallestNonzeroFloat64` (the `float64` minimum) and `math.MaxFloat64` (the `float64` maximum).
  - Conversely, the `float64` type has a **finite** number of bits: **64**.
  - Because making infinite values fit into a finite space isn’t possible, we have to work with **approximations**. Hence, we may lose **precision**. The same logic goes for the `float32` type.
- ⚠️ Using the `==` operator to **compare** two floating-point numbers can lead to **inaccuracies**:
  -  Instead, we should compare their difference to see if it is less than some small error value.
  -  For example, the [testify](https://github.com/stretchr/testify) library has an `InDelta` function to assert that two values are within a given delta of each other.
- ⚠️ The result of floating-point calculations depends on the actual **processor**.
  - Most processors have a *floating-point unit* (FPU) to deal with such calculations. There is no guarantee that the result executed on one machine will be the **same on another machine** with a **different FPU**.
  - Comparing two values using a delta can be a solution for implementing valid tests across different machines.
- ⚠️ Also note that the **error can accumulate** in a **sequence** of floating-point operations:
  - Keep in mind that the **order** of floating-point calculations can affect the **accuracy** of the result.
- ⚠️ When performing a **chain** of **additions** and **subtractions**, we should group the operations to add or subtract values with a **similar order of magnitude** before adding or subtracting those with magnitudes that aren’t close.
```go
func f1(n int) float64 {
    result := 10_000.
    for i := 0; i < n; i++ {
        result += 1.0001
    }
    return result
}
func f2(n int) float64 {
    result := 0.
    for i := 0; i < n; i++ {
        result += 1.0001
    }
    return result + 10_000. // Because f2 adds 10,000, in the end it produces more accurate results than f1.
}
```
- ⚠️ When performing floating-point calculations involving addition, subtraction, multiplication, or division, we have to complete the **multiplication** and **division** operations **first** to get better **accuracy**.

## #20: Not understanding slice length and capacity

- In Go, a **slice** is **backed** by an **array**. That means the slice’s data is stored **contiguously** in an array data structure.
- A slice also handles the logic of adding an element if the backing array is **full** or shrinking the backing array if it’s almost **empty**.
- The **length** is the number of elements the slice contains, whereas the **capacity** is the number of elements in the backing array.
- ⚠️ Accessing an element outside the length range is **forbidden**, even though it’s already allocated in memory.
- Adding an element to a full slice (length == capacity) leads to creating a new backing array with a **new capacity**, **copying** all the elements from the previous array, and updating the slice pointer to the new array.
- 💡 In Go, a slice grows by **doubling** its size until it contains *1,024* elements, after which it grows by *25%*.
- What happens with **slicing**? Slicing is an operation done on an array or a slice, providing a half-open range; the first index is included, whereas the second is excluded:
  ```go
    s1 := make([]int, 3, 6)
    s2 := s1[1:3]
    ```
  - When `s2` is created by slicing `s1`, both slices **reference** the s**ame backing array**. However, `s2` starts from a different index.
  - If we update `s1[1]` or `s2[0]`, the change is made to the same array, hence, visible in both slices.
  - If we append an element to `s2`, the shared backing array is modified, but only the length of `s2` changes.
  - if we keep appending elements to `s2` until the backing array is full, `s1` and `s2` will reference **two different arrays**. As `s1` is still a three-length, six-capacity slice, it still has some available buffer, so it keeps referencing the initial array. Also, the new backing array was made by copying the initial one from the first index of `s2`.

## #21: Inefficient slice initialization

- Consider the following example:
    ```go
    func convert(foos []Foo) []Bar {
        bars := make([]Bar, 0)
        for _, foo := range foos {
            bars = append(bars, fooToBar(foo))
        }
        return bars
    }
    ```
- This logic of creating another array because the current one is **full** is repeated multiple times when we add a third element, a fifth, a ninth, and so on.
  - Assuming the input slice has *1,000* elements, this algorithm requires allocating 10 backing arrays and copying more than 1,000 elements in total from one array to another.
  -  This leads to additional effort for the *GC* to clean all these temporary backing arrays.
- There are two different options for this:
  - The first option is to reuse the same code but allocate the slice with a **given capacity**:
    - Internally, Go preallocates an array of n elements. Therefore, adding up to n elements means reusing the **same backing array** and hence reducing the number of **allocations** drastically.
    ```go
    func convert(foos []Foo) []Bar {
        n := len(foos)
        bars := make([]Bar, 0, n)
        for _, foo := range foos {
            bars = append(bars, fooToBar(foo))
        }
        return bars
    }
    ```
  - The second option is to allocate bars with a given length:
    ```go
    func convert(foos []Foo) []Bar {
        n := len(foos)
        bars := make([]Bar, n)
        for i, foo := range foos {
            bars[i] = fooToBar(foo)
        }
        return bars
    }
    ```
    - Because we initialize the slice with a length, `n` elements are already allocated and initialized to the zero value of Bar. Hence, to set elements, we have to use, not `append` but `bars[i]`.
    - This approach is faster because we avoid **repeated calls** to the built-in `append` function, which has a small **overhead** compared to a **direct assignment**.

## #22: Being confused about nil vs. empty slices

```go
func main() {
    var s []string
    log(1, s)
    s = []string(nil)
    log(2, s)
    s = []string{}
    log(3, s)
    s = make([]string, 0)
    log(4, s)
}

func log(i int, s []string) {
    fmt.Printf("%d: empty=%t\tnil=%t\n", i, len(s) == 0, s == nil)
}
```
- This example prints the following:
    ```sh
    1: empty=true nil=true
    2: empty=true nil=true
    3: empty=true nil=false
    4: empty=true nil=false
    ```
- ▶️ All the slices are **empty**, meaning the length equals 0. Therefore, a `nil` slice is also an **empty** slice. However, only the first two are `nil` slices.
- 👍 If a function returns a slice, we **shouldn’t** do as in other languages and **return a non-nil** collection for **defensive** reasons.
  - Because a `nil` slice doesn’t require any **allocation**, we should favor returning a `nil` slice instead of an **empty** slice.
- A `nil` slice is (json) marshaled as a `null` element, whereas a **non-nil**, empty slice is marshaled as an **empty array**.
- `reflect.DeepEqual` returns false if we compare a nil and a non-nil empty slice
- All in all:
    - `var s []string` if we aren’t sure about the final length and the slice can be empty.
    - `[]string(nil)` as syntactic sugar to create a nil and empty slice.
    - `make([]string, length)` if the future length is known.
    - `[]string{}`, should be avoided if we initialize the slice without elements.

## #23: Not properly checking if a slice is empty

- We mentioned in the previous section that an empty slice has, by definition, a length of zero. Meanwhile, nil slices are always empty. Therefore, by checking the length of the slice, we cover all the scenarios:
    - If the slice is nil, `len(operations) != 0` is **false**.
    - If the slice isn’t nil but empty, `len(operations) != 0` is also **false**.
- Hence, checking the length is the best option to follow as we can’t always control the approach taken by the functions we call (by checking if the return != `nil`).
- 👍 When returning slices, it should make neither a semantic nor a technical **difference** if we return a `nil` or **empty** slice. Both **should mean the same thing** for the **callers**.
  - This principle is the same with **maps**. To check if a map is empty, check its length, not whether it’s `nil`.

## #24: Not making slice copies correctly

- To use `copy` effectively, it’s essential to understand that the number of elements copied to the destination slice corresponds to the **minimum** between:
    - The source slice’s length
    - The destination slice’s length
- If we want to perform a complete `copy`, the destination slice must have a length **greater** than or **equal** to the source slice’s length. Here, we set up a length based on the source slice:
    ```go
    src := []int{0, 1, 2}
    dst := make([]int, len(src))
    copy(dst, src)
    ```
- `copy` built-in function isn’t the only way to copy slice elements. There are different alternatives, the best known being probably the following, which uses `append`:
    ```go
    src := []int{0, 1, 2}               // using copy is more idiomatic and, therefore, easier to
    dst := append([]int(nil), src...)   // understand, even though it takes an extra line.

    ```

## #25: Unexpected side effects using slice append

- Consider the following example:
    ```go
    s1 := []int{1, 2, 3}
    s2 := s1[1:2]
    s3 := append(s2, 10)
    // All the slices are backed by the same array
    ```
<p align="center"><img src="./assets/append-slice-problem.png" width="300px" height="auto"></p>

- All the slices are backed by the same array 😮‍💨.
  - Because `s2` in not full, the `append` function adds the element by updating the backing array and returning a slice having a length incremented by 1.
- The `s1` slice’s content was **modified**, even though we did not update `s1[2]` or `s2[1]` **directly**. We should keep this in mind to avoid unintended consequences ⚠️.
- Therefore, if we print all the slices, we get this output: `s1=[1 2 10], s2=[2], s3=[2 10]`.
- If we want to protect against such side effects in function calls:
    - The first is to pass a **copy** of the slice and then construct the resulting slice.
      - 👎 Makes the code more complex to read and adds an extra copy.
    - The second option can be used to limit the range of potential side effects to the first two elements only. This option involves the so-called **full slice expression**: `s[low:high:max]`.
      - 👍 This statement creates a slice similar to the one created with `s[low:high]`, except that the resulting slice’s capacity is equal to `max - low`.
- ▶️ When using slicing, we must remember that we can face a situation leading to unintended side effects. If the resulting slice has a **length** **smaller** than its **capacity**, `append` can **mutate** the original slice.

## #26: Slices and memory leaks

### Leaking capacity

- Consider the example below:
    ```go
    func consumeMessages() {
        for {
            msg := receiveMessage()
            // Do something with msg
            storeMessageType(getMessageType(msg))
            // After a new loop iteration, msg is no longer used.
            // However, its backing array will still be used by msg[:5]
        }
    }
    func getMessageType(msg []byte) []byte {
        return msg[:5]
    }
    ```
- The slicing operation on msg using `msg[:5]` creates a five-length slice. However, its capacity remains the same as the initial slice. The remaining elements are still allocated in memory, even if eventually msg is **not referenced**.
- What can we do to solve this issue? We can make a slice **copy** instead of slicing msg:
    ```go
    func getMessageType(msg []byte) []byte {
        msgType := make([]byte, 5)
        copy(msgType, msg)
        return msgType
    }
    ```
- Because we perform a copy, `msgType` is a **five-length**, **five-capacity** slice regardless of the size of the message received. Hence, we only store 5 bytes per message type.
- ⚠️ Using the **full slice expression** isn’t a valid option (unless a future update of Go tackles this). The whole backing array still lives in memory 😢.
- 👍 As a rule of thumb, remember that slicing a large slice or array can lead to potential **high memory consumption**. The remaining space won’t be reclaimed by the GC, and we can keep a large backing array despite using only a few elements. Using a slice **copy** is the solution to prevent such a case.

### Slice and pointers

- Consider the example below:
    ```go
    type Foo struct {
        v []byte
    }

    func main() {
        foos := make([]Foo, 1_000)
        printAlloc()

        for i := 0; i < len(foos); i++ {
            foos[i] = Foo{ v: make([]byte, 1024*1024) }
        }
        printAlloc()

        two := keepFirstTwoElementsOnly(foos)
        runtime.GC()
        printAlloc()

        runtime.KeepAlive(two) // keep a ref to the two variable after the GC so that it won’t be collected
    }
    func keepFirstTwoElementsOnly(foos []Foo) []Foo {
        return foos[:2]
    }
    ```
- 🎯 It’s essential to keep this rule in mind when working with **slices**:
  - if the element is a **pointer** or a **struct** with pointer fields, the elements won’t be reclaimed by the GC.
- What can we do to solve this issue? We can create a **copy** of the slice. The second option if we want to keep the underlying capacity of 1,000 elements, which is to mark the slices of the remaining elements **explicitly** as `nil`.

## #27: Inefficient map initialization

- When a `map` grows, it doubles its number of **buckets**. What are the conditions for a map to grow?
    - The average number of items in the buckets (called the *load factor*) is greater than a constant value. This constant equals **6.5** (but it may change in future versions because it’s internal to Go).
    - Too many buckets have **overflowed** (containing more than **eight** elements).
<p align="center"><img src="./assets/maps-internals.png" width="400px" height="auto"></p>

- When a `map` **grows**, all the keys are dispatched again to all the buckets. This is why, in the worst-case scenario, inserting a key can be an *O(n)* operation, with `n` being the
total number of elements in the `map`.
- Like **slices**, we can use the make built-in function to provide an **initial size** when creating a `map`. For example, if we want to initialize a `map` that will contain 1 million elements, it can be done this way:
    ```go
    m := make(map[string]int, 1_000_000)` // ask Go runtime to allocate a map with room for at least 1m elements.
    ```
- By specifying a size, we provide a hint about the number of elements expected to go into the `map`. Internally, the map is created with an appropriate number of buckets to store 1 million elements. This saves a lot of **computation** time because the `map` won’t have to create buckets on the fly and handle **rebalancing buckets**.

## #28: Maps and memory leaks

- Consider the example below:
    ```go
    n := 1_000_000
    m := make(map[int][128]byte)
    printAlloc()

    for i := 0; i < n; i++ {
        m[i] = randBytes()
    }
    printAlloc()

    for i := 0; i < n; i++ {
        delete(m, i)
    }
    runtime.GC()
    printAlloc()

    runtime.KeepAlive(m)
    >>>>>>>>>
    0 MB
    461 MB
    293 MB
    ```
- At first, the heap size is minimal. Then it grows significantly after having added 1 million elements to the `map`. But if we expected the heap size to decrease after removing all the elements, this isn’t how maps work in Go.
- The reason is that the number of **buckets in a map cannot shrink**. Therefore, removing elements from a `map` doesn’t impact the number of existing buckets 😮‍💨; it just zeroes the slots in the buckets. 💡 A `map` can only grow and have more buckets; it never shrinks !
- One solution could be to **re-create a copy** of the current map at a regular pace. For example, every hour, we can build a new `map`, copy all the elements, and release the previous one. The main drawback of this option is that following the copy and until the next garbage collection, we may **consume twice** the current **memory** for a short period.
- Another solution would be to change the `map` type to store an array pointer: `map[int]*[128]byte`:
  - It doesn’t solve the fact that we will have a significant **number of buckets**; however, each bucket entry will reserve the **size of a pointer** for the value instead of 128 bytes.
  - Also as an optimization, if a key or a value is **over 128 bytes**, Go won’t store it directly in the `map` bucket. Instead, Go stores a pointer to reference the key or the value.

## #29: Comparing values incorrectly

- It’s essential to understand how to use `==` and `!=` to make comparisons effectively. We can use these operators on operands that are comparable:
    - **Booleans**: Compare whether two Booleans are equal.
    - **Numerics** (int, float, and complex types): Compare whether two numerics are equal.
    - **Strings**: Compare whether two strings are equal.
    - **Channels**: Compare whether two channels were created by the same call to make or if both are `nil`.
    - **Interfaces**: Compare whether two interfaces have identical dynamic types and equal dynamic values or if both are `nil`.
    - **Pointers**: Compare whether two pointers point to the same value in memory or if both are `nil`.
    - **Structs and arrays**: Compare whether they are composed of similar types.
- With these behaviors in mind, what are the options if we have to compare two **slices**, two **maps**, or two **structs** containing **non-comparable** types?
  - If we stick with the standard library, one option is to use run-time **reflection** with the `reflect` package.
  - `reflect.DeepEqual` reports whether two elements are deeply equal by **recursively** traversing two values.
- However, using `reflect.DeepEqual` has two catches ⚠️
  - It makes the distinction between an **empty** and a **nil** collection.
  - Because this function uses reflection, which introspects values at run time to discover how they are  formed, it has a **performance penalty** (is about 100 times slower than `==`).
- If performance is a crucial factor, another option might be to implement our **own comparison method**.
- In the context of unit tests, some other options are possible, such as using external libraries with [go-cmp](https://github.com/google/go-cmp) or [testify](https://github.com/stretchr/testify).
- The standard library has some existing comparison methods. For example, we can use the optimized `bytes.Compare` to compare two slices of bytes. Before implementing a custom method, we need to make sure we don’t reinvent the wheel 🧠.
