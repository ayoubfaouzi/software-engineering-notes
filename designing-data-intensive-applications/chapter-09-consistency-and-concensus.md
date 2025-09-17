# Chapter 9. Consistency and Consensus

- Distributed systems can fail in many ways (lost, delayed, or duplicated messages; clock issues; node pauses or crashes).
- Instead of letting services fail, fault tolerance aims to keep them functioning despite such problems.
- The best approach is to build **general-purpose abstractions** — like **transactions** or **consensus** — that hide failures from applications.
- Consensus, in particular, ensures **all nodes agree on decisions** (e.g., leader election), preventing issues like **split brain**.
- This chapter explores the guarantees and abstractions distributed systems can provide, along with their fundamental limits, offering an informal overview while pointing to deeper research in the literature.

## Consistency Guarantees

- Replication in databases leads to temporary inconsistencies between nodes, regardless of replication method. Most systems offer **eventual consistency** (better termed **convergence**), meaning replicas will agree eventually, but without guarantees on when.
- This weak model complicates application development, since it differs from the predictable behavior of variables in single-threaded programs. Subtle bugs often appear only under faults or high concurrency.
- The chapter explores **stronger consistency models**, which trade performance and fault tolerance for easier correctness. It introduces:
  - **Linearizability** (a strong model) and its pros/cons.
  - **Ordering guarantees** (causality, total order).
  - **Distributed transactions** and **consensus**, focusing on atomic commits.
- Consistency models, while related to transaction isolation levels, address different challenges: replica coordination vs. concurrency control.

## Linearizability

- Linearizability (also called **strong** or **atomic consistency**) is the guarantee that a distributed system behaves as if there were only **one copy of the data**. Every operation appears instantaneous and immediately visible to all clients.
- In a linearizable system, once a **write completes**, all subsequent reads return the **updated value** (no stale data).
- It provides a **recency guarantee**: reads always reflect the most recent completed write.
- Without linearizability (e.g., in **eventually consistent** systems), different replicas may return different answers, leading to anomalies.
- 👉 After Alice sees the final score of a match, Bob refreshes later but still sees an outdated result — a clear violation of linearizability.

### What Makes a System Linearizable?

- **Reads and Writes**:
  - Reads that happen strictly before a write must return the old value.
  - Reads that happen strictly after a write must return the new value.
  - Reads concurrent with a write may return either, but once one client sees the new value, all later reads must also return it (no “flip-flopping”).
- **Atomicity model**:
  - Each operation is pictured as a bar (request to response) with a point inside representing the moment it logically takes effect.
  - The sequence of these points must form a valid forward-moving history: once a value is written or read, all subsequent reads reflect it until overwritten again.
- **Concurrency quirks**:
  - Operations may complete in an **order different** from request arrival, and that’s okay as long as the resulting sequence is valid.
  - Responses can be **delayed**, so a read might return a new value before the writer has even received its own success acknowledgment.
  - CAS helps ensure updates aren’t lost to concurrent modifications.
- **Violations**:
  - If a client reads an older value after another client has already read a newer one, linearizability is broken (similar to the Alice & Bob example).
- **Testing linearizability**:
  - You can verify it (at high cost) by recording request/response timings and checking if operations can be ordered into a valid sequential history.
- 👉 Linearizability = the illusion of a single, atomic timeline of operations, preserving recency guarantees across concurrent clients.

#### Linearizability Versus Serializability

- TThe key difference between **linearizability** and **serializability** is often confused because both involve the idea of a **sequential order**. However, they are distinct guarantees:
  - Serializability is an **isolation property** for **transactions** (groups of operations). It guarantees that the result of multiple concurrent transactions is equivalent to if they had executed one after another (in some serial order). It is about the correctness of a group of operations.
  - Linearizability is a **recency guarantee** for a single object. It ensures that once a write completes, all subsequent reads will see that value until it is overwritten. It is about the freshness of a single operation.
A system can be:
  - Both: This combination is called **strict serializability**.
  - Serializable but not linearizable: For example, databases using SSI use consistent snapshots, so reads may not see the very latest write.
  - Linearizable but not serializable: A system can guarantee fresh reads on individual objects but not protect against multi-object transaction anomalies like **write skew**.
- 👉 Serializability is about transactions appearing atomic; linearizability is about reads and writes appearing instantaneous.

### Relying on Linearizability

