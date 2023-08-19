package ui

import (
	"context"
	"fmt"
	"net/http"

	"github.com/byte3org/oidc-orbi/internal/models"
	"github.com/gorilla/mux"
)

const (
	SIGN_UP_TEMPLATE = "sign-up"
)

type AuthRouter struct {
	authenticate authenticate
	router       *mux.Router
	callback     func(context.Context, string) string
}

func NewAuthRouter(authenticate authenticate, callback func(context.Context, string) string) *AuthRouter {
	l := &AuthRouter{
		authenticate: authenticate,
		callback:     callback,
	}
	l.createRouter()
	return l
}

func (l *AuthRouter) createRouter() {
	l.router = mux.NewRouter()
	l.router.Path("/sign-in").Methods("GET").HandlerFunc(l.signInHandler)
	l.router.Path("/sign-in").Methods("POST").HandlerFunc(l.checkSignInHandler)

	l.router.Path("/sign-up").Methods("GET").HandlerFunc(l.signUpHandler)
	l.router.Path("/sign-up").Methods("POST").HandlerFunc(l.checkSignUpHandler)
}

type authenticate interface {
	CheckUsernamePassword(username, password, id string) error
	CreateUser(user models.User) (*models.User, error)
}

func (l *AuthRouter) signInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	// the oidc package will pass the id of the auth request as query parameter
	// we will use this id through the login process and therefore pass it to the login page
	renderLogin(w, r.FormValue(queryAuthRequestID), nil)
}

func renderLogin(w http.ResponseWriter, id string, err error) {
	data := &struct {
		ID    string
		Error string
	}{
		ID:    id,
		Error: errMsg(err),
	}
	err = templates.ExecuteTemplate(w, "login", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (l *AuthRouter) checkSignInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	id := r.FormValue("id")
	err = l.authenticate.CheckUsernamePassword(username, password, id)

	if err != nil {
		renderLogin(w, id, err)
		return
	}
	http.Redirect(w, r, l.callback(r.Context(), id), http.StatusFound)
}

// sign up routers

func (l *AuthRouter) signUpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, SIGN_UP_TEMPLATE, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (l *AuthRouter) checkSignUpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}

	user := models.User{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Username:  r.FormValue("username"),
		Phone:     r.FormValue("phone"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
	}

	result, err := l.authenticate.CreateUser(user)

	fmt.Println(result)

	if err != nil {
		templates.ExecuteTemplate(w, SIGN_UP_TEMPLATE, err)
		return
	}

	w.Write([]byte("signed up success"))
}
