package main

/*
En este caso tenemos una libreria que nos sirve para poder deserializar los jsons
una que es la standar
una que es para los servidores que es la de github
una que se encarga de las peticiones http (get, post etc)
una que se encarga de los mutex que es sync
*/

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

/*
 Creamos una estructura como en el ejemplo de la pagina
 que nos sugirio  https://gowebexamples.com/forms/
*/

type User struct {
	Firstname string
	Lastname  string
	Age       string
}

/*
Creamos una estructura que se va encargar de almacernar
todos los usuarios dentro del server
*/

type ListUsers struct {
	users []User
}

/*
Creamos una funcion que se encargue de recibir los datos que vienen del
formulario del html que se encarga de enviar los datos en formato json
ya que no estoy seguro si no lo especifico no se si se envien como json
u otra cosa entonces el javascript lo hace.

Lo que hace esta funcion es deserializar los datos que vienen desde el
javascript del frontend como json y los convierte a una estructura,
cuando ya se haya creado el objeto y este apunto de escribirse
en una lista de Usuarios
*/

func (s *ListUsers) DecodeJson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recibiendo petición")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error al recibir los datos del json que vienen del html", http.StatusBadRequest)
		return
	}

	fmt.Printf("User: %+v\n", user)

	s.users = append(s.users, user)

	w.WriteHeader(http.StatusOK)
}

// Esta función solo se encarga de enviar el html al cliente en un navegador
func SentHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "forms.html")
}

/*
Esta otra funcion se hace para poder volver a serializar los datos que estan en estructuras y enviarlos como
json a una ruta que en este caso es /users
*/
func (s *ListUsers) ShowUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// main donde se declaran las rutas
func main() {
	list_users := &ListUsers{}

	r := mux.NewRouter()
	r.HandleFunc("/decode", list_users.DecodeJson).Methods(http.MethodPost)
	r.HandleFunc("/users", list_users.ShowUsers).Methods(http.MethodGet)
	r.HandleFunc("/", SentHtml).Methods(http.MethodGet)

	fmt.Println("ListUsers is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
