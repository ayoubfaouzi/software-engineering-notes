# Chapter 8. The Trouble with Distributed Systems

- Recent chapters discussed handling failures like replica failover, replication lag, and weak isolation in transactions.
- But those views were still too optimistic 🥹 — in reality, **anything that can go wrong will go wrong** in distributed systems.
- Distributed systems introduce many new failure modes beyond single-computer software.
- This chapter provides a pessimistic overview of problems: unreliable networks, faulty clocks/timing, and reasoning challenges about system state.
- The 🥅 for engineers: build reliable systems that meet guarantees despite inevitable failures 🚀.

## Faults and Partial Failures

- On a **single computer**, software is mostly **deterministic**: operations either work or fail, and hardware errors usually cause total crashes, **not partial flakiness**.
- Computers are deliberately designed to **fail-stop **(crash rather than return wrong results) and present a perfect, idealized model of computation.
- In distributed systems, this idealization breaks down — nodes and networks can fail independently, leading to **partial failures**.
- Partial failures are **nondeterministic**: some operations may succeed, others fail unpredictably, and it may be unclear whether an action completed.
- 👉 This **nondeterminism** and **uncertainty** is what makes distributed systems fundamentally hard to build and reason about.

### Cloud Computing and Supercomputing

- There are different philosophies for building large-scale computing systems:
  - **Supercomputers (HPC)**: Designed for scientific batch jobs (e.g., weather simulations). They checkpoint state and often crash the entire workload if one node fails, treating partial failure as **total failure**. Built with specialized, reliable hardware and fast interconnects.
  - **Cloud computing / internet services**: Built from commodity hardware, running online, user-facing applications that must stay available. **Partial failures** are common and expected, so systems must tolerate node failures, allow rolling upgrades, and recover without downtime. Networks are IP/Ethernet-based, often spanning geographic regions with unreliable links.
  - **Enterprise datacenters**: Sit somewhere in between.
- 🔑 point: Distributed systems must **assume partial failure is inevitable**. Unlike single-node systems or supercomputers, services can’t stop everything for repairs. Instead, software must be designed with **fault tolerance**, **pessimism**, and **paranoia** — anticipating failures, handling them gracefully, and testing failure scenarios to ensure reliability.
- > Reliability emerges from layering fault-tolerant mechanisms over unreliable components, but only up to a point (error-correcting code, TCP, ..).

## Unreliable Networks

- The internet and most internal networks in DC (often Ethernet) are **asynchronous** packet networks. 
  - Messages may be delayed, lost, duplicated, or dropped.
  - Many possible failure modes:
    - Request lost.
    - Request delayed in a queue.
    - Remote node crashed or shut down.
    - Remote node paused temporarily (e.g., GC).
    - Request processed, but response lost.
    - Request processed, but response delayed
- Core challenge:
  - The sender cannot tell why no response was received.
  - **All failures look the same** → “no response yet.”
- Typical solution:
  - Use **timeouts**: assume failure if no response in time.
  - But timeouts are **ambiguous** → the request may still eventually arrive or be processed after sender has given up.

### Network Faults in Practice

- Despite decades of experience, networks are still unreliable 🤷‍♂️.
- Studies show frequent faults:
  - `~12 faults/month` in a medium datacenter.
  - Half disconnect one machine, half an entire rack.
  - Redundant hardware helps little against **human errors** (misconfigurations).
- Clouds like EC2 often have transient glitches; even private datacenters are not immune.
- Failures can be surprising:
  - **Switch upgrades** causing >1 minute delays.
  - 🦈 damaging undersea cables.
  - Links working in one direction but not the other.
- 🔑 point: Any network communication can fail, and software must be prepared to handle it.
- If fault handling is missing or weak:
  - Systems can deadlock or stop serving requests.
  - Worst case: data loss.
- Approaches:
  - Sometimes acceptable to just show users an error during outages.
  - More importantly: ensure systems recover cleanly once network returns.
  - Proactively test by injecting failures (e.g., *Chaos Monkey*).

### Detecting Faults

- Detecting faulty nodes in distributed systems is **inherently difficult** due to network uncertainty
- While some feedback is possible — like `TCP` connection refusals, process crash notifications, ICMP errors, or switch management interfaces — none of these guarantees that the application has actually processed a request.
- Rapid feedback is helpful but unreliable; to be sure a request succeeded, you need a **positive response** from the **application itself**.
- In general, systems must handle the possibility of no response, using **retries** and **timeouts** to eventually declare a node dead.

### Timeouts and Unbounded Delays

- Choosing a timeout for detecting node failures involves a ⚖️.
- A **long timeout** delays **failure detection**, **slowing user response**, while a **short timeout** risks **falsely** declaring a slow but alive **node** as **dead**.
- **Prematurely** marking nodes as dead can cause duplicated actions and additional load on other nodes, potentially triggering **cascading failures**.
- In theory, if network delays and node response times were **bounded**, a precise timeout could be **calculated**, but in real systems, networks are **asynchronous** and node response times are unpredictable, so transient delays can easily lead to incorrect failure detection.

