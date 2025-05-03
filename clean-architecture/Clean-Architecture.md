# Clean Architecture

notes taken from reading the _Clean Architecture_ book by *Robert C.Martin*.

## Preface

The rules of software architecture are the same, they are independent of every other variable:
    - programming language
    - hardware
    - domain
    - application type.
    - ...

## Introduction

- Getting some code **to work** is not that hard, get it **right** is another matter entirely.
- When you get a software right:
    - 👍 Don't need hordes of programmers to **keep it working**.
    - 👍 Don't need a **massive requirements** doc and huge issue tracking systems.
    - 👍 Require a fraction of human resources to maintain, **efforts** are minimized, **functionality** and flexibility are maximized.

## Chapter 1. What is design and architecture

- The word **architecture** is often used in the context of something at a _high level_ that is divorced from the low-level details, whereas **design** more often seems to imply structures and decisions at a _low level_.
- The low-level details and the high-level structure are all part of the same whole, they form a continuous fabric that defines the shape of the system. You can't have one without the other; indeed **no clear dividing** lines separates them.
- The measure of **design quality** is simply the measure of the **effort required** to meet the needs if the customer.
- Developers buy into a familiar lie 😮:
    - We can **clean it up later**; we just need to get to **market first** but things never do get cleaned up later because market pressures never abate.
    - Making **messy code** makes them go **fast** in the short term and just slows them down in the long term.
- A fact is that making messes is **always slower** then staying clean, no matter which time scale you are using.
- Developers may think that the answer is to start over **from scratch** and redesign the whole system but that could be another mistake due to their **overconfidence**.

## Chapter 2. A tale of two values

- Every software system provides two different values to the stakeholders:
    - **Behavior** or **function**: program should work and behave in a way that makes or saves money for the stakeholders.
    - **Architecture**: software system should be easy to change.
- Both business managers and developers would tend more to advocate function over architecture, but it is the **wrong** attitude.
- Business managers and developers fail to separate those features that **urgent** but **not important**, from those features that truly are **urgent and important**. This failure leads to ignoring the important architecture of the system in favor of the unimportant features of the system.
- If **architecture** comes last, then the system will become ever **more costly** to develop, and eventually **change** wil become practically impossible for part or for all of the system. If that is allowed to happen, it means the software development team did not fight hard enough for what they knew was necessary 😳.

## Chapter 3. Paradigm overview

- **Structured programming**:
    - Shows that the use of **unrestrained** jumps (`goto` statements) is **harmful** to program structure.
    - Replaced those jumps with more familiar `if/then/else` and `do/while/until` constructs.
    - ➡️ _Structured programming imposes discipline on direct transfer of control_.
- **Object oriented programming**:
    - Engineers noticed that the function call stack frame in the *ALGOL* language could be moved to a **heap**, thereby allowing local variables declared by a function to **exist long** after the function returned.
    - The function became a constructor for a class, the local variables became instance variables, and the nested functions became methods. This led inevitably to the discovery of polymorphism through the disciplined use of function pointers.
    - ➡️ _Object-oriented programming imposes discipline on indirect transfer of control_.
- **Functional programming**:
    - Only recently begun to be adopted, was the first to be invented.
    - A foundational notion of lambda calculus is **immutability**, that is, the notion that the values of symbols do not change ▶️ **no assignment statement**.
    - Most functional languages do, in fact, have some means to alter the value of a variable, but only under very strict discipline.
    - ➡️ _Functional programming imposes discipline upon assignment_.
- Each of these paradigms **removes** capabilities from the programmer, none of them adds new ones. Each imposes some kind of extra discipline that is negative in its intent. The paradigms tell us **what not to do**, more than they tell us **what to do**. They removed `goto` statements, function pointers, and assignments.

## Chapter 4. Structured Programming

-  _Dijkstra_ discovered that certain uses of `goto` statements prevent modules from being **decomposed recursively** into smaller and smaller units, thereby preventing use of the **divide-and-conquer** approach necessary for reasonable proofs.
- Other uses of `goto`, however, did not have this problem. Dijkstra realized that these “good” uses of goto corresponded to simple **selection** and **iteration** control structures such as `if/then/else` and `do/while`. Modules that used only those kinds of control structures could be **recursively subdivided into provable units**.
- Nowadays we are all structured programmers, though not necessarily **by choice**. It’s just that our languages don’t give us the option to use **undisciplined direct transfer of control**.
- Some may point to named `breaks` in Java or exceptions as `goto` analogs. In fact, these structures are not the utterly unrestricted transfers of control that older languages like Fortran or COBOL once had. Indeed, even languages
that still support the `goto` keyword often restrict the target to within the **scope of the current function**.
- Structured programming allows modules to be recursively decomposed into **provable units**, which in turn means that modules can be **functionally decomposed**.
- That is, you can take a large-scale problem statement and decompose it into **high-level functions**. Each of those functions can then be decomposed into **lower-level functions**, ad infinitum. Moreover, each of those decomposed functions can be represented using the restricted control structures of structured programming.
- Dijkstra once said, _“Testing shows the presence, not the absence, of bugs.”_
    - A program **can be proven incorrect** by a test, but it **cannot be proven correct**.
    - All that tests can do, after sufficient testing effort, is allow us to deem a program to be **correct enough** for our purposes.
    - ➡️ Software development is **not** a **mathematical endeavor**, even though it seems to manipulate mathematical constructs.
    - ➡️ Rather, **software is like a science**. We show correctness by failing to prove incorrectness, despite our best efforts.
- Structured programming forces us to recursively decompose a program into a set of small provable functions. We can then use tests to try to prove those small provable functions incorrect. If such tests fail to prove incorrectness, then we deem the functions to be correct enough for our purposes.

## Chapter 5. Object-Oriented Programming

-  What is OO?
    - Combination of data and function.
        - :arrow_forward: **absurd**, because programmers were passing data structures into functions long before 1966, when Dahl and Nygaard moved the function call stack frame to the heap and invented OO.
    - A way to model the real world.
        - :arrow_forward: evasive answer at best. What does “modeling the real world” actually mean, and why is it something we would want to do?
    - Proper admixture of these three things: **encapsulation**, **inheritance**, and **polymorphism**.

### Encapsulation?

-  is certainly **not unique** to OO. Indeed, we had perfect encapsulation in C.
- `point.h`:
```C
struct Point;
struct Point* makePoint(double x, double y);
double distance (struct Point *p1, struct Point *p2);
```
- `point.c`:
```c
#include "point.h"
#include <stdlib.h>
#include <math.h>

struct Point {
    double x,y;
};
struct Point* makepoint(double x, double y) {
struct Point* p = malloc(sizeof(struct Point));
    p->x = x;
    p->y = y;
    return p;
}
double distance(struct Point* p1, struct Point* p2) {
    double dx = p1->x - p2->x;
    double dy = p1->y - p2->y;
    return sqrt(dx*dx+dy*dy);
}
```
- The users of `point.h` have no access whatsoever to the members of struct `Point`.
- But then came OO in the form of C++— and the perfect encapsulation of C was **broken**.
- The C++ compiler, for technical reasons (needs to know the size of the instances of each class.), needed the member variables of a class to be declared **in the header file of that class**. So our `Point` program
changed to look like this:
- `point.h`:
```c
class Point {
public:
    Point(double x, double y);
    double distance(const Point& p) const;
private:
    double x;
    double y;
};
```
- `point.cc`:
```c
#include "point.h"
#include <math.h>

Point::Point(double x, double y) : x(x), y(y) {}
double Point::distance(const Point& p) const {
    double dx = x-p.x;
    double dy = y-p.y;
    return sqrt(dx*dx + dy*dy);
}
```
- Clients of the header file `point.h` know about the member variables `x` and `y`! The compiler will prevent access to them, but the client still knows they exist.
- For example, if those member names are changed, the `point.cc` file must be recompiled! **Encapsulation has been broken**.
- Indeed, the way encapsulation is partially **repaired** is by introducing the `public`, `private`, and `protected` keywords into the language. This, however, was a _hack_ necessitated by the technical need for the compiler to see those variables in the header file.
- For these reasons, it is difficult to accept that OO depends on **strong encapsulation**. Indeed, many OO languages have little or no enforced encapsulation.

### Inheritance?

- Inheritance is simply the **redeclaration** of a group of variables and functions within an **enclosing scope**. This is something C programmers were able to do manually long before there was an OO languag.
- `namedPoint.h`:
```c
struct NamedPoint;
struct NamedPoint* makeNamedPoint(double x, double y, char* name);
void setName(struct NamedPoint* np, char* name);
char* getName(struct NamedPoint* np);
```
- `namedPoint.c`:
```c
#include "namedPoint.h"
#include <stdlib.h>

struct NamedPoint {
    double x,y;
    char* name;
};
struct NamedPoint* makeNamedPoint(double x, double y, char* name) {
    struct NamedPoint* p = malloc(sizeof(struct NamedPoint));
    p->x = x;
    p->y = y;
    p->name = name;
    return p;
}

void setName(struct NamedPoint* np, char* name) {
    np->name = name;
}

char* getName(struct NamedPoint* np) {
    return np->name;
}
```
- `main.c`:
```c
#include "point.h"
#include "namedPoint.h"
#include <stdio.h>

int main(int ac, char** av) {
    struct NamedPoint* origin = makeNamedPoint(0.0, 0.0, "origin");
    struct NamedPoint* upperRight = makeNamedPoint(1.0, 1.0, "upperRight");
    printf("distance=%f\n", distance(
        (struct Point*) origin, (struct Point*) upperRight));
}
```
- If you look carefully at the `main` program, you’ll see that the `NamedPoint` data structure acts as though it is a derivative of the `Point` data structure.
- This is because the **order** of the first two fields in `NamedPoint` is the **same** as `Point`.
- We had a trick, but it’s not nearly as **convenient** as **true inheritance**. Moreover, **multiple inheritance** is a considerably more difficult to achieve by such trickery.
- To recap: We can award no point to OO for encapsulation, and perhaps a half-point for inheritance. So far, that’s not such a great score :stuck_out_tongue:.

