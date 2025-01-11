package graphql

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

func StartServer() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:        Schema,
			RequestString: params.Query,
		})

		if len(result.Errors) > 0 {
			log.Printf("GraphQL errors: %v", result.Errors)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("GraphQL API is running on http://localhost:8081/graphql")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
