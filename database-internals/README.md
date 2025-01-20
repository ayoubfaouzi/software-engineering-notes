Study notes taken from reading *Database Internals:  A Deep Dive into How Distributed Data Systems Work* by *Alex Petrov*.

# Part I: Storage Engines

Since the term **database management system (DBMS)** is quite bulky, throughout this book we use more compact terms, **database system** and **database**, to refer to the same concept.
- **DBMS** are apps built on top of **storage engines**, offering a schema, a query language, indexing, transactions, and many other useful features.
- This clear separation is 👍 because:
  - It enabled database developers to bootstrap database systems using existing storage engines, and concentrate on the other subsystems.
  - It opens up an opportunity to switch between different engines, potentially better suited for particular use cases.

## Comparing Databases

- To reduce the risk of an **expensive migration**, you can invest some time before you decide on a specific database to build confidence in its ability to meet your application’s needs.
- Even a superficial understanding of how each database works and what’s inside it can help you land a more weighted conclusion then looking at DB [comparison](https://db-engines.com/en/ranking) websites.
- If you’re searching for a database that would be a good fit for the workloads you have, the best thing you can do is to **simulate these workloads** against different DB systems **measure the performance** metrics that are important for you, and compare results.
- Some issues, especially when it comes to performance and scalability, **start showing only after some time** or as the **capacity grows** 😼.
- To compare databases, it’s helpful to understand the use case in great detail and define the current and anticipated **variables**, such as:
  - Schema and record sizes
  - Number of clients
  - Types of queries and access patterns
  - Rates of the read and write queries
  - Expected changes in any of these variables
- Knowing these variables can help to answer the following questions:
  - Does the database support the required queries?
  - Is this database able to handle the amount of data we’re planning to store?
  - How many read and write operations can a single node handle?
  - How many nodes should the system have?
  - How do we expand the cluster given the expected growth rate?
  - What is the maintenance process?
- One of the popular tools used for benchmarking, performance evaluation, and comparison is **Yahoo! Cloud Serving Benchmark (YCSB)**.
- Also worth checking the **Transaction Processing Performance Council (TPC)**.

## Understanding Trade-Offs

- There are many different approaches to storage engine design, and every implementation has its own upsides and downsides.
- Some are optimized for **low read or write latency**, some try to **maximize density** (the amount of stored data per node), and some concentrate on **operational simplicity**.