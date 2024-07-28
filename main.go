package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Subject struct {
	Name  string
	Group string
}

type TopazResponse struct {
	Decisions []TopazDecision `json:"decisions"`
}
type TopazDecision struct {
	Decision string `json:"decision"`
	Is       bool   `json:"is"`
}

func main() {
	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/valid-agent", func(w http.ResponseWriter, r *http.Request) {

		jsonData := []byte(`
		{
			"identity_context":{
				"identity":"rick@the-citadel.com",
				"type":"IDENTITY_TYPE_SUB"
			},
			"policy_context":{
				"decisions":["allowed"],
				"path":"policies.hello"
			},
			"resource_context":{
				"object_id":"member.claims",
				"object_type":"file",
				"relation":"agent"}
			}
		`)

		topazURL := os.Getenv("TOPAZ_URL")

		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_URL environment variable is not set.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("post body:", string(jsonData))
		req, err := http.NewRequest("POST", topazURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Create a transport that doesn't verify certificates - this is needed because of topaz's invalid certs and this is a POC
		// For the real deal, we will need real certs
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			tresponse := TopazResponse{}
			fmt.Println("POST request successful!")

			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&tresponse)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			// fmt.Printf("t response: %+v\n", tresponse)

			if !tresponse.Decisions[0].Is {
				http.Error(w, "Gandalf: You shall not PPAAAAAASSS!!!!!", http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("it works!!"))
			return
		} else {
			fmt.Printf("POST request failed. Status: %s\n", resp.Status)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			fmt.Println("error response:" + string(body))
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Gandalf: You shall not PPAAAAAASSS!!!!!"))
			return
		}

	})
	r.Get("/invalid-agent", func(w http.ResponseWriter, r *http.Request) {

		jsonData := []byte(`
		{
			"identity_context":{
				"identity":"jerry@the-smiths.com",
				"type":"IDENTITY_TYPE_SUB"
			},
			"policy_context":{
				"decisions":["allowed"],
				"path":"policies.hello"
			},
			"resource_context":{
				"object_id":"member.claims",
				"object_type":"file",
				"relation":"agent"}
			}
		`)

		topazURL := os.Getenv("TOPAZ_URL")

		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_URL environment variable is not set.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("post body:", string(jsonData))
		req, err := http.NewRequest("POST", topazURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Create a transport that doesn't verify certificates - this is needed because of topaz's invalid certs and this is a POC
		// For the real deal, we will need real certs
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			tresponse := TopazResponse{}
			fmt.Println("POST request successful!")

			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&tresponse)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			fmt.Printf("t response: %+v\n", tresponse)

			if !tresponse.Decisions[0].Is {
				http.Error(w, "Gandalf: You shall not PPAAAAAASSS!!!!!", http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("it works!!"))
			return
		} else {
			fmt.Printf("POST request failed. Status: %s\n", resp.Status)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			fmt.Println("error response:" + string(body))
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Gandalf: You shall not PPAAAAAASSS!!!!!"))
			return
		}

	})

	r.Get("/external", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("external called!!")
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
		fmt.Println("external done!!")
		return
	})

	http.ListenAndServe(":8888", r)
}
