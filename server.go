package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const (
	ListenPort   = "PORT"
	DbUri        = "DB_URI"
	UserHeader   = "gs-user"
	ResourcesUrl = "/resources"
)

var db *sql.DB

func main() {
	var err error

	listPort := os.Getenv(ListenPort)
	dbUri := os.Getenv(DbUri)

	if listPort == "" || dbUri == "" {
		log.Fatal("Required env vars not found")
	}

	db, err = sql.Open("mysql", dbUri)
	if err != nil {
		log.Fatal("Can't open database")
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Can't connect to database")
	}
	log.Println("Database connection stablished")

	r := mux.NewRouter()
	r.HandleFunc(ResourcesUrl, GetResources).Methods("GET")
	r.HandleFunc(ResourcesUrl, AddResource).Methods("POST")
	r.HandleFunc(ResourcesUrl+"/{resourceId}", ResourceOptions).Methods("OPTIONS")
	r.HandleFunc(ResourcesUrl+"/{resourceId}", GetResource).Methods("GET")
	r.HandleFunc(ResourcesUrl+"/{resourceId}", UpdateResource).Methods("PUT")
	r.HandleFunc(ResourcesUrl+"/{resourceId}", PatchResource).Methods("PATCH")
	r.HandleFunc(ResourcesUrl+"/{resourceId}", DeleteResource).Methods("DELETE")
	r.HandleFunc(ResourcesUrl+"/{resourceId}", ResourceOptions).Methods("OPTIONS")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+listPort, Log(http.DefaultServeMux)))
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		log.Printf("REQ - %v - %v - %v - %v", r.RemoteAddr, r.Method, r.URL, r.Header.Get("gs-user"))
		handler.ServeHTTP(w, r)
	})
}