#### Network congestion and queueing

- Network delays are highly variable due to **queueing**, which occurs at multiple stages:
  - network switches during congestion,
  - the destination machine when CPU cores are busy,
  - and in virtualized environments when a VM is paused. 
- `TCP` adds additional variability through **flow control** and **retransmissions**.
- Queueing delays are especially large when systems are near **maximum capacity**.
- In public clouds and multi-tenant datacenters, shared resources and “**noisy neighbors**” further increase unpredictability.
- 👉 Because of this **variability**, timeouts for failure detection cannot be fixed a priori. They must be determined experimentally or dynamically adjusted based on observed network latency and jitter, as implemented in systems like *Phi Accrual failure detectors* or TCP retransmission algorithms.

### Synchronous Versus Asynchronous Networks

- Distributed systems would be simpler if networks could **guarantee fixed maximum delays** and never drop packets.
- Unlike datacenter networks, traditional **telephone networks** achieve this reliability through **circuit-switched** **synchronous** connections, where a **fixed bandwidth is reserved** along the **entire route**.
- This eliminates **queueing**, ensuring a **bounded**, predictable end-to-end latency for the duration of the call.
- In contrast, typical computer networks are asynchronous, shared, and subject to variable delays and packet loss.

#### Can we not simply make network delays predictable❓

- Telephone circuits reserve a fixed bandwidth for a connection, unlike TCP, which opportunistically uses available bandwidth.
- Datacenter and internet networks use **packet switching** instead of circuits because they are optimized for **bursty** traffic (web pages, emails, file transfers).
- **Circuits** would **waste capacity** or require guessing bandwidth. TCP adapts dynamically to available capacity.
- **Hybrid** approaches (e.g., ATM, InfiniBand with QoS) can reduce queueing, but public clouds and multi-tenant DCs cannot guarantee bounded delays or reliability.
- ▶️ network timeouts must be chosen experimentally, since congestion and unbounded delays are inevitable 🤷.

## Unreliable Clocks

- Clocks are crucial in distributed systems for measuring durations (e.g., request latency, throughput) and recording points in time (e.g., timestamps, cache expiration).
- However, in **distributed** systems, network delays make it hard to determine the exact **order of events**, and each **machine’s clock** can **drift** due to **hardware imperfections**.
- Clock synchronization, commonly via `NTP`, can partially align clocks using servers with more accurate time sources like `GPS`.

### Monotonic Versus Time-of-Day Clocks

- Modern computers have two main types of clocks:
- **Time-of-day clocks**:
  - Provide the current date and time (wall-clock time).
  - Can be synchronized across machines using `NTP`.
  - Unsuitable for measuring elapsed time because they can **jump backward** (e.g., NTP adjustments) and may have coarse resolution.
- **Monotonic clocks**:
  - Measure **durations** or **time intervals**.
  - Guaranteed to move **forward only**, so safe for timeouts or response measurements.
  - **Absolute value** is **meaningless**; cannot be compared across machines.
  - High resolution (microseconds or better) and generally reliable, though multiple CPU timers may require OS compensation.
  - NTP can slew the clock rate slightly but cannot cause jumps.
- 🔑 takeaway: use time-of-day clocks for timestamps and monotonic clocks for measuring elapsed time in distributed systems.


### Clock Synchronization and Accuracy

- Monotonic clocks don’t need synchronization, but time-of-day clocks rely on NTP or other external sources, which are often **unreliable**.
- Clocks can **drift** due to **quartz inaccuracy** (temperature-dependent, e.g., `6 ms` drift every `30s` or `17s/day`):
  - NTP sync may fail if drift is **too large**.
  - Servers may be unreachable, or network **delays** can introduce errors (`35 ms` typical, up to `1s`).
  - Some NTP servers are **misconfigured**, and **leap seconds** complicate timekeeping, sometimes crashing systems.
  - **VMs** add further inaccuracies since **pauses** cause time jumps.
  - On **untrusted** devices, hardware clocks may be wrong **intentionally** or **accidentally**.
- Very high accuracy is possible (e.g., required `100 μs` sync in financial trading via `GPS/PTP`), but achieving it requires specialized hardware, careful configuration, and continuous monitoring, since small misconfigurations (e.g., blocked NTP) can quickly lead to large errors.

### Relying on Synchronized Clocks

- Clocks seem simple but are full of pitfalls: days aren’t always `86,400` seconds, clocks may move **backward**, and nodes often **disagree on time**.
- Like networks, clocks usually work but can fail, so software must tolerate incorrect clocks. Faulty or misconfigured clocks are dangerous because they often go unnoticed — systems keep running, but clock drift or NTP misconfigurations can silently cause data loss.
- 👉 Systems that depend on synchronized clocks must monitor **clock offsets** across nodes and treat any machine with excessive drift as failed, removing it from the cluster before it causes harm.

#### Timestamps for ordering events

