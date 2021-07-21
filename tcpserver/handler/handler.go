package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	dummytcp "github.com/alex-a-renoire/tcp"
	"github.com/alex-a-renoire/tcp/storage"
)

type Handler struct {
	Storage storage.Storage
	Message chan string
}

func New(s storage.Storage) Handler {
	return Handler{
		Storage: s,
		Message: make(chan string),
	}
}

func (h *Handler) ProcessAction(action dummytcp.Action) string {
	var response string
	person := action.Parameters

	switch action.FuncName {
	case "AddPerson":
		id := h.Storage.AddPerson(person.Name)
		response = fmt.Sprintf("Person with id %d and name %s added \n", id, person.Name)
	case "UpdatePerson":
		p, err := h.Storage.UpdatePerson(person.Id, person.Name)
		if err != nil {
			response = fmt.Sprintf("Person with id %d not found \n", person.Id)
		} else {
			response = fmt.Sprintf("Person with id %d updated with name %s \n", p.Id, p.Name)
		}
	case "GetPerson":
		p, err := h.Storage.GetPerson(person.Id)
		if err != nil {
			response = fmt.Sprintf("Person with id %d not found \n", person.Id)
		} else {
			response = fmt.Sprintf("Person with id %d has name %s \n", p.Id, p.Name)
		}
	case "DeletePerson":
		if err := h.Storage.DeletePerson(person.Id); err != nil {
			response = fmt.Sprintf("Person with id %d not found \n", person.Id)
		} else {
			response = fmt.Sprintf("Person with id %d deleted \n", person.Id)
		}
	default:
		response = fmt.Sprintf("%s is not a valid command. Try again... \n", action.FuncName)
	}

	return response
}

func (h *Handler) HandleConnection(conn net.Conn) {
	//create connnection readwriter
	reader := bufio.NewReader(conn)
	log.Print("readwriter created")
	var response string

	for {
		log.Print("waiting for the client to send action")

		//read data from connection
		s, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				//TODO: Maybe add numbers to connected clients?
				log.Printf("client closed the connection")
			} else {
				log.Printf("failed reading from connection: %s", err)
			}
			return
		}

		//Unmarshall type of action and parameters
		action := dummytcp.Action{}
		if err := json.Unmarshal(s, &action); err != nil {
			response = fmt.Sprintf("unable to unmarshal action json, some fields are not valid: %s \n", err)
			h.Message <- response
			continue
		}

		if err := action.Validate(); err != nil {
			response = fmt.Sprintf("action json is not valid: %s \n", err)
			h.Message <- response
			continue
		}

		log.Printf("Command received: %s", s)

		//Select the correct action and perform it in the database
		response = h.ProcessAction(action)

		h.Message <- response
		log.Print("message sent to channel")
	}
}

func (h *Handler) WriterToServer(conn net.Conn, quit chan interface{}) {
	writer := bufio.NewWriter(conn)
	var response string

	for {
		select {
		case <-quit:
			response = "abort"
		case m := <-h.Message:
			response = m
		}

		log.Printf("Message received from channel")

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
			log.Print("action completed")
		}
	}
}
