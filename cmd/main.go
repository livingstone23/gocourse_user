package main

import (
	"fmt"
	"git/gocourse_user/internal/user"
	"git/gocourse_user/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	//realizamos el ruteo con el paquete de gorilla.mux
	router := mux.NewRouter()
	_ = godotenv.Load()
	//Habilitamos la funcion de nuestro paquete
	l := bootstrap.InitLogger()

	//Habilitamos la funcion de conexion
	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	//Especificamos el repositorio
	userRepo := user.NewRepo(l, db)
	//Especificamos el servicio
	userSrv := user.NewService(l, userRepo)
	//Importamos nuestro paquete de carpeta interna
	userEnd := user.MakeEndPoints(userSrv, user.Config{LimPageDef: pagLimDef})

	//Llamamos a nuestros endpoints
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	//Levantar el servidor, brindamos propiedades
	srv := &http.Server{
		Handler:           router,
		Addr:              address,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
