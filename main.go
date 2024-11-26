package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GDG-on-Campus-KHU/SC4_BE/config"
	"github.com/GDG-on-Campus-KHU/SC4_BE/db"
	"github.com/GDG-on-Campus-KHU/SC4_BE/handlers"
	"github.com/GDG-on-Campus-KHU/SC4_BE/services"
)

func main() {
	conf := config.GetDBConfig()

	err := db.InitDB(conf)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.DB.Close()

	r := mux.NewRouter()

	userService := &services.UserService{}
	userHandler := handlers.NewUserHandler(userService)

	//api
	r.HandleFunc("/api/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/api/register", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/user", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/user/{id}", userHandler.DeleteUser).Methods("DELETE")

	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
