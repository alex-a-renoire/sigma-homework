package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var ch = make(chan string)

type config struct {
	TCPAddr     string
}

func getCfg() config {
	TCPAddr := os.Getenv("TCP_ADDR")
	if TCPAddr == "" {
		TCPAddr = "127.0.0.1:8080"
	}

	return config {
		TCPAddr: TCPAddr,
	}
}

func ClientLoop(conn net.Conn) {
	//create readers for stdin and for connection
	writerConn := bufio.NewWriter(conn)
	readerIO := bufio.NewReader(os.Stdin)

	var processedData string

	for {
		//reading a command from the user
		fmt.Println("Please enter the value: ")

		text, err := readerIO.ReadString('\n')
		if err != nil {
			log.Printf("error readind from stdio: %s", err)
			return
		}

		//checking if entered data is a valid json
		if !json.Valid([]byte(text)) {
			fmt.Println("The input data is not a valid JSON. Try again...")
			continue
		}

		text = fmt.Sprint(text) // to add a newline

		//sending to server
		_, err = writerConn.Write([]byte(text))
		if err != nil {
			log.Printf("failed writing data to connection readwriter: %s", err)
			return
		}

		if err = writerConn.Flush(); err != nil {
			log.Printf("failed sending the command to server: %s", err)
			return
		}

		log.Print("Text sent to server")

		processedData = <-ch

		fmt.Printf("%s \n", processedData)
	}
}

func ResponseReadLoop(conn net.Conn) {
	//create reader
	readerConn := bufio.NewReader(conn)

	for {
		//reading the answer from the server
		processedData, _, err := readerConn.ReadLine()
		if err != nil {
			log.Printf("failed reading data from server: %s", err)
			continue
		}

		if string(processedData) == "abort" {
			return
		}

		ch <- string(processedData)
	}
}

func main() {
	cfg := getCfg()

	wg := sync.WaitGroup{}

	//create connection
	conn, err := net.Dial("tcp", cfg.TCPAddr)
	if err != nil {
		log.Fatalf("error dialing %s: %s", cfg.TCPAddr, err)
	}

	wg.Add(1)
	go ClientLoop(conn)
	go func() {
		ResponseReadLoop(conn)
		wg.Done()
	}()

	wg.Wait()

	log.Print("Connection closed by server")

	defer conn.Close()
}
