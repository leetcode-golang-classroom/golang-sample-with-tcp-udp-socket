package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-tcp-udp-socket/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-tcp-udp-socket/internal/logger"
)

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
	ulog.InfoContext(ctx, "server start", slog.String("port", config.AppConfig.Port))
	defer listener.Close()

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			ulog.ErrorContext(ctx, "failed to accept request", slog.Any("err", err))
			continue
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
