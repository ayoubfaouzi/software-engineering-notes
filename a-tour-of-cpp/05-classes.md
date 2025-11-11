# Classes

- A **class hierarchy** is a set of classes ordered in a lattice created by derivation (e.g., : public).
- Objects are constructed **bottom up** (base first) by constructors and destroyed **top down** (derived first) by destructors.

## Non Member Functions

- also known as a **free function**, is a function that is not a member of any class or namespace.
- It is a standalone function that can be declared and defined outside of any class. Non-member functions are not associated with any specific class and are not called using an object instance.

## Concrete Types

- The basic idea of concrete classes is that they behave *just like built-in types*.
- The defining characteristic of a concrete type is that **its representation is part of its definition**.
- They represent specific, fully defined types that provide complete implementations for all of their methods. In C++, a concrete type typically corresponds to a class that:
    - Has **no pure virtual functions**: All member functions are either fully implemented within the class or inherited and implemented from base classes.
    - Can be **instantiated**: You can create objects (instances) of a concrete type.
    - Represents **complete** functionality: Concrete types are intended to be used as-is, with no need for further specialization or extension.
- Example:
```cpp
    class complex:
        double re, im;
        public:
            complex(double r, double i) : re{i}, im{i} {}
            complex(double r) : re{r}, im{0} {}
            complex() : re{0}, im{0} {}
            complex(complex z) : re{z.r}, im{z.im} {} // copy constructor

            double real() const { return re; }
            double imag() const { return im; }
            void real(double r) { re = r };
            void imag(double i) { im = i };

            complex operator+=(complex z) {                        // define a += b
                re += z.real()
                im += z.imag()
                return *this;
            }
            complex operator-=(complex z) {
                re -= z.real()
                im -= z.imag()
                return *this;
            }

            complex operator +=(complex a, complex b) { return a+=b };      // define a + b
            complex operator -=(complex a, complex b) { return a-=b };
            complex operator -(complex a) return {-a.real(), -a.imag() };    // unary minus
            complex operator∗(complex a, complex b) { return a∗=b; }
            complex operator/(complex a, complex b) { return a/=b; }
            bool operator==(complex a, complex b) { return a.real()==b.real() && a.imag()==b.imag(); } // equal
            bool operator!=(complex a, complex b) { return !(a==b); } // not equal

complex a {2.3}; // construct {2.3,0.0} from 2.3
complex b {1/a};// means operator/(complex{1},a)
```

- The technique of acquiring resources in a constructor and releasing them in a destructor, known as **Resource Acquisition Is Initialization** or *RAII*.
- **Concrete classes** – especially classes with small representations – are much like built-in types:
  - We define them as local variables, access them using their names, copy them around, etc.
  - Classes in class hierarchies are different: we tend to allocate them on the **free store** using `new`, and we access them through **pointers** or **references**

## Initializer List

- `std::initializer_list` is a standard library template that provides a way to pass a **fixed-size** list of elements to a function or constructor.
- It allows you to initialize objects or containers using a list of values enclosed in curly braces `{}`. This feature was introduced in `C++11`.

```cpp
Vector(std::initializer_list<double>); // initialize with a list of doubles
Vector::Vector(std::initializer_list<double> lst) // initialize with a list
    : elem{new double[lst.size()]}, sz{static_cast<int>(lst.size())}
{
    copy(lst.begin(),lst.end(),elem); // copy from lst into elem (§13.5)
}
```

## Class Invariant

- Refer to certain conditions or properties that **must be true** for all **instances** (objects) of a class.
- These conditions define the state of the object and should be maintained throughout the object's lifetime.
- Class invariants are important for ensuring the correctness and consistency of the object's behavior.

## Virtual Functions

- Virtual functions enable dynamic **polymorphism**, allowing a derived (subclass) class to **override** the implementation of a function defined in a base (superclass) class.
- Without "**virtual**" you get "**early binding**". Which implementation of the method is used gets decided at compile time based on the type of the pointer that you call through.
- With "**virtual**" you get "**late binding**". Which implementation of the method is used gets decided at run time based on the type of the pointed-to object - what it was originally constructed as. This is not necessarily what you'd think based on the type of the pointer that points to that object.
- A class that provides the interface to a variety of other classes is often called a **polymorphic type**.
- `override` after member functions:
  - It shows the reader of the code that "this is a virtual method, that is **overriding** a virtual method of the base class."
  - The compiler also knows that it's an override, so it can "check" that you are not altering/adding new methods that you think are overrides.
- **Virtual Function Table** (vtbl):
  - When a class contains one or more virtual functions, the compiler automatically creates a `vtable` for that class.
  - The `vtable` is essentially an **array of pointers** to the virtual functions of the class.
  - Each entry in the `vtable` corresponds to a virtual function, and it points to the function's implementation for that specific class.
  - In every object of a class with virtual functions, there is a hidden pointer called the **vtable pointer** (often referred to as *vptr*). This pointer points to the `vtable` associated with the class of the object.
  - The `vptr` is automatically set up by the **constructor** of the class to point to the correct vtable.

## Pure Virtual Functions

```cpp
class Container {
    public:
        virtual double& operator[](int) = 0; // pure virtual function
        virtual int size() const = 0; // const member function
        virtual ~Container() {} // destructor
}

class Vector_container : public Container { // Vector_container implements Container
    public:
        Vector_container(int s) : v(s) { } // Vector of s elements
        ~Vector_container() {} //  ~Vector() is implicitly invoked by ~Vector_container()
        double& operator[](int i) override { return v[i]; }
        int size() const override { return v.size(); }
    private:
        Vector v;
};
```
- also known as an abstract function, is a virtual function that has no implementation in the base class.
- Instead, it is declared with the `= 0` syntax in its declaration.
- A class containing at least one pure virtual function is known as an **abstract class**, and cannot be instantiated directly. Instead, it serves as a base for other classes, providing an interface that derived classes must implement.
- As is common for abstract classes, `Container` does not have a **constructor**. After all, it does not have any data to initialize. On the other hand, Container does have a **destructor** and that destructor is virtual, so that classes derived from `Container` can provide implementations.
