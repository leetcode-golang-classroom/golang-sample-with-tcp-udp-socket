package main

import (
	"context"
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
