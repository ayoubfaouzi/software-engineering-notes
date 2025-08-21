# Chapter 3. Encoding and Evolution

- When a data format or schema changes, a corresponding change to application code often needs to happen. However, in a large application, code changes often cannot happen instantaneously:
  - With **server-side** apps you may want to perform a **rolling upgrade** (also known as a **staged rollout**).
  - With **client-side** apps you’re at the mercy of the user, who may not install the update for some time 🤷.
- This means that old and new versions of the code, and old and new data formats, may potentially all coexist in the system at the same time. In order for the system to continue running smoothly, we need to maintain compatibility in both directions:
    - **Backward compatibility**: Newer code can read data that was written by older code.
    - **Forward compatibility**: Older code can read data that was written by newer code.
- ⚠️ Forward compatibility can be trickier, because it requires **older code** to **ignore additions** made by a newer version of the code.

## Formats for Encoding Data

### Language-Specific Formats

- Many programming languages come with built-in support for encoding in-memory objects into byte sequences. For example:
  -  Java has `java.io.Serializable`
  -  Ruby has `Marshal`
  -  Python has `pickle` and so on.
- Many third-party libraries also exist, such as `Kryo` for Java.
- These encoding libraries are very **convenient**. However, they also have a number of deep problems:
  - Encoding is often tied to a particular programming language, and reading the data in another language is very difficult.
  - Decoding process needs to be able to **instantiate arbitrary classes**. This is frequently a source of **security** problems.
  - **Versioning** data is often an **afterthought** in these libraries.
  - **Efficiency** is also often **afterthought**.

### JSON, XML, and Binary Variants

- Moving to standardized encodings that can be written and read by many programming languages, *JSON* and *XML* are the obvious contenders. They are **widely known** & **supported**. *XML* is often criticized for being too verbose and unnecessarily complicated 🤷‍♂️. *JSON*’s popularity is mainly due to its built-in support in **web browsers**.
Besides the superficial syntactic
issues, they also have some subtle problems:
- There is a lot of ambiguity around the encoding of numbers:
  - In XML and CSV, you **cannot distinguish** between a **number** and a **string** that happens to consist of digits (except by referring to an external schema 🫤). *JSON* distinguishes strings and numbers, but it **doesn’t distinguish** integers and **floating-point numbers**, and it doesn’t specify a precision.
  - This is a problem when dealing with **large numbers**; for example, integers greater than `2^53` cannot be exactly represented in an *IEEE 754* double-precision floating-point number, so such numbers become inaccurate when parsed in a language that uses floating-point numbers (such as *JavaScript*).
  - JSON and XML have good support for Unicode character strings (i.e., human-readable text), but they don’t support binary strings (sequences of bytes without a character encoding).
  - There is **optional schema support** for both XML and JSON. These schema languages are quite powerful, and thus quite complicated to learn and implement.
  - CSV does not have any schema, so it is up to the application to define the meaning of each row and column. If an application change adds a new row or column, you have to handle that change manually. CSV is also a quite vague format (what happens if a value contains a comma or a *newline* character?). Although its **escaping rules** have been **formally specified**, not all parsers implement them correctly.

### Binary encoding

- JSON is less verbose than XML, but both still use a **lot of space** compared to binary formats 💯.
  - Since they don’t prescribe a schema, they need to include **all the object field names** within the encoded data.
- ➡️ This observation led to the development of a profusion of binary encodings for JSON (`MessagePack`, `BSON`, `BJSON`, `UBJSON`, `BISON`, and `Smile`, ...) and for XML (`WBXML` and `Fast Infoset`, for example).

## Thrift and Protocol Buffers

- Apache Thrift and (Google) Protocol Buffers (protobuf) are binary encoding libraries that are based on the same principle.
- Both Thrift and Protocol Buffers **require a schema** for any data that is encoded.
- Thrift and Protocol Buffers each come with a **code generation** tool that takes a **schema definition** like and produces **classes** that implement the schema in various programming languages.
- For Thrift **BinaryProtocol**, the big difference is that there are no **field names**. Instead, the encoded data contains **field tags**, which are numbers (1, 2, and 3). Those are the numbers that appear in the **schema definition**.
- Field tags are like aliases for fields—they are a compact way of saying what field we’re talking about, without having to spell out the field name.
- Each field has a **type annotation** and, where required, a length indication (length of a string, number of items in a list). <p align="center"><img src="assets/thrift-binary-protocol.png" width="400px" height="auto"></p>
- The Thrift **CompactProtocol** encoding is semantically equivalent to **BinaryProtocol** it packs the same information into only **34 bytes**. It does this by packing the **field type** and **tag number** into a **single byte**, and by using **variable-length integers**. Rather than using a full eight bytes for the number `1337`, it is encoded in two bytes, with the **top bit** of each byte used to indicate whether there are still more bytes to come. This means numbers between `64` and `63` are encoded in
1 byte, numbers between `8192` and `8191` are encoded in 2 bytes, etc. Bigger numbers use more bytes. <p align="center"><img src="assets/thrift-compact-protocol.png" width="400px" height="auto"></p>
- Finally, **Protocol Buffers** does the bit packing slightly differently, but is otherwise very similar to Thrift’s CompactProtocol. Protocol Buffers fits the same record in 33 bytes. <p align="center"><img src="assets/protocol-buffers.png" width="400px" height="auto"></p>

## Field tags and schema evolution