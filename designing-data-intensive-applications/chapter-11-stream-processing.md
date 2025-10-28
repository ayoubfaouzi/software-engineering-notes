# Chapter 11. Stream Processing

- Chapter 10 introduced **batch processing**, where a **finite** set of input files is processed to produce derived output files. This model is effective for building **search indexes**, **recommendation systems**, and **analytics pipelines**, since the data is **bounded** and processing knows when it’s complete.
- However, **real-world data** is often **unbounded** 🫤, it continuously arrives over time (e.g., user activity, logs, or sensor data). Because the dataset is never “finished,” batch systems must **artificially segment** the input into time-based chunks (e.g., hourly or daily batches).
- This segmentation introduces **latency**: if the batch runs daily, users only see updates once per day. Increasing the frequency (hourly, per-second, or continuously) reduces this delay, leading to the concept of **stream processing**, which processes data **as it arrives** rather than in discrete batches.
- A **stream** represents **incrementally available data** over time, a concept found across systems like Unix I/O streams, lazy lists in programming, file APIs, TCP connections, and multimedia delivery.

## Transmitting Event Streams

- In **batch processing**, jobs read input files and produce output files. In **stream processing**, the equivalent concept is an **event stream** — a continuous sequence of **records** (called **events**) representing things that happened over time.
- An **event** is a small, immutable object describing an occurrence (e.g., a page view, purchase, or sensor reading). It usually includes a **timestamp** indicating when it occurred.
- Like batch files, which can be read by multiple jobs, events in streaming systems can be **consumed** by multiple **subscribers**. The **producers** (*publishers*) send events to **topics** or **streams**, grouping related events together.
- While a file or database could, in theory, connect producers and consumers (with consumers polling for new data), **polling** becomes inefficient as latency requirements drop — frequent polls waste resources when no new data is available 🫤. Thus, **push-based** systems that **notify consumers** of new events are preferred.
- Traditional **databases** offer limited support for notifications (e.g., via **triggers**), which are inflexible and not designed for continuous event delivery. To fill this gap, specialized systems were developed for **event notification and streaming**, forming the backbone of modern stream processing architectures.

### Messaging Systems

- To evaluate messaging systems, two key questions arise:
  1. **Handling Producer-Consumer Speed Mismatch**
     - If producers send messages **faster** than consumers can handle, the system can:
       - **Drop messages**
       - **Buffer messages** in a queue
       - **Apply backpressure** (block the producer until the consumer catches up)
     - For buffered systems, it’s important to know what happens when queues grow too **large**:
       - Do they crash when memory runs out?
       - Do they spill to disk, and how does that affect performance?
  2. **Handling Failures and Message Loss**
     - Durability can be achieved through **disk writes** or **replication**, which add overhead.
     - ⚖️ Allowing message loss yields higher throughput and lower latency.
     - **Acceptability of loss** depends on context:
       - For periodic sensor data, occasional losses may be fine.
       - For event counting or critical data, message loss leads to incorrect results.

#### Direct messaging from producers to consumers

- Some messaging systems bypass intermediaries and use **direct network communication** between producers and consumers. Examples include:
  - **UDP Multicast**
    - Common in finance (e.g., stock market feeds) for **low-latency** message delivery.
    - UDP is unreliable, but higher-level protocols can **recover lost packets** by retransmission.
  - **Brokerless Messaging (ZeroMQ, nanomsg)**
    - Implement **publish/subscribe** models over **TCP or IP multicast**, avoiding centralized brokers.
  - **Metrics Collection (StatsD, Brubeck)**
    - Use **unreliable UDP** to send lightweight metrics across machines.
    - Since UDP can drop packets, results (e.g., counters) are **approximate** at best.
  - **Webhooks (HTTP/RPC Push Model)**
    - A producer pushes messages directly to a consumer’s **callback URL** whenever an event occurs.
- These **direct messaging systems** perform well for **low-latenc**y or lightweight applications but have **limited fault tolerance**:
  - Applications must handle **message loss** explicitly.
  - They assume **constant online availability** of producers and consumers.
  - If a consumer goes offline, it **misses messages** sent during downtime.
  - Some protocols retry failed deliveries, but if a **producer crashes**, any buffered messages waiting for retry are **lost**.

#### Message Brokers

- **Message brokers** (or **message queues**) act as intermediaries between producers and consumers, functioning like a specialized **database for message streams**.
- They run as **servers**, with producers and consumers connecting as **clients**:
  - **Producers** write messages to the broker.
  - **Consumers** read messages from it.
