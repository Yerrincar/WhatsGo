# WhatsGo

Simple real-time messaging system using TCP sockets in Go.

# Architecture overview

- TCP Server:
  1. Listens for incoming connections
  2. Tracks connected clients
  3. Receives messages from each client
  4. Broadcast messages to all other clients

- TCP Client:
  1. Connects to the Server
  2. Receives messages from other clients
  3. Listen and displays messages from the server

# Message framing: Length-Prefixed Protocol

- Since TCP is a stream-oriented protocol, it doesn't preserver message boundaries. To ensure
  full and clean messages everytime:
  1. 4-byte header (uint32 BigEndian): Message Length
  2. Message body: raw bytes of the message

  Before sending, the client prepends the message with its length.
  The server reads the 4-byte header first, then read the exact number of bytes for the message.

# Concurrency

- Each incoming client is handled in a dedicated goroutine.
- Worker Pool to manage and scale the handling of multiple connections efficiently.
