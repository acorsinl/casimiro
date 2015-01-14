package main

import (
	"net/http"
)

type Resource struct {
	Id   string `json:"-"`
	Href string `json:"href,omitempty"`
}

// GetResources retrieves all resources for the current logged user
func GetResources(w http.ResponseWriter, r *http.Request) {
	var resources []Resource
	var stmt string
	user := r.Header.Get(UserHeader)

	stmt = ""
	query, err := db.Prepare(stmt)
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	rows, err := query.Query()
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	defer rows.Close()

	for rows.Next() {
		resource := Resource{}

		if err := rows.Scan(); err != nil {
			APIReturn(http.StatusInternalServerError, err.Error(), w)
			return
		}
		resource.Href = ResourcesUrl + "/" + resource.Id
		resources = append(resources, resource)
	}

	if len(resources) == 0 {
		APIReturn(http.StatusNotFound, "Not found", w)
		return
	}

	output := APIMultipleOutput{}
	output.Data = make([]map[string]interface{}, len(wishlists))
	for index := range wishlists {
		output.Data[index] = make(map[string]interface{})
		output.Data[index]["id"] = wishlists[index].Hash
		output.Data[index]["href"] = wishlists[index].Href
	}
	APIMultipleResults(http.StatusOK, "OK", output, w)
}

// AddResource creates a new resource owned by the current user
func AddResource(w http.ResponseWriter, r *http.Request) {
	var resource Resource
	var stmt string
	var err error

	if err = DecodeJSON(r, &resource); err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	resource.Id = NewUUID()

	tx, err := db.Begin()
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	stmt = ""
	query, err := db.Prepare(stmt)
	if err != nil {
		tx.Rollback()
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	_, err = query.Exec()
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	tx.Commit()

	data := make(map[string]interface{})
	data["href"] = resource.Href
	data["id"] = wishlist.Id
	APISingleResult(http.StatusCreated, "Resource added", data, w)
}

// GetResource retrieves a resource owned by the current user given
// its resource Id.
func GetResource(w http.ResponseWriter, r *http.Request) {
	var resource Resource
	var stmt string
	user := r.Header.Get(UserHeader)
	resourceId := mux.Vars(r)["resourceId"]

	stmt = ""
	query, err := db.Prepare(stmt)
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	err = query.QueryRow().Scan()
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	resource.Href = ResourcesUrl + "/" + resource.Id
	data := make(map[string]interface{})
	data["href"] = wishlist.Href
	data["id"] = wishlist.Hash
	APISingleResult(http.StatusOK, "OK", data, w)
}

// UpdateResource allows to full update a record in the database
func UpdateResource(w http.ResponseWriter, r *http.Request) {
	var resource Resource
	var stmt string
	user := r.Header.Get(UserHeader)
	resourceId := mux.Vars(r)["resourceId"]

	err := DecodeJSON(r, &resource)
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	resource.Id = resourceId
	resource.Href = ResourcesUrl + "/" + resource.Id

	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	_, err = query.Exec()
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	data := make(map[string]interface{})
	data["href"] = resource.Href
	data["id"] = resource.Id
	APISingleResult(http.StatusOK, "Resource modified", data, w)
}

// PatchResource allows partial updates of a given resource owned
// by the current user
func PatchResource(w http.ResponseWriter, r *http.Request) {
	APIReturn(http.StatusNotImplemented, "Patch method not implemented yet", w)
}

// DeleteResource deletes a given resource owned by the current user
func DeleteResource(w http.ResponseWriter, r *http.Request) {
	var stmt string
	user := r.Header.Get(UserHeader)
	resourceId := mux.Vars(r)["resourceId"]

	stmt = ""
	query, err := db.Prepare(stmt)
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}
	_, err = query.Exec()
	if err != nil {
		APIReturn(http.StatusInternalServerError, err.Error(), w)
		return
	}

	APIReturn(http.StatusOK, "Resource deleted", w)
}

// ResourceOptions returns the Access-Control tier headers for this API resource
func ResourceOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+UserHeader)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
