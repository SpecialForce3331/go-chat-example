package main

import (
    "fmt"
    "net"
    "os"
    "strings"
    "bytes"
    "internal/pkg/safeconnmap"
)

const (
    SERVER_LISTEN = "localhost"
    SERVER_PORT = "9999"
    SERVER_PROTOCOL = "tcp"
)


var connections *safeconnmap.SafeConnMap

func main() {
    connections = safeconnmap.NewSafeConnMap()
    l, err := net.Listen(SERVER_PROTOCOL, SERVER_LISTEN+":"+SERVER_PORT)
    if err != nil {
        fmt.Println("Error, can't start listen: ", err)
    }
    defer l.Close()

    messages := make(chan []byte)

    fmt.Println("Listening on " + SERVER_LISTEN+":"+SERVER_PORT)
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err)
            os.Exit(1)
        }

        go handleConnection(conn, messages)
    }
}

func removeConnection(login string) {
    connections.Delete(login)
}

func readMessage(conn net.Conn, buff []byte) (bool) {
    _, err := conn.Read(buff)
    if err != nil {
        fmt.Println("Error reading: ", err)
        return true
    }
    return false
}

func handleConnection(conn net.Conn, messages chan []byte) {
    defer conn.Close()
    var login string

    buff := make([]byte, 1024)
    if login == "" {
        for {
            conn.Write([]byte("Login: "))
            readMessage(conn, buff)
            buffer := bytes.NewBuffer(buff)
            msg, _ := buffer.ReadString('\n')
            msg = strings.TrimSpace(msg)
            if msg != "" {
                login = msg
                if connections.Exists(login) {
                    conn.Write([]byte("Sorry login already in use, please try another!\n"))
                    continue
                }
                connections.Add(login, conn)
                break
            } else {
                conn.Write([]byte("You sent empty message, try again!\n"))
                continue
            }
        }
    }

    for {
        is_conn_closed := readMessage(conn, buff)
        if is_conn_closed {
            removeConnection(login)
            return
        }
        for user_login, c := range connections.Raw() {
            fmt.Println(login +" "+ user_login)
            if login != user_login {
                c.Write(buff)
            }
        }
    }
}
