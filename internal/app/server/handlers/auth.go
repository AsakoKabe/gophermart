package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"net/http"
	"time"
)

type UserHandler struct {
	userStorage storage.UserStorage
	tokenAuth   *jwtauth.JWTAuth
}

const tokenExpired = 7 * 24 * time.Hour

func NewUserHandler(userStorage storage.UserStorage, tokenAuth *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{userStorage: userStorage, tokenAuth: tokenAuth}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("error to read body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userStorage.CreateUser(r.Context(), &user)
	if errors.Is(err, models.LoginAlreadyExist) {
		slog.Error("login already existed", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		slog.Error("err to register user", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := h.makeToken(user.Login)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(tokenExpired),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("error to read body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := h.userStorage.IsUserValid(r.Context(), &user)
	if err != nil {
		slog.Error("err to check exist user", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		slog.Error("user not valid", slog.String("login", user.Login))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := h.makeToken(user.Login)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(tokenExpired),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) makeToken(login string) string {
	_, tokenString, _ := h.tokenAuth.Encode(map[string]interface{}{"login": login})
	return tokenString
}
