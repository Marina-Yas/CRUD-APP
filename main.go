package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

type Driver struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	RacingNumber int    `json:"racing_number"`
	PersonID     int    `json:"id"`
	Active       bool   `json:"active"`
}

var drivers = []Driver{
	{FirstName: "Lando", LastName: "Norris", RacingNumber: 4, PersonID: 0, Active: true},
	{FirstName: "Daniel", LastName: "Ricciardo", RacingNumber: 3, PersonID: 1, Active: true},
	{FirstName: "Max", LastName: "Verstappen", RacingNumber: 33, PersonID: 3, Active: true},
}

func indexByID(drivers []Driver, id string) int {

	for i := 0; i < len(drivers); i++ {
		id, _ := strconv.Atoi(id)
		if drivers[i].PersonID == id {
			return i
		}
	}
	return -1
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/drivers", getDrivers).Methods("GET")
	router.HandleFunc("/drivers", createDriver).Methods("POST")
	router.HandleFunc("/drivers/{id}", getDriverByID).Methods("GET")
	router.HandleFunc("/drivers/{id}", deleteDriver).Methods("DELETE")
	router.HandleFunc("/drivers/{id}", updateDriver).Methods("PATCH")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

func getDrivers(w http.ResponseWriter, r *http.Request) {

	res, _ := json.Marshal(drivers)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
func getDriverByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println("Error while parsing")
	}

	res, _ := json.Marshal(drivers[id])
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func createDriver(w http.ResponseWriter, r *http.Request) {
	driver := Driver{}
	json.NewDecoder(r.Body).Decode(&driver)

	drivers = append(drivers, driver)
	response, err := json.Marshal(&driver)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)

}

func deleteDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	driver := drivers[id]
	driver.Active = false
	// Put updates back to the list

	index := indexByID(drivers, params["id"])
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	drivers = append(drivers[:index], drivers[index+1:]...)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&driver)
}

func updateDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"]
	index := indexByID(drivers, params)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	dr := Driver{}
	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	drivers[index] = dr
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&dr); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}

}
