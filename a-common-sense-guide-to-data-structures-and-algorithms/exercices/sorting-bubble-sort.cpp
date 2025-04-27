#include <array>
#include <iostream>

void print_array(const std::array<int, 11> &arr) {
  for (const int &value : arr) {
    std::cout << value << " ";
  }

  std::cout << "\n";
}

void bubble_sort(std::array<int, 11> &arr) {

  size_t end = arr.size() - 1;
  bool swapped = false;
  do {
    swapped = false;
    for (size_t i = 0; i < end; i++) {
      if (arr[i] > arr[i + 1]) {
        std::swap(arr[i], arr[i + 1]);
        swapped = true;
      }
    }
    end--;
  } while (swapped);
}

int main() {

  std::array<int, 11> arr = {1, 5, 99, 14, 56, 4, 78, 100, 45, 87, 1};

  std::cout << "Original array" << std::endl;
  print_array(arr);
  bubble_sort(arr);
  std::cout << "Sorted array" << std::endl;
  print_array(arr);

  return 0;
}
