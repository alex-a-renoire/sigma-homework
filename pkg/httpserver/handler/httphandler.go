package httphandler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/service"
	"github.com/gorilla/mux"
	"github.com/jszwec/csvutil"
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
	router.HandleFunc("/persons/dump", s.DownloadPersonsCSV).Methods("GET")
	router.HandleFunc("/persons/upload", s.RenderTemplate).Methods("GET")
	router.HandleFunc("/persons/upload", s.UploadPersonsCSV).Methods("POST")
	router.HandleFunc("/persons/{id}", s.GetPerson).Methods("GET")
	router.HandleFunc("/persons/{id}", s.UpdatePerson).Methods("PATCH")
	router.HandleFunc("/persons/{id}", s.DeletePerson).Methods("DELETE")

	router.Use(s.loggingMiddleware)

	return router
}

var (
	BadRequestErr     = 1
	InternalServerErr = 2
)

type Error struct {
	Message string `json:"message"`
}

// TODO разобраться со статусами еггогов
func (s *HTTPHandler) reportError(w http.ResponseWriter, err error, errType int) {
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
	if errType == BadRequestErr {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

func (s *HTTPHandler) AddPerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command AddPerson received")

	//read new person in json format from the request body
	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Unmarshal person
	person := model.Person{}
	if err = json.Unmarshal(p, &person); err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//send the appropriate action to service
	id, err := s.service.AddPerson(person.Name)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//Marshal person to json
	person.Id = id
	p, err = json.Marshal(person)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/persons/%d", id))
	w.WriteHeader(http.StatusCreated)
	w.Write(p)
}

func (s *HTTPHandler) GetPerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command GetPerson received")

	//get the route variable ID of the person we want to retrieve
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Ask the service to process action
	person, err := s.service.GetPerson(id)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//Marshal person
	p, err := json.Marshal(person)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(p)
}

func (s *HTTPHandler) GetAllPersons(w http.ResponseWriter, req *http.Request) {
	log.Print("Command GetAllPersons received")

	//Ask the service to process action
	persons, err := s.service.GetAllPersons()
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//Marshal persons
	ps, err := json.Marshal(persons)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write(ps)
}

func (s *HTTPHandler) UpdatePerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command UpdatePerson received")

	//get the route variable ID of the person we want to retrieve
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Get the new person name from the json in the req body
	p, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//unmarshal person
	person := model.Person{}
	if err = json.Unmarshal(p, &person); err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Ask the service to process action
	updatedPerson, err := s.service.UpdatePerson(id, person.Name)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//Marshal updated person
	p, err = json.Marshal(updatedPerson)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write(p)
}

func (s *HTTPHandler) DeletePerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command DeletePerson received")

	//get the route variable ID of the person we want to delete
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Ask the service to process action
	if err := s.service.DeletePerson(id); err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

///////
//CSV//
///////

func (s *HTTPHandler) DownloadPersonsCSV(w http.ResponseWriter, req *http.Request) {
	log.Print("Command DownloadPersonsCSV received")

	//Ask the service to process action
	persons, err := s.service.GetAllPersons()
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//Marshal persons into csv
	ps, err := csvutil.Marshal(persons)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	log.Print("Transmitting file to client...")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=myfilename.csv")
	w.Write(ps)
}

func (s *HTTPHandler) RenderTemplate(w http.ResponseWriter, req *http.Request) {
	tmp, err := template.ParseFiles(filepath.Join("/templates", "upload.html"))
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	tmp.Execute(w, struct{}{})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (s *HTTPHandler) UploadPersonsCSV(w http.ResponseWriter, req *http.Request) {
	file, _, err := req.FormFile("uploadfile")
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	lines, err := csv.NewReader(file).ReadAll()
	persons := []model.Person{}

	for _, line := range lines[1:] {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			s.reportError(w, err, BadRequestErr)
			return
		}
		p := model.Person{
			Id:   id,
			Name: line[1],
		}

		persons = append(persons, p)
	}

	ps, err := json.Marshal(persons)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ps)
}

//////////////
//Middleware//
//////////////

func (s *HTTPHandler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received...")
		next.ServeHTTP(w, r)
	})
}
