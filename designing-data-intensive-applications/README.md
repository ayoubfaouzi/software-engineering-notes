# Designing Data Intensive Applications

notes taken from reading the _Designing Data Intensive Applications_ book by Martin Kleppmann.

## Part II

- There are various reasons why you might want to distribute a database across multiple machines:
  - Scalability
  - Fault Taulerance/High Availability
  - Latency
- If all you need is to scale to **higher load**, the simplest approach is to buy a more powerful machine ➡️ scale vertically.
  - In this kind of **shared-memory architecture**, all the components can be treated as a single machine.
  - The problem with a shared-memory approach is that the cost grows faster than **linearly**: a machine with twice as many CPUs, twice as much RAM, and twice as much disk capacity as another typically costs significantly more than twice as much 🤷.
  - Limited fault tolerance (hot-swappable components), limited to a **single geographic** location.
- Another approach is the **shared-disk architecture**, the problem is you can have **contention** and the overhead of **locking** limit the scalability of the shared-disk approach.
- By contrast, **shared-nothing architectures** (sometimes called **horizontal scaling** or scaling out)
  - Coordination between nodes is done at the software level, using a conventional network.
  - Part II focus on shared-nothing architectures not because they are necessarily the best choice for every use case, but rather because they require the most caution from you, the application developer.
