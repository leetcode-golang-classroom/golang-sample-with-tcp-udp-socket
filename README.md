# golang-sample-with-tcp-udp-socket

  This repository is demo how to use tcp and udp socket in golang

## tcp logic

```golang
func main() {
	// setup logger
	ctx := context.Background()
	ulog := logger.FromContext(ctx)
	// create tcp socket
	listener, err := net.Listen("tcp", config.AppConfig.ServerURI)
	if err != nil {
		ulog.ErrorContext(ctx, "failed to listen", slog.Any("addr", listener.Addr().String()))
		os.Exit(1)
	}
	defer listener.Close()

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			ulog.ErrorContext(ctx, "failed to accept request", slog.Any("err", err))
		}

		// Handle the connection
		go handleConnectionNonBlocking(ctx, conn)
	}
}

// handleConnectionNonBlocking - nonblock reading with setup conn read/write timeout
func handleConnectionNonBlocking(ctx context.Context, conn net.Conn) {
	hlog := logger.FromContext(ctx)
	defer conn.Close()
	for {
		// Set a deadline for reading data
		conn.SetReadDeadline(time.Now().Add(time.Second))
		// Read from client
		buf := make([]byte, 1024)
		// block util read timeout
		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// The read operation timeout, continue waiting for next read
				continue
			} else if !errors.Is(err, io.EOF) {
				// Other errors (client disconnected, etc)
				hlog.ErrorContext(ctx, "Connection closed", slog.Any("err", err))
				break
			}
		}
		if n > 0 {
			msg := string(buf[:n])
			hlog.InfoContext(ctx, "Received", slog.String("message", msg))

			// Set a deadline for writing data
			conn.SetWriteDeadline(time.Now().Add(time.Second))
			// Send a response to the client
			_, err = fmt.Fprintf(conn, "Echo - %s", msg)
			if err != nil {
				hlog.ErrorContext(ctx, "error writing to client", slog.Any("err", err))
				break
			}
		}
	}
}
```
## udp logic
```golang
func main() {
	// setup logger
	ctx := context.Background()
	ulog := logger.FromContext(ctx)
	// setup connection socket with udp
	addr, err := net.ResolveUDPAddr("udp", config.AppConfig.ServerURI)
	if err != nil {
		ulog.InfoContext(ctx, "failed to resolve addr", slog.Any("err", err))
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		ulog.ErrorContext(ctx, "failed to listen udp", slog.Any("err", err))
		os.Exit(2)
	}
	defer conn.Close()

	for {
		// Read a message from the client
		buf := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			ulog.ErrorContext(ctx, "Error reading", slog.Any("err", err))
			continue
		}

		if n > 0 {
			msg := string(buf[:n])
			// Print received message
			ulog.InfoContext(ctx, "Received message", slog.Any("clientAddr", clientAddr), slog.String("message", msg))
			// Send response to the client
			_, err = conn.WriteToUDP([]byte("Hey client!"), clientAddr)
			if err != nil {
				ulog.ErrorContext(ctx, "error writing", slog.Any("err", err))
			}
		}
	}
}
```