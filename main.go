package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/voutasaurus/env"
)

func main() {
	logger := log.New(os.Stderr, "hop: ", log.LstdFlags|log.LUTC|log.Llongfile)

	fatal := func(key string) {
		logger.Fatalf("required: %q", key)
	}

	addr := env.Get("HOP_LISTEN").WithDefault(":8080")
	remote := env.Get("HOP_TO").Required(fatal)

	logger.Printf("starting, listening on %q and forwarding to %q", addr, remote)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("listen failed: %v", err)
	}
	defer ln.Close()

	for {
		in, err := ln.Accept()
		if err != nil {
			logger.Fatalf("accept failed: %v", err)
		}

		out, err := net.Dial("tcp", remote)
		if err != nil {
			logger.Fatalf("dial failed: %v", err)
		}

		go func() {
			if _, err := io.Copy(out, in); err != nil {
				logger.Fatalf("copy failed: %v", err)
			}
		}()

		go func() {
			if _, err := io.Copy(in, out); err != nil {
				logger.Fatalf("copy failed: %v", err)
			}
		}()
	}
}