- 🔑 Characteristics:
  - **Centralized durability and reliability:**
    - The broker manages persistence and client disconnections.
    - Some brokers keep messages **in memory** only.
    - Others **persist to disk** to survive crashes.
  - **Handling slow consumers:**
    - Brokers typically use **unbounded queueing** (buffering messages) rather than dropping them or applying backpressure — though this behavior can be **configured**.
  - **Asynchronous delivery model:**
    - Producers only wait for **acknowledgment** that the broker has **buffered the message**, not for consumers to process it.
    - Messages are then delivered **later**, possibly immediately or after delays if **queues backlog**.
- 👉 Message brokers **decouple producers and consumers**, providing **fault tolerance**, **durability**, and **asynchronous communication**—but at the cost of **potential queuing delays**.

#### Message brokers compared to databases

- Some **message brokers** can participate in **two-phase commits** (via **XA** or **JTA**), making them somewhat similar to **databases**, though several key differences remain:
  - **Data Retention:**
    - **Databases** store data until explicitly deleted.
    - **Message brokers** usually **delete messages after delivery**, making them unsuitable for long-term storage.
  - **Working Set Size:**
    - Brokers assume **short queues** and **small working sets**.
    - When consumers are slow and queues grow large, performance and throughput **degrade** - especially if messages spill to disk.
  - **Data Access:**
    - Databases offer **secondary indexes** and **search queries**.
    - Brokers allow clients to **subscribe** to specific **topics or patterns**, providing a more limited filtering mechanism.
  - **Change Awareness:**
    - Database queries give **point-in-time snapshots**; clients must **poll** to detect updates.
    - Brokers **push notifications** when new messages arrive, providing **real-time updates**.
- This traditional model is defined by **JMS** and **AMQP** standards and implemented in systems such as: **RabbitMQ**, **ActiveMQ**, **Google Cloud Pub/Sub**.

#### Multiple consumers

- When multiple consumers read from the same topic, **two primary messaging patterns** are used:
- 1️⃣ **Load Balancing**
  - **Each message is delivered to only one consumer** within the group.
  - Used to **distribute workload** and **parallelize message processing**.
  - Ideal when messages are **expensive to process**.
  - Implemented as:
    - **Multiple clients consuming from the same queue** in **AMQP**.
    - **Shared subscriptions** in **JMS**.
- 2️⃣ **Fan-Out**
  - **Each message is delivered to all consumers**.
  - Enables independent consumers to **receive the same data stream**, similar to multiple batch jobs reading the same file.
  - Implemented via:
    - **Topic subscriptions** in **JMS**.
    - **Exchange bindings** in **AMQP**.
- 3️⃣ **Combined Pattern**
  - Groups of consumers can **each subscribe** to a topic:
    - Each group receives **all messages**.
    - Within each group, **only one node** processes each message.
<p align="center"><img src="assets/pubsub-patterns.png" width="500px" height="auto"></p>

#### Acknowledgments and redelivery

- Message brokers use **acknowledgments (acks)** to ensure messages aren’t lost when consumers crash or disconnect:
- **Acknowledgment mechanism**:
  - Consumers must **explicitly** confirm message processing.
  - If the broker doesn’t receive an ack (e.g., due to crash or timeout), it **redelivers the message** to another consumer.
  - This prevents message loss but can lead to **duplicate processing** if the ack was lost after successful processing.
- With **load balancing**, redelivery can cause **message reordering**:
  - Example: Consumer 2 crashes while processing message *m3*.
  - The broker reassigns *m3* to Consumer 1, which is already processing *m4*.
  - Result: Consumer 1 processes messages in order *m4 → m3 → m5*, breaking the original order.
- Message brokers like `JMS` and `AMQP` try to preserve order, but **load balancing + redelivery** can still reorder messages.
- To maintain strict ordering ▶️ Use a **dedicated queue per consumer** (no load balancing).
- Reordering is acceptable if messages are **independent**, but problematic if **causal dependencies** exist between messages.

### Partitioned Logs

#### **From Transient Messaging to Durable Logs**
- Traditional messaging systems (like AMQP/JMS) are **transient**:
  - Messages are deleted after being consumed.
  - Adding a new consumer only gives access to **future messages**.
  - Reprocessing past data is impossible.
- Databases, by contrast, keep data **durably**, enabling reprocessing and experimentation.
**Log-based message brokers** combine both ideas — durable, append-only storage with streaming semantics.

