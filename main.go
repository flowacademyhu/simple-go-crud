package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/flowacademyhu/simple-go-crud/storage"
	"github.com/go-chi/chi"
)

func main() {
	// creates a new Server
	server := Server{
		storage: storage.New(),
	}

	// start to listen on :8080
	server.Listen(":8080")
}

// Server is the main struct holds the storage and server related content
type Server struct {
	storage storage.Storage
}

// Get is the GET /resources/{id} endpoint implementation
// calls the storage and it will return content based on that
func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := s.storage.Get(id)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	s.response(w, res)
	return
}

// Create is the POST /resources/ endpoint implementation
// calls the storage and creates the content pased in the body
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	var resource storage.Resource
	err = json.Unmarshal(data, &resource)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	err = s.storage.Create(resource)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	s.success(w, http.StatusCreated)
}

// Delete is the DELETE /resources/{id} endpoint implementation
// calls the storage and delete the resource
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.storage.Delete(id)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	s.success(w, http.StatusAccepted)
	return
}

// Update is the PUT /resources/{id} endpoint implementation
// calls the storage and update the resource based on the ID
func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	var resource storage.Resource
	err = json.Unmarshal(data, &resource)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}
	resource.ID = chi.URLParam(r, "id")

	err = s.storage.Update(resource)
	if err != nil {
		s.errResponse(w, err.Error())
		return
	}

	s.success(w, http.StatusAccepted)
	return
}

// Listen setup the endpoints end start the server
func (s *Server) Listen(addr string) error {
	r := chi.NewRouter()
	r.Route("/resources", func(r chi.Router) {
		r.Get("/{id}", s.Get)
		r.Post("/", s.Create)
		r.Delete("/{id}", s.Delete)
		r.Put("/{id}", s.Update)
	})

	log.Printf("[INFO] Server listening on %s\n", addr)
	return http.ListenAndServe(addr, r)
}

// Response is the response container
type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// response setup the Data field of the Response
func (s *Server) response(w http.ResponseWriter, data interface{}) {
	resp := Response{
		Data: data,
	}

	respJSON, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(respJSON)

	return
}

// success return only a status code
func (s *Server) success(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)

	return
}

// errResponse setup the Error field of the Response
func (s *Server) errResponse(w http.ResponseWriter, err string) {
	resp := Response{
		Error: err,
	}

	respJSON, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(respJSON)

	return
}
