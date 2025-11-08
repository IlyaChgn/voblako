package delivery

import (
	authdel "github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/rest"
	filedel "github.com/IlyaChgn/voblako/internal/pkg/file/delivery/rest"

	"github.com/gorilla/mux"
)

func NewRouter(
	authHandler *authdel.AuthHandler,
	fileHandler *filedel.FileHandler,
	loginRequiredMiddleware mux.MiddlewareFunc,
) *mux.Router {
	router := mux.NewRouter()
	rootRouter := router.PathPrefix("/api").Subrouter()

	subrouterAuth := rootRouter.PathPrefix("/auth").Subrouter()
	subrouterAuth.HandleFunc("/signup", authHandler.Signup).Methods("POST")
	subrouterAuth.HandleFunc("/login", authHandler.Login).Methods("POST")
	subrouterAuth.HandleFunc("/check", authHandler.CheckAuth).Methods("GET")

	subrouterLogout := subrouterAuth.PathPrefix("/logout").Subrouter()
	subrouterLogout.Use(loginRequiredMiddleware)
	subrouterLogout.HandleFunc("", authHandler.Logout).Methods("POST")

	subrouterFiles := rootRouter.PathPrefix("/files").Subrouter()
	subrouterFiles.Use(loginRequiredMiddleware)
	subrouterFiles.HandleFunc("", fileHandler.UploadFile).Methods("POST")
	subrouterFiles.HandleFunc("/list", fileHandler.GetFilesList).Methods("POST")
	subrouterFiles.HandleFunc("/{id}", fileHandler.GetFile).Methods("GET")
	subrouterFiles.HandleFunc("/{id}/meta", fileHandler.GetMetadata).Methods("GET")
	subrouterFiles.HandleFunc("/{id}", fileHandler.UpdateFile).Methods("POST")
	subrouterFiles.HandleFunc("/{id}/name", fileHandler.UpdateFilename).Methods("POST")
	subrouterFiles.HandleFunc("/{id}", fileHandler.DeleteFile).Methods("DELETE")

	return router
}
