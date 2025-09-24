# NimbusFS

**NimbusFS** is a distributed, peer-to-peer Content Addressable Storage (CAS) system built in Go that enables decentralized file storage and retrieval across a network of nodes. Inspired by distributed file systems like IPFS and Hadoop's HDFS, NimbusFS provides a robust foundation for building decentralized applications that require reliable, scalable, and secure data storage.

## What is NimbusFS?

NimbusFS is designed to solve the challenges of traditional centralized storage systems by creating a distributed network where:

- **Files are stored across multiple nodes** for redundancy and fault tolerance
- **Content addressing ensures data integrity** - files are identified by their cryptographic hash rather than location
- **Peer-to-peer architecture eliminates single points of failure** 
- **Built-in encryption provides security** for data at rest and in transit
- **Automatic deduplication** reduces storage overhead by storing identical content only once

### Core Concepts

**Content Addressable Storage (CAS)**: Unlike traditional file systems that locate files by path, NimbusFS uses the cryptographic hash of the file content as its unique identifier. This ensures that:
- Identical files are automatically deduplicated
- Data integrity is guaranteed (any corruption changes the hash)
- Files can be retrieved from any node that has the content

**Peer-to-Peer Network**: Nodes in the network can both store and retrieve files, creating a self-sustaining distributed storage system where:
- New nodes can join and leave dynamically
- Data is replicated across multiple nodes for availability
- Network capacity scales automatically as more nodes join

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [How It Works](#how-it-works)
- [Resources](#resources)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Use Cases](#use-cases)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)


## Features

### 🔒 **Content Addressable Storage (CAS)**
- Files are identified by their cryptographic hash (SHA-1), ensuring data integrity
- Automatic deduplication - identical content is stored only once across the network
- Tamper detection - any modification to a file changes its hash identifier

### 🌐 **Peer-to-Peer Architecture**
- Decentralized network with no single point of failure
- Dynamic node discovery and connection management
- Load balancing across available nodes

### ⚡ **Distributed & Scalable**
- Horizontal scaling - add more nodes to increase storage capacity
- Built-in fault tolerance through data replication
- Efficient data retrieval from the nearest available node

### 🔐 **Advanced Security**
- AES encryption for data at rest and in transit
- Cryptographic integrity verification
- Secure peer-to-peer communication protocols

### 🚀 **High Performance**
- Efficient TCP-based transport layer with potential for protocol expansion
- GOB encoding for fast data serialization/deserialization  
- Optimized storage with configurable path transformation functions

### 🛠 **Developer Friendly**
- Simple Go API for file operations (Store, Get, Delete)
- Modular architecture with pluggable components
- Comprehensive test coverage

## Architecture

NimbusFS follows a modular, layered architecture designed for scalability and maintainability:

```
┌─────────────────────────────────────────────┐
│                Application Layer             │
│            (File Operations API)             │
├─────────────────────────────────────────────┤
│               FileServer Core                │
│        (Message Handling & Peer Mgmt)       │
├─────────────────────────────────────────────┤
│          Storage Layer (Store)               │
│     (CAS Path Transform & Encryption)       │
├─────────────────────────────────────────────┤
│           Transport Layer (P2P)              │
│         (TCP Transport & Peers)              │
├─────────────────────────────────────────────┤
│            Network Layer                     │
│         (TCP/IP Communication)               │
└─────────────────────────────────────────────┘
```

### Key Components

#### 1. **FileServer** (`fileserver.go`)
The central orchestrator that manages all file operations and peer interactions:
- **Message Handling**: Processes `MessageStoreFile` and `MessageGetFile` requests
- **Peer Management**: Maintains connections to other nodes in the network
- **Request Routing**: Forwards requests to appropriate nodes when files aren't available locally

#### 2. **Store** (`store.go`)
Handles local file storage with advanced features:
- **Path Transformation**: Uses `CASPathTransformFunc` to create directory structures from file hashes
- **Encryption**: Transparent encryption/decryption using AES
- **File Operations**: Create, read, delete operations with integrity checks
- **Example**: File with hash `4295b3c8f0a940a89eb10af8415c9d8ff3234234` becomes:
  ```
  storage_root/node_id/4295b/3c8f0/a940a/89eb1/0af84/15c9d/8ff32/34234/4295b3c8f0a940a89eb10af8415c9d8ff3234234
  ```

#### 3. **P2P Transport Layer** (`p2p/`)
Manages network communication between nodes:
- **TCPTransport**: Handles TCP connections, listening, and dialing
- **TCPPeer**: Represents individual peer connections with bi-directional communication
- **Message Encoding**: Uses GOB encoding for efficient data serialization
- **Connection Management**: Automatic peer discovery and connection lifecycle management

#### 4. **Cryptographic Layer** (`crypto.go`)
Provides security features:
- **Key Generation**: Random encryption keys for each node
- **AES Encryption**: Stream cipher for encrypting file content
- **Hash Functions**: SHA-1 for content addressing, MD5 for key hashing

## How It Works

### File Storage Process
1. **Client calls `Store(key, data)`**
2. **Content is hashed** to create unique identifier
3. **Data is encrypted** using node's encryption key  
4. **File is stored locally** using CAS path structure
5. **Store message is broadcast** to connected peers
6. **Peers can replicate** the file for redundancy

### File Retrieval Process
1. **Client calls `Get(key)`**
2. **Local storage is checked** first for the file
3. **If not found locally**, a `GetFile` message is sent to connected peers
4. **First peer with the file** streams it back
5. **Data is decrypted** and returned to client
6. **File can optionally be cached** locally

### Network Bootstrap Process  
1. **First node starts** with empty peer list (bootstrap node)
2. **Subsequent nodes connect** to one or more bootstrap nodes
3. **Peer information propagates** through the network
4. **Nodes maintain** optimal connections for performance

## Use Cases

NimbusFS is designed for applications that require distributed, reliable file storage:

### 🌍 **Distributed Web Applications**
- **Content Delivery Networks**: Distribute static assets across geographical regions
- **Decentralized Social Media**: Store user-generated content without central servers
- **Collaborative Platforms**: Share files across distributed teams

### 📊 **Data Archival & Backup**
- **Long-term Storage**: Archive data across multiple locations for durability  
- **Disaster Recovery**: Maintain copies of critical data in geographically distributed nodes
- **Personal Cloud Storage**: Build your own distributed cloud storage system

### 🔬 **Research & Scientific Computing**
- **Dataset Distribution**: Share large datasets across research institutions
- **Collaborative Research**: Distribute compute results and intermediate data
- **Reproducible Science**: Ensure data availability for research validation

### 🏢 **Enterprise Applications**
- **Microservices Storage**: Shared storage backend for containerized applications
- **Edge Computing**: Distribute data to edge nodes for low-latency access
- **Multi-site Deployments**: Synchronize data across multiple office locations

## Security

NimbusFS implements multiple layers of security to protect your data:

### 🔐 **Encryption at Rest**
- **AES Encryption**: All stored files are encrypted using industry-standard AES
- **Unique Keys**: Each node generates its own random encryption key
- **Stream Cipher**: Efficient CTR mode for large file encryption

### 🛡️ **Data Integrity**
- **Content Addressing**: Files are identified by cryptographic hash (SHA-1)
- **Tamper Detection**: Any modification to content changes the hash identifier
- **Verification**: Automatic integrity checking during retrieval

### 🔒 **Network Security**
- **Encrypted Transport**: All peer-to-peer communication is encrypted
- **Handshake Protocol**: Secure connection establishment between peers
- **Identity Verification**: Nodes maintain verified peer relationships

### 🚫 **Privacy Features**
- **No Metadata Leakage**: Only content hashes are shared across the network
- **Selective Sharing**: Nodes only share content they choose to store
- **Access Control**: Files are only accessible via their cryptographic hash

### ⚠️ **Security Considerations**

While NimbusFS provides strong security features, consider these points:

- **Key Management**: Secure storage of encryption keys is critical
- **Network Exposure**: Ensure nodes are properly firewalled in production
- **Hash Algorithms**: SHA-1 is used for compatibility; consider SHA-256 for high-security applications
- **Access Patterns**: Monitor network access to detect unusual activity

## Resources

### Educational Videos

These educational videos provided foundational knowledge and inspiration for building NimbusFS:

- **Distributed Computing - DFS** by [Perfect CE](https://www.youtube.com/@perfectcomputerengineer): [Watch on YouTube](https://youtu.be/Xhi3hqbiXNM)
- **Big Data Analysis - DFS** by [Perfect CE](https://www.youtube.com/@perfectcomputerengineer): [Watch on YouTube](https://youtu.be/McTWc6N-pBg)
- **Architecture of a DFS** by [Knowledge Hub](https://www.youtube.com/@knowledgehub9741): [Watch on YouTube](https://youtu.be/QmNlluPbEEk)
- **Hadoop - HDFS** by [SimpliLearn](https://www.youtube.com/@SimplilearnOfficial): [Watch on YouTube](https://youtu.be/6apXsm_25s0)
- **[Talk on p2p implementation](https://youtu.be/waVtYYSXkXU?si=UFP2YRSx0dxZ1fRc)**
- **[Additional](https://youtu.be/eRndYq8iTio?si=5XuYlcs6FgDIbkxC)**

### Related Technologies

- **[IPFS](https://ipfs.io/)** - The InterPlanetary File System
- **[Hadoop HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html)** - Hadoop Distributed File System
- **[BitTorrent](https://www.bittorrent.org/)** - Peer-to-peer file sharing protocol

## Contributing

Contributions are welcome! Please read the [contribution guidelines](CONTRIBUTING.md) before submitting a pull request.

### Development

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and add tests
4. Run the test suite (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the Apache-2.0 license - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ by [MonalBarse](https://github.com/MonalBarse)**

*NimbusFS - Bringing distributed storage to the cloud and beyond* ☁️