### Polymorphism?

```c
#include <stdio.h>

void copy() {
    int c;
    while ((c=getchar()) != EOF)
        putchar(c);
}
```
- The function `getchar()` reads from `STDIN`. But which device is STDIN?
- The `putchar()` function writes to `STDOUT`. But which device is that?
- These functions are polymorphic—their behavior depends on the type of `STDIN` and `STDOUT`.
- So how does the call to `getchar()` actually get delivered to the device driver that reads the character?
- The UNIX operating system **requires** that every IO device driver provide five standard functions: `open, close, read, write, and seek`. The signatures of those functions must be **identical** for every IO driver.
- The `FILE` data structure contains five pointers to functions. In our example, it might look like this:
```c
struct FILE {
    void (*open)(char* name, int mode);
    void (*close)();
    int (*read)();
    void (*write)(char);
    void (*seek)(long index, int mode);
};
```
- The IO driver for the console will define those functions and load up a `FILE` data structure with their addresses—something like this:
```c
#include "file.h"

void open(char* name, int mode) {/*...*/}
void close() {/*...*/};
int read() {int c;/*...*/ return c;}
void write(char c) {/*...*/}
void seek(long index, int mode) {/*...*/}

struct FILE console = {open, close, read, write, seek};
```
- Now if `STDIN` is defined as a `FILE*`, and if it points to the console data structure, then `getchar()` might be implemented this way:
```c
extern struct FILE* STDIN;

int getchar() {
    return STDIN->read();
}
```
- In other words, `getchar()` simply calls the function pointed to by the read pointer of the `FILE` data structure pointed to by `STDIN`.
- This simple trick is the basis for all polymorphism in OO. In C++, for example, every **virtual function** within a class has a pointer in a table called a _vtable_, and all calls to virtual functions go through that table. Constructors of derivatives simply load their versions of those functions into the vtable of the object being created.
- The bottom line is that polymorphism is an **application of pointers to functions**.
- Programmers have been using pointers to functions to achieve polymorphic behavior since _Von Neumann_ architectures were first implemented in the late 1940s.
- In other words, OO has provided nothing new. Ah, but that’s not quite correct. OO languages may not have given us polymorphism, but they have made it **much safer and much more convenient**.
- OO is the ability, through the use of polymorphism, to gain **absolute control** over every **source code dependency** in the system. It allows the architect to create a **plugin** architecture, in which modules that contain **high-level policies** are **independent** of modules that contain **low-level details**. The low-level details are relegated to plugin modules that can be deployed and developed independently from the modules that contain **high-level policies**.

## Chapter 6: Functional Programming

- In many ways, the concepts of functional programming predate programming itself. This paradigm is strongly based on the **lambda-calculus** invented by `Alonzo Church` in the 1930s.
- Variables in functional languages **do not vary**.
- Why ? All **race conditions**, **deadlock conditions**, and **concurrent update** problems are due to mutable variables.
- **Immutability** can be practicable, if certain compromises are made.
    - **Segregation of Mutability**: The immutable components perform their tasks in a purely functional way, without using any mutable variables. The immutable components communicate with one or more other components that are not purely functional, and allow for the state of variables to be mutated. <p align="center"><img src="assets/segregation-of-mutability.png" width="400px" height="auto"></p>
        - Since mutating state exposes those components to all the problems of concurrency, it is common practice to use some kind of **transactional** memory to protect the mutable variables from concurrent updates and race conditions. Transactional memory simply treats variables in memory the same way a database treats records on disk. It protects those variables with a transactionor retry-based scheme.
    - **Event Sourcing**: is a strategy wherein we store the **transactions, but not the state**. When state is required, we simply apply all the transactions from the beginning of time.
        - The more memory we have, and the faster our machines are, the less we need mutable state.
        - If we have enough storage and enough processor power, we can make our applications entirely immutable—and, therefore, entirely functional.

## Chapter 7 SRP: The Single Responsibility Principle

- :eyes: does not mean that every module should do just one thing.
- But means instead that _a module should be responsible to one, and only one, actor_.
- My favorite example is the `Employee` class from a payroll application. It has three methods: `calculatePay()`, `reportHours()`, and `save()`. <p align="center"><img src="assets/spr.png"></p>
- This class violates the SRP because those three methods are responsible to three very different actors.
    - The `calculatePay()` method is specified by the accounting department, which reports to the CFO.
    - The `reportHours()` method is specified and used by the human resources department, which reports to the COO.
    - The `save()` method is specified by the database administrators (DBAs), who report to the CTO.
- Many problems occur because we put code that **different actors** depend on into **close proximity**. The SRP says to
separate the code that different actors depend on.

## Chapter 8 OCP: The Open-Closed Principle

- The Open-Closed Principle (OCP) was coined in 1988 by _Bertrand Meyer_. It says: _A software artifact should be open for extension but closed for modification_.
- In other words, the behavior of a software artifact ought to be extendible, without having to modify that artifact.

### A Thought Experiment

- Imagine, for a moment, that we have a system that displays a financial summary on a web page. The data on the page is scrollable, and negative numbers are rendered in red.
- Now imagine that the stakeholders ask that this same information be turned into a report to be printed on a black-and-white printer. The report should be properly paginated, with appropriate page headers, page footers, and column labels. Negative numbers should be surrounded by parentheses.
- Clearly, some new code must be written. But how much old code will have to change?
- A good software architecture would reduce the amount of changed code to the barest minimum. **Ideally, zero**.
- The essential insight here is that generating the report involves **two separate responsibilities**:
    - the calculation of the reported data;
    - and the presentation of that data into a web and printer-friendly form. <p align="center"><img src="assets/ocp.png"></p>
- Classes marked with `<I>` are **interfaces**; those marked with `<DS>` are **data structures**.
- Open arrowheads are using **relationships**. Closed arrowheads are implements or **inheritance** relationships.
- An arrow pointing from class `A` to class `B` means that the source code of class `A` mentions the name of class `B`, but class `B` mentions **nothing about** class `A`. Thus, `FinancialDataMapper` knows about `FinancialDataGateway` through an implements relationship, but `FinancialGateway` knows nothing at all about `FinancialDataMapper`.
- The next thing to notice is that each double line is crossed in **one direction only**. This means that all component relationships are **unidirectional**. These arrows point toward the components that we want to **protect from change**. <p align="center"><img src="assets/components-relationships-uni.png"></p>
- We want to protect the _Controller_ from changes in the _Presenters_. We want to protect the _Presenters_ from changes in the _Views_. We want to protect the _Interactor_ from changes in—well, anything.
- The _Interactor_ is in the position that best conforms to the OCP. Changes to the Database, or the Controller, or the Presenters, or the Views, will have **no impact** on the Interactor. Why should the _Interactor_ hold such a privileged position? Because it contains the **business rules**. The Interactor contains the **highest-level policies** of the application. All the other components are dealing with peripheral concerns. The _Interactor_ deals with the **central concern**.
- This is how the OCP works at the architectural level. Architects separate functionality based on how, why, and when it changes, and then organize that separated functionality into a hierarchy of components. **Higher-level** components in that hierarchy are **protected** from the changes made to **lower-level** components.

## Chapter 9 LSP: The Liskov Substitution Principle

- In 1988, Barbara Liskov wrote the following as a way of defining subtypes:
    - :bulb: _What is wanted here is something like the following substitution property: If for each object o1 of type S there is an object o2 of type T such that for all programs P defined in terms of T, the behavior of P is unchanged when o1 is substituted for o2 then S is a subtype of T_.
- Example: Imagine that we have a class named `License`. This class has a method named `calcFee()`, which is called by the `Billing` application. There are two “subtypes” of License: `PersonalLicense` and `BusinessLicense`. They use different algorithms to calculate the license fee. <p align="center"><img src="assets/lsp-inheritance.png" width="400px" height="auto"></p>
- This design conforms to the LSP because the behavior of the `Billing` application does not depend, in any way, on which of the two subtypes it uses. Both of the subtypes are **substitutable** for the `License` type.
- The canonical example of a violation of the LSP is the famed (or infamous, depending on your perspective) square/rectangle problem: <p align="center"><img src="assets/lsp-square-rectangle.png" width="400px" height="auto"></p>
- In this example, `Square` is not a proper **subtype** of `Rectangle` because the height and width of the `Rectangle` are **independently mutable**; in contrast, the height and width of the `Square` must **change together**. Since the `User` believes it is communicating with a `Rectangle`, it could easily get confused.
- :arrow_forward: A simple violation of **substitutability**, can cause a system’s architecture to be polluted with a significant amount of extra mechanisms.

## Chapter 10 ISP: The Interface Segregation Principle

- The Interface Segregation Principle (ISP) derives its name from the diagram below: <p align="center"><img src="assets/isp.png"></p>
- In the situation illustrated in above, there are several users who use the operations of the `OPS` class. Let’s assume that `User1` uses only `op1`, `User2` uses only `op2`, and `User3` uses only `op3`.
- Now imagine that `OPS` is a class written in a language like Java. Clearly, in that case, the source code of `User1` will **inadvertently depend** on `op2` and `op3`, even though it doesn’t call them. This dependence means that a change to the source code of `op2` in `OPS` will force `User1` to be **recompiled and redeployed**, even though nothing that it cared about has actually changed.
- This problem can be resolved by **segregating the operations into interfaces**: <p align="center"><img src="assets/isp-segregated-ops.png"></p>

