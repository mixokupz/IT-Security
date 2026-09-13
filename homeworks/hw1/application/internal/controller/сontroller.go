package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"nsu.secure/internal/database"
	"nsu.secure/internal/model"
)

type ProductController struct {
	Repo *database.Database
}

func NewProductController(r *database.Database) *ProductController {
	return &ProductController{Repo: r}
}

// GET /Products/{id}
func (c *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	Product, err := c.Repo.FindByID(r.Context(), id)
	if errors.Is(err, database.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, Product)
}

// POST /Products
func (c *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var o model.Product
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if o.Name == "" {
		http.Error(w, "title обязателен", http.StatusBadRequest)
		return
	}

	id, err := c.Repo.Save(r.Context(), &o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	o.ID = id
	writeJSON(w, http.StatusCreated, o)
}

// GET /Products?limit=10&offset=0
func (c *ProductController) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit := parseIntDefault(r.URL.Query().Get("limit"), 50)
	offset := parseIntDefault(r.URL.Query().Get("offset"), 0)

	Products, err := c.Repo.FindAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, Products)
}

// DELETE /Products/{id}
func (c *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := c.Repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
