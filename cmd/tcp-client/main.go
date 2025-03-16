package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/leetcode-golang-classroom/golang-sample-with-tcp-udp-socket/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-tcp-udp-socket/internal/logger"
)

func main() {
	// setup logger
	ctx := context.Background()
	ulog := logger.FromContext(ctx)
	// Create a TCP socket
	conn, err := net.Dial("tcp", config.AppConfig.ServerURI)
	if err != nil {
		ulog.ErrorContext(ctx, "failed to connect tcp", slog.Any("err", err))
		os.Exit(1)
	}

	// Send a message
	fmt.Fprint(conn, "Hello!")

	// Read a response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		ulog.ErrorContext(ctx, "failed to read", slog.Any("err", err))
	}
	msg := string(buf[:n])
	ulog.InfoContext(ctx, "server reply", slog.String("message", msg))
	conn.Close()
}
