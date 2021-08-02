package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/service"
	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	service service.PersonService
}

func New(service service.PersonService) HTTPHandler {
	return HTTPHandler{
		service: service,
	}
}

func (s *HTTPHandler) GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/persons", s.AddPerson).Methods("POST")
	router.HandleFunc("/persons", s.GetAllPersons).Methods("GET")
	router.HandleFunc("/persons/{id}", s.GetPerson).Methods("GET")
	router.HandleFunc("/persons/{id}", s.UpdatePerson).Methods("PATCH")
	router.HandleFunc("/persons/{id}", s.DeletePerson).Methods("DELETE")

	return router
}

type Error struct {
	Message string `json:"message"`
}

// TODO разобраться со статусами еггогов
func (s *HTTPHandler) reportError(w http.ResponseWriter, err error) {
	httperr := Error{
		Message: err.Error(),
	}

	data, err := json.Marshal(httperr)
	if err != nil {
		log.Println("failed to marshal error message: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest) // TODO handle different errors differently
	w.Write(data)
}

func (s *HTTPHandler) AddPerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command AddPerson received")

	//read new person in json format from the request body
	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//Unmarshal person
	item := model.Person{}
	if err = json.Unmarshal(p, &item); err != nil {
		s.reportError(w, err)
		return
	}

	//send the appropriate action to service
	res, err := s.service.ProcessAction(model.Action{
		FuncName:   "AddPerson",
		Parameters: item,
	},
	)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (s *HTTPHandler) GetPerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command GetPerson received")

	//get the route variable ID of the person we want to retrieve
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.reportError(w, err)
		return
	}

	//Ask the service to process action
	res, err := s.service.ProcessAction(model.Action{
		FuncName:   "GetPerson",
		Parameters: model.Person{Id: id},
	},
	)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func (s *HTTPHandler) GetAllPersons(w http.ResponseWriter, req *http.Request) {
	log.Print("Command GetAllPersons received")

	//Ask the service to process action
	res, err := s.service.ProcessAction(model.Action{
		FuncName:   "GetAllPersons",
		Parameters: model.Person{},
	},
	)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func (s *HTTPHandler) UpdatePerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command UpdatePerson received")

	//get the route variable ID of the person we want to retrieve
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.reportError(w, err)
		return
	}

	//Get the new person name from the json in the req body
	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//unmarshal person
	item := model.Person{}
	if err = json.Unmarshal(p, &item); err != nil {
		s.reportError(w, err)
		return
	}

	//add ID to the person object.

	//In fact, Id could also be transferred through json, but common logic tells that
	//ID of the modified object should be a path variable
	item.Id = id

	//Ask the service to process action
	res, err := s.service.ProcessAction(model.Action{
		FuncName:   "UpdatePerson",
		Parameters: item,
	},
	)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func (s *HTTPHandler) DeletePerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command DeletePerson received")

	//get the route variable ID of the person we want to delete
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.reportError(w, err)
		return
	}

	//Ask the service to process action
	res, err := s.service.ProcessAction(model.Action{
		FuncName: "DeletePerson",
		Parameters: model.Person{
			Id: id,
		},
	},
	)
	if err != nil {
		s.reportError(w, err)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}
