package main

import (
	"net/http"

	"github.com/jritsema/gotoolbox/web"
)

// // Delete -> DELETE /company/{id} -> delete, companys.html
//
// // Edit   -> GET /company/edit/{id} -> row-edit.html
// // Save   ->   PUT /company/{id} -> update, row.html
// // Cancel ->	 GET /company/{id} -> nothing, row.html
//
// // Add    -> GET /company/add/ -> companys-add.html (target body with row-add.html and row.html)
// // Save   ->   POST /company -> add, companys.html (target body without row-add.html)
// // Cancel ->	 GET /company -> nothing, companys.html

// GET /
func index(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, html, "index.html", data, nil)
	res.Write(w)
}

// GET /company/add
func companyAdd(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, html, "company-add.html", nil, nil)
	res.Write(w)
}

// GET /company/edit/{id}
func companyEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := getCompanyByID(id)
	res := web.HTML(http.StatusOK, html, "row-edit.html", row, nil)
	res.Write(w)
}

// GET /company
func companies(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, html, "companies.html", data, nil)
	res.Write(w)
}

// GET /company/{id}
func getCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := getCompanyByID(id)
	res := web.HTML(http.StatusOK, html, "row.html", row, nil)
	res.Write(w)
}

// POST /company
func createCompany(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	row := Company{}
	row.Company = r.Form.Get("company")
	row.Contact = r.Form.Get("contact")
	row.Country = r.Form.Get("country")
	addCompany(row)
	res := web.HTML(http.StatusOK, html, "companies.html", data, nil)
	res.Write(w)
}

// PUT /company/{id}
func saveCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := getCompanyByID(id)
	r.ParseForm()
	row.Company = r.Form.Get("company")
	row.Contact = r.Form.Get("contact")
	row.Country = r.Form.Get("country")
	updateCompany(row)
	res := web.HTML(http.StatusOK, html, "row.html", row, nil)
	res.Write(w)
}

func removeCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	deleteCompany(id)
	res := web.HTML(http.StatusOK, html, "companies.html", data, nil)
	res.Write(w)
}
