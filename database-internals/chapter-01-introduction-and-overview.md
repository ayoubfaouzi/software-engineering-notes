# Chapter 1: Introduction and Overview

- DBMS can serve different purposes: some are used primarily for **temporary hot data**, some serve as a **long-lived cold storage**, some allow **complex analytical queries**, some only allow accessing values by the **key**, some are optimized to store **time-series** data, and some store **large blobs** efficiently.
- There are many ways DBMSs can be classified. For example, in terms of a **storage medium** (*Memory vs Disk-Based*) or **layout** (*“Column- vs Row-Oriented*). Some sources group them into three major categories:
    - *Online transaction processing (OLTP) databases* - These handle a large number of user-facing requests and transactions. Queries are often **predefined** and **short-lived**.
    - *Online analytical processing (OLAP) databases* - These handle **complex aggregations**. OLAP databases are often used for **analytics** and data **warehousing**, and are capable of handling complex, **long-running ad hoc queries**.
    - **Hybrid transactional and analytical processing (HTAP)**  - These databases combine properties of both OLTP and OLAP stores.
- There are many other terms and classifications: key-value stores, relational databases, document-oriented stores, and graph databases.

## DBMS Architecture

- DBMS use a **client/server** model, where database system instances (nodes) take the role of servers, and application instances take the role of clients.
<p align="center"><img src="assets/dbms-architecture.png" width="300px"></p>

- Client requests arrive through the **transport subsystem**. Requests come in the form of queries, most often expressed in some query language. The transport subsystem is also responsible for **communication with other nodes** in the database cluster.
- Upon receipt, the transport subsystem hands the query over to a **query processor**, which parses, interprets, and validates it.
- The parsed query is passed to the **query optimizer**, which first eliminates impossible and redundant parts of the query, and then attempts to find the most efficient way to execute it based on internal statistics and data placement.
- The query is usually presented in the form of an **execution plan** (or **query plan**): a sequence of operations that have to be carried out for its results to be considered complete. Since the same query can be satisfied using different execution plans that
can vary in efficiency, the optimizer picks the best available plan.
- The execution plan is handled by the **execution engine**, which collects the results of the execution of local and remote operations. **Remote** execution can involve writing and reading data to and from **other nodes** in the cluster, and replication. **Local** queries (coming directly from clients or from other nodes) are executed by the storage engine. The storage engine has several components with dedicated responsibilities:
  - **Transaction manager** - This manager schedules transactions and ensures they cannot leave the database in a logically inconsistent state.
  - **Lock manager** - This manager locks on the database objects for the running transactions, ensuring that concurrent operations do not violate physical data integrity.
  - **Access methods** (storage structures) - These manage access and organizing data on disk. Access methods include heap files and storage structures such as B-Trees or LSM Trees.
  - **Buffer manager** - This manager caches data pages in memory.
  - **Recovery manager** - This manager maintains the operation log and restoring the system state in case of a failure.

## Memory- Versus Disk-Based DBMS