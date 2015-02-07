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

package models

import (
	"database/sql"
	"github.com/acorsinl/casimiro/system"
)

type Resource struct {
	Id   string `json:"id"`
	Href string `json:"href"`
}

func InsertResource(db *sql.DB, resource *Resource) error {
	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(resource)
	if err != nil {
		return err
	}

	return nil
}

func InsertResourceWithTransaction(db *sql.DB, resource *Resource) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	_, err = query.Exec(resource)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func GetResourceById(db *sql.DB, userId, resourceId string) (*Resource, error) {
	var resource Resource

	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return &Resource{}, err
	}
	defer query.Close()

	err = query.QueryRow(userId, resourceId).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return &Resource{}, err
		}
		return &Resource{}, err
	}

	return &resource, nil
}

func GetResources(db *sql.DB, userId string, offset, limit int) ([]Resource, error) {
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
		resource.Href = system.ResourcesUrl + "/" + resource.Id
		resources = append(resources, resource)
	}

	return resources, nil
}

func ResourceExists(db *sql.DB, resourceId string) (bool, error) {
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

func DeleteResourceById(db *sql.DB, userId, resourceId string) error {
	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(userId, resourceId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateResource(db *sql.DB, resource *Resource, userId string) error {
	stmt := ""
	query, err := db.Prepare(stmt)
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(resource, userId)
	if err != nil {
		return err
	}

	return nil
}
