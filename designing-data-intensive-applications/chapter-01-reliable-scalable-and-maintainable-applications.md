# Chapter 1. Reliable, Scalable and Maintainable Applications

- A data-intensive application is typically built from standard building blocks which provide commonly needed functionality.
  - Store data so that they, or another application, can find it again later (**databases**),
  - Remember the result of an expensive operation, to speed up reads (**caches**),
  - Allow users to search data by keyword or filter it in various ways (**search indexes**),
  - Send a message to another process, to be handled asynchronously (**message queues**),
  - Observe what is happening, and act on events as they occur (**stream processing**),
  - Periodically crunch a large amount of accumulated data (**batch processing**).

## Thinking About Data Systems

- Even though each category of these systems serves a specific purpose, many new tools for data storage and processing have emerged that are optimized for a variety of different use cases:
  - For example, there are data stores that are also used as **message queues (Redis)**,
  - and there are **message queues with database-like durability guarantees (Kafka)**,
  - ▶️ the boundaries between the categories are becoming **blurred**.
- If you are designing a data system or service, a lot of tricky questions arise:
  - How do you ensure that the data remains correct and complete, even when things go wrong internally?
  - How do you provide consistently good performance to clients, even when parts of your system are degraded?
  - How do you scale to handle an increase in load? What does a good API for the service look like?
- We focus on three concerns that are important in most software systems:
  - Reliability
  - Scalability
  - Maintainability

## Reliability

- We understand reliability as meaning, roughly: _continuing to work correctly, even when things go wrong_.
- The things that can go wrong are called **faults**, and systems that anticipate faults and can cope with them are called **fault-tolerant**.
- Note that a fault is not the same as a **failure** for an overview of the terminology. A **fault** is usually defined as **one** component of the system **deviating from its spec**, whereas a **failure** is when the system as a **whole** stops providing the required service to the user.

### Hardware faults

- Hard disks crash, RAM becomes faulty, the power grid has a blackout, someone unplugs the wrong network cable. Anyone who has worked with **large data centers** can tell you that these things happen **all the time** when you have a lot of machines.
- Hard disks are reported as having a **mean time to failure (MTTF)** of about 10 to 50 years. Thus, on a storage cluster with 10,000 disks, we should expect on average **one disk to die per day** :smiley:
- There is a move towards systems that can tolerate the **loss of entire machines**, by using software fault-tolerance techniques in preference to hardware redundancy. Such systems also have operational advantages:
  - A single-server system requires planned downtime if you need to reboot the machine (to apply security patches, for example).
  - Whereas a system that can tolerate machine failure can be patched one node at a time, without downtime of the entire system.

### Software errors

- Harder to anticipate because they are correlated across nodes, they tend to cause many more system failures than uncorrelated hardware faults. Examples include:
  - A software **bug** that causes every instance of an application server to crash when given a particular bad
  input. For example, consider the leap second on _June 30, 2012_ that caused many applications to hang simultaneously, due to a **bug in the Linux kernel**.
  - A **runaway process** uses up some shared resource—CPU time, memory, disk space or network bandwidth.
  - A **service** that the system depends on slows down, becomes **unresponsive** or starts returning corrupted responses.
  - **Cascading failures**, where a small fault in one component triggers a fault in another component, which in turn triggers further faults.

## Human errors

- How do we make our system reliable, in spite of unreliable humans?:
  - Design systems in a way that **minimizes opportunities for error**. For example, **well-designed abstractions, APIs and admin interfaces** make it easy to do “the right thing”, and discourage “the wrong thing”.
  - Provide fully-featured **non-production** sandbox environments where people can explore and experiment safely, using real data, without affecting real users.
  - Test thoroughly at all levels, from **unit tests** to whole-system **integration tests** and **manual tests**.
  - Allow quick and easy recovery from human errors, to minimize the impact in the case of a failure. For example, make it **fast to roll back** configuration changes, **roll out new code gradually** (so that any unexpected bugs affect only a small subset of users), and provide tools to recompute data (in case it turns out that the old computation was incorrect).
  - Set up **detailed and clear monitoring and telemetry**, such as performance metrics and error rates.
