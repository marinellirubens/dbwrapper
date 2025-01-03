package app

import (
	"fmt"
	"net/http"

	"github.com/marinellirubens/dbwrapper/database"
	"github.com/marinellirubens/dbwrapper/internal/utils"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes(mux *http.ServeMux, app *database.App) (http.Handler, error) {
	// TODO: need to create some treatment on the path variable to understand how to do path handling without any framework
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Log.Debug(fmt.Sprintf("Requested ping by %v", utils.ReadUserIP(r)))
		app.Log.Debug(fmt.Sprintf("Headers %s", r.Header))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("I'm alive.\n"))

		if err != nil {
			app.Log.Error(fmt.Sprintf("Error trying to get server. %v", err))
			panic(err)
		}
	}))

	mux.HandleFunc("/databases", app.GetDatabasesRequest)
	mux.HandleFunc("/database", app.ProcessGenericRequest)

	return corsMiddleware(mux), nil
}
