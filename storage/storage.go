package storage

import (
	"errors"
	"log"

	"github.com/rs/xid"
)

// Storage is the storage interface of the in-memory storage
type Storage interface {
	Get(id string) (Resource, error)
	Create(r Resource) error
	Delete(id string) error
	Update(r Resource) error
}

// New creates a new Storage with a pre-defined test data
func New() Storage {
	s := &storage{
		resources: make(map[string]Resource),
	}

	s.Create(Resource{
		Name:          "Test Feri",
		DOB:           "1934.04.01",
		Phone:         "+36305544554",
		NumOfChickens: 1222,
	})

	for id := range s.resources {
		log.Printf("[INFO] Test ID %s\n", id)
	}

	return s
}

// storage the main handler of the storage
type storage struct {
	resources map[string]Resource
}

// Get the data based on the ID
func (s *storage) Get(id string) (Resource, error) {
	for resourceID, resource := range s.resources {
		if id == resourceID {
			return resource, nil
		}
	}
	return Resource{}, errors.New("resource not found")
}

// Create the resource
func (s *storage) Create(r Resource) error {
	id := xid.New().String()
	r.ID = id
	s.resources[id] = r
	log.Printf("[INFO] Data created with ID: %s\n", id)
	return nil
}

// Delete the resource
func (s *storage) Delete(id string) error {
	for resourceID := range s.resources {
		if id == resourceID {
			delete(s.resources, resourceID)
		}
	}
	return nil
}

// Update the resource
func (s *storage) Update(r Resource) error {
	if _, ok := s.resources[r.ID]; r.ID == "" || !ok {
		return errors.New("resource not found")
	}
	s.resources[r.ID] = r
	return nil
}

// Resource is the main representation of the stored data
type Resource struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	DOB           string `json:"dob"`
	Phone         string `json:"phone"`
	NumOfChickens int    `json:"num_of_chickens"`
}