- How important is reliability? Bugs in business applications cause **lost productivity** (and legal risks if figures are reported incorrectly), and outages of e-commerce sites can have **huge costs** in terms of lost revenue and reputation.

## Scalability

- Scalability is the term we use to describe a system’s ability to **adapt to increased load**.

### Describing load

- Load can be described with a few numbers which we call **load parameters**. Examples:
  - It’s requests per second;
  - Ratio of reads to write;
  - The number of simultaneously active users;
  - Number of messages in the queue ...

### Describing performance

- In a batch-processing system such as _Hadoop_, we usually care about **throughput**: the number of records we can process per second, or the total time it takes to run a job on a dataset of a certain size.
-  In online systems, what’s usually more important is the service’s **response time** — that is, the time between a client sending a request and
receiving a response.
- In practice, in a system handling a variety of requests, the latency per request can **vary a lot**. We therefore need to think of latency not as a **single number**, but as a **probability distribution**.
- Even in a scenario where you’d think all requests should take the same time, you get variation: random additional latency could be introduced by:
  - A **context switch** to a background process;
  - The **loss of a network packet** and TCP retransmission;
  - A **GC** pause, a **page fault** forcing a read from disk, or many other things.
- It’s common to see the **average** response time of a service reported (arithmetic mean).
- However, the mean is not a very good metric if you want to know your “typical” response time, because it is easily **biased by outliers**.
- Usually it is better to use **percentiles**. If you take your list of response times and sort it, from fastest to slowest, then the median is the half-way point: for example, if your median response time is 200 ms, that means half your requests return in less than 200 ms, and half your requests take longer than that.
- This makes the median a **good metric** if you want to know how long users typically have to wait. The median is also known as **50th percentile**, and sometimes abbreviated as **p50**.
    <p align="center"><img src="assets/response-time.png"></p>
- High percentiles of response times, also known as **tail latencies**, are important because they directly affect users’ experience of the service 🤷‍♀️.

### Approaches for coping with load

- good architectures usually involve a pragmatic mixture of approaches (vertical scaling and horizontal scaling).
- there is no such a thing as **magic scaling sauce**. The problem may be:
  - The volume of reads, the volume of writes, the volume of data to store;
  - The complexity of the data, the latency requirements, the access patterns;
  - Or (usually) some mixture of all of these plus many more issues.
- In an **early-stage startup** or an unproven product it’s usually more important to be able to **iterate quickly** on product features, than it is to scale to some **hypothetical future load**.

## Maintainability

- It is well-known that the majority of the cost of software is not in its initial development, but in its ongoing maintenance.
- We will pay particular attention to three design principles for software systems:
  - **Operability**: Make it **easy** for operations teams to keep the system **running smoothly**.
  - **Simplicity**: Make it **easy** for new engineers to **understand the system**, by removing as much complexity as possible from the system.
  - **Plasticity**: Make it **easy** for engineers in future to **make changes** to the system, adapting it for unanticipated use cases as requirements change. Also known as extensibility, modifiability or malleability.

### Operability: making life easy for operations

- _Good operations can often work around the limitations of bad (or incomplete) software, but good software cannot run reliably with bad operations._

## Simplicity: managing complexity

- *Moseley and Marks* define complexity as _accidental_ if it is not inherent in the problem that the software solves (as seen by the users), but arises only from the implementation.
- One of the best tools we have for removing accidental complexity is **abstraction**. A good abstraction can hide a great deal of implementation detail behind a clean, simple-to-understand facade.
- High-level programming languages are abstractions that hide machine code, CPU registers and syscalls.
- SQL is an abstraction that hides complex on-disk and in-memory data structures, concurrent requests from other clients, and inconsistencies after crashes.

## Plasticity: making change easy

- In terms of organizational processes, _agile_ working patterns provide a framework for adapting to change. The agile community has also developed technical tools and patterns that are helpful when developing software in a frequently-changing environment, such as **test-driven development (TDD)** and **refactoring**.
