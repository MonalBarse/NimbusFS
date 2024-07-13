# NimbusFS

NimbusFS is a robust, distributed, peer-to-peer Content Addressable Storage (CAS) system implemented in Go. It provides a scalable and efficient solution for storing and retrieving data across a network of nodes.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Resouces](#resources)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)


## Features

- **Content Addressable Storage (CAS)**: Ensures data integrity and deduplication.
- **Peer-to-Peer Architecture**: Enables decentralized data storage and retrieval.
- **Distributed System**: Scales horizontally for increased storage capacity and fault tolerance.
- **Encryption**: Provides data security with AES encryption.
- **Flexible Transport Layer**: Supports TCP with potential for easy expansion to other protocols.
- **Custom Encoding**: Uses GOB encoding for efficient data serialization.

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

### Prerequisites

- Go 1.15 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MonalBarse/NimbusFS
   cd NimbusFS
   ```
2. Build
   ```bash
   go mod tidy
   make build 
   ```
## Usage

NimbusFS comes with a Makefile for easy building, running, and testing.

### Building the Application
To build the application, run:
```bash
make run
```
This command will build the application and then execute it.

### Running Tests

To run the tests for all packages, use:
```bash
make test
```
### Understanding the Default Behavior

The `main.go` file sets up a network of three file servers:

1. Server 1 listens on port 3000
2. Server 2 listens on port 7000
3. Server 3 listens on port 5000 and connects to the other two servers

The program then demonstrates the functionality by:

1. Storing 20 files (named "picture_0.png" to "picture_19.png") on Server 3
2. Immediately deleting these files from Server 3's local storage
3. Retrieving these files from the network
4. Printing the contents of each retrieved file

This showcases the distributed nature of NimbusFS, where files can be retrieved from the network even if they're not present in the local storage of the requesting node.

### Customizing the Network

- To create your own network of NimbusFS nodes, you can modify the `main.go` file. Use the `makeServer` function to create new server instances:

  ```go
  server := makeServer(listenAddr, bootstrapNodes...)
  ```
  Where:

  - listenAddr is the address the server will listen on (e.g., ":3000")
  - bootstrapNodes are the addresses of existing nodes to connect to (can be empty for the first node)

- Then start each server with:
  ```go
  go server.Start()
  ```

- Storing and Retrieving Files

  - To store a file:
    ```go
    goCopydata := bytes.NewReader([]byte("file contents"))
    server.Store("filename.txt", data)
    ```
  - To retrieve a file:
    ```go
    goCopyr, err := server.Get("filename.txt")
    if err != nil {
      log.Fatal(err)
    }
    contents, err := io.ReadAll(r)
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(string(contents))
    ```
## Contributing

Contributions are welcome! Please read the [contribution guidelines](CONTRIBUTING.md) before submitting a pull request.

## License

This project is licensed under the Apache-2.0 license - see the [LICENSE](LICENSE) file for details.
