package main

import (
	"bookstore_go/app/internal/config"
	"bookstore_go/app/internal/controllers"
	"bookstore_go/app/internal/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg := config.MustLoad()
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s)/test?parseTime=true",
		cfg.DataBase.DbUser, cfg.DataBase.DbPassword, cfg.DataBase.DbHost)
	db, err := gorm.Open("mysql", dbConfig)
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Book{})
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	authMiddleware := middleware.NewAuthMiddleWare(db)
	protectedRoutes := r.PathPrefix("/book").Subrouter()
	protectedRoutes.Use(authMiddleware.AuthMiddleware)
	bookController := controllers.NewBookController(db)
	authController := controllers.NewAuthorizationController(db)
	r.HandleFunc("/authorization/", authController.Authorize).Methods("POST")
	protectedRoutes.HandleFunc("/", bookController.CreateBook).Methods("POST")
	protectedRoutes.HandleFunc("/", bookController.GetBook).Methods("GET")
	protectedRoutes.HandleFunc("/{bookId}", bookController.GetBookById).Methods("GET")
	protectedRoutes.HandleFunc("/{bookId}", bookController.UpdateBook).Methods("PUT")
	protectedRoutes.HandleFunc("/{bookId}", bookController.DeleteBook).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Address, r))
}
