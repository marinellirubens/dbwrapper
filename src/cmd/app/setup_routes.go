package app

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/marinellirubens/dbwrapper/src/database"
	"github.com/marinellirubens/dbwrapper/src/internal/utils"
)

func basicAuthMiddleware(next http.Handler, app *database.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "API key is required", http.StatusUnauthorized)
			return
		}

		dbID := r.Header.Get("database")
		if dbID == "" {
			http.Error(w, "Database must be informed", http.StatusUnauthorized)
			return
		}

		// Retrieve the list of valid API keys from the configuration
		validAPIKeys := app.Config.APIKeys
		isValidKey := false
		for _, key := range validAPIKeys {
			if apiKey == key.Key {
				isValidKey := slices.Contains(key.AllowedDbs, dbID)

				if isValidKey {
					break
				} else {
					http.Error(w, "Database not allowed for this key", http.StatusUnauthorized)
					return
				}
			}
		}

		if !isValidKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		_, ok := app.DbConns[dbID]
		if !ok {
			http.Error(w, "Database not located", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
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

	mux.Handle("GET /databases", http.HandlerFunc(app.GetDatabasesRequest))
	mux.Handle("GET /database", basicAuthMiddleware(http.HandlerFunc(app.GetQueryFromDatabase), app))
	mux.Handle("/database", basicAuthMiddleware(http.HandlerFunc(app.ProcessGenericRequest), app))

	return corsMiddleware(mux), nil
}