### ISP and Language

- Clearly, the previously given description depends critically on language type. **Statically typed languages** like Java force programmers to create declarations that users must `import`, or use, or otherwise `include`. It is these included
declarations in source code that create the source code dependencies that **force recompilation and redeployment**.
- In **dynamically typed** languages like Ruby and Python, such declarations don’t exist in source code. Instead, they are inferred at **runtime**. Thus there are no source code dependencies to force recompilation and redeployment. This is
the primary reason that dynamically typed languages create systems that are **more flexible and less tightly coupled** than statically typed languages.
- :arrow_forward: This fact could lead you to conclude that the ISP is a language issue, rather than an architecture issue.

### ISP and Architecture

- In general, it is **harmful** to depend on modules that contain **more than you need**. This is obviously true for source code dependencies that can force unnecessary recompilation and redeployment—but it is also true at a much higher, **architectural level**.
- Consider, for example, an architect working on a system, S. He wants to include a certain framework, F, into the system. Now suppose that the authors of F have bound it to a particular database, D. So S depends on F. which depends on D. <p align="center"><img src="assets/isp-problematic-architecture.png"></p>.
-  Now suppose that D contains features that F does not use and, therefore, that S does not care about. Changes to those features within D may well force the redeployment of F and, therefore, the redeployment of S. Even worse, a failure of one of the features within D may cause failures in F and S.

## Chapter 11 DIP: Dependency Inversion Principle

- The Dependency Inversion Principle (DIP) tells us that the most flexible systems are those in which source code dependencies refer only to **abstractions, not to concretions**.
- In a statically typed language, like Java, this means that the `use, import,` and `include` statements should **refer only to source modules containing interfaces, abstract classes**, or some kind of abstract declaration.
- The same rule apply for dynamically typed languages, like Ruby and Python. **Source code dependencies should not refer to concrete modules**. However, in these languages it is a bit harder to define what a concrete module is, in particular, it is any module which the functions being called are implemented.
- Clearly treating this idea as a rule is unrealistic, because software systems must depend on many concrete facilities.
    - We tend to ignore the **stable background of operating system** and platform facilities because we know we can rely on them to not change.
    - Example: `String` class in Java is **very stable**.
- Is the the **volatile concrete** elements of our system that we want to **avoid depending on**. Those are the modules that are actively developing, and that are undergoing frequent change.

### Stable Abstractions

- Every change to an abstract interface corresponds to a change to its concrete implementations. Conversely, changes to concrete implementations do not always, or even usually, require changes to the interface that they implement
    - :arrow_forward: Therefore interfaces are less volatile than implementations.
    - A good software architect work hard to **reduce the volatility of interfaces** (stable abstract interfaces).
    - Coding practices:
        - :white_check_mark: Don’t refer to volatile concrete classes.
        - :white_check_mark: Don’t derive doom volatile concrete classes.
        - :white_check_mark: Don’t override concrete functions.

### Factories

- To comply with these rules, the creation of volatile concrete objects requires special handling. In most object-oriented languages, such a Java, we would use an *Abstract Factory* to manage this undesirable factory.
- The diagram below shows the structure:
    - The `Application` uses the `ConcreteImpl` though the `Service` interface.
    - However the `Application` must somehow create instances of the `ConcreteImpl`. To achieve this without creating a source code dependency on the `ConcreteImpl`, the `Application` calls the `makeSvc` method of the `ServiceFactory` interface. This method is implemented by the `ServiceFactoryImpl` class, which derives from `ServiceFactory`. The implementation instantiates the `ConcreteImpl` and returns it as a `Service`.
    - The curved line is an architectural boundary, it separates the abstract from the concrete.
    <p align="center"><img src="assets/factories.png" width="400px" height="auto"></p>
- The way the dependencies cross that curved line in one direction, and toward more abstract entities, will become a new rule that we will call the **Dependency Rule**.

## Chapter 12 Components

- Components are the **unit of deployment**. The smallest entities that can be deployed as part of a system.
    - In Java, they are jar files. In Ruby, they are gem files. In .NET, they are DLLs.
- In the early days of software development, programmers controllers the memory locations and layout of their programs
    - *origin* statement tells the compiler to generate code that will be loaded at address 200.
    - Libraries where kept in source, not in a binary.
    - Devices were slow, memory was expensive and, therefore was limited.
    - :arrow_forward: Compilation took a long time.
    - To shorten compile times, the source code of the function library was separated from the apps.
        - Programmers compiled the function library separately and loaded the binary at a known address.
        - Memory layout looks like this: <p align="center"><img src="assets/memory-layout-old-days.png" width="400px" height="auto"></p>
    - This works fine as long as the app could fit between addresses 0000 and 1777. But soon apps grow to be larger than the space allocated for them

### Relocatability

- The compiler was changed to output binary code that could be relocated in memory by a smart loader.
- The loader would be told where to load the relocatable code (`.reloc` in PE files :smile: ).
- The relocatable code was instrumented with flags that told the loader which parts of the loaded data had to be altered to be loaded at the selected address.
- The compiler was also changed to emit the names of the functions as metadata in the relocatable binary. If a program called a library function, the compiler would emit that name as an **external reference**. If a program defined a library function, the compiler would emit that name as an **external definition**. Then the loader could link the external references to the external definitions once it had determined where it had loaded those definitions.
	- :cool: and the linking loader was born.

### Linkers

The linking loader allowed programmers to divide their programs up onto **separately compilable and loadable segments**. This worked well when relatively small programs were being linked with relatively small libraries.
- Eventually, the linking loaders were **too slow to tolerate**. They had to read dozens, if not hundreds, of binary libraries (stored on slow devices) to resolve the external references.
- As a result, the loading and the linking were separated into two phases. Programmers took the slow part — the part that did that linking — and put it into a separate application called the linker. The output of the linker was a linked relocatable that a relocating loader could load very quickly.
- Compiling each individual module was relatively fast, but compiling all the modules took a bit of time. The linker would then take even more time. Turnaround had gain grown to an hour or more in many cases.
- It seemed as if programmers were doomed to endlessly chase their tail :hushed:.
- By the mid-1990s, the time spent linking had begun to shrink faster than our
ambitions could make programs grow (**thanks to Moore’s law**). In many cases, link time decreased to a matter of seconds.

## Chapter 13: Component Cohesion

### The Reuse/Release Equivalence Principle

- REP is a principle that seems obvious, at least in hindsight. People who want to reuse software components cannot, and will not, do so unless those components are **tracked through a release process** and are given **release numbers**.

### The Common Closure Principle

- Gather into components those classes that **change for the same reasons** and at the **same times**. Separate into different components those classes that **change at different times** and for **different reasons**.
- This is the Single Responsibility Principle restated for **components**.

### The Common Reuse Principle

- Don’t force users of a component to depend on things they don’t need.
- When we depend on a component, we want to make sure we depend **on every class** in that component. Put another way, we want to make sure that the classes that we put into a component are **inseparable** — that it is impossible to depend on some and not on the others. Otherwise, we will be redeploying more components than is necessary, and wasting significant effort.
- The CRP is the generic version of the ISP. The ISP advises us not to depend on classes that have **methods** we don’t use. The CRP advises us not to depend on **components** that have classes we don’t use.

### The Tension Diagram for Component Cohesion

- The edges of the diagram describe the **cost of abandoning** the principle on the opposite vertex.
<p align="center"><img src="assets/cohesion-principles-tension.png" width="400px" height="auto"></p>

## Chapter 14: Component Coupling

### The Acyclic Dependencies Principle

- The _morning after syndrome_ occurs in development environments where many developers are modifying the same source files. Some code that used to work are now no longer working because someone modified something you depends on.
- Two solutions to this problem have evolved: **The weekly build** and the **Acyclic Dependencies Principle (ADP)**.

#### The Weekly Build

- Used to be common in medium-sized projects.
- All the developers ignore each other for the first four days of the week. Then, on Friday, they integrate all their changes and build the system.
- :white_check_mark: Allowing the developers to live in an isolated world for four days out of five.
- :thumbsdown: Large integration penalty that is paid on Friday.
- :arrow_forward: As the project size grows, integration and testing become increasingly harder to do.

#### Eliminating Dependency Cycles

- The solution to this problem is to **partition** the development environment into **releasable components**.
- To make it work successfully, however, you must manage the **dependency structure** of the components. There can be no cycles. If there are cycles in the dependency structure, then the _morning after syndrome_ cannot be avoided.

<p align="center"><img src="assets/typical-components-diagram.png" width="500px" height="auto"></p>

- Regardless of which component you begin at, it is impossible to follow the dependency relationships and wind up back at that component. This structure has no cycles. It is a **directed acyclic graph (DAG)**.
- When it is time to release the whole system, the process proceeds from the bottom up. First the `Entities` component is compiled, tested, and released. Then the same is done for `Database` and `Interactors`. These components are followed by `Presenters`, `View,` `Controllers`, and then `Authorizer`. `Main` goes last. This process is very clear and easy to deal with. We know how to build the system because we understand the dependencies between its parts.

#### The Effect of a Cycle in the Component Dependency Graph

- For example, let’s say that the `User` class in `Entities` uses the `Permissions` class in `Authorizer`. This creates a **dependency cycle**. When such cycles in the dependency graph are created:
 - :bangbang: All of the developers working on any of those components will experience the dreaded _morning after syndrome_.
 - :bangbang: Unit testing and releasing become very difficult and error prone.

#### Breaking the Cycle

- There are two primary mechanisms for doing so:
  - Apply the Dependency Inversion Principle (DIP).
  - Create a new component that both `Entities` and `Authorizer` depend on. Move the class(es) that they both depend on into that new component.

#### The “Jitters”

