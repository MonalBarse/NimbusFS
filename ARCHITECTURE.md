# NimbusFS Architecture

This document provides a detailed overview of the NimbusFS architecture, explaining the role of each file and how they interact. The document is organized chronologically and logically to help developers understand how the entire codebase works together.

## 1. Entrypoint: main.go

**Purpose**: This is the starting point of the application and demonstrates how to set up and use the NimbusFS distributed file system.

**Functionality**: 
- Creates and initializes multiple FileServer instances to simulate a distributed network
- Sets up the peer-to-peer network connections between servers
- Demonstrates basic use cases of storing and retrieving files across the network
- Shows how nodes can bootstrap into an existing network

**Key Components**:
- `makeServer()`: A factory function that creates a FileServer with all necessary configurations including TCP transport, encryption keys, storage paths, and bootstrap nodes
- `main()`: The entry function that creates three servers (s1, s2, s3), starts them as goroutines, and demonstrates file operations

**How it works**: The main function creates three servers listening on different ports (:3000, :7000, :5000), where s3 bootstraps by connecting to s1 and s2. It then demonstrates storing 20 files, deleting them locally, and retrieving them from the network.

## 2. Core Logic: fileserver.go

**Purpose**: This file contains the primary logic for the distributed file server and serves as the coordination layer between storage and networking.

**Key Components**:

### FileServer struct
The main struct that holds the server's state, including:
- `FileServerOpts`: Configuration options (ID, encryption key, storage root, path transform function, transport, bootstrap nodes)
- `peers`: A map of connected peer nodes
- `store`: Reference to the local storage system
- `quitCh`: Channel for graceful shutdown

### Core Methods

**Store() method**: 
- Handles storing files both locally and across the P2P network
- Uses encryption to secure data before storage
- Broadcasts file availability to connected peers
- Coordinates between local disk storage and network replication

**Get() method**: 
- Retrieves files from local storage or fetches from network peers
- First checks local storage, then queries network if not found locally
- Handles decryption of retrieved data
- Provides a unified interface regardless of data location

**broadcast() method**: 
- Sends messages to all connected peers in the network
- Used for coordinating file operations across the distributed system
- Handles message serialization and transmission

## 3. Local Storage: store.go

**Purpose**: This file manages the storage of files on the local disk using a sophisticated content-addressable storage system.

**Key Concepts**:

### Content-Addressable Storage (CAS)
- Uses a `PathTransformFunc` to determine storage paths based on file content
- The `CASPathTransformFunc` creates unique paths using SHA-1 hashes of file content
- Provides deduplication - identical content is stored only once

### Path Transformation
Example: A file with key "hello.txt" might be stored as:
```
storage_root/user_id/68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff/hello.txt
```

Where `68044/29f74/...` is derived from the SHA-1 hash of "hello.txt".

### Storage Operations

**Write()**: 
- Creates directory structure based on path transformation
- Writes file content to the computed path
- Returns number of bytes written

**Read()**: 
- Locates file using path transformation
- Returns an io.Reader for the file content

**Delete()**: 
- Removes files and cleans up empty directories
- Handles both individual files and entire directory structures

**Has()**: 
- Checks if a file exists without reading its content
- Used for quick existence verification

## 4. Peer-to-Peer Networking: The p2p Directory

The p2p directory contains all networking and communication components, designed with clean interfaces for modularity.

### transport.go
**Purpose**: Defines core interfaces for network communication.

**Interfaces**:
- `Peer`: Represents a remote node with methods for sending data and managing connections
- `Transport`: Abstracts the underlying network protocol (TCP, UDP, WebSocket, etc.)

### tcp_transport.go
**Purpose**: Implements the Transport interface using TCP for reliable network communication.

**Key Components**:
- `TCPPeer`: Concrete implementation of the Peer interface for TCP connections
- `TCPTransport`: Manages TCP listeners, connections, and peer discovery
- Handles both incoming connections (server) and outgoing connections (client)

**Features**:
- Connection management with proper cleanup
- Concurrent handling of multiple peer connections
- Integration with handshake and encoding systems

### message.go
**Purpose**: Defines the structure for inter-node communication.

**RPC Struct**: 
- `From`: Identifier of the sending node
- `Payload`: The actual data being transmitted
- `Stream`: Boolean indicating if this is a data stream or a message

**Constants**:
- `IncomingMessage`: Identifies regular messages
- `IncomingStream`: Identifies file/data streams

### encoding.go
**Purpose**: Handles serialization and deserialization of messages between nodes.

**Decoder Interface**: Provides abstraction for different encoding schemes.

**DefaultDecoder**: 
- Distinguishes between regular messages and file streams
- Reads the first byte to determine message type
- Handles binary data transmission efficiently

**GOBDecoder**: Alternative decoder using Go's GOB encoding for complex data structures.

### handshake.go
**Purpose**: Defines the handshake process when peers connect.

**HandshakeFunc**: Function type that defines how peers authenticate/validate each other.

**NOPHandshakeFunc**: A no-operation handshake used for testing and simple scenarios where no authentication is needed.

## 5. Cryptography: crypto.go

**Purpose**: Provides all cryptographic functions to ensure data security and privacy.

**Key Functions**:

### newEncryptionKey()
- Generates a new 256-bit AES encryption key using cryptographically secure random number generation
- Each FileServer instance gets its own unique key

### copyEncrypt() and copyDecrypt()
- Implement AES encryption in CTR (Counter) mode for stream encryption
- `copyEncrypt()`: Reads from source, encrypts data, and writes to destination
- `copyDecrypt()`: Reads encrypted data, decrypts it, and writes to destination
- Handles initialization vectors (IV) for security

### Utility Functions
- `generateID()`: Creates unique identifiers for nodes
- `hashKey()`: Creates MD5 hashes for key derivation
- `copyStream()`: Helper function for streaming data through cipher operations

**Security Features**:
- Uses AES-256 encryption for data protection
- CTR mode allows for efficient streaming encryption/decryption
- Proper IV handling prevents cryptographic attacks
- Secure random number generation for keys and IVs

## How It All Works Together

1. **Initialization**: `main.go` creates FileServer instances with TCP transports, encryption keys, and storage configurations.

2. **Network Setup**: TCP transports listen for connections and establish peer relationships using the handshake mechanism.

3. **File Storage**: When a file is stored, `fileserver.go` coordinates with `store.go` to save it locally using CAS, then broadcasts its availability to peers.

4. **File Retrieval**: When a file is requested, the system first checks local storage, then queries network peers if needed.

5. **Communication**: All network communication uses the RPC system defined in the p2p package, with messages encoded/decoded appropriately.

6. **Security**: All data operations use the cryptographic functions to ensure confidentiality and integrity.

This architecture provides a robust, scalable, and secure distributed file system that can handle multiple nodes, automatic replication, and efficient content-addressable storage.