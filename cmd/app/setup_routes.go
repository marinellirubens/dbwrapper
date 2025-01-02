package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	pg "github.com/marinellirubens/dbwrapper/database"
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

func SetupRoutes(mux *http.ServeMux, app *pg.App) (http.Handler, error) {
	// TODO: need to create some treatment on the path variable to understand how to do path handling without any framework
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Log.Debug(fmt.Sprintf("Requested ping by %v", ReadUserIP(r)))
		app.Log.Debug(fmt.Sprintf("Headers %s", r.Header))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("I'm alive.\n"))

		if err != nil {
			app.Log.Error(fmt.Sprintf("Error trying to get server. %v", err))
			panic(err)
		}
	}))

	mux.Handle("/databases", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch method := r.Method; method {
		case http.MethodGet:
			app.Log.Debug(fmt.Sprintf("Requested ping by %v", ReadUserIP(r)))
			app.Log.Debug(fmt.Sprintf("Headers %s", r.Header))

			w.WriteHeader(http.StatusOK)
			jsonResponse, _ := json.Marshal(app.DbHandlers)
			_, err := w.Write(jsonResponse)
			if err != nil {
				app.Log.Error(fmt.Sprintf("Error trying to get server. %v", err))
				panic(err)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}
	}))

	mux.HandleFunc("/database", app.ProcessGenericRequest)

	// process base request for postgresl
	mux.HandleFunc("/pg", app.ProcessPostgresRequest)

	// process requests with path arguments
	mux.Handle(
		"/pg/",
		http.StripPrefix("/pg/", http.HandlerFunc(app.ProcessPostgresRequestHandlePath)),
	)

	mux.HandleFunc("/oracle", app.ProcessOracleRequest)
	mux.HandleFunc("/mongodb", app.ProcessMongoRequest)
	return corsMiddleware(mux), nil
}
