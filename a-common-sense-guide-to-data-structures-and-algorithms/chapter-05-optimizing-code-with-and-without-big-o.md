# Optimizing Code with and Without Big O

## Selection Sort

Selection Sort is a comparison-based sorting algorithm. It sorts an array by repeatedly selecting the **smallest** (or largest) element from the unsorted portion and **swapping** it with the first unsorted element. This process continues until the entire array is sorted.

Here’s a JavaScript implementation of Selection Sort:

```js
function selectionSort(array) {
    for(let i = 0; i < array.length - 1; i++) {
        let lowestNumberIndex = i;
        for(let j = i + 1; j < array.length; j++) {
            if(array[j] < array[lowestNumberIndex]) {
                lowestNumberIndex = j;
            }
        }
        if(lowestNumberIndex != i) {
            let temp = array[i];
            array[i] = array[lowestNumberIndex];
            array[lowestNumberIndex] = temp;
        }
    }
return array;
}
```

## The Efficiency of Selection Sort

To put it in a way that works for arrays of all sizes, we’d say that for `N` elements, we make `(N - 1) + (N - 2) + (N - 3) … + 1` comparisons.

As for swaps, though, we only need to make a maximum of one swap per passthrough.

Here’s a side-by-side comparison of **Bubble Sort** and **Selection Sort**:

| N Elements | Max # of Steps in Bubble Sort | Max # of Steps in Selection Sort   |
| ---------- | ----------------------------- | ---------------------------------- |
| 5          | 20                            | 14 (10 comparisons + 4 swaps)      |
| 10         | 90                            | 54 (45 comparisons + 9 swaps)      |
| 20         | 380                           | 199 (180 comparisons + 19 swaps)   |
| 40         | 1560                          | 819 (780 comparisons + 39 swaps)   |
| 80         | 6320                          | 3239 (3160 comparisons + 79 swaps) |

Because Selection Sort takes **roughly half** of N^2 steps, it would seem reasonable that we’d describe the efficiency of Selection Sort as being `O(N2 / 2)`. That is, for N data elements, there are `N2 / 2` steps.

In reality, however, Selection Sort is described in Big O as `O(N2)` 🤒, just like **Bubble Sort**. This is because of a major rule of Big O that I’m now introducing for the first time: *Big O Notation ignores constants* 📌.

This is simply a mathematical way of saying that Big O Notation never includes regular numbers that aren’t an **exponent**. We simply drop these regular numbers from the expression. In our case, then, even though the algorithm takes `N2 / 2` steps, we drop the `“/ 2”` because it’s a regular number, and express the efficiency as `O(N2)`.

## Big O Categories

All the types of Big O we’ve encountered, whether it’s `O(1)`, `O(log N)`, `O(N)`, `O(N2)`, or the types we’ll encounter later in this book, are **general categories** of Big O that are widely different from each other. **Multiplying** or **dividing** the number of steps by a regular number doesn’t make them change to another category 💁.

However, when two algorithms fall under the **same classification** of Big O, it doesn’t necessarily mean that both algorithms have the **same speed**. After all, Bubble Sort is twice as slow as Selection Sort even though both are `O(N2)`. So, while Big O is perfect for contrasting algorithms that fall under different classifications of Big O, when two algorithms fall under the same classification, **further analysis** is required to determine which algorithm is **faster**.

### Significant Steps

In the previous chapters, I alluded to the fact that you’d learn how to determine which steps are significant enough to be counted when expressing the Big O of an algorithm. In our case, then, which of these steps are considered **significant** ❓ Do we care about the **comparisons**, the **printing**, or the **incrementing** of number?

The answer is that **all steps are significant**. It’s just that when we express the steps in Big O terms, we drop the constants and thereby simplify the expression 💁.

Let’s apply this here. If we count all the steps, we have **N comparisons**, **N incrementings**, and **N / 2 printings**. This adds up to **2.5N steps**. However, because we eliminate the constant of **2.5**, we express this as `O(N)`. So, which step was significant? They all were, but by dropping the constant, we effectively focus more on the number of times the loop runs, rather than the exact details of what happens within the loop.

## Exercises

> 1. Use Big O Notation to describe the time complexity of an algorithm that takes 4N + 16 steps.

O(N).

> 2. Use Big O Notation to describe the time complexity of an algorithm that takes 2N2.

o(N^2)

> 3. Use Big O Notation to describe the time complexity of the following function, which returns the sum of all numbers of an array after the numbers have been doubled:

```py
def double_then_sum(array)
    doubled_array = []
    array.each do |number|
        doubled_array << number *= 2
    end
    sum = 0
    doubled_array.each do |number|
        sum += number
    end
return sum
end
```

O(N).

> 4. Use Big O Notation to describe the time complexity of the following function, which accepts an array of strings and prints each string in multiple cases:

```py
def multiple_cases(array)
    array.each do |string|
        puts string.upcase
        puts string.downcase
        puts string.capitalize
    end
end
```

O(N).

> 5. The next function iterates over an array of numbers, and for each number whose index is even, it prints the sum of that number plus every number in the array. What is this function’s efficiency in terms of Big O Notation?

```py
def every_other(array)
    array.each_with_index do |number, index|
        if index.even?
            array.each do |other_number|
                puts number + other_number
            end
        end
    end
end
```

O(N^2).
