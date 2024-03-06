package app

import (
	"net/http"
	"stdrifa/app/company"

	"github.com/jritsema/gotoolbox/web"
)

// GET /company
func (a *App) CompanyIndex(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, a.Templ, "companies.html", company.GetCompanies(), nil)
	res.Write(w)
}

// GET /company/{id}
func (a *App) GetCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := company.GetCompanyByID(id)
	res := web.HTML(http.StatusOK, a.Templ, "row.html", row, nil)
	res.Write(w)
}

// GET /company/add
func (a *App) AddCompany(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, a.Templ, "company-add.html", company.GetCompanies(), nil)
	res.Write(w)
}

// GET /company/{id}/edit
func (a *App) EditCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := company.GetCompanyByID(id)
	res := web.HTML(http.StatusOK, a.Templ, "row-edit.html", row, nil)
	res.Write(w)
}

// POST /company
func (a *App) CreateCompany(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	row := company.Company{}
	row.Company = r.Form.Get("company")
	row.Contact = r.Form.Get("contact")
	row.Country = r.Form.Get("country")
	company.AddCompany(row)
	res := web.HTML(http.StatusOK, a.Templ, "companies.html", company.GetCompanies(), nil)
	res.Write(w)
}

// PUT /company/{id}
func (a *App) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := company.GetCompanyByID(id)
	r.ParseForm()
	row.Company = r.Form.Get("company")
	row.Contact = r.Form.Get("contact")
	row.Country = r.Form.Get("country")
	company.UpdateCompany(row)
	res := web.HTML(http.StatusOK, a.Templ, "row.html", row, nil)
	res.Write(w)
}

// DELETE /company/{id}
func (a *App) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	company.DeleteCompany(id)
	res := web.HTML(http.StatusOK, a.Templ, "companies.html", company.GetCompanies(), nil)
	res.Write(w)
}
