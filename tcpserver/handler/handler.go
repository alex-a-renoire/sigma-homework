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
}

func New(s storage.Storage) Handler {
	return Handler{
		Storage: s,
	}
}

func (h *Handler) HandleConnection(conn net.Conn) {
	//create connnection readwriter
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	log.Print("readwriter created")

	for {
		//read data from connection
		s, err := rw.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				//TODO: Maybe add numbers to connected clients?
				log.Printf("client closed the connection")
			} else {
				log.Printf("failed reading from connection: %s", err)
			}
			return
		}
		log.Printf("Command received: %s", s)

		/* --------------------------Logic goes here ---------------------------*/
		//Unmarshall type of action and parameters
		action := dummytcp.Action{}
		if err := json.Unmarshal(s, &action); err != nil {
			log.Printf("unable to unmarshal action json: %s", err)
			return
		}

		//Select the correct action and perform it in the database
		person := action.Parameters
		var response string

		switch action.FuncName {
		case "AddPerson":
			id := h.Storage.AddPerson(person.Name)
			response = fmt.Sprintf("Person with id %d added \n", id)
		case "UpdatePerson":
			p := h.Storage.UpdatePerson(person.Id, person.Name)
			response = fmt.Sprintf("Person with id %d updated with name %s \n", p.Id, p.Name)
		case "GetPerson":
			p := h.Storage.GetPerson(person.Id)
			response = fmt.Sprintf("Person with id %d has name %s \n", p.Id, p.Name)
		case "DeletePerson":
			h.Storage.DeletePerson(person.Id)
			response = fmt.Sprintf("Person with id %d deleted \n", person.Id)
		}

		/*----------------------EndLogic--------------------------*/
		_, err = rw.Write([]byte(response))
		if err != nil {
			log.Printf("failed sending data back to client: %s", err)
		}

		//Send the response from database to client
		rw.Flush()
		log.Print("action completed")
	}
}
