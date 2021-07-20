package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	dummytcp "github.com/alex-a-renoire/tcp"
)

func main() {
	//create connection
	conn, err := net.Dial("tcp", dummytcp.TCP_ADDR)
	if err != nil {
		log.Fatalf("error dialing %s: %s", dummytcp.TCP_ADDR, err)
	}

	defer conn.Close()

	//create readers for stdin and for connection
	readerIO := bufio.NewReader(os.Stdin)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

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
		_, err = rw.Write([]byte(text))
		if err != nil {
			log.Printf("failed writing data to connection readwriter: %s", err)
			return
		}

		if err = rw.Flush(); err != nil {
			log.Printf("failed sending the command to server: %s", err)
			return
		}

		log.Print("Text sent to server")

		//reading the answer from the server
		processedData, _, err := rw.ReadLine()
		if err != nil {
			log.Printf("failed reading data from server: %s", err)
			return
		}

		fmt.Printf("The output is: %s \n", processedData)
	}
}
