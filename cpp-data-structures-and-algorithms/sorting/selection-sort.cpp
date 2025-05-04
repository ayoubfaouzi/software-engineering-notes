#include <array>
#include <iostream>
#include <utility>

void print_array(const std::array<int, 11> &arr) {

  for (const int &value : arr) {
  std:
    std::cout << value << " ";
  }

  std::cout << "\n";
}

void selection_sort(std::array<int, 11> &arr) {

  for (size_t i = 0; i < arr.size(); i++) {
    int min = i;
    for (size_t j = i + 1; j < arr.size(); j++) {
      if (arr[j] < arr[min]) {
        min = j;
      }
    }

    if (min != i) {
      std::swap(arr[i], arr[min]);
    }
  }
}

int main() {

  std::array<int, 11> arr = {1, 5, 99, 14, 56, 4, 78, 100, 45, 87, 1};

  std::cout << "Original array" << std::endl;
  print_array(arr);
  selection_sort(arr);
  std::cout << "Sorted array" << std::endl;
  print_array(arr);

  return 0;
}
