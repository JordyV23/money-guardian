package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// APIServer es una estructura que representa un servidor de API.
type APIServer struct {
	linstenAddress string  // El campo linstenAddress almacena la dirección en la que el servidor escuchará las solicitudes.
	store          Storage // El campo store almacena una instancia de Storage que se utilizará para almacenar y recuperar datos.
}

// NewApiServer crea una nueva instancia de APIServer y la inicializa con la dirección IP proporcionada.
// Parámetros:
// - IpAddress: La dirección IP que se utilizará como dirección de escucha del servidor.
// Devuelve:
// - Un puntero a la instancia de APIServer creada.
func NewAPIServer(IpAddress string, store Storage) *APIServer {
	return &APIServer{
		linstenAddress: IpAddress,
		store:          store,
	}
}

// Run inicia el servidor de la API.
// Esta función configura las rutas del enrutador y comienza a escuchar las solicitudes en la dirección IP y puerto especificados en el campo linstenAddress de la estructura APIServer.
// La función utiliza el enrutador mux.NewRouter() para definir las rutas y utiliza makeHttpHandleFunc para manejar las solicitudes.
func (s *APIServer) Run() {
	// Crea un nuevo enrutador utilizando la biblioteca mux.
	router := mux.NewRouter()

	// Define una ruta para el endpoint "/account" y la asocia con la función handleAccount de este servidor.
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))

	// Define una ruta para el endpoint "/account" que recibe un parametro llamado id y la asocia con la función handleGetAccount de este servidor.
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleGetAccountByID))

	// Define una ruta para el endpoint "/transfer" y la asocia con la función handleTransfer del servidor.
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))

	// Registra un mensaje de inicio en el registro de logs.
	log.Print("Starting API server on port: ", s.linstenAddress)

	// Inicia el servidor HTTP y comienza a escuchar en la dirección IP y puerto especificados en linstenAddress.
	http.ListenAndServe(s.linstenAddress, router)
}

// handleAccount funcion para la ruta "/account" en el servidor de la API.
// Esta función inspecciona el método HTTP de la solicitud entrante y enrutará la solicitud a las funciones de manejo correspondientes
// según el método (GET, POST o DELETE).
// Parámetros:
// - w: El objeto http.ResponseWriter para escribir la respuesta.
// - r: El objeto *http.Request que contiene la solicitud HTTP entrante.
// Devuelve:
// - Un error si ocurre algún problema durante el procesamiento de la solicitud, de lo contrario, devuelve nil.
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// Inspecciona el método HTTP de la solicitud entrante.
	switch r.Method {
	case "GET":
		// En caso de solicitud GET, llama a la función handleGetAccountByID para manejar la solicitud.
		return s.handleGetAccount(w, r)
	case "POST":
		// En caso de solicitud POST, llama a la función handleCreateAccount para manejar la solicitud.
		return s.handleCreateAccount(w, r)
	default:
		// Si el método no coincide con ninguno de los casos anteriores, devuelve nil para indicar que la solicitud no se manejará.
		return fmt.Errorf("accion no permitida %s", r.Method)
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// handleGetAccountByID funcion para obtener detalles de una cuenta por su ID o eliminar una cuenta.
// Parámetros:
// - w: El objeto http.ResponseWriter para escribir la respuesta.
// - r: El objeto *http.Request que contiene la solicitud HTTP entrante.
// Devuelve:
// - Un error en caso de que ocurra algún problema durante el procesamiento de la solicitud, de lo contrario, devuelve nil.
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		// Obtiene el ID de la solicitud.
		id, err := getID(r)

		if err != nil {
			return err
		}
		// Obtiene los detalles de la cuenta por su ID.
		account, err := s.store.GetAccountById(id)
		if err != nil {
			return err
		}
		// Responde con un JSON que contiene los detalles de la cuenta.
		return WriteJSON(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		// Si la solicitud es un DELETE, llama a la función handleDeleteAccount para eliminar la cuenta.
		return s.handleDeleteAccount(w, r)
	}
	// Si el método HTTP no es ni GET ni DELETE, devuelve un error indicando que la acción no está permitida.
	return fmt.Errorf("acción no permitida %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	account := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// handleDeleteAccount funcion para eliminar una cuenta por su ID.
// Parámetros:
// - w: El objeto http.ResponseWriter para escribir la respuesta.
// - r: El objeto *http.Request que contiene la solicitud HTTP entrante.
// Devuelve:
// - Un error en caso de que ocurra algún problema durante el procesamiento de la solicitud, de lo contrario, devuelve nil.
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// Obtiene el ID de la solicitud.
	id, err := getID(r)
	if err != nil {
		return err
	}
	// Llama a la función s.store.DeleteAccount para eliminar la cuenta por su ID.
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	// Responde con un JSON que contiene un mensaje indicando que la cuenta ha sido eliminada.
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "cuenta eliminada"})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

