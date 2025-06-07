## Program Flow by Function Calls

### 1. `net.Dial("tcp", "localhost:8080")`

- Establishes TCP connection to server
- Returns connection object and error

### 2. `log.Fatal(err)` (if error)

- Terminates program if connection fails
- Program exits here if server unavailable

### 3. `defer conn.Close()`

- Schedules connection closure for function exit
- Ensures cleanup happens automatically

### 4. `fmt.Fprintf(conn, "GET /index.html\n")`

- Writes request string to connection
- Sends data to server over network

### 5. `make([]byte, 1024)`

- Creates 1024-byte buffer
- Prepares memory for incoming response

### 6. `conn.Read(bs)`

- Reads server response into buffer
- Returns bytes read count and error

### 7. `log.Fatal(err)` (if read error)

- Terminates program if read operation fails
- Program exits here if read problems occur

### 8. `fmt.Println(string(bs[:n]))`

- Converts buffer bytes to string
- Prints server response to console

### 9. `conn.Close()` (deferred)

- Executes automatically when main() ends
- Closes network connection
