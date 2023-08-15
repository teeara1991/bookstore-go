package controllers

import (
	"bookstore_go/app/internal/models"
	"bookstore_go/app/internal/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

var User models.User

type AuthorizationController struct {
	userModel *models.UserModel
}

func NewAuthorizationController(db *gorm.DB) *AuthorizationController {
	userModel := models.NewUserModel(db)
	return &AuthorizationController{userModel: userModel}
}

type Token struct {
	AuthToken string `json:"authtoken"`
}

func (ac *AuthorizationController) Authorize(w http.ResponseWriter, r *http.Request) {
	GotUser := &models.User{}
	utils.ParseBody(r, GotUser)
	user, db := ac.userModel.GetUser(models.UserSearchOptions{Login: GotUser.Login, Password: GotUser.Password})
	token := Token{}

	if user != nil {
		if user.AuthToken == "" || user.AuthTokenExpired*1000 < time.Now().UnixMilli() {
			tokenString, tokenExpired := CreateToken()
			user.AuthToken = tokenString
			user.AuthTokenExpired = tokenExpired
			db.Save(&user)
			token.AuthToken = tokenString

		} else {
			token.AuthToken = user.AuthToken
		}
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(token)
	w.Write(jsonData)
}

func CreateToken() (string, int64) {
	var mySigningKey = []byte("secret")
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, _ := token.SignedString(mySigningKey)
	tokenExpired := time.Now().Add(time.Hour * 24).Unix()

	return tokenString, tokenExpired

}
