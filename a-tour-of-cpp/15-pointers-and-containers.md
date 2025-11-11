# Start Pointers

## auto_ptr

- introduced in the C++98 standard library
- automatically manage the memory allocated for an object on the heap. When an `auto_ptr` goes out of scope, it automatically deletes the object it owns.
- has a unique ownership model.
- designed for managing single objects and did not support arrays.
- had several limitations and drawbacks. For example, it did not provide strong safety guarantees, especially when multiple pointers were involved.
  -  STL containers require that their contents exhibit “normal” copying behavior, so containers of auto_ptr aren’t allowed
- was deprecated in C++11 and subsequently removed from the C++ Standard Library in C++17. The C++ community recommended using safer smart pointers like std::unique_ptr, std::shared_ptr, and std::weak_ptr introduced in C++11.