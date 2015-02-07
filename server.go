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
	"github.com/acorsinl/casimiro/controllers/api"
	"github.com/acorsinl/casimiro/models"
	"github.com/acorsinl/casimiro/system"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const (
	ListenPort = "PORT"
	DbUri      = "DB_URI"
)

//var db *sql.DB

func main() {
	listPort := os.Getenv(ListenPort)
	dbUri := os.Getenv(DbUri)

	if listPort == "" || dbUri == "" {
		log.Fatal("Required env vars not found")
	}

	model := models.Model{}
	model.InitDB(dbUri)

	r := mux.NewRouter()
	r.HandleFunc(system.ResourcesUrl, api.GetResources).Methods("GET")
	r.HandleFunc(system.ResourcesUrl, api.AddResource).Methods("POST")
	r.HandleFunc(system.ResourcesUrl+"/{resourceId}", api.ResourceOptions).Methods("OPTIONS")
	r.HandleFunc(system.ResourcesUrl+"/{resourceId}", api.GetResource).Methods("GET")
	r.HandleFunc(system.ResourcesUrl+"/{resourceId}", api.UpdateResource).Methods("PUT")
	r.HandleFunc(system.ResourcesUrl+"/{resourceId}", api.PatchResource).Methods("PATCH")
	r.HandleFunc(system.ResourcesUrl+"/{resourceId}", api.DeleteResource).Methods("DELETE")
	r.HandleFunc(system.ResourcesUrl+"/{resourceId}", api.ResourceOptions).Methods("OPTIONS")
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
