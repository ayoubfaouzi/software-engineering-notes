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

- We can replicate the Unix log analysis using a simple Ruby script.

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

- While it’s less concise than the Unix pipeline, the Ruby version is **easier to read**. However, the key point is that — beyond syntax — the **execution flow differs** significantly, especially when processing large files, which impacts **performance** and **efficiency**.

##### Sorting versus in-memory aggregation

- The **Ruby script** uses a **hash table** that keeps a counter for each unique URL. Its **memory** usage depends on the number of **distinct** URLs, not total log lines. This works well when all unique URLs fit comfortably in memory (e.g., within 1 GB).
- The **Unix pipeline** avoids keeping everything in memory by **sorting** the repeated URLs. Sorting can efficiently handle large datasets using **disk-based mergesort**, which relies on sequential I/O — ideal for disks.
- GNU sort automatically handles datasets larger than memory by spilling to disk and **parallelizing** sorting across **CPU cores**, allowing it to scale smoothly.
- 👉 Hash tables are faster for smaller datasets, while sorting pipelines scale better for very large datasets that exceed available memory.

### The Unix Philosophy
