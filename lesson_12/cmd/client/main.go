package main

import (
	"bufio"
	"fmt"
	"lesson_12/internal/protocol"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: client <addr> (e.g. localhost:8080)")
	}
	addr := os.Args[1]

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	stdin := bufio.NewReader(os.Stdin)

	fmt.Fprintf(os.Stderr, "Connected to %s. Enter JSON commands (one per line).\n", addr)
	for {
		fmt.Fprint(os.Stderr, "> ")
		line, err := stdin.ReadBytes('\n')
		if err != nil {
			break
		}
		line = trimNewline(line)
		if len(line) == 0 {
			continue
		}

		if err := protocol.WriteMessage(w, line); err != nil {
			log.Println("write:", err)
			break
		}
		if err := w.Flush(); err != nil {
			log.Println("flush:", err)
			break
		}

		respLine, err := protocol.ReadMessage(r)
		if err != nil {
			log.Println("read:", err)
			break
		}
		resp, err := protocol.DecodeResponse(trimNewline(respLine))
		if err != nil {
			fmt.Println(string(respLine))
			continue
		}
		printResponse(resp)
	}
}

func trimNewline(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
		b = b[:len(b)-1]
	}
	return b
}

func printResponse(resp *protocol.Response) {
	if resp.OK {
		if len(resp.Names) > 0 {
			fmt.Println("names:", resp.Names)
		} else if resp.Doc != nil {
			fmt.Printf("doc: %+v\n", resp.Doc)
		} else if len(resp.Docs) > 0 {
			fmt.Printf("docs: %d document(s)\n", len(resp.Docs))
			for i, d := range resp.Docs {
				fmt.Printf("  [%d] %+v\n", i, d)
			}
		} else {
			fmt.Println("ok")
		}
	} else {
		fmt.Println("error:", resp.Err)
	}
}
