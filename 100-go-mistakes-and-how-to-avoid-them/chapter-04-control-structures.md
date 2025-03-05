# Chapter 4: Control structures

## #30: Ignoring the fact that elements are copied in range loops

- In Go, everything we **assign** is a **copy**:
    - If we assign the result of a function returning a **struct**, it performs a **copy** of that struct.
    - If we assign the result of a function returning a **pointer**, it performs a **copy** of the memory address.
- It’s crucial ⚠️ to keep this in mind to avoid common mistakes, including those related to `range` loops. Indeed, when a `range` loop iterates over a data structure, it performs a
**copy** of each element to the value variable.
- So, what if we want to update the slice elements? There are two main options:
    ```go
    for i := range accounts {
        accounts[i].balance += 1000
    }
    for i := 0; i < len(accounts); i++ {
        accounts[i].balance += 1000
    }
    ```
- Another option is to keep using the `range` loop and access the value but modify the slice type to a slice of account **pointers**:
    ```go
    accounts := []*account{
        {balance: 100.},
        {balance: 200.},
        {balance: 300.},
    }
    for _, a := range accounts {
        a.balance += 1000
    }
    ```
    - 👎 Iterating over a slice of pointers may be **less efficient** for a CPU because of the lack of **predictability** (CPU caches).

## #31: Ignoring how arguments are evaluated in range loops

- Consider the example below:
    ```go
    s := []int{0, 1, 2}
    for range s {
        s = append(s, 10)
    }
    ```
- When using a `range` loop, the provided expression is evaluated only once, **before** the **beginning** of the loop.
- In this context, *evaluated* means the provided expression is copied to a **temporary** variable, and then `range` iterates over this variable. In this example, when the `s` expression is evaluated, the result is a **slice copy**:
<p align="center"><img src="./assets/range-slice-copy.png" width="300px" height="auto"></p>

- The behavior is **different** with a classic for `loop`.
- The same logic applies to **channels** regarding how the `range` expression is evaluated.
    ```go
    ch1 := make(chan int, 3)
    go func() {
        ch1 <- 0
        ch1 <- 1
        ch1 <- 2
        close(ch1)
    }()

    ch2 := make(chan int, 3)
    go func() {
        ch2 <- 10
        ch2 <- 11
        ch2 <- 12
        close(ch2)
    }()

    ch := ch1
    for v := range ch {
        fmt.Println(v)
        ch = ch2
    }
    ```
  - The expression provided to `range` is a `ch` channel pointing to `ch1`. Hence, `range` evaluates `ch`, performs a **copy to a temporary** variable, and iterates over elements from this channel. Despite the `ch = ch2` statement, range keeps iterating over `ch1`, **not** `ch2.`
- In **arrays**, the `range` expression is also evaluated **before** the **beginning** of the loop, what is assigned to the temporary loop variable is a **copy** of the array.
- Let’s see this principle in action with the following example that updates a specific array index during the iteration:
    ```go
    a := [3]int{0, 1, 2}
    for i, v := range a {
        a[2] = 10
        if i == 2 {
            fmt.Println(v)
        }
    }
    ```
    - This code updates the last index to `10`. However, if we run this code, it does not print `10`; it prints `2`, instead.
    - The loop doesn’t update the copy; it updates the **original** array ‼️
    - If we want to print the actual value of the last element, we can do so in two ways:
      - By accessing the element from its **index**: `fmt.Println(a[2])`.
      - Using an array pointer: `for i, v := range &a`.
        - We assign a copy of the array pointer to the temporary variable used by `range`. But because both pointers **reference** the **same array**.
        - Doesn’t lead to copying the whole array, which may be something to keep in mind in case the array is **significantly large** 💡.

## #32: Ignoring the impact of using pointer elements in range loops

- 💡 If we store **large** structs, and these structs are **frequently mutated**, we can use pointers instead to **avoid a copy** and an insertion for each mutation.
- We will consider the following two structs:
    ```go
    // A Store that holds a map of Customer pointers
    type Store struct {
        m map[string]*Customer
    }

    // A Customer struct representing a customer
    type Customer struct {
        ID string
        Balance float64
    }
    ```
