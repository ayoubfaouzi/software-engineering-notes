# Big O in Everyday Code

## Mean Average of Even Numbers

```ruby
def average_of_even_numbers(array)
    sum = 0.0
    count_of_even_numbers = 0

    # Since the loop iterates over each of the N elements, we know that the algorithm takes at least N steps.
    array.each do |number| # In worst case (all the numbers are even), we need 3N steps.
        if number.even? # 1 step.
            sum += number # 1 more step.
            count_of_even_numbers += 1 # 1 more step.
        end
    end

    # Outside of the loop as well. Before the loop, we initialize the two variables and set them to 0 + the division: 3 steps in total.
    return sum / count_of_even_numbers
end
```

The total number of steps is **3N + 3** ➡️ `O(N)`.

## Word Builder

he next example is an algorithm that collects every combination of two character strings built from an array of single characters. Here is a JS implementation:

```js
function wordBuilder(array) {
    let collection = [];
    for(let i = 0; i < array.length; i++) {
        for(let j = 0; j < array.length; j++) {
            if (i !== j) {
                collection.push(array[i] + array[j]);
            }
        }
    }
    return collection;
}
```