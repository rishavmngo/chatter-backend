package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rishavmngo/chat-backend-v2/entities"
	intrf "github.com/rishavmngo/chat-backend-v2/interfac"
	jwts "github.com/rishavmngo/chat-backend-v2/jwt"
	"github.com/rishavmngo/chat-backend-v2/responses"
)

type controller struct {
	db intrf.Store
}

func (controller *controller) Register(w http.ResponseWriter, r *http.Request) {

}

func (controller *controller) GetUserByUID(w http.ResponseWriter, r *http.Request) {

	uid, err := jwts.ExtractTokenID(r)

	if err != nil {
		log.Println(err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// var user *entities.User
	user := entities.User{}
	user.UID = uid

	err = controller.db.GetUserByUID(&user)

	if err != nil {

		responses.ERROR(w, http.StatusUnauthorized, err)
		return

	}

	log.Println("GetUserByUID", user.Username, user.CreatedAt)

	responses.JSON(w, http.StatusOK, user)

	// w.Write([]byte("hello"))
}
func (controller *controller) Login(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := entities.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	log.Println(user.Username, user.Password)

	err = controller.db.GetUser(&user)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("username or password wrong!"))
		return
	}

	token, err := jwts.CreateToken(user.UID)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"token": token})
}
