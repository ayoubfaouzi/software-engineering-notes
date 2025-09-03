# Chapter 7. Transactions

- Data systems face many possible failures: software/hardware crashes, application crashes mid-operation, network interruptions, concurrent writes, partial updates, and race conditions.
- To handle these, systems need fault tolerance, which is complex and requires careful design and testing.
- **Transactions simplify** error handling by **grouping reads/writes** into a **single logical unit** that either **fully succeeds** (commit) or **fully fails** (rollback). This frees applications from dealing with partial failures or concurrency issues, since the database provides safety guarantees.
- Transactions are not always necessary; sometimes **weaker guarantees** are chosen for **performance** or **availability** 🤷‍♂️.
- To decide whether transactions are needed, one must understand what guarantees they provide and their trade-offs ⚖️.

## The Slippery Concept of a Transaction

- Most relational databases today (MySQL, PostgreSQL, Oracle, SQL Server, etc.) implement transactions in a style originating from *IBM’s System R* in the 1970s. While details have evolved, the core ideas remain consistent.
- When NoSQL databases rose in the late 2000s, they introduced new data models and emphasized replication and partitioning. Transactions were **often dropped** or **weakened**, fueling a belief that transactions **limit scalability**, while traditional vendors promoted them as essential for critical applications 🤷.
- In reality, both claims are exaggerated. Transactions are neither universally necessary nor inherently incompatible with scalability — they are a trade-off ⚖️, offering **strong safety guarantees** but with **performance** and a**vailability costs**.
- To reason about these trade-offs, one must understand exactly what guarantees transactions provide under normal and failure conditions 🤔.

### The Meaning of ACID

- The concept of `ACID` (**Atomicity**, **Consistency**, **Isolation**, **Durability**) was introduced in 1983 to define database fault-tolerance guarantees. However, implementations differ greatly — especially around **isolation** —so “*ACID compliance*” has become more of a **marketing** term than a precise guarantee 🫤.
- As a contrast, `BASE` (**Basically Available**, **Soft state**, **Eventual consistency**) emerged to describe *non-ACID* systems, but it’s even vaguer and often just means “*not ACID*.” 🫨.

#### Atomicity

- Atomicity has slightly different meanings depending on context:
  - In **multithreading**: An atomic operation cannot be observed in a **partial state** — only before or after completion.
  - In **ACID transactions**: Atomicity ensures that if a transaction fails partway (due to crash, network error, full disk, or constraint violation), **all writes** are **discarded/undone**.
- Without atomicity, partial changes make it unclear which updates succeeded, and retries risk duplicates or corruption. With atomicity, aborted transactions guarantee no changes, making retries safe.
- 👉 ACID atomicity = **all-or-nothing** execution of a transaction (better thought of as *abortability*).

#### Consistency

- The word consistency is terribly overloaded 🫣:
  - In Chapter 5 we discussed **replica consistency** and the issue of **eventual consistency** that arises in **asynchronously** replicated systems.
  - **Consistent hashing** is an approach to partitioning that some systems use for rebalancing.
  - In the **CAP** theorem, the word consistency is used to mean **linearizability** (see Chapter 9).
  - In the context of **ACID**, consistency refers to an application-specific notion of the database being in a “**good state**”.
- It’s unfortunate that the same word is used with at least four different meanings 🤷‍♂️.
- ACID consistency means that application-defined invariants (rules that must always hold true, e.g., credits = debits in accounting) are preserved across transactions.
- The **application** is responsible for ensuring its transactions maintain these **invariants**.
- The database cannot guarantee consistency on its own —it only **enforces limited constraints** (like foreign keys or uniqueness).
- In contrast, *atomicity*, *isolation*, and *durability* are properties provided by the **database**.
- 👉 Therefore, consistency in ACID is really an **application concern**, not a **true database property** — some argue the “`C`” doesn’t even belong in ACID 😐.

#### Isolation

- Isolation in ACID ensures that **concurrent transactions** don’t **interfere with each other**, avoiding **race conditions** (e.g., two clients incrementing the same counter but ending up with the wrong result).
- Formally, isolation is defined as **serializability**: transactions behave as if they ran **one after another**, even if they actually run concurrently.
- In practice, full serializable isolation is rarely used due to performance costs.
- Many databases (e.g., *Oracle 11g*) use **weaker guarantees** like **snapshot isolation**, even when labeled “serializable.”
- 👉 Isolation = transactions **don’t step on each other’s toes**, ideally serializable, but often weaker in practice.

