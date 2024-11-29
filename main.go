package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	ADDR = "0.0.0.0:8000"
)

var (
	conn_id  = 0
	mu       = &sync.Mutex{}
	conn_map = &sync.Map{}
)

type TCPConn struct {
	id   int
	conn net.Conn
}

func main() {
	l, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatal(err)
	}

	serverContext, serverCancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		l.Close()
		conn_map.Range(func(id, tcpConn any) bool {
			log.Println("closing connection with id", id.(int))
			tcpConn.(TCPConn).conn.Close()
			return true
		})
		serverCancel()
	}()

	go func() {
		log.Println("listening on", ADDR)
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println(err)
				break
			}

			go func(conn net.Conn) {
				mu.Lock()
				tcpConn := TCPConn{
					id:   conn_id,
					conn: conn,
				}
				conn_id += 1
				mu.Unlock()

				remoteAddr := conn.RemoteAddr().String()
				log.Printf("new connection from %s with id %d\n", remoteAddr, tcpConn.id)

				conn_map.Store(tcpConn.id, tcpConn)
				if err := handleConn(tcpConn.conn); err != nil {
					if err != net.ErrClosed {
						log.Println(err)
					}
				}
				conn_map.Delete(tcpConn.id)
				log.Printf("connection from %s and id %d closed", remoteAddr, tcpConn.id)
			}(conn)
		}
	}()

	<-serverContext.Done()
}

func handleConn(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}
		conn.Write(buf[:n])
	}

	return nil
}
