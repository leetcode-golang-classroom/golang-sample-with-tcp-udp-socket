package main

import (
	"context"
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
	// Create UDP connection
	conn, err := net.Dial("udp", config.AppConfig.ServerURI)
	if err != nil {
		ulog.ErrorContext(ctx, "failed to dial", slog.Any("err", err))
		os.Exit(1)
	}
	defer conn.Close()

	// Send message to the server
	message := "Hello UDP Server"
	_, err = conn.Write([]byte(message))
	if err != nil {
		ulog.ErrorContext(ctx, "failed to sending message", slog.Any("err", err))
		return
	}
	ulog.InfoContext(ctx, "Sent:", slog.String("message", message))
	// Set read deadline to avoid blocking indefinitely
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Receive response from server
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		ulog.InfoContext(ctx, "errror reading from server", slog.Any("err", err))
		return
	}
	if n > 0 {
		receivedMsg := string(buf[:n])
		ulog.InfoContext(ctx, "Received:", slog.String("receivedMsg", receivedMsg))
	}
}