### **How Log-Based Brokers Work**
- A **log** is an append-only sequence of records on disk.
- **Producers** append messages to the log’s end.
- **Consumers** read sequentially and wait for new messages (like `tail -f`).
- Logs are **partitioned** for scalability, allowing parallelism across machines.
- Each partition has **monotonically increasing offsets** that uniquely identify messages.
- Examples: **Apache Kafka**, **Amazon Kinesis Streams**, **Twitter DistributedLog**, **Google Pub/Sub**.
<p align="center"><img src="assets/partitioned-logs.png" width="500px" height="auto"></p>

👉 These systems achieve **high throughput (millions of msgs/sec)** via partitioning and replication.

### **Log-Based vs. Traditional Messaging**

| Feature | Traditional Brokers (AMQP/JMS) | Log-Based Brokers (Kafka-style) |
|----------|-------------------------------|--------------------------------|
| **Message retention** | Deleted after consumption | Retained for configured time |
| **Fan-out** | Requires duplication | Natural — multiple consumers can read same log |
| **Ordering** | Can reorder under load balancing | Strict within each partition |
| **Reprocessing** | Not possible (messages deleted) | Possible via offsets |
| **Consumer tracking** | Per-message acks | Simple per-partition offset tracking |
| **Throughput** | Degrades when queues grow | Constant (disk-based append only) |

### **Consumer Offsets and Recovery**
- Each consumer tracks a **current offset**.
- Broker doesn’t need per-message **acknowledgments** — just records offsets periodically.
- If a consumer crashes:
  - Another node can resume from the last recorded offset.
  - Some messages may reprocess if offsets weren’t committed yet (duplicates).

This mirrors **database replication**: offsets act like **log sequence numbers**.

### **Disk Space and Retention**
- Logs are **divided into segments**; old ones are periodically **deleted or archived**.
- Acts as a **bounded disk buffer** (circular buffer):
  - Example: A 6 TB drive at 150 MB/s can hold ~11 hours of messages.
  - Usually, retention is configured to keep **days or weeks** of data.
- Throughput remains **constant** regardless of history size (unlike memory-based brokers).

### **Handling Slow Consumers**
- The log-based approach uses **disk buffering** instead of backpressure or dropping.
- If a consumer falls too far behind (past retention window), it **misses messages**.
- Other consumers remain unaffected — great for **fault isolation** and **experimentation**.
- Operators can monitor consumer lag and take corrective action before loss.

### **Replaying Messages**
- Consuming from a log is **non-destructive** — the log remains intact.
- Consumers can **reset offsets** to reprocess historical data (e.g., replay yesterday’s data).
- Enables:
  - Rebuilding derived datasets.
  - Experimentation with new logic.
  - Easy recovery from bugs or data corruption.

### **Key Advantages**
- Durable storage + low-latency streaming.
- Simplified bookkeeping via offsets.
- Natural fan-out for multiple consumers.
- Replay and reprocessing support (like batch jobs).
- Operational robustness — consumers can come and go independently.

## Databases and Streams

- **Writes = Events** → Every database write is an event that can be captured and streamed.  
- **Replication Logs = Event Streams** → Leaders produce a stream of writes; followers consume it to stay in sync.  
- **State Machine Replication** → If replicas process the same events in the same order, they reach the same state.  
- **Key Insight** → Databases store **current state**; streams record **state changes** (history).  
- **Unified View** → Databases and event streams are two sides of the same coin — one shows the end result, the other shows how it got there.

### Keeping Systems in Sync

- Modern applications use multiple systems (OLTP databases, caches, search indexes, data warehouses), each with its own optimized copy of the data. These copies must be kept **synchronized** when data changes.
- **Common Synchronization Methods:**
  - **Batch ETL Processes:** Effective for data warehouses, involving periodic full dumps and bulk loading.
  - **Dual Writes:** The application code explicitly writes to all systems (e.g., database, then search index, then cache) when data changes.
- **Problems with Dual Writes:**
  - **Race Conditions:** Concurrent clients can cause systems to see writes in different orders, leading to permanent inconsistency (e.g., the database ends with value B while the search index ends with value A).
  - **Fault Tolerance:** If one write succeeds and another fails, the systems become inconsistent. Solving this requires an expensive atomic commit protocol.
- **The Core Issue:** Dual writes fail because there is no single system determining the order of writes across the different technologies (like having multiple leaders).
- **Proposed Solution:** A better approach is to have a single leader (e.g., the database) and make the other systems (like the search index) followers that consume its stream of changes.