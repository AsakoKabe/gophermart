package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"net/http"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/service"
	userService "github.com/AsakoKabe/gophermart/internal/app/service/user"
)

type UserHandler struct {
	userService service.UserService
	tokenAuth   *jwtauth.JWTAuth
}

const tokenKey = "login"

func NewUserHandler(userService service.UserService, tokenAuth *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{userService: userService, tokenAuth: tokenAuth}
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

	err = h.userService.Add(r.Context(), &user)
	if errors.Is(err, userService.ErrLoginAlreadyExist) {
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

	ok, err := h.userService.IsValidUser(r.Context(), &user)
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
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) makeToken(login string) string {
	_, tokenString, _ := h.tokenAuth.Encode(map[string]interface{}{"login": login})
	return tokenString
}

func (h *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userLogin, ok := claims[tokenKey].(string)
	if !ok {
		slog.Error("error to get user login")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := h.userService.GetBalance(r.Context(), userLogin)
	if err != nil {
		slog.Error(
			"error to get balance",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawal, err := h.userService.GetSumWithdrawal(r.Context(), userLogin)
	if err != nil {
		slog.Error(
			"error to get sum withdrawal",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := balanceResponse{Current: balance - withdrawal, Withdrawn: withdrawal}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error to create response get balance", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

}
