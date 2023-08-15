package middleware

import (
	"bookstore_go/app/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type AuthMiddleWare struct {
	userModel *models.UserModel
}

func NewAuthMiddleWare(db *gorm.DB) *AuthMiddleWare {
	userModel := models.NewUserModel(db)
	return &AuthMiddleWare{userModel: userModel}
}

func (am *AuthMiddleWare) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationToken := r.Header.Get("auth-token")
		if authorizationToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Print("Token is valid")
		user, _ := am.userModel.GetUser(models.UserSearchOptions{AuthToken: authorizationToken})
		log.Print(user)
		if user != nil {
			if user.AuthTokenExpired*1000 >= time.Now().UnixMilli() {
				// Token is valid, call the next handler
				log.Print("Token is valid")
				next.ServeHTTP(w, r)
				return
			} else {
				http.Error(w, "Auth token has expired", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

	})
}
