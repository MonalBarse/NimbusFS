# NimbusFS

NimbusFS is a robust, distributed, peer-to-peer Content Addressable Storage (CAS) system implemented in Go. It provides a scalable and efficient solution for storing and retrieving data across a network of nodes.

## Conceptual Overview

NimbusFS is a distributed Content-Addressable Storage system where data is replicated across multiple nodes without requiring centralized coordination. Each file is identified by its cryptographic hash (SHA-1), ensuring content integrity and deduplication. The system distributes encrypted copies of files across participating nodes, enabling resilient data retrieval even when individual nodes become unavailable.

**The Problem Addressed:** Traditional storage systems rely on centralized servers that create single points of failure, potential for censorship, and centralized control. NimbusFS provides a peer-to-peer alternative where data resilience is achieved through distribution across network participants rather than dependence on any single authority.

## Potential Use Cases

- **Securely sharing research data** across a team of academics without relying on a single university server that could go offline or restrict access
- **A decentralized photo backup service** where your encrypted photos are stored across your own devices (and trusted friends' devices) instead of being locked into one company's cloud
- **Building censorship-resistant applications** where content can't be taken down by shutting off one server — the network keeps it alive
- **Collaborative software development** where code and assets are distributed across team members' machines, ensuring no single point of failure for critical project files
- **Medical record sharing** between hospitals and clinics without depending on a central database that could be compromised or become unavailable

## Table of Contents

- [Conceptual Overview](#conceptual-overview)
- [Potential Use Cases](#potential-use-cases)
- [How It Works](#how-it-works)
- [System Guarantees & Implementation](#system-guarantees--implementation)
- [Architecture](#architecture)
- [Resources](#resources)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)


## How It Works

NimbusFS operates through a distributed protocol handling file storage and retrieval across network nodes:

### Storing a File
When storing a file, the local NimbusFS node generates a SHA-1 hash fingerprint (`store.go`), encrypts the content using AES encryption (`crypto.go`), and distributes encrypted copies to peer nodes via TCP connections (`p2p/tcp_transport.go`). The node announces the file's availability by broadcasting its hash identifier to the network.

### Retrieving a File
File retrieval begins with a hash-based lookup on the local node (`store.go`). If the file is not available locally, the node broadcasts a retrieval request to the network (`fileserver.go`). Peer nodes respond with encrypted file data, which is automatically decrypted using the appropriate encryption key.

## System Guarantees & Implementation

### Data Integrity & Deduplication
**Guarantee:** Files maintain integrity and identical content is never duplicated across the system.

**Implementation:** Content-addressing via SHA-1 hashing (`store.go`) ensures each file receives a unique fingerprint based on its content. Any modification changes the hash completely, enabling immediate detection of alterations. Identical files produce identical hashes, preventing storage duplication.

### Fault Tolerance
**Guarantee:** Data remains accessible despite multiple node failures.

**Implementation:** The peer-to-peer architecture (`p2p/tcp_transport.go`) distributes files across multiple nodes without central coordination. Each node operates independently for storage, retrieval, and serving operations, eliminating single points of failure.

### Privacy & Security
**Guarantee:** File contents remain confidential during network transmission.

**Implementation:** AES encryption (`crypto.go`) secures data before network transmission. Each node maintains independent encryption keys, ensuring only authorized participants can decrypt and access file contents.

### Efficient Network Usage
**Guarantee:** File sharing minimizes network overhead through optimized data structures.

**Implementation:** GOB encoding provides efficient data serialization, while hierarchical path transformation functions organize files for fast lookups and minimal network requests.

## Architecture

NimbusFS is built on several key components:

1. **FileServer**: The core component managing file operations and peer interactions.
2. **Store**: Handles local file storage with a customizable path transform function.
3. **TCPTransport**: Manages network communications between peers.

## Resources

### Educational Videos

These educational videos provided foundational knowledge and inspiration for building NimbusFS:

- **Distributed Computing - DFS** by [Perfect CE](https://www.youtube.com/@perfectcomputerengineer): [Watch on YouTube](https://youtu.be/Xhi3hqbiXNM)
- **Big Data Analysis - DFS** by [Perfect CE](https://www.youtube.com/@perfectcomputerengineer): [Watch on YouTube](https://youtu.be/McTWc6N-pBg)
- **Architecture of a DFS** by [Knowledge Hub](https://www.youtube.com/@knowledgehub9741): [Watch on YouTube](https://youtu.be/QmNlluPbEEk)
- **Hadoop - HDFS** by [SimpliLearn](https://www.youtube.com/@SimplilearnOfficial): [Watch on YouTube](https://youtu.be/6apXsm_25s0)

- **[Talk on p2p implementation](https://youtu.be/waVtYYSXkXU?si=UFP2YRSx0dxZ1fRc)**
- **[Additional](https://youtu.be/eRndYq8iTio?si=5XuYlcs6FgDIbkxC)**

---

## Getting Started

For developers who want to run a node in the distributed network, here's how to get started:

### Prerequisites

- Go 1.22 or higher
- Make (for building)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MonalBarse/NimbusFS
   cd NimbusFS
   ```

2. Build the application:
   ```bash
   go mod tidy
   make build 
   ```
## Usage

NimbusFS comes with a Makefile for easy building, running, and testing.

### Running the Demo

To see NimbusFS in action, run:
```bash
make run
```
This command will build and execute the application, demonstrating the distributed storage system!

### Running Tests

To run the tests for all packages, use:
```bash
make test
```

### Understanding the Default Demo

The `main.go` file sets up a network of three file servers:

1. **Server 1** listens on port 3000
2. **Server 2** listens on port 7000  
3. **Server 3** listens on port 5000 and connects to the other two servers

The demo demonstrates distributed storage capabilities by:

1. **Storing 20 files** (named "picture_0.png" to "picture_19.png") on Server 3
2. **Deleting Server 3's local copies** (simulating local storage failure)
3. **Successfully retrieving all files from the network** — the data remains available from Servers 1 and 2
4. **Printing the contents** of each retrieved file

This demonstrates the distributed nature of the system: when one node loses its local data, the network maintains availability through other participating nodes.

### Customizing Your Own Network

To create a custom network of NimbusFS nodes, modify the `main.go` file:

**Creating Servers:**
```go
server := makeServer(listenAddr, bootstrapNodes...)
```
Where:
- `listenAddr` is the address for the server to listen on (e.g., ":3000")
- `bootstrapNodes` are addresses of existing servers to connect to (can be empty for the first server)

**Starting Servers:**
```go
go server.Start()
```

**File Operations:**

- To store a file:
  ```go
  data := bytes.NewReader([]byte("file contents"))
  server.Store("filename.txt", data)
  ```
  
- To retrieve a file:
  ```go
  r, err := server.Get("filename.txt")
  if err != nil {
    log.Fatal(err)
  }
  contents, err := io.ReadAll(r)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(string(contents))
  ```

## Architecture

NimbusFS is built on several key components that handle different aspects of the distributed storage system:

1. **FileServer** (`fileserver.go`): The central coordination component that manages file operations, maintains peer connections, and handles incoming storage and retrieval requests. It orchestrates communication between the storage layer and network layer, processing messages via GOB decoding and routing them to appropriate handlers.

2. **Store** (`store.go`): The local storage management component that implements content-addressable storage through configurable path transformation functions. It handles file persistence, retrieval, and integrity checking using SHA-1 hashing for content identification and hierarchical directory organization for efficient lookups.

3. **TCPTransport** (`p2p/tcp_transport.go`): The network communication component that manages TCP connections between peers. It handles connection establishment, peer discovery, message routing, and provides abstractions for sending data streams between nodes in the distributed network.

4. **Crypto** (`crypto.go`): The security component that provides AES encryption and decryption capabilities for data in transit. It generates encryption keys, handles secure data transformation during network operations, and ensures confidentiality of stored and transmitted content.

## Resources

### Educational Videos

These educational videos provided foundational knowledge and inspiration for building NimbusFS:

- **Distributed Computing - DFS** by [Perfect CE](https://www.youtube.com/@perfectcomputerengineer): [Watch on YouTube](https://youtu.be/Xhi3hqbiXNM)
- **Big Data Analysis - DFS** by [Perfect CE](https://www.youtube.com/@perfectcomputerengineer): [Watch on YouTube](https://youtu.be/McTWc6N-pBg)
- **Architecture of a DFS** by [Knowledge Hub](https://www.youtube.com/@knowledgehub9741): [Watch on YouTube](https://youtu.be/QmNlluPbEEk)
- **Hadoop - HDFS** by [SimpliLearn](https://www.youtube.com/@SimplilearnOfficial): [Watch on YouTube](https://youtu.be/6apXsm_25s0)
- **[Talk on p2p implementation](https://youtu.be/waVtYYSXkXU?si=UFP2YRSx0dxZ1fRc)**
- **[Additional](https://youtu.be/eRndYq8iTio?si=5XuYlcs6FgDIbkxC)**

## Contributing

Contributions are welcome! Please read the [contribution guidelines](CONTRIBUTING.md) before submitting a pull request.

## License

This project is licensed under the Apache-2.0 license - see the [LICENSE](LICENSE) file for details.
