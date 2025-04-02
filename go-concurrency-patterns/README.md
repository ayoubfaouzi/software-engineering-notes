# Go Concurrency Patterns

- Go is the latest on the _Newsqueak-Alef-Limbo_ branch, distinguished by **first-class** channels.
- The Go approach _Don't communicate by sharing memory, share memory by communicating._

## Generator

- Generator: function that **returns a channel**.
- Channels are **first-class values**, just like strings or integers.
- Channels as a handle on a service.
  > Our boring function returns a channel that lets us communicate with the boring service it provides.

## Multiplexing

- The generator pattern makes _Joe_ and _Ann_ count in lockstep.
- We can instead use a fan-in function to let whosoever is ready talk.
- We stitch the two channel into a **single one**, and the fan-in function forwards the messages to the output channel.