#### Durability

- Durability in ACID means that **once a transaction commits**, its data will **not be lost** — even after crashes or hardware faults.
- **Single-node** databases: Ensure durability by writing to nonvolatile storage (HDD/SSD) and often use a WAL for recovery.
- **Replicated** databases: Ensure durability by replicating data to multiple nodes **before confirming commit**.
- A database must wait for these writes/replications to finish before declaring success.
- Perfect durability is impossible — if all disks and backups are destroyed, data is lost 🤓.
- 👉 Durability = **committed data survives crashes**, but only as strong as your **storage + replication setup**.

### Single-Object and Multi-Object Operations

- 👉 Atomicity = rollback on failure; Isolation = no half-visible states in concurrent access.
- Example (email app):
  - Denormalized unread counter can get inconsistent (message inserted but counter not updated).
  - Isolation prevents anomalies (users see consistent state).
  - Atomicity ensures if counter update fails, the message insert is rolled back.
- Implementation:
  - Relational DBs: Use `BEGIN TRANSACTION` … `COMMIT` (tied to a client connection).
  - Many nonrelational DBs: Don’t support true transactions — multi-object operations may partially succeed, leaving inconsistencies.

#### Single-object writes

- When writing a single object (e.g., 20KB JSON doc), DBs must avoid:
  - Partial writes (cut-off JSON fragments).
  - Corrupted values (old + new spliced together).
  - Reads seeing half-updated data.
- Storage engines ensure **atomic single-object writes** (via crash-recovery logs) and isolation (via per-object locks).
- Extra features:
  - Some DBs provide atomic operations like:
    - increment (avoids read-modify-write race conditions).
    - compare-and-set (update only if value hasn’t changed).
  - These help prevent lost updates in concurrent writes.
- ⚠️ But:
  - These are **not true transactions** (multi-object, grouped operations).
  - Calling them “*lightweight transactions*” or “ACID” is **marketing**, not accurate 🤷‍♀️.
- 👉 Atomic writes & isolation are almost always guaranteed per object, but real transactions cover **multiple objects together**.

#### The need for multi-object transactions

- Many distributed datastores dropped them for simplicity, availability, and performance, but they are still possible to implement.
- Single-object operations (insert, update, delete) are sometimes enough, but many cases need coordinated multi-object writes:
  - Relational DBs: enforcing **foreign keys and references** across rows/tables.
  - Document DBs: usually fine with single-object updates, but denormalization (due to lack of joins) often requires updating **multiple documents** at once.
  - **Secondary indexes**: must be updated with the base record; without transactions, you risk inconsistent index states.
- Without multi-object transactions:
  - Apps can still work, but error handling is much harder.
  - Lack of isolation causes concurrency anomalies.
  - Transactions simplify correctness by handling these automatically.
- 👉 Transactions aren’t strictly required, but they greatly **reduce complexity** and prevent data inconsistencies in relational, document, and indexed databases.

### Weak Isolation Levels

- Transactions that don’t access the same data can run in parallel safely.
- Concurrency issues happen when:
  - One transaction reads data modified by another.
  - Two transactions try to modify the same data at once.
- Concurrency 🐛 are:
  - Rare and timing-dependent → hard to reproduce in testing.
  - Difficult to reason about in large apps with many users.
- Role of isolation:
  - Goal: hide concurrency from developers by making execution look serial.
  - Serializable isolation: strongest guarantee, but has performance costs.
  - Many databases use weaker isolation levels → fewer guarantees, easier performance, but risk subtle and dangerous bugs.
- Real-world impact:
  - Weak isolation has caused:
    - Financial losses.
    - Auditor investigations.
    - Customer data corruption.
  - Even “ACID” databases may use weak isolation, so ACID ≠ full safety 😲.
- Key takeaway:
  - Don’t blindly trust tools.
  - Understand concurrency problems and isolation levels.
  - Choose the right isolation level for your app’s needs.

#### Read Committed

- The most basic level of transaction isolation is read committed. It makes two guarantees:
  1. When reading from the database, you will only see data that has been committed (**no dirty reads**).
  2. When writing to the database, you will only overwrite data that has been committed (**no dirty writes**).

##### No dirty reads