- The second solution implies that the component structure is **volatile** in the presence of changing requirements.
  - As the app grow, the component dependency structure **jitters and grows**. Thus the dependency structure must always be monitored for cycles.
  - When cycles occur, they must be **broken** somehow. Sometimes this will mean creating new components, making the dependency structure grow.

### Top-Down Design

The component structure cannot be designed from the **top down**.
- As more and more modules accumulate in the early stages of implementation and design:
    - :+1: Keep changes as **localized** as possible, so we start paying attention to the **SRP** and **CCP** and collocate classes that are likely to change together.
    - :+1: Isolate volatile components. We don’t want components that change frequently and for capricious reasons to affect components that otherwise ought to be stable.
- The component dependency structure grows and evolves with the logical design of the system.

### The Stable Dependencies Principle

- Designs cannot be completely **static**. Some volatility is necessary if the design is to be maintained.
- Any component that we expect to be **volatile** should not be depended on by a component that is difficult to change. Otherwise, the volatile component will also be difficult to change.

#### Stability

- One sure way to make a software component difficult to change, is to make lots of other software components depend on it. A component with lots of incoming **dependencies** is **very stable** because it requires a great deal of work to reconcile any changes with all the dependent components.
- _X_, is a **stable** component. Three components depend on _X_, so it has three good reasons **not to change**.
- We say that _X_ is **responsible** to those three components.
- Conversely, _X_ depends on nothing, so it has no external influence to make it change. We say it is **independent**.
```mermaid
flowchart TD
    a[[ ]] --> x[[x]]
    b[[ ]] --> x
    c[[ ]] --> x
```
- _Y_, is a very **unstable** component. No other components depend on _Y_, so we say that it is **irresponsible**.
- _Y_ also has three components that it depends on, so changes may come from three external sources. We say that _Y_ is **dependent**.

```mermaid
flowchart TD
    y[[y]] --> a[[ ]]
    y --> b[[ ]]
    y --> c[[ ]]
```

#### Stability Metrics

:question: How can we measure the stability of a component? Count the number of dependencies that enter and leave that component.
- **Fan-in**: Incoming dependencies. This metric identifies the number of classes outside this component that depend on classes within the component.
- **Fan-out**: Outgoing dependencies. This metric identifies the number of classes inside this component that depend on classes outside the component.
- **I**: Instability: `I = Fan-out % (Fan-in + Fan-out)`. This metric has the range [0, 1]. I = 0 indicates a maximally stable component. I = 1 indicates a maximally unstable component.
>  The SDP says that the I metric of a component should be larger than the I metrics of the components that it depends on. That is, I metrics should decrease in the direction of dependency.

#### Not All Components Should Be Stable

- `Flexible` is a component that we have designed to be easy to change. We want `Flexible` to be unstable.
- However, some developer, working in the component named `Stable`, has hung a dependency on `Flexible`.
- :bangbang: This violates the SDP because the I metric for `Stable` is much smaller than the I metric for `Flexible`.
- :arrow_forward: `Flexible` will no longer be easy to change. A change to `Flexible` will force us to deal with Stable and all its
dependents. <p align="center"><img src="assets/violating-sdp.png" width="300px" height="auto"></p>.
- To fix this problem, we somehow have to **break** the dependence of `Stable` on `Flexible` (DIP for the win).

### The Stable Abstractions Principle

> A component should be as abstract as it is stable.

#### Where Do We Put the High-Level Policy?

- Some software in the system should not change very often. This software represents **high-level architecture and policy decisions**. We don’t want these business and architectural decisions to be **volatile**.
- However, if the high-level policies are placed into stable components, then the source code that represents those policies will be **difficult to change**.
:arrow_forward: This could make the overall architecture inflexible. How can a component that is maximally stable (I = 0) be flexible enough to withstand change :question: ==> The answer is found in the OCP.

#### Introducing the Stable Abstractions Principle

- The Stable Abstractions Principle (SAP) sets up a relationship between **stability** and **abstractness**.
- On the one hand, it says that a **stable** component should also be **abstract** so that its stability does not prevent it from being **extended**.
- On the other hand, it says that an **unstable** component should be **concrete** since it its instability allows the concrete code within it to be easily changed.
- :arrow_forward: If a component is to be stable, it should consist of interfaces and abstract classes so that it can be extended.
- The SAP and the SDP combined amount to the DIP for components. This is true because the SDP says that dependencies should run in the direction of stability, and the SAP says that stability implies abstraction. Thus _dependencies run in the direction of abstraction_.

#### Measuring Abstraction

- The `A` metric is a measure of the abstractness of a component. Its value is simply the ratio of interfaces and abstract classes in a component to the total number of classes in the component.
- The A metric ranges from 0 to 1. A value of 0 implies that the component has no abstract classes at all. A value of 1 implies that the component contains nothing but abstract classes.

#### The Main Sequence

- We are now in a position to define the relationship between stability (I) and abstractness (A).
- If we plot the two “good” kinds of components on this graph, we will find the components that are maximally stable and abstract at the upper left at (0, 1). The components that are maximally unstable and concrete are at the lower right at (1, 0). <p align="center"><img src="assets/zone-of-exclusion.png" width="300px" height="auto"></p>.
- :+1: Strive to be close to the **main sequence** and far from the **zones of exclusion**.

#### The Zone of Pain

- Consider a component in the area of (0, 0). This is a highly stable and concrete component. Such a component is not desirable because it is rigid. It cannot be extended because it is not abstract, and it is very difficult to change because of its stability.
- :arrow_forward: An example of that would be a database schema.

#### The Zone of Uselessness

- Consider a component near (1, 1). This location is undesirable because it is maximally abstract, yet has no dependents.
- :-1: Such components are useless.
- :arrow_forward: An example of that are leftover abstract classes that no one ever implemented. We find them in systems from time to time, sitting in the code base, unused.

## Chapter 15 What Is Architecture?

- The architecture of a software system is the shape given to that system by those who build it. The form of that shape is in the division of that system into components, the arrangement of those components, and the ways in which those components communicate with each other.
- :+1: Good architecture makes the system easy to **understand**, easy to **develop**, easy to **maintain**, and easy to **deploy**. The ultimate goal is to minimize the lifetime **cost** of the system and to maximize programmer **productivity**.

### Development

- ‼️ A software system that is hard to develop is not likely to have a long and healthy lifetime. So the architecture of a system should make that system easy to develop, for the team(s) who develop it.

### Deployment

- ‼️ Unfortunately, deployment strategy is seldom considered during initial development. This leads to architectures that may make the system easy to develop, but leave it very difficult to deploy.

### Operation

- 🛸 The fact that **hardware is cheap** and **people are expensive** means that architectures that impede operation are not
as costly as architectures that impede development, deployment, and maintenance.

### Maintenance

- ‼️ Of all the aspects of a software system, **maintenance** is the most **costly**. The never-ending parade of new features and the inevitable trail of defects and corrections consume vast amounts of human resources.
- The primary cost of maintenance is in **spelunking** and risk. Spelunking is the cost of digging through the existing software, trying to determine the best place and the best strategy to add a new feature or to repair a defect. While
making such changes, the likelihood of creating inadvertent defects is always there, adding to the cost of risk.

### Keeping Options Open

- The way you keep software soft is to **leave as many options open as possible**, for as long as possible. What are the options that we need to leave open? ▶️ They are the details that don’t matter.
All software systems can be decomposed into two major elements: **policy** and **details**.
 - The policy element embodies all the business rules and procedures.  The policy is where the true value of the system lives.
 - The details are those things that are necessary to enable humans, other systems, and programmers to communicate with the policy, but that do not impact the behavior of the policy at all. They include IO devices, databases, web systems, servers, frameworks, communication protocols, and so forth.
- The goal of the architect is to create a shape for the system that recognizes policy as the most essential element of the system while **making the details irrelevant to that policy**. This allows decisions about those details to be delayed and deferred. The longer you wait to make those decisions, the more information you have with which to make them properly.
- The longer you leave options open, the more experiments you can run, the more things you can try, and the more information you will have when you reach the point at which those decisions can no longer be deferred.

## Chapter 16 Independence

### Decoupling Layers

- Example of layers we can decouple:
    - **User interfaces** change for reasons that have nothing to do with **business rules**.
    - **Business rules** themselves may be closely tied to the application, or they may be more general.
    - **The database**, the **query language**, and even the **schema** are technical details that have nothing to do with the business rules or the UI.

### Decoupling Use Cases

- The use case for _adding an order_ to an order entry system almost certainly will change at a different rate, and for different reasons, than the use case that *deletes an order* from the system.
- Each use case uses some UI, some application-specific business rules, some application-independent business rules, and some database functionality. Thus, as we are dividing the system in to **horizontal layers**, we are also dividing the system into **thin vertical** use cases that cut through those layers.

### Decoupling Mode

- If the different aspects of the use cases are separated, then those that must run at a high throughput are likely already separated from those that must run at a low throughput.
- If the UI and the database have been separated from the business rules, then they can run in different servers. Those that require higher bandwidth can be replicated in many servers.
- ▶️ In short, the decoupling that we did for the sake of the use cases also helps with operations.
- The point being made here is that sometimes we have to separate our components all the way to the **service level**.

### Independent Develop-ability

▶️ So long as the layers and use cases are decoupled, the architecture of the system will support the organization of the teams, irrespective of whether they are organized as **feature** teams, **component** teams, **layer** teams, or some other variation.

### Independent Deployability

▶️ The decoupling of the use cases and layers also affords a high degree of flexibility in deployment.

### Duplication

- There are different kinds of duplication:
    - **True duplication**, in which every change to one instance necessitates the same change to every duplicate of that instance.
    - **False or accidental duplication**: If two apparently duplicated sections of code evolve along different paths—if they change at different rates, and for different reasons—then they are not true duplicates.
