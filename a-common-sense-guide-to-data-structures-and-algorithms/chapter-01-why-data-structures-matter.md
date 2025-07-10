# Why Data Structures Matter

- Depending on how you choose to **organize** your data, your program may run faster or slower by orders of magnitude.

## The Array: The Foundational Data Structure

- One of the most basic data structures in computer science.
- Many data structures are used in four basic ways, which we refer to as **operations**. These operations are:
  - Read
  - Insert
  - Delete
  - Search

## Measuring Speed

- When we measure how “fast” an operation takes, we do **not** refer to how fast the operation takes in terms of **pure time**, but instead in how many **steps** it takes 🎯.
- Measuring the speed of an operation in terms of time is undependable, since the time will always change depending on the **hardware** it is run on.
- Measuring the speed of an operation is also known as measuring its **time complexity**.

## Reading

- Reading from an array is an efficient operation, since the computer can read any index by jumping to any memory address in **one step**.

## Searching

- Searching, though, is tedious, since the computer has no way to jump to a particular value.
- For N cells in an array, linear search would take a maximum of **N steps**.

## Insertion

- The efficiency of inserting a new piece of data into an array depends on **where within the array** you’re inserting it.
- Inserting at the end of the array takes just one step. But there’s one hitch, we need an extra **memory allocation**.
- Inserting to the middle of the array need to shift pieces of data to make room for what we’re inserting, leading to additional steps.
- The **worst-case** scenario for insertion into an array - that is, the scenario in which insertion takes the most steps - is when we insert data at the **beginning** of the array.
- We can say that insertion in a worst-case scenario can take **N + 1** steps for an array containing `N` elements.

## Deletion

- Like insertion, the **worst-case** scenario of deleting an element is deleting the **very first** element of the array. This is because index 0 would become empty, and we’d have to shift all the remaining elements to the left to fill the gap.
- We can say then, that for an array containing `N` elements, the maximum number of steps that deletion would take is **N steps**.

## Sets: How a Single Rule Can Affect Efficiency

- A set is a data structure that does **not allow duplicate** values to be contained within it.
- Reading / searching a set is exactly the same as reading / searching an array.
- Insertion, however, is where arrays and sets diverge.
  - Every insertion into a set first requires a search.
  - In the worst-case scenario, where we’re inserting a value at the beginning of a set, the computer needs to search N cells to ensure that the set doesn’t already contain that value, another N steps to shift all the data to the right, and another final step to insert the new value. ▶️ That’s a total of 2N + 1 steps.

## Exercises:

> 1. For an array containing 100 elements, provide the number of steps the following operations would take:

    a. Reading -> 1
    b. Searching for a value not contained within the array -> 100
    c. Insertion at the beginning of the array -> 101
    d. Insertion at the end of the array -> 1
    e. Deletion at the beginning of the array -> 100
    f. Deletion at the end of the array -> 1

> 2. For an array-based set containing 100 elements, provide the number of steps the following operations would take:

    a. Reading -> 1
    b. Searching for a value not contained within the array -> 100
    c. Insertion of a new value at the beginning of the set -> 201
    d. Insertion of a new value at the end of the set -> 101
    e. Deletion at the beginning of the set -> 100
    f. Deletion at the end of the set -> 1

> 3. Normally the search operation in an array looks for the first instance of a given value. But sometimes we may want to look for every instance of a given value. For example, say we want to count how many times the value “apple” is found inside an array. How many steps would it take to find all the “apples”? Give your answer in terms of N

    N
