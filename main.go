package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"./connect"
	"./structures"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	connect.InitializeDatabase()
	defer connect.CloseConnection()
	// crear enrutador
	r := chi.NewRouter()
	r.Use(middleware.Timeout(10000 * time.Millisecond)) // 10 Seg
	r.Use(middleware.NoCache)                           // deshabilitar cache
	r.Use(middleware.Logger)                            // registros de solicitudes
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/user", func(r chi.Router) {
		r.Get("/{id}", GetUser)
		r.Post("/new", NewUser)
		r.Patch("/update/{id}", UpdateUser)
		r.Delete("/delete/{id}", DeleteUser)
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "id")
	user := connect.GetUser(user_id)
	var response structures.Response
	if !user.IsValid() {
		response = ServerResponseError("User not found.")
	} else {
		response = ServerResponseOk(user, "")
	}
	json.NewEncoder(w).Encode(response)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	userRequest, err := GetUserRequest(r)
	var response structures.Response
	if err != nil {
		response = ServerResponseError("No se han podido leer los datos")
	} else {
		user := connect.CreateUser(userRequest)
		response = ServerResponseOk(user, "Usuario creado correctamente")
	}
	json.NewEncoder(w).Encode(response)
}

func GetUserRequest(r *http.Request) (structures.User, error) {
	var user structures.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return structures.User{}, errors.New("NO se pudo crear el usuario")
	}
	return user, nil
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "id")
	userRequest, err := GetUserRequest(r)
	var response structures.Response
	if err != nil {
		response = ServerResponseError("No se han podido leer los datos")
	} else {
		user := connect.GetUser(user_id)
		if !user.IsValid() {
			response = ServerResponseError("User not found.")
		} else {
			user := connect.UpdateUser(user_id, userRequest)
			response = ServerResponseOk(user, "")
		}
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "id")
	var response structures.Response

	user := connect.GetUser(user_id)
	if !user.IsValid() {
		response = ServerResponseError("User not found.")
	} else {
		connect.DeleteUser(user_id)
		response = ServerResponseOk(structures.User{}, "User delete.")
	}
	json.NewEncoder(w).Encode(response)
}

func ServerResponseOk(data structures.User, message string) structures.Response {
	return ServerResponse(
		http.StatusBadRequest,
		data,
		message,
	)
}

func ServerResponseError(message string) structures.Response {
	return ServerResponse(
		http.StatusOK,
		structures.User{},
		message,
	)
}

func ServerResponse(status int, data structures.User, message string) structures.Response {
	return structures.Response{Status: status, Data: data, Message: message}
}