- Linearizability is essential in distributed systems for three main use cases where strict agreement on the most recent state of a single object is necessary:
  - **Locking and Leader Election**: Systems that use a distributed lock or single-leader replication must have a linearizable lock to prevent **split-brain** scenarios. All nodes must agree on which node currently holds the lock. Coordination services like `ZooKeeper` and `etcd` provide this linearizable foundation.
  - **Enforcing Uniqueness Constraints**: Hard constraints, such as ensuring a username or email address is unique, require linearizability. This is effectively like acquiring a lock on a value. If two users try to create the same username concurrently, a linearizable system guarantees that only one will succeed. This also applies to constraints like preventing a bank account from going negative.
  - **Cross-Channel Timing Dependencies**: Linearizability prevents race conditions when two different communication channels depend on the same data. For example, if a web server writes a photo to storage (channel 1) and then sends a message via a queue (channel 2) for another service to process it, a non-linearizable storage service might cause the processor to read a stale or missing file, leading to inconsistencies.

### Implementing Linearizable Systems

- The naive way is a single copy of data, but that isn’t fault-tolerant 🤷.
- So, replication methods are used — but not all support linearizability:
- **Single-leader replication**:
  - Can be linearizable if **reads** go to the **leader** or **synchronously** updated followers.
  - Risks: leader uncertainty (split brain), async replication causing lost writes, or databases intentionally using weaker models.
- **Consensus algorithms** (e.g., `Raft`, `Paxos`, `ZooKeeper`, `etcd`)
  - Provide linearizability by preventing split brain, ensuring durability, and keeping replicas consistent.
  - Safest way to implement linearizable storage 👍.
- **Multi-leader replication**:
  - Not linearizable: multiple leaders accept concurrent writes → conflicts, async replication → no single-copy illusion.
- **Leaderless replication** (`Dynamo`-style):
  - Often claimed to be “strong” with quorum reads/writes (`w + r > n`), but still not strictly linearizable.
  - LWW with clocks → broken by **clock skew**.
  - **Sloppy quorums/hinted handoff** break linearizability further.
  - Even strict quorums can allow anomalies.

#### Linearizability and quorums

- Quorums (`r + w > n`) in `Dynamo`-style systems may seem to guarantee linearizability, but network delays can cause race conditions.
- Example:
  - Writer updates `x = 1` across 3 replicas (w = 3).
  - Reader A (r = 2) sees `1`.
  - Later, Reader B (r = 2) still sees 0.
  - Even though quorum conditions are satisfied, this violates linearizability.
- To achieve linearizability:
  - Readers must perform **synchronous read repair**.
  - Writers must read the latest quorum state before writing.
- ⚖️:
  - Riak skips synchronous read repair → better performance, but no linearizability.
  - `Cassandra` does synchronous repair for quorum reads, but still breaks linearizability under concurrent writes (LWW).
- Limitation: Only reads and writes can be made linearizable this way — CAS requires full consensus.
- Conclusion: Dynamo-style leaderless replication should be assumed not linearizable.

### The Cost of Linearizability

- Multi-leader replication:
  - Each datacenter can keep accepting writes independently.
  - Writes are queued and later exchanged once connectivity is restored.
  - System stays available, but may have conflicts to resolve later.
- Single-leader replication:
  - All writes and linearizable reads must go through the leader.
  - If the leader’s DC is unreachable, follower DC clients cannot perform writes or linearizable reads (only stale reads).
  - This means **outages** for clients that can’t reach the leader until the network recovers.
- ⚖️ Multi-leader offers higher availability under network partitions, while single-leader ensures strict consistency but risks unavailability.

#### The CAP theorem

- Linearizable systems: If replicas are disconnected, they cannot safely process requests → they must block or return errors → unavailability.
- Non-linearizable systems (e.g., multi-leader): Replicas can keep working independently during disconnections → higher availability, but weaker consistency.
- This trade-off is the essence of the CAP theorem: during a **network partition**, you must choose between **Consistency** (linearizability) and **Availability**.
- CAP is often misrepresented as “*pick 2 of 3,*” but the real meaning is:
  - In normal conditions: you can have both consistency + availability.
  - Under partition: you must sacrifice one.
- CAP was influential (helped inspire `NoSQL` systems), but:
  - It only covers linearizability + partitions, not **other faults** (delays, crashes, etc.).
  - Its definitions of availability are confusing 🤷‍♀️.
- It has limited practical value today, mostly of historical interest, as more precise results now exist.

#### Linearizability and network delays

- Few systems are truly linearizable — even RAM on modern **multi-core CPUs** isn’t.
  - Each CPU core has caches and store buffers. Writes are propagated asynchronously to main memory.
  - This boosts performance but breaks linearizability unless memory **fences** are used.
- Reason for dropping linearizability:
  - Not CAP or fault tolerance.
  - Purely performance — linearizability is always slower, not just during faults.
- Distributed databases often avoid linearizability for the same reason: to improve speed and latency.
- Theory (*Attiya & Welch*):
  - Linearizable storage is fundamentally slow.
  - Response times are at least proportional to network delay uncertainty.
  - No faster algorithm exists for linearizability.
- ⚖️:
  - Linearizability = correctness but high latency.
  - Weaker consistency = much faster, better for latency-sensitive systems.

## Ordering Guarantees