# Chapter 11. The Future of Data Systems

The final chapter shifts focus from how systems **are** to how they **should be**, proposing ideas to improve the design of reliable, scalable, and maintainable applications. It synthesizes the book’s main themes—**fault tolerance, scalability, and maintainability**—and explores how to build future systems that are **robust, correct, evolvable**, and **beneficial to humanity**.

## Data Integration

- There are **multiple valid solutions** for any problem, each with trade-offs (e.g., log-structured vs. B-tree vs. column stores; single-leader vs. multi-leader replication).
- No universal solution fits all use cases 🤷‍♀️; **different workloads require specialized tools**.
- Real-world applications often need to **combine multiple systems**, leading to complex integration challenges.

### Combining Specialized Tools Through Derived Data

- Example: integrating an **OLTP database** with a **full-text search index**.
- As systems multiply (databases, search engines, caches, ML models, etc.), synchronization and consistency become more **difficult**.
- Data integration requires clear reasoning about **dataflows** — where data is written first and which systems derive from it.

#### Reasoning About Dataflows

- **CDC** or **event sourcing** can help maintain derived data in sync with a system of record.
- A single system deciding the **total order of writes** simplifies synchronization and prevents permanent inconsistencies.
- This mirrors the **state machine replication** model: process writes in a consistent order to ensure deterministic results.

#### Derived Data vs. Distributed Transactions

| Aspect | Distributed Transactions | Derived Data (Log-based) |
|--------|---------------------------|---------------------------|
| Ordering | Locks & 2PL | Log ordering |
| Commit Model | Atomic commit | Deterministic retry + idempotence |
| Consistency | Linearizable | Eventually consistent |
| Performance | High coordination cost | Asynchronous, scalable |

- Distributed transactions (e.g., XA/2PC) provide strong guarantees but **poor fault tolerance and performance**.
- Log-based derived data systems are **more promising** for integration, though they lack immediate consistency guarantees.
- The goal: find a **middle ground**—asynchronous systems with stronger correctness properties.

#### The Limits of Total Ordering

Constructing a **totally ordered event log** works for small systems but faces limits as scale increases:
1. **Throughput**: a single leader cannot handle all events.
2. **Geographic Distribution**: multiple data centers introduce ambiguous orderings.
3. **Microservices**: independent services lack shared durable state.
4. **Offline Clients**: different event orders between client and server.

Total ordering = **consensus problem**.  
Scaling consensus across partitions or datacenters remains an **open research challenge**.

#### Ordering Events to Capture Causality

When total order is infeasible, causal relationships must still be preserved:
- Example: a user unfriends someone, then posts a message. If causality is lost, the ex-partner might still get notified.
- The challenge: maintain **causal dependencies** between related events stored in different systems.

#### Potential Solutions

- **Logical timestamps**: give order without coordination, but require extra metadata.
- **Causal references**: events can reference prior events that influenced them.
- **Conflict resolution algorithms**: handle reordering but not side effects.

Over time, application patterns may evolve to **efficiently capture causality** and maintain derived state **without total ordering bottlenecks**.

### Batch and Stream Processing

The main goal of **data integration** is to ensure that data is transformed into the right form and delivered to the right places.  
This involves:
- **Consuming inputs**
- **Transforming, joining, filtering, aggregating**
- **Training and evaluating models**
- **Writing outputs**

**Batch** and **stream processing** systems are the core tools that achieve these transformations.

#### Derived Datasets

The outputs of these processing systems are **derived datasets**, such as:
- Search indexes  
- Materialized views  
- Recommendations  
- Aggregate metrics  

Batch and stream processing share similar principles; the key distinction is:
- **Batch processing** → finite datasets  
- **Stream processing** → unbounded, continuous data

Modern systems blur this line:
- **Apache Spark**: stream processing via *microbatches*  
- **Apache Flink**: batch processing built on top of a streaming model  

#### Maintaining Derived State

- Both batch and stream processing emphasize **deterministic, functional operations**:
  - Pure functions (output depends only on input)
  - Immutable inputs and append-only outputs
  - No side effects other than explicit outputs
- Stream processors extend this model with **managed, fault-tolerant state**.
  - 👍 Improves **fault tolerance** through **idempotent and deterministic** processing  
  - 👍 Simplifies **reasoning about dataflows** across an organization  
  - 👍 Makes derived data (e.g., indexes, caches, models) easier to maintain through **functional pipelines**
- Maintaining derived systems **asynchronously** increases robustness:
  - 👍 Failures are isolated to local components  
  - 👍 Avoids failure amplification common in **distributed transactions**

Cross-partition indexes (e.g., term/document partitioning) are most scalable when updated asynchronously.

#### Reprocessing data for application evolution

- Batch and stream processing support **system evolution**:
  - **Stream processing** → low-latency updates  
  - **Batch processing** → large-scale reprocessing of historical data  
- Benefits of Reprocessing:
  - Enables major **schema and model changes**, not just incremental ones  
  - Supports **gradual evolution**: maintain old and new views side-by-side
  - Allows **canary migrations** (testing new views with limited users)
  - Each step is **reversible**, reducing migration risk and improving confidence

#### The Lambda Architecture

The **Lambda Architecture** combines batch and stream processing:
- **Stream layer**: processes recent data for quick, approximate updates  
- **Batch layer**: reprocesses historical data for accurate, corrected results  
- Data is stored as **immutable, append-only events**
- **Advantages**:
  - Encourages **event-sourced** dataflows and **derived views**
  - Promotes **fault tolerance** through immutability and reprocessing
- **Limitations**:
  - Duplicated logic between batch and stream systems  
  - Operational complexity (two frameworks to maintain)  
  - Difficulty merging outputs from batch and stream pipelines  
  - Expensive to frequently reprocess large datasets  
  - Incremental batch updates add **temporal complexity** and blur distinctions between layers

#### Unifying Batch and Stream Processing

Modern architectures overcome Lambda’s downsides by **unifying batch and streaming**:
- One engine handles **both historical reprocessing** and **live event processing**
- Simplifies code reuse, operations, and consistency guarantees
- **Required Features**:
  1. **Event replay**: ability to reprocess historical data through the same pipeline, e.g., Kafka’s log replay or reading from distributed filesystems (HDFS)
  2. **Exactly-once semantics**: ensures consistent outputs despite failures  
  3. **Event-time windowing**: processes data based on event timestamps, not processing time. Supported by frameworks like **Apache Beam**, **Flink**, and **Google Cloud Dataflow**
