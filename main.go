package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
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

type AddUserImpersonator struct {
	UserId         string `json:"user_id"`
	ImpersonatorId string `json:"impersonator_id"`
}

type UserRelation struct {
	ObjectType  string    `json:"object_type"`
	ObjectId    string    `json:"object_id"`
	Relation    string    `json:"relation"`
	SubjectType string    `json:"subject_type"`
	SubjectId   string    `json:"subject_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type UserImpersonatorResp struct {
	Type      string         `json:"type"`
	Id        string         `json:"id"`
	Relations []UserRelation `json:"relations"`
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
				"object_type":"capability",
				"relation":"agent"}
			}
		`)

		topazURL := os.Getenv("TOPAZ_IS_URL")

		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_IS_URL environment variable is not set.")
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
				"object_type":"capability",
				"relation":"agent"}
			}
		`)

		topazURL := os.Getenv("TOPAZ_IS_URL")

		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_IS_URL environment variable is not set.")
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

	r.Get("/check-external", func(w http.ResponseWriter, r *http.Request) {

		jsonData := []byte(`
		{
			"identity_context":{
				"identity":"jerry@the-smiths.com",
				"type":"IDENTITY_TYPE_SUB"
			},
			"policy_context":{
				"decisions":["allowed"],
				"path":"policies.helloexternal"
			},
			"resource_context":{
				"object_id":"member.claims",
				"object_type":"capability",
				"relation":"agent"}
			}
		`)

		topazURL := os.Getenv("TOPAZ_IS_URL")

		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_IS_URL environment variable is not set.")
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

	// This endpoint is a way for the authz engine to call in external function when
	// making a decision. See hello-with-external.rego to understand how the authorizer uses it.
	r.Get("/external", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("external called!!")

		// simulate some processing
		time.Sleep(1 * time.Second)

		w.WriteHeader(http.StatusOK)
		fmt.Println("external done!!")
		return
	})

	r.Get("/valid-agent-end-user-check", func(w http.ResponseWriter, r *http.Request) {
		// Step 1. we parse JWT to get subject - the impersonated
		// and actor  - the Impersonator.

		//  We first check if the Impersonator has access to the Capability.

		// TODO: Add parsing logic in JWT to fetch the payload shown below.
		// For now, assume this payload is automagically constructed.

		// this payload would check if the User in Identity Context can access The resource
		jsonData := []byte(`
		{
			"identity_context":{
				"identity":"beth@the-smiths.com",
				"type":"IDENTITY_TYPE_SUB"
			},
			"policy_context":{
				"decisions":["allowed"],
				"path":"policies.hello"
			},
			"resource_context":{
				"object_id":"member.claims",
				"object_type":"capability",
				"relation":"can_read",
				"impersonated_id" : "homer@the-simpsons.com" 
				} 
			}
		`) // relation can be "agent | can_read | can_write"

		topazURL := os.Getenv("TOPAZ_IS_URL")

		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_IS_URL environment variable is not set.")
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

			// TODO - Now that we know Step 1, passed, we need to programmatically do:
			// - Add the Impersonated into the Users DB.
			// - The relationship between the Impersonated and Impersonator in the relations DB for Users.
			// - We create some cron job to delete that Impersonated user and relationship using Directory APIs after some set timeout.

			// Now assume these are done, we begin step 2 verification
			// Assume we added Homer Simpson as an Impersonated into the mix via Directory APIs
			// Add we added a Impersonator who is Beth Smiths and linked it to Homer.

			impersonatedJsonData := []byte(`
			{
				"identity_context": {
					"type": "IDENTITY_TYPE_SUB",
					"identity": "beth@the-smiths.com"
				},
				"resource_context": {
					"object_type": "user",
					"object_id": "homer@the-simpsons.com",
					"relation": "impersonator"
				},
				"policy_context": {
					"decisions": [
					"allowed"
					],
					"path": "policies.hello"
				}
			}
			`)
			impreq, err := http.NewRequest("POST", topazURL, bytes.NewBuffer(impersonatedJsonData))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			impreq.Header.Set("Content-Type", "application/json")
			// Create a transport that doesn't verify certificates - this is needed because of topaz's invalid certs and this is a POC
			// For the real deal, we will need real certs
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			impresp, err := client.Do(impreq)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer impresp.Body.Close()
			if impresp.StatusCode == http.StatusOK {
				imptresponse := TopazResponse{}
				fmt.Println("POST request successful!")

				decoder := json.NewDecoder(impresp.Body)
				err := decoder.Decode(&imptresponse)
				if err != nil {
					fmt.Println("Error reading response body:", err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				// fmt.Printf("t response: %+v\n", imptresponse)

				if !imptresponse.Decisions[0].Is {
					http.Error(w, "Gandalf: You shall not PPAAAAAASSS!!!!!", http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("it works!!"))
			}
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

	r.Get("/valid-agent-end-user-check-compact", func(w http.ResponseWriter, r *http.Request) {
		// Step 1. we parse JWT to get subject - the impersonated
		// and actor  - the Impersonator.

		//  We first check if the Impersonator has access to the Capability.

		// TODO: Add parsing logic in JWT to fetch the payload shown below.
		// For now, assume this payload is automagically constructed.

		// this payload would check if the User in Identity Context can access The resource
		jsonData := []byte(`
		{
			"identity_context":{
				"identity":"beth@the-smiths.com",
				"type":"IDENTITY_TYPE_SUB"
			},
			"policy_context":{
				"decisions":["allowed"],
				"path":"policies.impersonator"
			},
			"resource_context":{
				"object_id":"member.claims",
				"object_type":"capability",
				"relation":"can_read",
				"impersonated_id" : "homer@the-simpsons.com" 
				} 
			}
		`) // relation can be "agent | can_read | can_write"

		topazURL := os.Getenv("TOPAZ_IS_URL")
		if topazURL == "" {
			// Handle the case where the environment variable is not set
			fmt.Println("TOPAZ_IS_URL environment variable is not set.")
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
		} else {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
	})

	r.Post("/add-user-with-impersonator", func(w http.ResponseWriter, r *http.Request) {
		// We parse user id and impersonator id
		// We then check if user exists with impersonator using GET Object API
		// If exists and if exp is valid, we simply return without creating
		// Otherwise, we add user with impersonator relation using Object and Relation API

		// we parse request body
		var userImpersonator AddUserImpersonator

		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&userImpersonator)
		if err != nil {
			http.Error(w, "Invalid body in POST", http.StatusBadRequest)
			return
		}

		impersonatedUserId := template.HTMLEscapeString(userImpersonator.UserId)
		userUrl := fmt.Sprintf(
			"https://topaz:9393/api/v3/directory/object/user/%s?with_relations=true", impersonatedUserId,
		)
		// Check if user exists with relations
		req, err := http.NewRequest("GET", userUrl, nil)

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
			decoder := json.NewDecoder(resp.Body)
			var userResponse UserImpersonatorResp
			err := decoder.Decode(&userResponse)
			if err != nil {
				http.Error(w, "Unable to decode User Relations Response", http.StatusInternalServerError)
				return
			}

			fmt.Printf("user relation resp: %+v", userResponse)
			// Now look for impersonator and the updatedAt time
			for _, relation := range userResponse.Relations {
				if relation.SubjectId == userImpersonator.ImpersonatorId {
					createdTS := relation.CreatedAt
					now := time.Now().UTC()
					if now.Sub(createdTS) <= 15*time.Minute {
						// this means we have not yet expired. so we dont need to add
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("relation still valid....user added with impersonator relation"))
						return
					}
					break
				}

			}
			// Found but relation expired. we Delete and Post relation again
			// Delete url - curl -X 'DELETE' \
			//   'https://localhost:9393/api/v3/directory/relation?object_type=user&object_id=homer%40the-simpsons.com&relation=impersonator&subject_type=user&subject_id=beth%40the-smiths.com' \
			//   -H 'accept: application/json'

			deleteUrl := fmt.Sprintf(
				"https://topaz:9393/api/v3/directory/relation?object_type=user&object_id=%s&relation=impersonator&subject_type=user&subject_id=%s",
				template.HTMLEscapeString(userImpersonator.UserId), template.HTMLEscapeString(userImpersonator.ImpersonatorId))
			delreq, err := http.NewRequest("DELETE", deleteUrl, nil)

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			delresp, err := client.Do(delreq)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			if delresp.StatusCode == http.StatusOK {
				// Now we add the new relation back
				// curl -X 'POST' \
				//   'https://localhost:9393/api/v3/directory/relation' \
				//   -H 'accept: application/json' \
				//   -H 'Content-Type: application/json' \
				//   -d '{
				//   "relation": {
				//     "object_id": "homer@the-simpsons.com",
				//     "object_type": "user",
				//     "relation": "impersonator",
				//     "subject_id": "beth@the-smiths.com",
				//     "subject_type": "user"
				//   }
				// }'
				postRelationBody := `{"relation":{ "object_id": "` + template.HTMLEscapeString(userImpersonator.UserId) + `","object_type": "user","relation": "impersonator","subject_id": "` + template.HTMLEscapeString(userImpersonator.ImpersonatorId) + `","subject_type": "user"}}`
				postRelationBodyBytes := []byte(postRelationBody)
				postRelationUrl := "https://topaz:9393/api/v3/directory/relation"

				postRelationReq, err := http.NewRequest("POST", postRelationUrl, bytes.NewBuffer(postRelationBodyBytes))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return
				}
				postRelationReq.Header.Set("Content-Type", "application/json")
				postRelationReq.Header.Set("Accept", "application/json")

				postRelationResp, err := client.Do(postRelationReq)
				if err != nil {
					fmt.Println("Error sending request:", err)
					return
				}

				if postRelationResp.StatusCode == http.StatusOK {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("added new relation with user and impersonator"))
					return
				} else {
					http.Error(w, "Oops", http.StatusInternalServerError)
					return
				}

			} else {
				http.Error(w, "Oops", http.StatusInternalServerError)
				return
			}

			return
		} else {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

	})

	http.ListenAndServe(":8888", r)
}
