# Chapter 10. Batch Processing

- There are 3️⃣ major styles of data processing — **services**, **batch processing**, and **stream processing**/
- Most modern systems follow a **request/response model**, where a client sends a request and quickly receives a reply. This is typical of online services like web servers, APIs, and databases, where **response time** and **availability** are **critical** because users are actively waiting for results.
- However, not all systems work this way. Two other models are common:
  - **Batch Processing Systems** (Offline Systems):
    - These handle large volumes of data in scheduled jobs that may take **minutes to days**.
    - There’s no immediate user interaction; instead, the focus is on **throughput** — how efficiently data can be processed.
    - Technologies like `MapReduce` (and its open-source implementations such as **Hadoop**) exemplify this model. Though MapReduce is now less dominant, it remains important for understanding scalable data processing.
  - **Stream Processing Systems** (Near-Real-Time Systems):
    - These process continuous streams of data as **events** occur, achieving lower latency than batch jobs.
    - Stream processing blends elements of online and batch systems and is discussed as an evolution of batch processing.

## Batch Processing with Unix Tools

