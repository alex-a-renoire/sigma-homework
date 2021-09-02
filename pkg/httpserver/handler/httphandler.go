package httphandler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/service/authservice"
	"github.com/alex-a-renoire/sigma-homework/service/csvservice"
	"github.com/alex-a-renoire/sigma-homework/service/personservice"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//TODO: write a local interface

type HTTPHandler struct {
	service      personservice.PersonService
	csvprocessor csvservice.CsvProcessor
	authservice  authservice.AuthService
}

func New(srv personservice.PersonService, csvprocessor csvservice.CsvProcessor, authservice authservice.AuthService) HTTPHandler {
	return HTTPHandler{
		service:      srv,
		csvprocessor: csvprocessor,
		authservice:  authservice,
	}
}

func (s *HTTPHandler) GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/persons", s.AddPerson).Methods("POST")
	router.HandleFunc("/persons", s.GetAllPersons).Methods("GET")
	router.HandleFunc("/persons/dump", s.DownloadPersonsCSV).Methods("GET")
	router.HandleFunc("/persons/upload", s.RenderTemplate).Methods("GET")
	router.HandleFunc("/persons/upload", s.UploadPersonsCSV).Methods("POST")
	router.HandleFunc("/persons/me", s.MyUser).Methods("GET")
	router.HandleFunc("/persons/{id}", s.GetPerson).Methods("GET")
	router.HandleFunc("/persons/{id}", s.UpdatePerson).Methods("PUT")
	router.HandleFunc("/persons/{id}", s.DeletePerson).Methods("DELETE")

	router.HandleFunc("/login/{id}", s.Login).Methods("GET")

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
	person := model.AddUpdatePerson{}
	if err = json.Unmarshal(p, &person); err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//send the appropriate action to service
	id, err := s.service.AddPerson(person)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	person.Id = id

	//Marshal person to json
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
	id, err := uuid.Parse(vars["id"])
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
	id, err := uuid.Parse(vars["id"])
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
	person := model.AddUpdatePerson{}
	if err = json.Unmarshal(p, &person); err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Ask the service to process action
	err = s.service.UpdatePerson(id, person)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//write the data to response
	w.WriteHeader(http.StatusOK)
}

func (s *HTTPHandler) DeletePerson(w http.ResponseWriter, req *http.Request) {
	log.Print("Command DeletePerson received")

	//get the route variable ID of the person we want to delete
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
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

	ps, err := s.csvprocessor.DownloadPersonsCSV()
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
	//Read a file from the form
	file, _, err := req.FormFile("uploadfile")
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	if err := s.csvprocessor.ProcessCSV(*csv.NewReader(file)); err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

////////
//AUTH//
////////

func (s *HTTPHandler) Login(w http.ResponseWriter, req *http.Request) {
	//get the route variable ID of the person we want to retrieve
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.reportError(w, err, BadRequestErr)
		return
	}

	//Find the person
	person, err := s.service.GetPerson(id)
	if err != nil {
		s.reportError(w, err, InternalServerErr)
		return
	}

	//If the person exists, issue a token
	token, err := s.authservice.GenerateSessionToken(person)
	if err != nil {
		s.reportError(w, fmt.Errorf("failed to generate jwt token for the user: %w", err), InternalServerErr)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (s *HTTPHandler) MyUser(w http.ResponseWriter, req *http.Request) {
	tokenHeader := req.Header.Get("Authorization")

	p, err := s.authservice.MyUser(tokenHeader)
	if err != nil {
		s.reportError(w, fmt.Errorf("failed to authenticate: %w", err), BadRequestErr)
	}

	person, err := json.Marshal(p)
	if err != nil {
		s.reportError(w, fmt.Errorf("failed to marshal person: %w", err), InternalServerErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(person)
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
