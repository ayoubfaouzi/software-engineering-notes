# Go Concurrency Patterns

- Go is the latest on the _Newsqueak-Alef-Limbo_ branch, distinguished by **first-class** channels.
- The Go approach _Don't communicate by sharing memory, share memory by communicating._

## Generator

- Generator: function that **returns a channel**.
- Channels are **first-class values**, just like strings or integers.
- Channels as a handle on a service.
  > Our boring function returns a channel that lets us communicate with the boring service it provides.

## Fan-In (Multiplexing)

- The generator pattern makes _Joe_ and _Ann_ count in lockstep.
- We can instead use a fan-in function to let whosoever is ready talk.
- We stitch the two channel into a **single one**, and the fan-in function forwards the messages to the output channel.
- Fan In is used when a **single function** reads from **multiple inputs** and proceeds until all are closed. This is made possible by multiplexing the input into a single channel.
- What: Combine results from multiple goroutines into a single channel.
- Why: Aggregate results or wait for all goroutines to complete.

## Fan-Out

- What: Distribute work across multiple goroutines to run in parallel.
- Why: Improve throughput by utilizing multiple CPU cores.
- Example: Multiple workers pulling tasks from the same job queue.

## Daisy Chain

- Goroutines and channels are chained together to pass data along a series of steps.
- It’s often used to illustrate the power and simplicity of Go’s concurrency model and how cheap is a goroutine compared to a thread!
