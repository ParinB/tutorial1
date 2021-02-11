package hello

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/handlers/hello/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}
func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJson(rw)
	//	err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	//rw.Write(d)
}
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle  post request")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	// prod := &data.Product{}
	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw,"Unable to unmarshal json",http.StatusBadRequest)
	// }
	// p.l.Printf("Prod : %#v",prod)
	data.AddProduct(&prod)
}
func (p Products) UpdateProducts(rw http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	id,err := strconv.Atoi(vars["id"])
	if err !=nil {
		http.Error(rw,"Unable to covert id",http.StatusBadRequest)
		return
	}
	p.l.Println("Handle  PUT  Product",id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id,&prod)
	if err == data.ErrProductNotFound {
		http.Error(rw,"Product  not  found",http.StatusNotFound)
		return
	}
	if err !=nil {
		http.Error(rw,"Product not found",http.StatusInternalServerError)
	}
}

type KeyProduct struct {}

// func (p Products) MiddlewareProductsValidation(next http.Handler) http.Handler {
// 	return http.HandlerFunc(rw http.ResponseWriter, r*http.Request) {
// 		prod := &data.Product{}
// 		err  :=  prod.FromJson(r.Body)
// 		if err != nil {
// 			http.Error(rw,"Unable  to  marshal json",http.StatusBadRequest)
// 		}
// 		ctx := r.Context().WithValue(KeyProduct{},prod)
// 		req := r.WithContext(ctx)
// 		next.ServeHTTP(rw,r)
// 	}
//}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product",err)
			http.Error(rw,"Error reading product",http.StatusBadRequest)
			return
		}
		//validate the  product 
		err = prod.Validate()
		if err !=nil {
			p.l.Println("[ERROR] validating product",err)
			http.Error(
				rw ,
				fmt.Sprintf("Error validating product: %s",err),
				http.StatusBadRequest,
			)
		}
		ctx := context.WithValue(r.Context(),KeyProduct{},prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw,r)
	})
}