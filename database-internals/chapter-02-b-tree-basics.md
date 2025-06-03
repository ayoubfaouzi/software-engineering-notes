# Chapter 2: B-Tree Basics

## Binary Search Tree

- Insertion might lead to the situation where the tree is **unbalanced**. The worst-case scenario is where we end up with a **pathological** tree, which looks more like a linked list, and instead of desired logarithmic complexity `O(log2 N)`, we get linear `O(N)`.
- One of the ways to keep the tree **balanced** is to perform a **rotation** step after nodes are added or removed.
  - If the insert operation leaves a branch unbalanced, we can rotate nodes around the middle one.
  - In the example below, during rotation the middle `node (3)`, known as a **rotation pivot**, is promoted one level higher, and its parent becomes its right child.
<p align="center"><img src="assets/tree-balancing.png" width="300px"></p>.

- At the same time, due to **low fanout** (fanout is the maximum allowed number of children per node), we have to perform balancing, relocate nodes, and update pointers rather frequently. Increased maintenance costs make BSTs impractical as on-disk data structures.
- If we wanted to maintain a BST on disk, we’d face several problems:
  - locality: node child pointers may span across several disk pages, since elements are added in random order (we can improve the situation by modifying the tree layout and using **paged binary trees**).
  - tree height: since BT has a fanout of just two, height is a binary logarithm of the number of the elements in the tree, and we have to perform O(log2 N) **seeks** to locate the searched element and, subsequently, perform the same number of **disk transfers**.
- Considering these factors, a version of the tree that would be better suited for **disk implementation** has to exhibit the following properties:
  - High fanout to improve locality of the neighboring keys.
  - Low height to reduce the number of seeks during traversal.

## Disk-Based Structures

On-disk data structures are often used when the amounts of data are so large that keeping an entire dataset in memory is impossible or not feasible. Only a fraction of the data can be cached in memory at any time, and the rest has to be stored on disk in a manner that allows efficiently accessing it.

### Hard Disk Drives

On spinning disks, seeks increase costs of random reads because they require disk rotation and mechanical head movements to position the read/write head to the desired location. However, once the expensive part is done, reading or writing contiguous bytes (i.e., sequential operations) is relatively cheap.

Head positioning is the most expensive part of an operation on the HDD. This is one of the reasons we often hear about the positive effects of sequential I/O: reading and writing contiguous memory segments from disk.

### Solid State Drives

Since in both device types (HDDs and SSDs) we are addressing chunks of memory rather than individual bytes (i.e., accessing data block-wise), most operating systems have a block device abstraction. It hides an internal disk structure and buffers I/O operations internally, so when we’re reading a single word from a block device, the whole block containing it is read. This is a constraint we cannot ignore and should always take into account when working with disk-resident data structures.

In SSDs, we don’t have a strong emphasis on random versus sequential I/O, as in HDDs, because the difference in latencies between random and sequential reads is not as large. There is still some difference caused by prefetching, reading contiguous
pages, and internal parallelism.

Even though garbage collection is usually a background operation, its effects may negatively impact write performance, especially in cases of random and unaligned write workloads. Writing only full blocks, and combining subsequent writes to the same block, can help to reduce the number of required I/O operations.

### On-Disk Structures

Besides the cost of disk access itself, the main limitation and design condition for building efficient on-disk structures is the fact that the smallest unit of disk operation is a block. To follow a pointer to the specific location within the block, we have to fetch an entire block. Since we already have to do that, we can change the layout of the data structure to take advantage of it.

In summary, on-disk structures are designed with their target storage specifics in mind and generally optimize for fewer disk accesses. We can do this by improving locality, optimizing the internal representation of the structure, and reducing the
number of out-of-page pointers.

We came before to the conclusion that high fanout and low height are desired properties for an optimal on-disk data structure. We’ve also just discussed additional space overhead coming from pointers, and maintenance over‐
head from remapping these pointers as a result of balancing. B-Trees combine these ideas: increase node fanout, and reduce tree height, the number of node pointers, and the frequency of balancing operations.

## Ubiquitous B-Trees

B-Trees build upon the foundation of balanced search trees and are different in that they have higher fanout and smaller height.