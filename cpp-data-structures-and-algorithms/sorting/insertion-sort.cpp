#include <array>
#include <iostream>

void print_array(const std::array<int, 11> &arr) {
  for (const int &value : arr) {
    std::cout << value << " ";
  }

  std::cout << "\n";
}

void insertion_sort(std::array<int, 11> &arr) {

  for (size_t i = 0; i < arr.size() - 1; i++) {
    int temp = arr[i + 1];
    size_t j = i;
    while (j >= 0 && arr[j] > temp) {
      arr[j + 1] = arr[j];
      j--;
    }
    arr[j + 1] = temp;
  }
}

int main() {

  std::array<int, 11> arr = {1, 5, 99, 14, 56, 4, 78, 100, 45, 87, 1};

  std::cout << "Original array" << std::endl;
  print_array(arr);
  insertion_sort(arr);
  std::cout << "Sorted array" << std::endl;
  print_array(arr);

  return 0;
}
