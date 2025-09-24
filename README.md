# NimbusFS

NimbusFS is a robust, distributed, peer-to-peer Content Addressable Storage (CAS) system implemented in Go. It provides a scalable and efficient solution for storing and retrieving data across a network of nodes.

## Conceptual Overview

Imagine NimbusFS as a magical, self-organizing library. Instead of one central building, every member holds a piece of the collection. When you add a book, the library gives it a unique fingerprint (a hash). Copies of the book are then cleverly and securely distributed among all the members. If you want to read that book later, you just ask the library for it by its fingerprint, and the nearest members who have a copy deliver it to you instantly.

**The Problem We Solve:** Traditional storage relies on centralized servers — single points of failure that can go down, be censored, or be controlled by one entity. NimbusFS offers a peer-to-peer alternative where your data is resilient, distributed, and truly belongs to the network of participants rather than any single authority.

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
- [Our Guarantees & How We Keep Them](#our-guarantees--how-we-keep-them)
- [Architecture](#architecture)
- [Resources](#resources)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)


## How It Works

Continuing our magical library analogy, here's what happens when you interact with NimbusFS:

### Storing a File (Adding a Book)
When you want to add a new file, you bring it to the "librarian" (your NimbusFS node). The librarian creates a unique fingerprint for it using SHA-1 hashing (`store.go`), locks it in a secure, encrypted box using AES encryption (`crypto.go`), and then sends copies of this box to other members of the library network via TCP connections (`p2p/tcp_transport.go`). Your librarian also announces to everyone, "I have a new book with this fingerprint!"

### Retrieving a File (Requesting a Book)
When you need a file, you ask the network, "Does anyone have the book with this fingerprint?" Your local librarian first checks its own shelves (`store.go`). If it's not there, it broadcasts the request to the network (`fileserver.go`). The first librarian to respond sends you their encrypted copy, which your librarian automatically decrypts for you using the shared encryption key.

## Our Guarantees & How We Keep Them

### Data Integrity & No Duplicates
**Our Guarantee:** Your files will never be secretly corrupted, and we never waste space storing the same file twice.

**How We Keep It:** We use a technique called "Content-Addressing" (`store.go`), where every file gets a unique fingerprint (SHA-1 hash). If even a single bit changes, the fingerprint changes completely, so we immediately know the file is different. This also means identical files always get the same fingerprint, preventing duplicates.

### No Central Point of Failure
**Our Guarantee:** Your data is safe even if multiple computers go offline.

**How We Keep It:** This is a true Peer-to-Peer system (`p2p/tcp_transport.go`). Files are distributed across many nodes in the network. There is no "main" server to attack, control, or fail. Each node can store, retrieve, and serve files independently.

### Privacy & Security
**Our Guarantee:** No one can snoop on your files as they travel across the network.

**How We Keep It:** We use strong AES encryption (`crypto.go`) to lock your data before it ever leaves your machine. Each node has its own encryption key, ensuring that only authorized participants can decrypt and access the files.

### Efficient Network Usage
**Our Guarantee:** Files are shared efficiently without unnecessary network overhead.

**How We Keep It:** We use custom GOB encoding for efficient data serialization and smart path transformation functions that organize files in a hierarchical structure, making lookups fast and network requests minimal.

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

For developers who want to run a node in our distributed network, here's how you can get started:

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
This command will build and execute the application, demonstrating the magical library in action!

### Running Tests

To run the tests for all packages, use:
```bash
make test
```

### Understanding the Default Demo

The `main.go` file sets up a network of three file servers (like three interconnected library branches):

1. **Library Branch 1** listens on port 3000
2. **Library Branch 2** listens on port 7000  
3. **Library Branch 3** listens on port 5000 and connects to the other two branches

The demo then shows the magic by:

1. **Adding 20 "books"** (files named "picture_0.png" to "picture_19.png") to Branch 3
2. **Destroying Branch 3's local copy** (simulating a local disk failure)
3. **Successfully retrieving all books from the network** — they're still available from Branches 1 and 2!
4. **Printing the contents** of each retrieved file

This perfectly demonstrates the distributed nature of our magical library: even when one branch loses its books, the network remembers and can deliver them from other branches.

### Customizing Your Own Network

To create your own network of NimbusFS nodes, you can modify the `main.go` file:

**Creating Library Branches:**
```go
server := makeServer(listenAddr, bootstrapNodes...)
```
Where:
- `listenAddr` is the address your branch will listen on (e.g., ":3000")
- `bootstrapNodes` are addresses of existing branches to connect to (can be empty for the first branch)

**Starting Your Branches:**
```go
go server.Start()
```

**Adding and Retrieving Books:**

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

NimbusFS is built on several key components that work together like the staff of our magical library:

1. **FileServer** (`fileserver.go`): The main "librarian" that manages file operations, coordinates with other library branches (peers), and handles requests from visitors.

2. **Store** (`store.go`): The "filing system" that handles local book storage with smart organization. It uses the Content-Addressable System to give each file a unique location based on its fingerprint.

3. **TCPTransport** (`p2p/tcp_transport.go`): The "communication network" that manages conversations between different library branches, handling connections and message passing.

4. **Crypto** (`crypto.go`): The "security system" that encrypts books before sending them to other branches and decrypts them when they're needed locally.

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
