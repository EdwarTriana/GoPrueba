package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// se llama una sola tarea
type tarea struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}



type allTareas []tarea



// se llaman varias tareas
var tareas = allTareas{
	{
		ID:      1,
		Name:    "Primer Nombre",
		Content: "nuevos contenidos",
	},
}



func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi API!")
}


// funciono metodo para crear una Tarea
func createTarea(w http.ResponseWriter, r *http.Request) {
	var newTarea tarea
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte un dato valido")
	}

	json.Unmarshal(reqBody, &newTarea)
	newTarea.ID = len(tareas) + 1
	tareas = append(tareas, newTarea)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTarea)

}


// funcion o metodo para consultar una Tarea
func getTareas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tareas)
}

func getOneTarea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tareaID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	for _, tarea := range tareas {
		if tarea.ID == tareaID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tarea)
		}
	}
}


//funcion o metodo para actualizar una tarea
func updateTarea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tareaID, err := strconv.Atoi(vars["id"])
	var updatedTarea tarea

	if err != nil {
		fmt.Fprintf(w, "ID invalido")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Por favor ingrese un dato valido")
	}
	json.Unmarshal(reqBody, &updatedTarea)

	for i, t := range tareas {
		if t.ID == tareaID {
			tareas = append(tareas[:i], tareas[i+1:]...)

			updatedTarea.ID = t.ID
			tareas = append(tareas, updatedTarea)

			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(updatedTask)
			fmt.Fprintf(w, "la tarea con el ID %v se ha actualizado con exito", tareaID)
		}
	}

}

//funcion o metodo para eliminar una tarea
func deleteTarea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tareaID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "ID de usuario invalido")
		return
	}

	for i, t := range tareas {
		if t.ID == tareaID {
			tareas = append(tareas[:i], tareas[i+1:]...)
			fmt.Fprintf(w, "la tarea con el ID %v ha sido removida con exito", tareaID)
		}
	}
} 


//rutas  y aplicacion de los metodos
func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tareas", createTarea).Methods("POST")
	router.HandleFunc("/tareas", getTareas).Methods("GET")
	router.HandleFunc("/tareas/{id}", getOneTarea).Methods("GET")
	router.HandleFunc("/tareas/{id}", deleteTarea).Methods("DELETE")
	router.HandleFunc("/tareas/{id}", updateTarea).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
