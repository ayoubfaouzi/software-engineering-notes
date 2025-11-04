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

## Unbundling Databases

At a high level, **databases**, **Hadoop**, and **operating systems** share a common goal: managing data storage and processing. While Unix provides low-level abstractions (files, pipes), relational databases offer high-level abstractions (SQL, transactions) that hide complexity such as concurrency and recovery. The philosophical tension between Unix simplicity and database abstraction continues today - manifesting in movements like **NoSQL**, which adopt Unix-style low-level flexibility for distributed systems.

### Composing Data Storage Technologies

Databases internally implement mechanisms like:
- **Secondary indexes** (for efficient lookups)
- **Materialized views** (cached query results)
- **Replication logs** (for data consistency across nodes)
- **Full-text search indexes**

These functions parallel how **batch** and **stream processors** manage derived data systems.  
For example, creating a database index resembles setting up a **new replica** or **bootstrapping change data capture**—it involves scanning existing data and keeping updates synchronized.

Thus, the entire dataflow of an organization can be viewed as one large “meta-database,” where batch and stream processors maintain various derived data views—akin to indexes or materialized views—across multiple systems.

### Two Integration Philosophies

#### 1. **Federated Databases (Unifying Reads)**
- Provide a **single query interface** across diverse storage engines.
- Example: PostgreSQL’s **Foreign Data Wrappers**.
- Follows the relational tradition — high-level unified querying over heterogeneous systems.
- Suitable for combining data for read-only purposes.

#### 2. **Unbundled Databases (Unifying Writes)**
- Focus on **keeping multiple systems in sync**.
- Instead of distributed transactions, rely on **asynchronous event logs** and **idempotent writes**.
- Inspired by Unix’s philosophy: small, composable tools communicating via uniform APIs (like pipes).

#### Benefits of Log-Based (Unbundled) Integration
1. **System robustness** — asynchronous event logs decouple components, preventing local failures from escalating system-wide.
2. **Team autonomy** — each service or data system can evolve independently, connected via durable, ordered logs that maintain consistency.

### Unbundled vs Integrated Systems

- **Databases remain essential** for maintaining local state and serving queries from batch/stream outputs.
- **Specialized engines** (e.g., MPP warehouses) will continue to exist for niche workloads.
- Integrated systems can offer better performance and easier management for specific needs.
- **Unbundling’s goal** is not outperforming single databases, but enabling **breadth** — combining multiple specialized systems for a wider range of workloads.

Use integrated systems if one tool meets your needs; embrace unbundling when no single system fits all requirements.

### What’s Missing: The “Unix Shell” for Data Systems

We still lack a **declarative, high-level language** to compose unbundled data systems as easily as Unix commands.  
Ideally, one could define: `mysql | elasticsearch` to automatically replicate and index MySQL data into Elasticsearch, handling changes transparently—without custom glue code.

Future systems could extend this idea to **declarative caching and materialized views**, possibly using innovations like **differential dataflow**, bridging the gap between low-level unbundled tools and the declarative power of databases.
