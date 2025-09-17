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