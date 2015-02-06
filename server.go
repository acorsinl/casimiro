/*
Copyright (c) 2015, Alberto Corsín Lafuente
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
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
	PagingOffset = 0
	PagingLimit  = 10
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

	log.Println("Server listening on port " + listPort)
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