- When you are **vertically separating** use cases from one another, you will run into this issue, and your temptation will be to couple the use cases because they have similar screen structures, or similar algorithms, or similar database queries and/or schemas. ‼️ Be careful. Resist the temptation to commit the sin of knee-jerk elimination of duplication. Make sure the duplication is real.
- By the same token, when you are separating **layers horizontally**, you might notice that the data structure of a particular database record is very similar to the data structure of a particular screen view. You may be tempted to simply pass the database record up to the UI, rather than to create a view model that looks the same and copy the elements across. Be careful: ‼️ This duplication is almost certainly accidental.

### Decoupling Modes (Again)

- Back to modes. There are many ways to decouple layers and use cases.:
    - **Source level**. We can control the dependencies between source code modules so that changes to one module do not force changes or recompilation of others (e.g., Ruby Gems).
    -  **Deployment level**. We can control the dependencies between deployable units such as jar files, DLLs, or shared libraries, so that changes to the source code in one module do not force others to be rebuilt and redeployed.
    - **Service level**. We can reduce the dependencies down to the level of data structures, and communicate solely through network packets such that every execution unit is entirely independent of source and binary changes to others (e.g., services or micro-services).

What is the best mode to use ❓

> A good architecture will allow a system to be born as a monolith, deployed in a single file, but then to grow into a set of independently deployable units, and then all the way to independent services and/or micro-services. Later, as things change, it should allow for reversing that progression and sliding all the way back down into a monolith.

## Chapter 17 Boundaries: Drawing Lines

- Software architecture is the art of drawing lines that I call *boundaries*. Those boundaries separate software elements from one another, and restrict those on one side from knowing about those on the other.

### A Couple of Sad Stories

- Company *P* used a *three-tiered architecture* for implementing a web solution. The programmers decided, very early on, that all domain objects would have three instantiations: one in the **GUI** tier, one in the **middleware** tier, and one in the **database** tier.
- The architect did not think of all the object instantiations, all the serializations, all the marshaling and de-marshaling, all the building and parsing of messages, all the socket communications, timeout managers, retry scenarios, and all the other extra stuff that you have to do just to get one simple thing done: adding a new field to an existing record 😵.
- The irony is that company *P* never sold a system that required a server farm. Every system they ever deployed was a single server.
- Similarly, company *W* **prematurely** adopted and enforced a suite of tools that promised *SOA*.
- To test anything you had to fire up all the necessary services, one by one, and fire up the message bus, and the BPel server, and ... And then, there were the propagation delays as these messages bounced from service to service, and waited in queue after queue.
- And then if you wanted to add a new feature—well, you can imagine the coupling between all those services, and the sheer volume of WSDLs that needed changing, and all the redeployments those changes necessitated ...

### FitNesse

- The idea was to create a simple wiki that wrapped Ward Cunningham’s FIT tool for writing acceptance tests.
- One of the first decisions was to **write our own web server**, specific to the needs of `FitNesse`.
    - ▶️ Postpone any web framework decision until much later.
- Also, avoid thinking about a database. We had `MySQL` in the back of our minds, but we purposely **delayed** that decision by employing a design that made the decision irrelevant. That design was simply to put an interface named `WikiPage` between all data accesses and the data repository itself.
- Early in the development of `FitNesse`, we drew a **boundary line between business rules and databases**. That line prevented the business rules from knowing anything at all about the database, other than the simple data access methods.
- The fact that we did not have a database running for early stages of development meant that, we did not have schema issues, query issues, database server issues, password issues, connection time issues, and all the other nasty issues that raise their ugly heads when you fire up a database 🤦.

### Which Lines Do You Draw, and When Do You Draw Them?

- You draw lines between things that matter and things that don’t.
	- The GUI doesn’t matter to the business rules, so there should be a line between them.
	- The database doesn’t matter to the GUI, so there should be a line between them.
	- The database doesn’t matter to the business rules, so there should be a line between them.
- The `BusinessRules` use the `DatabaseInterface` to load and save data. The `DatabaseAccess` implements the interface and directs the operation of the actual Database. <p align="center"><img src="assets/db-behind-interface.png" width="300px" height="auto"></p>
- Where is the boundary line? The boundary is drawn across the inheritance relationship, just below the `DatabaseInterface`.<p align="center"><img src="assets/boundary-line.png" width="300px" height="auto"></p>
- Note the direction of the arrow. The Database knows about the `BusinessRules`. The `BusinessRules` do not know about the Database. This implies that the `DatabaseInterface` classes live in the `BusinessRules` component, while the `DatabaseAccess` classes live in the Database component. <p align="center"><img src="assets/business-rules-and-database-components.png" width="300px" height="auto"></p>
- The direction of this line is important. It shows that the `Database` does not matter to the `BusinessRules`, but the `Database` cannot exist without the BusinessRules.

### What About Input and Output?

- Developers often think that the **GUI** is the system. They define a system in terms of the GUI, so they believe that they should see the GUI start working immediately. They fail to realize a critically important principle: **The IO is irrelevant**.
- And so, once again, we see the GUI and `BusinessRules` components separated by a boundary line. <p align="center"><img src="assets/boundary-gui-business-rules.png" width="300px" height="auto"></p>

### Plugin Architecture

- Taken together, these two decisions about the database and the GUI create a kind of pattern for the addition of other components. That pattern is the same pattern that is used by systems that allow third-party **plugins**.
- Because the UI in this design is considered to be a plugin, we have made it possible to plug in many different kinds of user interfaces (web based, client/server based, SOA based, Console based, etc...).
- THe same is true of the database. Since we have chosen to treat it as a plugin, we can replace it with any of the various `SQL` databases, or a `NOSQL` or a file system-based database, or any other kind of database.
<p align="center"><img src="assets/plugging-in-to-business-rules.png" width="300px" height="auto"></p>

## Chapter 18 Boundary Anatomy

- The architecture of a system is defined by a set of software components and the boundaries that separate them. Those boundaries come in many different forms. In this chapter we’ll look at some of the most common.

### Boundary Crossing

- At runtime, a boundary crossing is nothing more than a function on one side of the boundary calling a function on the other side and passing along some data. The trick to creating an appropriate boundary crossing is to **manage the source code dependencies**.
- The simplest possible boundary crossing is a function call from a low-level client to a higher-level service. Both the runtime dependency and the compile-time dependency point in the same direction, toward the higher-level component.
<p align="center"><img src="assets/low-high-level-boundary-crossing.png" width="400px" height="auto"></p>

- When a high-level client needs to invoke a lower-level service, dynamic polymorphism is used to invert the dependency against the flow of control.
<p align="center"><img src="assets/boundary-crossing-against-control-flow.png" width="400px" height="auto"></p>

### Deployment Components

- The simplest physical representation of an architectural boundary is a **dynamically linked library** like a .Net DLL, a Java jar file, a Ruby Gem, or a UNIX shared library.

### Threads

- Both monoliths and deployment components can make use of threads. Threads are **not architectural boundaries** or units of deployment, but rather a way to organize the schedule and order of execution. They may be wholly contained within a component, or spread across many components.

### Local Processes

- A much stronger physical architectural boundary is the local process.
- Local processes run in a separate address spaces but communicate with each other using sockets, or some other kind of operating system communications facility such as mailboxes or message queues.
- The segregation strategy between local processes is the same as for monoliths and binary components. Source code dependencies point in the same direction across the boundary, and always toward the higher-level component.

### Services

- The strongest boundary is a service. The services assume that all communications take place over the network.
- Otherwise, the same rules apply to services as apply to local processes. Lower-level services should “plug in” to higher-level services. The source code of higher-level services must not contain any specific physical knowledge (e.g., a URI) of any lower-level service.

## Chapter 19: Policy and Level

- Software systems are statements of **policy**. A computer program is a detailed description of the policy by which inputs are transformed into outputs.
- Part of the art of developing a software architecture is carefully separating those policies from one another, and regrouping them based on the ways that they change. Policies that change for the **same reasons**, and at the **same times**, are at the **same level** and belong together in the **same component**. Policies that change for **different reasons**, or at **different times**, are at **different levels** and should be **separated** into **different components**.
- Low-level components are designed so that they depend on high-level components.

### Level

- The **farther** a policy is from both the inputs and the outputs of the system, the **higher** its level.
<p align="center"><img src="assets/class-diagram-low-high-level-policy.png" width="300px" height="auto"></p>

- *ConsoleReader* and *ConsoleWriter* are shown here as classes. They are low level because they are close to the inputs and outputs.
- This structure decouples the high-level encryption policy from the lower-level input/output policies. This makes the *encryption* policy usable in a wide range of contexts. When changes are made to the input and output policies, they are not likely to affect the encryption policy.
- Keeping these policies separate, with all source code dependencies pointing in the **direction of the higher-level** policies, **reduces the impact of change**. Trivial but urgent changes at the lowest levels of the system have little or no impact on the higher, more important, levels.

## Chapter 20: Business Rules

- Strictly speaking, **business rules** are rules or procedures that make or save the business money.
- Very strictly speaking, these rules would make the business money, **irrespective** of whether they were implemented on a **computer**. They would make money even if they were executed **manually**.
- **Critical Business Rules** are the rules that are critical to the business itself, and would exist even if there were no system to automate them (bank charges N% interest for a loan is a critical business rule).
- Critical Business Rules usually require some data to work with. For example, our loan requires a loan balance, an interest rate, and a payment schedule. We shall call this data **Critical Business Data**. This is the data that would exist even if the system were not **automated**.

### Entities

