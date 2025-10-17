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

- A sample nginx access log line contains multiple fields, including client IP, user, timestamp, request, status, bytes sent, referrer, and user agent.
- For instance, the line shows that on `Feb 27, 2015, at 17:55:11 UTC`, a client at `216.58.210.78` requested `/css/typography.css`.
- The user wasn’t authenticated, the request succeeded with status `200`, the response size was `3,377 bytes`, the referrer was `http://martin.kleppmann.com/`.

### Simple Log Analysis

- By chaining commands like `cat`, `awk`, `sort`, `uniq`, and `head`, you can quickly find the five most popular pages on a website.
- The pipeline extracts the requested URLs, counts how often each appears, sorts them by frequency, and displays the top results.
  ```sh
  cat /var/log/nginx/access.log |
      awk '{print $7}' |
      sort |
      uniq -c |
      sort -r -n |
      head -n 5
  ```
- Despite looking complex, this method is **fast**, **flexible**, and **powerful**, capable of processing **gigabytes** of logs efficiently. Small tweaks to the `awk` command can easily change what’s analyzed (e.g., ignoring CSS files or counting client IPs).
- 👉 Unix text-processing tools like `awk`, `sed`, `grep`, `sort`, `uniq`, and `xargs` are valuable for quick, effective data analysis.

#### Chain of commands versus custom program

We can replicate the Unix log analysis using a simple Ruby script:
  ```ruby
  counts = Hash.new(0)
  File.open('/var/log/nginx/access.log') do |file|
    file.each do |line|
      url = line.split[6]
      counts[url] += 1
    end
  end

  top5 = counts.map{|url, count| [count, url] }.sort.reverse[0...5]
  top5.each{|count, url| puts "#{count} #{url}" }
  ```
While it’s less concise than the Unix pipeline, the Ruby version is **easier to read**. However, the key point is that — beyond syntax — the **execution flow differs** significantly, especially when processing large files, which impacts **performance** and **efficiency**.

##### Sorting versus in-memory aggregation

- The **Ruby script** uses a **hash table** that keeps a counter for each unique URL. Its **memory** usage depends on the number of **distinct** URLs, not total log lines. This works well when all unique URLs fit comfortably in memory (e.g., within 1 GB).
- The **Unix pipeline** avoids keeping everything in memory by **sorting** the repeated URLs. Sorting can efficiently handle large datasets using **disk-based mergesort**, which relies on sequential I/O — ideal for disks.
- GNU sort automatically handles datasets larger than memory by spilling to disk and **parallelizing** sorting across **CPU cores**, allowing it to scale smoothly.
- 👉 Hash tables are faster for smaller datasets, while sorting pipelines scale better for very large datasets that exceed available memory.

### The Unix Philosophy

- *Doug McIlroy* introduced the concept of **Unix pipes** in `1964`, comparing them to “garden hoses” that can be connected to process data in flexible ways. This idea evolved into the Unix philosophy, summarized as:
  - Make each program do **one thing well**.
  - Ensure each **program’s output** can serve as another **program’s input**.
  - Build **quickly**, **iterate**, and **rebuild** when necessary.
  - Use tools to **automate** and **simplify** tasks.
- This philosophy emphasizes automation, modularity, and experimentation, ideas that strongly resemble today’s **Agile** and **DevOps** principles.
- 👉 The Unix shell enables this **composability**, allowing independently built tools to interoperate seamlessly and form complex, efficient data processing pipelines.

#### A uniform interface

- Unix programs interoperate so well because they all share a simple, uniform interface — the **file descriptor**, which represents an ordered sequence of bytes. This abstraction allows diverse things (**files**, **devices**, **sockets**, **pipes**) to communicate seamlessly.
- Most Unix tools conventionally treat this byte stream as ASCII text, typically organized as lines separated by `\n` and fields split by whitespace. This shared convention enables tools like `awk`, `sort`, `uniq`, and `head` to work together easily, even if the text-based interface **isn’t elegant** or **strongly** structured.
-👉 While this simplicity sacrifices **readability** and **rich data semantics**, it provides extraordinary **composability**: Unix programs can be chained together flexibly — something rare in modern software ecosystems, where systems are often fragmented and data exchange between them is difficult.

#### Separation of logic and wiring

- Standard input (*stdin*) and standard output (*stdout*) are key to Unix’s flexibility and composability.
- By default, **stdin** comes from the **keyboard** and **stdout** goes to the **screen**, but users can redirect them to or from files — or connect programs together using **pipes**, which stream data directly between **processes** without writing to **disk**.
- Programs that simply read from stdin and write to stdout don’t need to know where their data comes from or goes, enabling **loose coupling** and easy composition in pipelines. Users can integrate their own tools just as easily as system ones, since everything shares the same interface.
- However, this model has limitations: it’s less suited for programs needing **multiple inputs/outputs** or **complex I/O** like network connections or subprocesses, which must be handled within the program itself — reducing the shell’s flexibility 🤷.

#### Transparency and experimentation

- Unix tools are powerful because they make it easy to observe and **experiment** with data processing:
  - Input files are **immutable**, so you can safely rerun commands without altering data.
  - You can inspect **intermediate** output at any stage (e.g., using `less`) for debugging.
  - You can save intermediate results to files, allowing **partial restarts** without reprocessing everything.
- This simplicity and transparency make Unix tools ideal for experimentation, despite lacking the sophistication of database query optimizers.
- Their main limitation, however, is that they operate on a **single machine**, which is where **distributed** systems like `Hadoop` extend the model.

## MapReduce and Distributed Filesystems
