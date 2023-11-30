package main

import (
	"fmt"
	"jwt/config"
	"jwt/routes/auth"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config.ConnectDB()
	// defer config.Database.Close();
	router := mux.NewRouter()
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/logindata", auth.LoginData).Methods("GET")

	c := cors.AllowAll()
	fmt.Print("Running....")
	handler := c.Handler(router)
	http.Handle("/", handlers.LoggingHandler(os.Stdout, handler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}
