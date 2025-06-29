package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// Data structures corresponding to data.json

type Specification struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type CarModel struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	ManufacturerID int           `json:"manufacturerId"`
	CategoryID     int           `json:"categoryId"`
	Year           int           `json:"year"`
	Specifications Specification `json:"specifications"`
	Image          string        `json:"image"`
}

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Data struct {
	Manufacturers []Manufacturer `json:"manufacturers"`
	Categories    []Category     `json:"categories"`
	CarModels     []CarModel     `json:"carModels"`
}

// Global dataset
var db Data

func loadData() error {
	f, err := os.Open(path.Join("api", "data.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&db)
}

func main() {
	if err := loadData(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/api", apiRoot)
	mux.HandleFunc("/api/models", handleModels)
	mux.HandleFunc("/api/models/", handleModelByID)
	mux.HandleFunc("/api/manufacturers", handleManufacturers)
	mux.HandleFunc("/api/manufacturers/", handleManufacturerByID)
	mux.HandleFunc("/api/categories", handleCategories)
	mux.HandleFunc("/api/categories/", handleCategoryByID)
	mux.HandleFunc("/api/recommendations", handleRecommendations)
	mux.HandleFunc("/api/models/compare", handleCompare)
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("api/img"))))

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", logRequest(mux)); err != nil {
		fmt.Println(err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, path.Join("web", "index.html"))
}

// /api root
func apiRoot(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"models":        "/api/models",
		"categories":    "/api/categories",
		"manufacturers": "/api/manufacturers",
	}
	writeJSON(w, resp)
}

// helpers
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// Models
func handleModels(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/models" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	models := db.CarModels
	// filtering
	if s := r.URL.Query().Get("search"); s != "" {
		s = strings.ToLower(s)
		filtered := make([]CarModel, 0)
		for _, m := range models {
			if strings.Contains(strings.ToLower(m.Name), s) {
				filtered = append(filtered, m)
			}
		}
		models = filtered
	}

	if mid := r.URL.Query().Get("manufacturerId"); mid != "" {
		if id, err := strconv.Atoi(mid); err == nil {
			filtered := make([]CarModel, 0)
			for _, m := range models {
				if m.ManufacturerID == id {
					filtered = append(filtered, m)
				}
			}
			models = filtered
		}
	}

	if cid := r.URL.Query().Get("categoryId"); cid != "" {
		if id, err := strconv.Atoi(cid); err == nil {
			filtered := make([]CarModel, 0)
			for _, m := range models {
				if m.CategoryID == id {
					filtered = append(filtered, m)
				}
			}
			models = filtered
		}
	}

	writeJSON(w, models)
}

func handleModelByID(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api/models/") {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/models/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for _, m := range db.CarModels {
		if m.ID == id {
			writeJSON(w, m)
			return
		}
	}
	http.Error(w, "model not found", http.StatusNotFound)
}

// Manufacturers
func handleManufacturers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/manufacturers" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, db.Manufacturers)
}

func handleManufacturerByID(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api/manufacturers/") {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/manufacturers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for _, m := range db.Manufacturers {
		if m.ID == id {
			writeJSON(w, m)
			return
		}
	}
	http.Error(w, "manufacturer not found", http.StatusNotFound)
}

// Categories
func handleCategories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/categories" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, db.Categories)
}

func handleCategoryByID(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api/categories/") {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for _, c := range db.Categories {
		if c.ID == id {
			writeJSON(w, c)
			return
		}
	}
	http.Error(w, "category not found", http.StatusNotFound)
}

// Compare endpoint: /api/models/compare?ids=1,2,3
func handleCompare(w http.ResponseWriter, r *http.Request) {
	idsStr := r.URL.Query().Get("ids")
	if idsStr == "" {
		http.Error(w, "ids required", http.StatusBadRequest)
		return
	}
	parts := strings.Split(idsStr, ",")
	result := make([]CarModel, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			continue
		}
		for _, m := range db.CarModels {
			if m.ID == id {
				result = append(result, m)
				break
			}
		}
	}
	writeJSON(w, result)
}

// Recommendations: based on manufacturerId or categoryId
func handleRecommendations(w http.ResponseWriter, r *http.Request) {
	var result []CarModel
	if mid := r.URL.Query().Get("manufacturerId"); mid != "" {
		if id, err := strconv.Atoi(mid); err == nil {
			for _, m := range db.CarModels {
				if m.ManufacturerID == id {
					result = append(result, m)
				}
			}
		}
	} else if cid := r.URL.Query().Get("categoryId"); cid != "" {
		if id, err := strconv.Atoi(cid); err == nil {
			for _, m := range db.CarModels {
				if m.CategoryID == id {
					result = append(result, m)
				}
			}
		}
	}

	if len(result) == 0 {
		// default recommendations: first 3 models
		limit := 3
		if len(db.CarModels) < limit {
			limit = len(db.CarModels)
		}
		result = db.CarModels[:limit]
	}
	writeJSON(w, result)
}
