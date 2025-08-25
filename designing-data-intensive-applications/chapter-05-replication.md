# Chapter 5. Replication

- Replication means keeping a **copy** of the same data on **multiple machines** that are connected via a network. 
  - To keep data **geographically close** to your users (and thus reduce latency)
  - Increase **availability**.
  - To scale out the number of machines that can serve read queries (and thus increase **read throughput**).
- If the data that you’re replicating **does not change** over time, then replication is easy 🤹‍♂️:
  - You just need to copy the data to every node once, and you’re done. 
  - All of the **difficulty** in replication lies in handling **changes** to replicated data, and that’s what this chapter is about.

## Leaders and Followers

- Each node that stores a copy of the database is called a **replica**.
- How do we ensure that all the data ends up on all the replicas ❓
- The most common solution for this is called **leader-based replication** (also known as **active/passive** or **master–slave** replication). It works as follows:
  1. One of the replicas is designated the **leader** (also known as **master** or **primary**). When clients want to write to the database, they must send their requests to the leader, which first writes the new data to its local storage.
  2. The other replicas are known as **followers** (read **replicas**, **slaves**, **secondaries**, or **hot standbys**). Whenever the leader writes new data to its local storage, it also sends the data change to all of its followers as part of a **replication log** or **change stream**. Each follower takes the log from the leader and updates its local copy of the data‐base accordingly, by applying all writes in the same order as they were processed on the leader.
  3. When a client wants to read from the database, it can query either the leader or any of the followers. However, writes are only accepted on the leader (the followers are read-only from the client’s point of view).

### Synchronous Versus Asynchronous Replication

- In the example of Figure 5-2, the replication to *follower 1* is **synchronous**: the leader waits until follower 1 has **confirmed** that it received the **write** before reporting success to the user, and before making the write visible to other clients. The replication to *follower 2* is **asynchronous**: the leader sends the message, but doesn’t wait for a response from the follower.
<p align="center"><img src="assets/replication-sync-async.png" width="500px" height="auto"></p>

- Normally, replication is **quite fast**: most database systems apply changes to followers in **less than a second**.
- However, there is no guarantee of how long it might take 🤓. There are circumstances when followers might fall behind the leader by several minutes or more, for example:
  - If a follower is recovering from a **failure**,
  - If the system is operating near **maximum capacity**,
  - or if there are **network problems** between the nodes.