- Using **time-of-day** clocks to order events across distributed nodes is dangerous because **clock skew**, even as small as a few milliseconds, can cause incorrect **event ordering**.
  - For example, in multi-leader replication with **LWW conflict resolution**, a later write can be discarded if its timestamp is slightly earlier due to clock differences, leading to silent data loss.
- Problems with LWW include:
  - Writes mysteriously disappearing due to clock skew.
  - Inability to distinguish between sequential and concurrent writes.
  - Collisions when two writes share the same timestamp.
  - Even with tight NTP synchronization, network delays and clock drift make reliable ordering via physical clocks impossible.
- 👉 use **logical clocks** (e.g., **Lamport clocks** or **version vectors**), which track causal ordering of events without depending on physical time.

#### Clock readings have a confidence interval

- Even if a system clock provides microsecond or nanosecond resolution, it is not truly accurate at that scale due to quartz drift, NTP limits, and network delays.
- In practice, public NTP often yields accuracy in the **10 to 100 of ms**, making fine-grained digits in timestamps meaningless 🤷‍♂️.
- A clock reading should be treated as a **range of possible times** (a confidence interval), not an exact point. The uncertainty depends on the time source (GPS/atomic clocks are far more precise than NTP over the internet).
- Most systems don’t expose this uncertainty — for example, `clock_gettime()` only returns a single value. **Google’s TrueTime API** (used in *Spanner*) is a rare exception: it explicitly returns an interval `[earliest, latest]`, ensuring the system accounts for uncertainty in distributed coordination.

#### Synchronized clocks for global snapshots

- On a single node, a simple increasing transaction ID works.
- In distributed databases, generating a **global monotonically** increasing ID is hard because it requires **coordination**.
- Using synchronized **time-of-day** clocks as **transaction IDs** seems attractive, but clock **uncertainty** makes it tricky.
- **Spanner’s approach**:
  - Uses *Google’s TrueTime* API, which gives each timestamp as an interval `[earliest,latest]` with uncertainty bounds.
  - If two intervals **don’t overlap**, their order is **guaranteed**.
  - To ensure causality, Spanner **waits** out the uncertainty **window** before committing, so future reads see consistent order.
- To minimize this wait, *Google* deploys GPS receivers or atomic clocks in each DC, keeping uncertainty around `7 ms`.
- 👉 Spanner leverages tightly synchronized clocks plus TrueTime intervals to assign consistent transaction timestamps across datacenters, but this technique hasn’t been widely adopted outside *Google*.

### Process Pauses

- In a leader-based distributed database, leases are often used to ensure only one node acts as leader at a time. A leader periodically renews its lease; if it fails, another node can take over.
- A naïve implementation checks the lease against the local system clock before processing requests, but this is **dangerous**:
  1. **Clock sync issues** – Lease expiry is calculated on one machine and compared to another’s clock. Even small clock skews can break correctness.
  2. **Unexpected pauses** – Even with a monotonic local clock, the system assumes little time passes between checking the lease and processing a request. But a thread may be paused unexpectedly for seconds or minutes, during which the lease could expire.
- Pauses can happen due to:
  - *Stop-the-world* garbage collection.
  - VM suspension or live migration.
  - Laptop sleep/resume.
  - OS/hypervisor **context switching** or CPU steal time.
  - Disk or network I/O delays.
  - Swapping/thrashing under memory pressure.
  - Signals like `SIGSTOP`.
- These pauses mean a node may process requests after its lease has expired, even though another leader is active ▶️ unsafe behavior.
- 👉 In distributed systems, you must assume any node can be paused arbitrarily long, while the rest of the system continues, so timing assumptions are unsafe.

#### Response time guarantees

- Some environments require software to respond within strict deadlines — these are hard real-time systems (e.g., aircraft control, car airbags, robots). In such cases, even brief pauses (like a GC pause) can be catastrophic.
- To achieve real-time guarantees, the entire software stack must cooperate:
  - **Real-time OS** (RTOS): ensures processes get **CPU time** at **fixed intervals**.
  - **Libraries**: must document **worst-case execution times**.
  - **Memory allocation**: often restricted; real-time garbage collectors exist but need strict limits.
  - **Testing/measurement**: extensive validation required.
- This makes real-time system development expensive and restrictive, limiting language, libraries, and tools. `Real-time ≠ high-performance` — **throughput** is often sacrificed for **predictability** 🤷‍♀️.
- Most server-side data systems don’t justify these costs, so they operate in non-real-time environments and must tolerate pauses and clock instability.

#### Limiting the impact of garbage collection

- The impact of process pauses (like GC) can be reduced without full real-time guarantees.
- **Planned GC handling**: Treat GC pauses as short outages — warn the app, stop new requests, finish current ones, then pause safely. Clients see no delay.
- **Selective GC & restarts**: Use GC mainly for short-lived objects and periodically restart processes before long-lived objects force a heavy GC. Restart one node at a time, shifting traffic elsewhere (like rolling upgrades).
- 👉 These methods don’t eliminate GC pauses but make their effects much less harmful.

## Knowledge, Truth, and Lies