- An **Entity** is an object within our computer system that embodies a small set of critical business rules operating on Critical Business Data. The Entity object either contains the Critical Business Data or has very easy access to that data. The interface of the Entity consists of the functions that implement the Critical Business Rules that operate on that data.
<p align="center"><img src="assets/loan-entity.png" width="300px" height="auto"></p>

- This class stands alone as a representative of the business. It is **unsullied** with **concerns about databases**, **user interfaces**, or **third-party frameworks**. It could serve the business in any system, irrespective of how that system was presented, or how the data was stored, or how the computers in that system were arranged. :bulb: The Entity is pure business and nothing else.

### Use Cases

- Not all business rules are as pure as `Entities`.
- A use case is a description of the way that an automated system is used. They make sense only as part of an **automated** system.
- A use case describes **application-specific** business rules as opposed to the Critical Business Rules within the Entities.
- :bangbang: Use cases do not describe how the system **appears to the user**. Instead, they describe the application-specific rules that govern the interaction between the users and the Entities. How the data gets in and out of the system is **irrelevant** to the use cases.
- Entities have **no knowledge** of the use cases that control them. This is another example of the direction of the dependencies following the `DIP`. High-level concepts, such as Entities, **know nothing** of lower-level concepts, such as use cases. Instead, the lower-level use cases know about the higher-level Entities.

### Request and Response Models

- Use cases expect input data, and they produce output data. However, a well-formed use case object should have no inkling about the way that data is communicated to the user, or to any other component. We certainly don’t want the code within the use case class to know about HTML or SQL!
- This lack of dependencies is **critical**. If the request and response models are **not independent**, then the use cases that depend on them will be indirectly bound to whatever dependencies the models carry with them.

## Chapter 21: Screaming Architecture

- What does the architecture of your application **scream**? When you look at the **top-level directory** structure, and the source files in the highest-level package, do they scream “*Health Care System*” or “*Accounting System*” or “*Inventory Management System*”? Or do they scream *“Rails*” or “*Spring/Hibernate*” or “*ASP*”?

### The Theme of an Architecture

- Just as the plans for a house or a library scream about the use cases of those buildings, so should the architecture of a software application scream about the use cases of the application.
- Architectures are not (or should not be) about **frameworks**. Architectures should not be supplied by frameworks. Frameworks are **tools** to be used, not architectures to be conformed to. If your architecture is based on frameworks, then it cannot be based on your use cases.

### The Purpose of an Architecture

- Good architectures are centered on use cases so that architects can safely describe the structures that support those use cases without committing to frameworks, tools, and environments.

### But What About the Web?

- The web is a delivery mechanism — an **IO device**— and your application architecture should treat it as such. The fact that your application is delivered over the web is a **detail** and should not dominate your system structure.

### Frameworks Are Tools, Not Ways of Life

- Frameworks can be very powerful and very useful.
- However, look at each framework with a jaded eye. View it **skeptically**. Yes, it might help, but at what **cost**? Ask yourself how you should use it, and how you should **protect** yourself from it. Think about how you can preserve the use-case emphasis of your architecture.
- :+1: Develop a strategy that prevents the framework from taking over that architecture.

### Testable Architectures

- If your system architecture is all about the use cases, and if you have kept your frameworks at arm’s length, then you should be able to unit-test all those use cases **without any of the frameworks** in place. You shouldn’t need the **web server** running to run your tests. You shouldn’t need the **database** connected to run your tests. Your Entity objects should be plain old objects that have **no dependencies** on frameworks or databases or other complications.

## Chapter 22: The Clean Architecture

- Over the last several decades we’ve seen a whole range of ideas regarding the architecture of systems. These include:
	- Hexagonal Architecture (also known as Ports and Adapters)
	- DCI (Data, Context, Interaction)
	- BCE (Boundary-Control-Entity)
- Although these architectures all vary somewhat in their details, they are very similar. They all have the same objective, which is the **separation of concerns**. They all achieve this separation by dividing the software into **layers**. Each has at least **one layer for business rules**, and **another layer for user and system interfaces**.
- Each of these architectures produces systems that have the following characteristics:
	- Independent of frameworks
	- Testable
	- Independent of the UI
	- Independent of the database
	- Independent of any external agency
<p align="center"><img src="assets/clean-architecture.png" width="300px" height="auto"></p>

### The Dependency Rule

- The further in you go, the higher level the software becomes. The outer circles are **mechanisms**. The inner circles are **policies**.
- Source code dependencies must point only inward, toward higher-level policies.

### Interface Adapters

- The software in the interface adapters layer is a set of adapters that **convert** data from the format most convenient for the use **cases and entities**, to the format most convenient for some **external agency** such as the database or the web.
- It is this layer, for example, that will wholly contain the `MVC` architecture of a GUI.

### Only Four Circles?

- Dependency Rule always applies. Source code **dependencies** always point **inward**. As you move inward, the level of **abstraction** and **policy** **increases**. The outermost circle consists of **low-level** concrete **details**. As you move **inward**, the software grows more **abstract** and **encapsulates** **higher-level** policies. The innermost circle is the most general and highest level.

### Crossing Boundaries

- For example, suppose the use case needs to call the *presenter*. This call must not be **direct** because that would violate the **Dependency Rule**: No name in an outer circle can be mentioned by an inner circle. So we have the use case call an **interface** (shown in Figure above as “use case output port”) in the inner circle, and have the presenter in the outer circle **implement** it.

### Which Data Crosses the Boundaries

- Typically the data that crosses the boundaries consists of **simple** data structures. You can use basic structs or simple data transfer objects if you like. Or the data can simply be arguments in **function calls**. Or you can pack it into a **hashmap**, or construct it into an **object**.
- When we pass data across a boundary, it is always in the form that is **most convenient for the inner circle**.

## Chapter 23 Presenters and Humble Objects

### The Humble Object Pattern

- Is a design pattern that was originally identified as a way to help unit testers to **separate behaviors** that are hard to test from behaviors that are easy to test.
- The idea is very simple: **split the behaviors into two modules or classes**.
	- One of those modules is humble; it contains all the hard-to-test behaviors stripped down to their barest essence.
	- The other module contains all the testable behaviors that were stripped out of the humble object.
- For example, GUIs are hard to unit test because it is very difficult to write tests that can see the screen and check that the appropriate elements are displayed there. However, most of the behavior of a GUI is, in fact, easy to test. Using the Humble Object pattern, we can separate these two kinds of behaviors into two different classes called the **Presenter** and the **View**.

### Presenters and Views

- The View is the humble object that is **hard to test**. The code in this object is kept as simple as possible. It moves data into the GUI but does not process that data.
- The Presenter is the **testable object**. Its job is to accept data from the application and format it for presentation so that the View can simply move it to the screen.

### Testing and Architecture

- The Humble Object pattern is a good example, because the separation of the behaviors into testable and non-testable parts often defines an **architectural boundary**. The **Presenter/View** boundary is one of these boundaries, but there are many others.

- Conclusion: At each architectural boundary, we are likely to find the Humble Object pattern lurking somewhere nearby. The communication across that boundary will almost always involve some kind of simple data structure, and the boundary will frequently divide something that is hard to test from something that is easy to test. The use of this pattern at architectural boundaries vastly increases the testability of the entire system.

## Chapter 24: Partial Boundaries

- Full-fledged architectural boundaries are **expensive**. They require reciprocal polymorphic Boundary interfaces, *Input* and *Output* data structures, and all of the **dependency management** necessary to **isolate** the two sides into independently **compilable** and **deployable** components.
- A potential solution is to create a **partial boundary**.

### Skip the Last Step

- One way to construct a partial boundary is to do all the work necessary to create independently compilable and deployable components, and then simply **keep them together in the same component**.
- The reciprocal interfaces are there, the input/output data structures are there, and everything is all set up—but we compile and deploy all of them as a **single component**.
- Obviously, this kind of partial boundary requires the same amount of code and preparatory design work as a full boundary. However, it does not require the **administration** of **multiple** components. There’s **no version number tracking** or **release management burden**. That difference should not be taken lightly.

### One-Dimensional Boundaries

- A simpler structure that serves to hold the place for later extension to a full- fledged boundary is shown below. It exemplifies the traditional **Strategy** pattern. A `ServiceBoundary` interface is used by clients and implemented by `ServiceImpl` classes.
<p align="center"><img src="assets/strategy-pattern.png" width="300px" height="auto"></p>

### Facades

- An even simpler boundary is the **Facade** pattern, illustrated below. In this case, even the **dependency inversion** is **sacrificed**. The boundary is simply defined by the Facade class, which lists all the services as methods, and deploys the service calls to classes that the client is not supposed to access.
<p align="center"><img src="assets/facade-pattern.png" width="300px" height="auto"></p>

## Chapter 25: Layers and Boundaries

- It is easy to think of systems as being composed of three components: UI, business rules, and database.
- For some simple systems, this is sufficient. For most systems, though, the number of components is larger than that.
- Consider, for example, a simple computer game (Hunt The Wumpus). It is easy to imagine the three components. The UI handles all messages from the player to the game rules. The game rules store the state of the game in some kind of persistent data structure. But is that all there is?
- Of course, no ! It should be clear that we could easily apply the clean architecture approach in this context:
<p align="center"><img src="assets/clean-architecture-hunt-the-wumpus.png" width="300px" height="auto"></p>

