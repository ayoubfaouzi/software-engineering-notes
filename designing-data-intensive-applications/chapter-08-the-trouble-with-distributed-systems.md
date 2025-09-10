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

