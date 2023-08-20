package ui

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/byte3org/oidc-orbi/internal/models"
	"github.com/byte3org/oidc-orbi/internal/storage"
	"github.com/byte3org/oidc-orbi/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/mux"
)

const (
	SIGN_UP_TEMPLATE           = "sign-up"
	REGISTRATION_SESSION_STORE = "registration-session"
	REGISTRATION_STORE_SECRET  = "secret"
)

type AuthRouter struct {
	authenticate authenticate
	router       *mux.Router
	callback     func(context.Context, string) string
	webAuthn     *webauthn.WebAuthn
	store        *storage.SessionStore
}

func NewAuthRouter(authenticate authenticate, callback func(context.Context, string) string) *AuthRouter {
	store := storage.NewSessionStore()
	res, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Foobar Corp.",     // Display Name for your site
		RPID:          "localhost",        // Generally the domain name for your site
		RPOrigin:      "http://localhost", // The origin URL for WebAuthn requests
	})

	if err != nil {
		panic(err)
	}

	l := &AuthRouter{
		authenticate: authenticate,
		callback:     callback,
		webAuthn:     res,
		store:        store,
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

	l.router.Path("/webauthn/sign-in").Methods("POST").HandlerFunc(l.WebauthnSignInHandler)
	l.router.Path("/webauthn/sign-up").Methods("POST").HandlerFunc(l.WebauthnSignUpHandler)
	l.router.Path("/webauthn/sign-up/complete").Methods("POST").HandlerFunc(l.WebauthnSignUpCompleteHandler)
}

func (l *AuthRouter) WebauthnSignUpCompleteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := l.authenticate.FindUserByUserName("test")

	if err != nil {
		log.Panicln(err)
	}

	response, err := protocol.ParseCredentialCreationResponseBody(r.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	session, ok := l.store.GetSessionByUserID(user.ID.String())

	if !ok {
		fmt.Println("No session found")
		return
	}

	credential, err := l.webAuthn.CreateCredential(user, session, response)

	if err != nil {
		fmt.Println(err)
		return
	}

	utils.JsonResponse(w, "Registration Success", http.StatusOK) // Handle next steps

	fmt.Println(credential)
	fmt.Println("save creds on user")
}

func (l *AuthRouter) WebauthnSignUpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email: r.FormValue("username"),
	}

	result, err := l.authenticate.CreateUser(user)
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	if err != nil {
		log.Panic(err)
	}

	options, sessionData, err := l.webAuthn.BeginRegistration(
		user,
		registerOptions,
	)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Print(result)
	l.store.StoreSession(user.ID.String(), *sessionData)
	utils.JsonResponse(w, options, http.StatusOK)
}

func (l *AuthRouter) WebauthnSignInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")

	user, err := l.authenticate.FindUserByUserName(username)

	if err != nil {
		renderLogin(w, r.FormValue(queryAuthRequestID), err)
	}

	fmt.Println(user)

	if err != nil {
		log.Println(err)
		return
	}

	renderLogin(w, r.FormValue(queryAuthRequestID), nil)
}

type authenticate interface {
	CheckUsernamePassword(username, password, id string) error
	CreateUser(user models.User) (*models.User, error)
	FindUserByUserName(username string) (*models.User, error)
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