1. Decouple text based UI from the game rules so that our version can use different languages in different markets. The game rules will communicate with the UI component using a language-independent API, and the UI will translate the API into the appropriate human language.
2. The state of the game is maintained on some persistent store—perhaps in flash, or perhaps in the cloud, or maybe just in RAM, we’ll create an API that the game rules can use to communicate with the data storage component.
3. We also might want to vary the mechanism by which we communicate the text. For example, we might want to use a normal shell window, or text messages, or a chat application. We should construct an API that crosses that boundary and isolates the language from the communications mechanism.
- Conclusion: On the one hand, some very smart people have told us, over the years, that we should not **anticipate the need for abstraction**. This is the philosophy of `YAGNI`: “*You aren’t going to need it.*” There is wisdom in this message, since o**ver-engineering is often much worse than under-engineering**. On the other hand, when you discover that you truly do need an architectural boundary where none exists, the costs and risks can be very high to add such a boundary.
- You must weigh the costs and determine where the architectural boundaries lie, and which should be fully implemented, and which should be partially implemented, and which should be ignored. This is an **on-going process** that every architect should do as the system evolves.

## Chapter 26: The Main Component

### The Ultimate Detail

- The **Main** component is the ultimate detail—the **lowest-level policy**. It is the initial entry point of the system. Nothing, other than the OS, depends on it. Its job is to create all the Factories, Strategies, and other global facilities, and then hand control over to the high-level abstract portions of the system.
- It is in this `Main` component that **dependencies** should be **injected** by a **Dependency Injection framework**. Once they are injected into `Main`, `Main` should distribute those dependencies normally, without using the framework.
- Think of `Main` as the **dirtiest** of all the dirty components.
- Think of `Main` as a **plugin** to the application—a plugin that sets up the initial conditions and configurations, gathers all the outside resources, and then hands control over to the high-level policy of the application. Since it is a plugin, it is possible to have many `Main` components, one for each configuration of your application.

## Chapter 27: Services: Great and Small

### Service Architecture?

- Services that simply separate application behaviors are little more than expensive function calls, and are not necessarily **architecturally significant**.
- This is not to say that all services should be architecturally significant. There are often substantial benefits to creating services that separate functionality across processes and platforms—whether they obey the Dependency Rule or not. It’s just that services, in and of themselves, **do not define an architecture**.

### Service Benefits?

- The question mark in the preceding heading indicates that this section is going to challenge the current popular orthodoxy of service architecture.

#### The Decoupling Fallacy

- One of the big supposed benefits of breaking a system up into services is that services are strongly decoupled from each other.
- Yes, services are decoupled at the level of **individual variables**. However, they can still be coupled by shared resources within a processor, or on the network. What’s more, they are strongly **coupled by the data they share**.

### The Fallacy of Independent Development and Deployment

- Another of the supposed benefits of services is that they can be owned and operated by a dedicated team.
- That team can be responsible for writing, maintaining, and operating the service as part of a dev-ops strategy. This independence of development and deployment is presumed to be *scalable*.
- The decoupling fallacy means that services cannot always be independently developed, deployed, and operated. To the extent that they are **coupled by data or behavior**, the development, deployment, and operation must be **coordinated**.

### The Kitty Problem

- As an example of these two fallacies, let’s look at our taxi aggregator system again.
<p align="center"><img src="assets/taxi-service-arch.png" width="400px" height="auto"></p>

- Company decided to add a new feature so a user can order kittens to be delivered to their homes or to their places of business.
- Such feature will require all services to be coordinated and to change.
- :arrow_forward: This is the problem with cross-cutting concerns. Every software system must face this problem, whether service oriented or not. Functional decompositions, of the kind depicted in the service diagram above, are **very vulnerable to new features** that cut across all those functional behaviors.

### Objects to the Rescue

- How would we have solved this problem in a component-based architecture?
    - Careful consideration of the **SOLID** design principles would have prompted us to create a set of classes that could be **polymorphically extended** to handle new features.

<p align="center"><img src="assets/taxi-object-oriented-arch.png" width="400px" height="auto"></p>

- Much of the logic of the original services is preserved within the base classes of the object model.
    - However, that portion of the logic that was specific to rides has been extracted into a `Rides` component.
    - The new feature for kittens has been placed into a `Kittens` component.
    - These two components **override** the abstract base classes in the original components using a pattern such as **Template Method** or **Strategy**.
- Note also that the classes that implement those features are created by **factories** under the control of the UI.

### Component-Based Services

- Can we do that for services? And the answer is, Yes!
- The services can be designed using the SOLID principle, but each has its own internal component design, allowing new features to be added as new derivative classes. Those derivative classes live within their own components.
<p align="center"><img src="assets/component-based-taxi-arch.png" width="400px" height="auto"></p>

### Cross-Cutting Concerns

- To deal with the **cross-cutting concerns** that all significant systems face, services must be designed with **internal component** architectures that follow the **Dependency Rule**. Those services **do not define the architectural boundaries** of the system; instead, the **components** within the services do.
<p align="center"><img src="assets/cross-cutting-concerns.png" width="400px" height="auto"></p>

- :arrow_forward: As useful as services are to the scalability and develop-ability of a system, they are not, in and of themselves, **architecturally** significant elements. The architecture of a system is defined by the **boundaries** drawn within that system, and by the **dependencies** that cross those boundaries. That architecture is not defined by the physical mechanisms by which elements communicate and execute.

## Chapter 28: The Test Boundary

- tests are part of the system, and they participate in the architecture just like every other part of the system does.

### Tests as System Components

- Tests, by their very nature, follow the Dependency Rule; they are very detailed and concrete; and they always depend inward toward the code being tested.
- In fact, you can think of the tests as the outermost circle in the architecture. Nothing within the system depends on the tests, and the tests always depend inward on the components of the system.
- Tests are the most **isolated** system component. They are not necessary for system operation. No user depends on them. Their role is to support **development**, not **operation**.

### Design for Testability

- The extreme isolation of the tests, combined with the fact that they are not usually deployed, often causes developers to think that tests fall outside of the design of the system. :warning: This is a catastrophic point of view.
- *Fragile Tests Problem* occur where a simple change would break hundreds or thousands of tests.
- :arrow_forward: The solution is to design for testability. The first rule of software design— whether for testability or for any other reason—is always the same: *Don’t depend on volatile things*.

### The Testing API

- The way to accomplish this goal is to create a specific API that the tests can use to verify all the business rules.
- The purpose of the testing API is to decouple the tests from the application. This decoupling encompasses more than just detaching the tests from the UI: The goal is to decouple the **structure of the tests** from the **structure of the application**.
- :arrow _forward: Tests are not outside the system; rather, they are parts of the system that must be well designed if they are to provide the desired benefits of stability and regression. Tests that are not designed as part of the system tend to be fragile and difficult to maintain. Such tests often wind up on the maintenance room floor—discarded because they are too difficult to maintain.

## Chapter 29: Clean Embedded Architecture

- According to the author's definition: *Firmware does not mean code lives in ROM. It’s not firmware because of where it is stored; rather, it is firmware because of what it depends on and how hard it is to change as hardware evolves*.
-  For engineers and programmers, the message is clear: Stop writing so much firmware and give your code a chance at a long useful life. Of course, demanding it won’t make it so. Let’s look at how we can keep embedded software architecture clean to give the software a fighting chance of having a long and useful life.

### App-titude Test

- Getting an app to work is what I call the **App-titude** test for a programmer. Programmers, embedded or not, who just concern themselves with getting their app to work are doing their products and employers a disservice. There is much more to programming than just getting an app to work.
    - “First make it work.” You are out of business if it doesn’t work.
    - “Then make it right.” Refactor the code so that you and others can understand it and evolve it as needs change or are better understood.
    - “Then make it fast.” Refactor the code for “needed” performance.

### The Target-Hardware Bottleneck

- One of the special embedded problems is the **target-hardware bottleneck**.
- When embedded code is structured without applying clean architecture principles and practices, you will often face the scenario in which you can test your code only on the target. If the target is the only place where testing is possible, the target-hardware bottleneck will slow you down.

### A Clean Embedded Architecture Is a Testable Embedded Architecture

- Layering comes in many flavors. Let’s start with three layers:
- due to technology advances and Moore’s law, the hardware will change. Parts become obsolete, and new parts use less power or provide better performance or are cheaper. Whatever the reason, as an embedded engineer, I don’t want to have a bigger job than is necessary when the inevitable hardware change finally happens.
<p align="center"><img src="assets/three-layers.png" width="400px" height="auto"></p>

- Software and firmware intermingling is an anti-pattern. Code exhibiting this anti-pattern will resist changes. In addition, changes will be dangerous, often leading to unintended consequences.
<p align="center"><img src="assets/mix-soft-firm-anti-pattern.png" width="400px" height="auto"></p>

- The line between software and firmware is a bit fuzzier than the line between code and hardware:
<p align="center"><img src="assets/hardware-is-a-detail.png" width="400px" height="auto"></p>

- One of your jobs as an embedded software developer is to firm up that line. The name of the boundary between the software and the firmware is the **hardware abstraction layer** (HAL). This is not a new idea: It has been in PCs since the days before Windows.
- The HAL exists for the software that sits on top of it, and its API should be tailored to that software’s needs.

### Don’t Reveal Hardware Details to the User of the HAL

#### The Processor Is a Detail

- When your embedded application uses a specialized tool chain, it will often provide header files to help you.
- These compilers often take liberties with the C language, adding new keywords to access their processor features. The code will **look like C**, but it is **no longer C**, because we lost **portability**.
- :light: A clean embedded architecture would use these device access registers directly in very few places and confine them totally to the **firmware**.
- If you use a micro-controller like this, your firmware could isolate these low-level functions with some form of a **processor abstraction layer (PAL)**. Firmware above the PAL could be tested off-target, making it a little less firm.

#### The Operating System Is a Detail

- To give your embedded code a good chance at a long life, you have to treat the operating system as a detail and protect against OS dependencies.
- A clean embedded architecture isolates software from the OS, through an **operating system abstraction layer** (OSAL). In some cases, implementing this layer might be as simple as changing the name of a function. In other cases, it might involve wrapping several functions together.
<p align="center"><img src="assets/os-abstraction-layer.png" width="400px" height="auto"></p>

