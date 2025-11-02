package delivery

import (
	authdel "github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/rest"

	"github.com/gorilla/mux"
)

func NewRouter(authHandler *authdel.AuthHandler) *mux.Router {
	router := mux.NewRouter()
	rootRouter := router.PathPrefix("/api").Subrouter()

	subrouter := rootRouter.PathPrefix("/auth").Subrouter()

	subrouter.HandleFunc("/signup", authHandler.Signup).Methods("POST")
	subrouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	subrouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	subrouter.HandleFunc("/check", authHandler.CheckAuth).Methods("GET")

	return router
}