- The following method iterates over a slice of `Customer` elements and stores them in the `m` map:
    ```go
    func (s *Store) storeCustomers(customers []Customer) {
        for _, customer := range customers {
            s.m[customer.ID] = &customer
        }
    }
    ```
- Iterating over the customers slice using the `range` loop, regardless of the number of elements, creates a **single** customer variable with a **fixed** address ⚠️. We can verify this by printing the pointer address during each iteration:
    ```go
    func (s *Store) storeCustomers(customers []Customer) {
        for _, customer := range customers {
            fmt.Printf("%p\n", &customer)
            s.m[customer.ID] = &customer
        }
    }
    >>>
    0xc000096020
    0xc000096020
    0xc000096020
    ```
- We can overcome this issue by: forcing the creation of a **local variable** in the loop’s scope (`current := customer`) or **creating a pointer** referencing a slice element via its **index** (`customer := &customers[i]`).
- Both solutions are fine. Also note that we took a slice data structure as an input, but the problem
would be similar with a map.

## #33: Making wrong assumptions during map iterations

### Ordering

- Regarding ordering, we need to understand a few fundamental behaviors of the map data structure:
    - It doesn’t keep the data **sorted by key** (a map isn’t based on a binary tree).
    - It doesn’t **preserve the order** in which the data was added.
- But can we at least expect the code to print the keys in the order in which they are currently stored in the map ? No, not even this 😮‍💨.
- However, let’s note that using packages from the **standard library** or **external libraries** can lead to different behaviors. For example, when the `encoding/json` package **marshals** a map into `JSON`, it reorders the data **alphabetically** by keys, regardless of the insertion order.

### Map insert during iteration

- Consider the following example:
    ```go
    m := map[int]bool{
        0: true,
        1: false,
        2: true,
    }
    for k, v := range m {
        if v {
            m[10+k] = true
        }
    }
    fmt.Println(m) // The result of this code is unpredictable
    ```
To understand the reason, we have to read what the Go specification says about a new map entry during an iteration:

> If a map entry is created during iteration, it may be produced during the iteration or skipped. The choice may vary for each entry created and from one iteration to the next.

Hence, when an element is added to a map during an iteration, it may be produced during a follow-up iteration, or it may not ⚠️.

👍 One solution is to create a copy of the map, like so: `m2 := copyMap(m)` and update `m2` instead.

## 34: Ignoring how the break statement works

- One essential rule to keep in mind is that a `break` statement terminates the execution of the **innermost** `for`, `switch`, or `select` statement.
- So how can we write code that breaks the loop instead of the `switch` statement? The most idiomatic way is to use a label:
    ```go
    loop:
        for i := 0; i < 5; i++ {
            fmt.Printf("%d ", i)
            switch i {
                default:
                case 2:
                    break loop // Not a fancy goto statement !
            }
        }
    ```
- 📔 We can also use `continue` with a label to go to the next iteration of the labeled loop.

## #35: Using defer inside a loop

- Consider the following example:
    ```go
    func readFiles(ch <-chan string) error {
        for path := range ch {
            file, err := os.Open(path)
            if err != nil {
                return err
            }
            defer file.Close()
            // Do something with file
        }
        return nil
    ```
- The `defer` calls are executed not during each loop iteration but when the `readFiles` function returns. If `readFiles` doesn’t return, the file descriptors will be kept open forever, causing **leaks**.
- So, what are the options if we want to keep using `defer`?
    1. We have to **create another surrounding function** around `defer` that is called during each iteration. For example, we can implement a `readFile` function holding the logic for each new file path received:
        ```go
        func readFile(path string) error {
            file, err := os.Open(path)
            if err != nil {
                return err
            }
            defer file.Close()
            // Do something with file
            return nil
        }
        ```
    2. Another approach could be to make the `readFile` function a **closure**:
        ```go
        func readFiles(ch <-chan string) error {
            for path := range ch {
                err := func() error {
                    // ...
                    defer file.Close()
                    // ...
                }()
                if err != nil {
                    return err
                }
            }
            return nil
        }
        ````
