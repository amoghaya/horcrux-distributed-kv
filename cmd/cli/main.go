package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	input := bufio.NewReader(os.Stdin)

	fmt.Println("Horcrux CLI connected. Type commands (SET/GET/DEL).")

	for {
		fmt.Print("> ")
		text, _ := input.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			break
		}

		fmt.Fprintf(conn, "%s\n", text)

		resp, _ := reader.ReadString('\n')
		fmt.Println(resp)
	}
}