#### Programming to Interfaces and Substitutability

- In addition to adding a HAL and potentially an OSAL inside each of the major layers (software, OS, firmware, and hardware), you can—and should— apply the principles described throughout this book. These principles encourage separation of concerns, programming to interfaces, and substitutability.

### DRY Conditional Compilation Directives

- One use of substitutability that is often overlooked relates to how embedded C and C++ programs handle different targets or OSs.
- There is a tendency to use **conditional compilation** to turn on and off segments of code.
- I see `#ifdef BOARD_V2` once, it’s not really a problem. Six thousand times is an extreme problem. :happy:

## Chapter 30: The Database Is a Detail

- From an architectural point of view, the database is a **non-entity**, it is a detail that does not rise to the level of an architectural element.
- The database is a utility that provides access to the data. From the architecture’s point of view, that utility is **irrelevant** because it’s a low-level detail —a mechanism. And a good architect does not allow low-level mechanisms to pollute the system architecture.

### Relational Databases

- no matter how brilliant, useful, and mathematically sound a technology (such as relational DBs) it is, it is still just a **technology**. And that means it’s a **detail**.
- :man_shrugging: Many data access frameworks (ORMs) allow database rows and tables to be passed around the system as **objects**. Allowing this is an **architectural error**. It couples the use cases, business rules, and in some cases even the UI to the relational structure of the data.

### Why Are Database Systems So Prevalent?

- What accounts for the preeminence of Oracle, MySQL, and SQL Server? In a word: **disks**.
- The rotating magnetic disk was the mainstay of data storage for five decades.

### What If There Were No Disk?

- As prevalent as disks once were, they are now a dying breed 🫨 Soon they will have gone the way of tape drives, floppy drives, and CDs. They are being replaced by **RAM**.
- When all the disks are gone, and all your data is stored in RAM, how will you organize that data? ▶️ You’ll organize it into **linked lists, trees, hash tables, stacks, queues**, or any of the other myriad data structures, and you’ll access it using pointers or references— because that’s what programmers do.

### But What about Performance?

- When it comes to data storage, performance is a concern that can be entirely **encapsulated** and **separated** from the **business** rules.
- It has nothing whatsoever to do with the overall architecture of our systems.

### Conclusion:

- The organizational structure of data, the data model, is architecturally significant. The technologies and systems that move data on and off a rotating magnetic surface are not.
- Relational database systems that force the data to be organized into tables and accessed with SQL have much more to do with the latter than with the former.
- ▶️ The data is significant. The database is a detail.

## Chapter 31: The Web Is a Detail

- As architects of the apps, we should keep the **UI** and **business rules** **isolated** from each other, because there are always marketing geniuses 😄 out there just waiting to pounce on the next little bit of coupling you create.
- The GUI is a detail. The web is a GUI. ▶️ the web is a detail. And, as an architect, you want to put details like that behind boundaries that keep them separate from your core business logic.
- The interaction between the app and the GUI is *chatty* in ways that are quite specific to the kind of GUI you have. The dance between a browser and a web application is different from the dance between a desktop GUI and its application. Trying to abstract out that dance, the way **devices** are abstracted out of *UNIX*, seems unlikely to be possible.
- But another boundary between the UI and the app can be abstracted. The business logic can be thought of as a suite of use cases, each of which performs some function on behalf of a user. Each use case can be described based on the input data, the processing preformed, and the output data.

### Conclusion

- This kind of abstraction is **not easy**, and it will likely take several iterations to get just right. But it is possible. And since the world is full of marketing geniuses, it’s not hard to make the case that it’s often very necessary.

## Chapter 32: Frameworks Are Details

- Frameworks are not architectures — though some try to be.

### Framework Authors

- Framework authors know their own problems, and the problems of their coworkers and friends. And they write their frameworks to solve those problems — not yours.

### Asymmetric Marriage

- The relationship between you and the framework author is extraordinarily asymmetric. You must make a huge **commitment** to the **framework**, but the framework author makes **no commitment to you** whatsoever.
- The author urges you to **couple** your application to the framework as **tightly** as possible.
- 💯 Nothing feels more **validating** to a framework author than a bunch of users willing to inextricably derive from the author’s base classes.

### The Risks

- The architecture of the framework is often not very clean. Frameworks tend to **violate** the **Dependency Rule**.
- As your product matures, it may outgrow the facilities of the framework. You’ll find the framework fighting you more and more as time passes.
- The framework may evolve in a direction that you don’t find helpful. You may be stuck upgrading to new versions that don’t help you.
- A new and better framework may come along that you wish you could switch to 🤷‍♂️.

### The Solution

- You can use the framework — just don’t couple to it. Keep it at arm’s length. Treat the framework as a detail that belongs in one of the outer circles of the architecture. Don’t let it into the **inner circles**.

### Conclusion

- When faced with a framework, try not to marry it **right away**. See if there aren’t ways to date it for a while before you take the plunge.
- Keep the framework behind an architectural boundary if at all possible, for as long as possible. Perhaps you can find a way to get the milk without buying the 🐮.

## Chapter 33 : Case Study: Video Sales

- The architecture diagram below includes two dimensions of separation. The first is the **separation of actors** based on the **SIP** the second is the Dependency Rule.
<p align="center"><img src="assets/video-sales-uml.png" width="400px" height="auto"></p>

- The goal of both is to **separate components** that change for **different reasons**, and at **different rates**. The different reasons correspond to the **actors**; the different rates correspond to the different levels of **policy**.
- Once you have structured the code this way, you can mix and match how you want to actually deploy the system. You can group the components into deployable deliverables in any way that makes sense, and easily change that grouping when conditions change.
<p align="center"><img src="assets/video-sales-arch.png" width="400px" height="auto"></p>

## Chapter 34: The Missing Chapter

### Package by Layer

The first, and perhaps simplest, design approach is the traditional horizontal layered architecture, where we separate our code based on what it does from a technical perspective.
In this typical layered architecture, we have one layer for the web code, one layer for our “business logic,” and one layer for persistence. <p align="center"><img src="assets/package-by-layer.png" width="400px" height="auto"></p>
- 👎 Problems:
    - Once your software grows in scale and complexity, you will quickly find that having three large buckets of code isn’t sufficient, and you will need to think about modularizing further.
    - A layered architecture doesn’t **scream** anything about the **business domain**.

### Package by Feature

- This is a vertical slicing, based on related features, domain concepts, or aggregate roots (to use domain-driven design terminology). In the typical implementations that I’ve seen, all of the types are placed into a single Java package, which is named to reflect the concept that is being grouped.
- Same interfaces and classes as before, but they are all placed into a single Java package rather than being split among three packages. <p align="center"><img src="assets/package-by-feature.png" width="400px" height="auto"></p>
- This is a very simple refactoring from the “package by layer” style, but the top-level organization of the code now screams something about the business domain.

### Ports and Adapters

- Approaches such as “ports and adapters,” the “hexagonal architecture,” “boundaries, controllers, entities,” and so on
aim to create architectures where business/domain-focused code is independent and separate from the technical implementation details such as frameworks and databases. To summarize, you often see such code bases being composed of an “inside” (domain) and an “outside” (infrastructure).
- The “inside” region contains all of the domain concepts, whereas the “outside” region contains the interactions with the outside world (e.g., UIs, databases, third-party integrations). The major rule here is that the “outside” depends on the “inside”—never the other way around. <p align="center"><img src="assets/ports-adapters.png" width="400px" height="auto"></p>

### Package by Component

I’s a hybrid approach to everything we’ve seen so far, with the goal of bundling all of the responsibilities related to a single coarse-grained component into a single Java package.
- It’s about taking a service-centric view of a software system, which is something we’re seeing with micro-service architectures as well. In the same way that ports and adapters treat the web as just another delivery mechanism, “package by component” keeps the user interface separate from these coarse-grained components. <p align="center"><img src="assets/package-by-component.png" width="400px" height="auto"></p>
- A key benefit of the “package by component” approach is that if you’re writing code that needs to do something with orders, there’s just one place to go—the `OrdersComponent`. Inside the component, the separation of concerns is still maintained, so the business logic is separate from data persistence, but that’s a component implementation detail that consumers don’t need to know about.

### The Devil Is in the Implementation Details

- Marking all of your types as `public` means you’re not taking advantage of the facilities that your programming language provides with regard to **encapsulation**. In some cases, there’s literally nothing preventing somebody from writing some code to instantiate a concrete implementation class directly, **violating** the intended architecture style.

### Organization versus Encapsulation

- Looking at this issue another way, if you make all types in your Java application `public`, the packages are simply an organization mechanism (a grouping, like folders), rather than being used for encapsulation.
- The net result is that if you ignore the packages (because they don’t provide any means of encapsulation and hiding), it doesn’t really matter which architectural style you’re aspiring to create 🤷.
- what I’ve described here relates to a monolithic app, where all of the code resides in a single source code tree. If you are building such an application (and many people are), I would certainly encourage you to lean on the **compiler** to enforce your architectural principles, rather than relying on **self-discipline** and **post-compilation** tooling.

### Conclusion: The Missing Advice

- :+1: The whole point of this chapter is to highlight that your best design intentions can be destroyed in a flash if you don’t consider the intricacies of the implementation strategy.
- :+1: Think about how to map your desired design on to code structures, how to organize that code, and which decoupling modes to apply during runtime and compile-time.
- :+1: Leave options open where applicable, but be pragmatic, and take into consideration the size of your team, their skill level, and the complexity of the solution in conjunction with your time and budgetary constraints.
- :+1: Also think about using your compiler to help you enforce your chosen architectural style, and watch out for coupling in other areas, such as data models. The devil is in the implementation details.