// WriteJSON es una función que toma un objeto http.ResponseWriter, un código de estado HTTP, y un valor 'v' que se va a convertir a JSON y escribir en la respuesta.
// Parámetros:
// - w: El objeto http.ResponseWriter en el que se escribirá la respuesta.
// - status: El código de estado HTTP que se establecerá en la respuesta.
// - v: El valor que se convertirá a JSON y se escribirá en la respuesta.
// Devuelve:
// - Un error si ocurre algún problema durante la escritura de la respuesta, de lo contrario, devuelve nil.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	// Establece el encabezado Content-Type para indicar que la respuesta es JSON.
	w.Header().Add("Content-Type", "application/json")

	// Establece el código de estado HTTP en la respuesta.
	w.WriteHeader(status)

	// Utiliza el codificador JSON para convertir el valor 'v' en JSON y escribirlo en la respuesta.
	return json.NewEncoder(w).Encode(v)
}

// apiFunc es un tipo de función que representa una función que maneja solicitudes HTTP.
// Esta función toma un objeto http.ResponseWriter, un puntero a http.Request y devuelve un error.
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError es una estructura que representa un error en una API.
// Contiene un campo "Error" que almacena una descripción o mensaje de error.
type ApiError struct {
	Error string `json:"error"` // El mensaje de error o descripción.
}

// makeHttpHandleFunc toma una función apiFunc y devuelve una función http.HandlerFunc que actúa como un manejador de solicitudes HTTP.
// Parámetros:
// - f: Una función apiFunc que manejará la solicitud HTTP.
// Devuelve:
// - Una función http.HandlerFunc que procesa las solicitudes HTTP utilizando la función apiFunc proporcionada y maneja los errores devueltos por la función.
func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	// Devuelve una función http.HandlerFunc que toma un objeto http.ResponseWriter y un puntero a http.Request.
	return func(w http.ResponseWriter, r *http.Request) {
		// Llama a la función apiFunc proporcionada y verifica si devuelve un error.
		if err := f(w, r); err != nil {
			// Si hay un error, utiliza la función WriteJSON para responder con un código de estado HTTP 400 (Bad Request) y un mensaje de error JSON.
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// getID extrae el ID de un parámetro en la ruta de la solicitud HTTP y lo convierte en un entero.
// Parámetros:
// - r: El objeto *http.Request que contiene la solicitud HTTP entrante.
// Devuelve:
// - El ID extraído como un entero y un posible error si el ID no se puede convertir en un entero.
func getID(r *http.Request) (int, error) {
	// Extrae el ID del parámetro en la ruta.
	idStr := mux.Vars(r)["id"]
	// Convierte el ID de cadena a un entero.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si la conversión falla, devuelve un error indicando que el ID es inválido.
		return id, fmt.Errorf("ID inválido: %s", idStr)
	}
	// Devuelve el ID como un entero y un error nil, indicando que la extracción y conversión fueron exitosas.
	return id, nil
}
