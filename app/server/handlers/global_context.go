package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"plandex-server/db"

	shared "plandex-shared"
)

type GlobalContextRequest struct {
	Content string `json:"content"`
}

type GlobalContextResponse struct {
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func GetGlobalContextHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for GetGlobalContextHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	orgUserConfig, err := db.GetOrgUserConfig(auth.User.Id, auth.OrgId)
	if err != nil {
		log.Printf("Error getting org user config: %v\n", err)
		http.Error(w, "Error getting org user config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if orgUserConfig == nil || orgUserConfig.GlobalContext == "" {
		http.Error(w, "No global context set", http.StatusNotFound)
		return
	}

	resp := GlobalContextResponse{
		Content: orgUserConfig.GlobalContext,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		http.Error(w, "Error marshalling response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func UpdateGlobalContextHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for UpdateGlobalContextHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req GlobalContextRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		http.Error(w, "Error unmarshalling request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get existing config or create new one
	orgUserConfig, err := db.GetOrgUserConfig(auth.User.Id, auth.OrgId)
	if err != nil {
		log.Printf("Error getting org user config: %v\n", err)
		http.Error(w, "Error getting org user config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if orgUserConfig == nil {
		orgUserConfig = &shared.OrgUserConfig{}
	}

	orgUserConfig.GlobalContext = req.Content

	err = db.UpdateOrgUserConfig(auth.User.Id, auth.OrgId, orgUserConfig)
	if err != nil {
		log.Printf("Error updating org user config: %v\n", err)
		http.Error(w, "Error updating org user config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully updated global context")
	w.WriteHeader(http.StatusOK)
}

func DeleteGlobalContextHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for DeleteGlobalContextHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	// Get existing config
	orgUserConfig, err := db.GetOrgUserConfig(auth.User.Id, auth.OrgId)
	if err != nil {
		log.Printf("Error getting org user config: %v\n", err)
		http.Error(w, "Error getting org user config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if orgUserConfig == nil {
		// No config exists, already deleted
		w.WriteHeader(http.StatusNoContent)
		return
	}

	orgUserConfig.GlobalContext = ""

	err = db.UpdateOrgUserConfig(auth.User.Id, auth.OrgId, orgUserConfig)
	if err != nil {
		log.Printf("Error updating org user config: %v\n", err)
		http.Error(w, "Error updating org user config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully deleted global context")
	w.WriteHeader(http.StatusNoContent)
}