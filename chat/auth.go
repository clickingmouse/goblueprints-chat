package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

// authHandler method
// todo: make it cope with empty value & empty cookie case
// if cookie, err := r.Cookie("auth"); err ==
// http.ErrNoCookie || cookie.Value == ""
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		//not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		//some other erro
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth function
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

//loginHandler handles the third-party login process
//format : /auth/{action}/{provider}
//TODO - code will panic if segs[2] & segs[3] does not exists

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("received login handler - %s ", r.URL.Path)
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		//	log.Println("TODO handle login for", provider)
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get Provider %s:%s", provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			//http.Error(w, fmt.Sprintf("Error when trying to get user from %s:%s", provider, err), http.StatusInternalServerError)
			log.Fatalln("Error when trying to get user from", provider, "-", err)
		}
		chatUser := &chatUser{User: user}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("Error when trying to get AvatarURL", "-", err)
		}
		//return
		//}

		// authCookieValue := objx.New(map[string]interface{}{
		// 	"name":       user.Name(),
		// 	"avatar_url": user.AvatarURL(),
		// 	"email":      user.Email(),
		// }).MustBase64()

		//m := md5.New()
		//io.WriteString(m, strings.ToLower(user.Email()))
		//userId := fmt.Sprintf("%x", m.Sum(nil))
		authCookieValue := objx.New(map[string]interface{}{
			//"userid":      userId,
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,

			//"avatar_url": user.AvatarURL(),
			//"email":       user.Email(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
