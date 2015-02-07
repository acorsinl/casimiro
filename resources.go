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
	"github.com/acorsinl/casimiro/models"
	"github.com/acorsinl/casimiro/system"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Resource struct {
	Id   string `json:"-"`
	Href string `json:"href,omitempty"`
}

/**********************Controler Methods ************************/

// GetResources retrieves all resources for the current logged user
func GetResources(w http.ResponseWriter, r *http.Request) {
	var offset, limit int
	userId := r.Header.Get(UserHeader)
	queryParams, err := system.GetQueryParameters(r.RequestURI)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if queryParams.Get("$offset") == "" || queryParams.Get("$limit") == "" {
		offset = PagingOffset
		limit = PagingLimit
	} else {
		offset, _ = strconv.Atoi(queryParams.Get("$offset"))
		limit, _ = strconv.Atoi(queryParams.Get("$limit"))
	}

	resources, err := getResources(userId, offset, limit)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	if resources == nil {
		system.APIReturn(http.StatusNotFound, err.Error(), w)
		return
	}

	output := system.APIMultipleOutput{}
	output.Data = make([]map[string]interface{}, len(resources))
	for index := range resources {
		output.Data[index] = make(map[string]interface{})
		output.Data[index]["id"] = resources[index].Id
		output.Data[index]["href"] = resources[index].Href
	}
	output.Paging = make(map[string]interface{})
	output.Paging["offset"] = PagingOffset
	output.Paging["limit"] = PagingLimit
	system.APIMultipleResults(http.StatusOK, "OK", output, w)
}

// AddResource creates a new resource owned by the current user
func AddResource(w http.ResponseWriter, r *http.Request) {
	var resource *models.Resource
	var err error

	if err = system.DecodeJSON(r, &resource); err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	resource.Id = system.NewUUID()

	err = models.InsertResource(db, resource)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	data := make(map[string]interface{})
	data["href"] = resource.Href
	data["id"] = resource.Id
	system.APISingleResult(http.StatusCreated, "Resource added", data, w)
}

// GetResource retrieves a resource owned by the current user given
// its resource Id.
func GetResource(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(UserHeader)
	resourceId := mux.Vars(r)["resourceId"]

	/*resource, err := getResource(userId, resourceId)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	if resource == nil {
		system.APIReturn(http.StatusNotFound, "Not found", w)
		return
	}*/

	resource, err := models.GetResourceById(db, userId, resourceId)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	resource.Href = ResourcesUrl + "/" + resource.Id
	data := make(map[string]interface{})
	data["href"] = resource.Href
	data["id"] = resource.Id
	system.APISingleResult(http.StatusOK, "OK", data, w)
}

// UpdateResource allows to full update a record in the database
func UpdateResource(w http.ResponseWriter, r *http.Request) {
	var resource Resource
	userId := r.Header.Get(UserHeader)
	resourceId := mux.Vars(r)["resourceId"]

	err := system.DecodeJSON(r, &resource)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	resource.Id = resourceId
	resource.Href = ResourcesUrl + "/" + resource.Id

	resourceUpdated, err := updateResource(resource, userId)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	if resourceUpdated == false {
		system.APIReturn(http.StatusInternalServerError, "Resource not updated", w)
		return
	}

	data := make(map[string]interface{})
	data["href"] = resource.Href
	data["id"] = resource.Id
	system.APISingleResult(http.StatusOK, "Resource modified", data, w)
}

// PatchResource allows partial updates of a given resource owned
// by the current user
func PatchResource(w http.ResponseWriter, r *http.Request) {
	system.APIReturn(http.StatusNotImplemented, "Patch method not implemented yet", w)
}

// DeleteResource deletes a given resource owned by the current user
func DeleteResource(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(UserHeader)
	resourceId := mux.Vars(r)["resourceId"]

	deletedResource, err := deleteResource(userId, resourceId)
	if err != nil {
		system.APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	if deletedResource == false {
		system.APIReturn(http.StatusInternalServerError, "Resource not deleted", w)
		return
	}

	system.APIReturn(http.StatusOK, "Resource deleted", w)
}

// ResourceOptions returns the Access-Control tier headers for this API resource
func ResourceOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+UserHeader)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

/**********************Model Methods ************************/

func getResources(userId string, offset, limit int) ([]Resource, error) {
	var resources []Resource

	stmt := "? LIMIT ?, ?"
	query, err := db.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query(userId, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resource := Resource{}

		if err := rows.Scan(); err != nil {
			return nil, err
		}
		resource.Href = ResourcesUrl + "/" + resource.Id
		resources = append(resources, resource)
	}

	return resources, nil
}

func deleteResource(userId, resourceId string) (bool, error) {

	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return false, err
	}
	defer query.Close()

	_, err = query.Exec(userId, resourceId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func updateResource(resource Resource, userId string) (bool, error) {
	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return false, err
	}
	defer query.Close()

	_, err = query.Exec(resource, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceExists(resourceId string) (bool, error) {
	var resource Resource

	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return false, err
	}
	defer query.Close()

	err = query.QueryRow(resourceId).Scan(&resource)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, err
	}

	return true, nil
}
