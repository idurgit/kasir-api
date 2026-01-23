package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int     `json:"id"`
	Nama  string  `json:"nama"`
	Harga float64 `json:"harga"`
	Stok  int     `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Mie Sedap Goreng", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 600ml", Harga: 6000, Stok: 20},
	{ID: 3, Nama: "Kecap ABC 275ml", Harga: 12000, Stok: 15},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	// get data dari request body
	var updatedProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updatedProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// loop produk cari id yg sesuai, ganti sesuai request body
	for i := range produk {
		if produk[i].ID == id {
			updatedProduk.ID = id
			produk[i] = updatedProduk
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduk)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func deleteProdukByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// ganti id ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	// loop produk cari id yg sesuai, hapus data
	for i := range produk {
		if produk[i].ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index i
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Produk deleted successfully"})
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func main() {
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProdukByID(w, r)
		} else if r.Method == "DELETE" {
			deleteProdukByID(w, r)
		}
	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			// masukkan data ke variabel produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(produkBaru)
		}
	})
	// localhost:8080 / health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API is running"})
	})
	fmt.Println("Starting server on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server")
	}
}
