package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/alex-a-renoire/sigma-homework/model"
	tcpcontroller "github.com/alex-a-renoire/sigma-homework/pkg/tcpserver/controller"
)

type Handler struct {
	controller tcpcontroller.PersonControllerTCP
}

func New(c tcpcontroller.PersonControllerTCP) Handler {
	return Handler{
		controller: c,
	}
}

func (h *Handler) HandleConnection(conn net.Conn, message chan string, connNumber int) {
	//create connnection readwriter
	reader := bufio.NewReader(conn)
	log.Print("readwriter created")
	var response string

	for {
		log.Printf("waiting for the client #%d to send action", connNumber)

		//read data from connection
		s, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("client #%d closed the connection", connNumber)
			} else {
				log.Printf("failed reading from connection %d: %s", connNumber, err)
			}
			close(message)
			return
		}

		//Unmarshall type of action and parameters
		action := model.Action{}
		if err := json.Unmarshal(s, &action); err != nil {
			response = fmt.Sprintf("unable to unmarshal action json, some fields are not valid: %s \n", err)
			message <- response
			continue
		}

		if err := action.Validate(); err != nil {
			response = fmt.Sprintf("action json is not valid: %s \n", err)
			message <- response
			continue
		}

		log.Printf("Command received: %s", s)

		//Select the correct action and perform it in the database
		response, err = h.controller.ProcessAction(action)
		if err != nil {
			response = fmt.Sprintf("error processing action: %s \n", err)
			message <- response
			continue
		}

		message <- response
		log.Print("message sent to server writer loop")
	}
}

func (h *Handler) WriterToServer(conn net.Conn, message chan string, quit chan interface{}, connNumber int) {
	writer := bufio.NewWriter(conn)
	var response string

	for {
		select {
		case <-quit:
			response = "abort"
		case m := <-message:
			response = m
		}

		log.Printf("Server Writer Loop for conn %d: Message received from handler", connNumber)
		log.Print(response)

		_, err := writer.Write([]byte(response))
		if err != nil {
			log.Printf("failed sending data back to client: %s", err)
			continue
		}

		//Send the response from database to client
		if err := writer.Flush(); err != nil {
			log.Print("error flushing the data")
			continue
		}

		if response == "abort" {
			return
		} else {
			log.Printf("action for conn %d completed", connNumber)
		}
	}
}
