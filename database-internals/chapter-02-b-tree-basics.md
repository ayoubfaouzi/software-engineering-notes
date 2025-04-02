# Chapter 2: B-Tree Basics

## Binary Search Tree

- Insertion might lead to the situation where the tree is **unbalanced**. The worst-case scenario is where we end up with a **pathological** tree, which looks more like a linked list, and instead of desired logarithmic complexity `O(log2 N)`, we get linear `O(N)`.
- One of the ways to keep the tree balanced is to perform a **rotation** step after nodes are added or removed.
  - If the insert operation leaves a branch unbalanced, we can rotate nodes around the middle one.
  - In the example below, during rotation the middle `node (3)`, known as a **rotation pivot**, is promoted one level higher, and its parent becomes its right child.
<p align="center"><img src="assets/tree-balancing.png" width="300px"></p>