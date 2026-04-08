package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "product-api/internal/db"

	_ "github.com/lib/pq"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Image string  `json:"image"`
}

var idCounter = 1

func getProduct(w http.ResponseWriter, r *http.Request) {
	list, err := queries.GetProducts(r.Context())
	if err != nil {
		http.Error(w, "DB lỗi", 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Image string  `json:"image"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "JSON lỗi", 400)
		return
	}

	result, err := queries.CreateProduct(r.Context(), db.CreateProductParams{
		Name:  body.Name,
		Price: body.Price,
		Image: body.Image,
	})
	if err != nil {
		http.Error(w, "DB lỗi", 500)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Sai method", 405)
		return
	}
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lệ", 400)
		return
	}

	var body struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Image string  `json:"image"`
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "JSON lỗi", 400)
		return
	}
	// cập nhật dữ liệu
	result, err := queries.UpdateProduct(r.Context(), db.UpdateProductParams{
		Name:  body.Name,
		Price: body.Price,
		Image: body.Image,
		ID:    int32(id),
	})

	if err != nil {
		http.Error(w, "DB lỗi", 500)
		return
	}

	// trả về kết quả
	json.NewEncoder(w).Encode(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Sai method", 405)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lê", 400)
		return
	}
	err = queries.DeleteProduct(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "DB lỗi", 500)
		return
	}
	fmt.Fprintln(w, "Xóa thành công")
}
func checkPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Sai method", 405)
		return
	}
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID lỗi", 400)
		return
	}

	// tìm sản phẩm
	p, err := queries.GetProductByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Không tìm thấy", 404)
		return
	}
	// chỉ tra về giá
	result := map[string]interface{}{
		"id":    p.ID,
		"name":  p.Name,
		"price": p.Price,
		"image": p.Image,
	}
	json.NewEncoder(w).Encode(result)
}

// tim kiem san pham
func searchProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Sai method", 405)
		return
	}
	keyword := r.URL.Query().Get("name")
	list, err := queries.SearchProducts(r.Context(), keyword)
	if err != nil {
		http.Error(w, "DB lỗi", 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// func main() {
// 	http.HandleFunc("/products", getProduct)
// 	http.HandleFunc("/product", createProduct)
// 	http.HandleFunc("/product/update", updateProduct)
// 	http.HandleFunc("/product/delete", deleteProduct)
// 	http.HandleFunc("/product/price", checkPrice)
// 	http.HandleFunc("/product/search", searchProduct)

// 	fmt.Println("Server đang chạy ở http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }

var queries *db.Queries

func main() {
	conn, err := sql.Open("postgres",
		"postgres://postgres:170206@localhost:5433/testdb?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	queries = db.New(conn)

	http.HandleFunc("/products", getProduct)
	http.HandleFunc("/product", createProduct)
	http.HandleFunc("/product/update", updateProduct)
	http.HandleFunc("/product/delete", deleteProduct)
	http.HandleFunc("/product/price", checkPrice)
	http.HandleFunc("/product/search", searchProduct)
	fmt.Println("Server đang chạy ở http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
