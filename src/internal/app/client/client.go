package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
)

const (
    SERVER_HOST = "localhost"
    SERVER_PORT = "9999"
    SERVER_PROTOCOL = "tcp"
)

func main() {
    conn, err := net.Dial(SERVER_PROTOCOL, SERVER_HOST+":"+SERVER_PORT)
    if err != nil {
        fmt.Println("Error connect to server: ", err)
        os.Exit(1)
    }

    defer conn.Close()
    go readMessages(conn)

    reader := bufio.NewReader(os.Stdin)
    for {
        text, _ := reader.ReadString('\n')
        conn.Write([]byte(text))
    }
}

func readMessages(conn net.Conn) {
    buff := make([]byte, 1024)
    for {
        _, err := conn.Read(buff)
        if err != nil {
            fmt.Println("Error read from connection: ", err)
            os.Exit(1)
        }
        fmt.Println(string(buff[:]))
    }
}
