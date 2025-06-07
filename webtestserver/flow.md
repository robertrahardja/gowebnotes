This Go program implements a simple TCP server that listens for connections and processes custom commands.

## Program Flow

**1. Server Setup (`main` function)**

- `net.Listen("tcp", ":8080")` - Creates a TCP listener on port 8080
- `defer ln.Close()` - Ensures the listener is closed when main() exits
- Enters an infinite loop to continuously accept connections
- `ln.Accept()` - Blocks until a client connects, then returns a connection object
- `go handleConnection(conn)` - Launches each connection in a separate goroutine for concurrent handling

**2. Connection Handling (`handleConnection` function)**

- `defer conn.Close()` - Ensures the connection is closed when the function exits
- `bufio.NewReader(conn)` - Creates a buffered reader for efficient reading from the connection
- `reader.ReadString('\n')` - Reads from the connection until it encounters a newline character
- `strings.TrimSpace(line)` - Removes leading/trailing whitespace from the received line
- `strings.SplitN(..., " ", 2)` - Splits the input into exactly 2 parts: command and resource
- Uses a switch statement to route different commands to appropriate handlers
- `log.Printf()` - Logs the received command and resource for debugging

**3. Command Processing (`handleGet` function)**

- `fmt.Fprintf(conn, ...)` - Writes the response back to the client through the connection
- Currently just echoes back what resource was requested (placeholder implementation)

## Key Functions and Their Roles

- **`net.Listen`** - Creates the server socket
- **`net.Accept`** - Waits for and accepts client connections
- **`bufio.NewReader`** - Provides buffered reading capabilities
- **`strings.SplitN`** - Parses the command protocol
- **`fmt.Fprintf`** - Sends responses back to clients
- **`log.Printf`** - Provides server-side logging

The server expects clients to send commands in the format "COMMAND RESOURCE" (like "GET user123") and currently only implements a GET command handler. Each client connection runs in its own goroutine, allowing the server to handle multiple clients simultaneously.
