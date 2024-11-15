package api

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Dendyator/AntiBF/internal/core"   //nolint
	"github.com/Dendyator/AntiBF/internal/logger" //nolint
)

var appLogger *logger.Logger

func InitLogger(log *logger.Logger) {
	appLogger = log
}

// AuthRequest структура запроса авторизации
type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

// HandleAuth godoc
// @Summary      Авторизация
// @Description  Обработка попытки авторизации
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body AuthRequest true "Auth data"
// @Success      200 {object} map[string]bool
// @Failure      400 {string} string "Invalid request"
// @Router       /auth [post]
func HandleAuth(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		appLogger.Warnf("Invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, _, err := net.ParseCIDR(request.IP); err != nil {
		appLogger.Warnf("Invalid IP address format: %v", request.IP)
		http.Error(w, "Invalid IP address format", http.StatusBadRequest)
		return
	}

	ok := core.CheckAuthorization(request.Login, request.Password, request.IP)
	response := map[string]bool{"ok": ok}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		appLogger.Errorf("Error encoding response: %v", err)
	}
}

func HandleManageList(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ListType string `json:"listType"`
		Subnet   string `json:"subnet"`
		Add      bool   `json:"add"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		appLogger.Warnf("Invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, _, err := net.ParseCIDR(request.Subnet); err != nil {
		appLogger.Warnf("Invalid subnet format: %v", request.Subnet)
		http.Error(w, "Invalid subnet format", http.StatusBadRequest)
		return
	}

	var success bool
	if request.ListType == "white" {
		success = core.ManageList(request.Subnet, core.Whitelist, request.Add)
	} else if request.ListType == "black" {
		success = core.ManageList(request.Subnet, core.Blacklist, request.Add)
	} else {
		http.Error(w, "Invalid list type", http.StatusBadRequest)
		return
	}

	response := map[string]bool{"success": success}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		appLogger.Errorf("Error encoding response: %v", err)
	}
}

func HandleCheckList(w http.ResponseWriter, r *http.Request) {
	appLogger.Debug("HandleCheckList called")

	var request struct {
		Subnet   string `json:"subnet"`
		ListType string `json:"listType"`
	}

	appLogger.Debug("Decoding request body")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		appLogger.Warnf("Invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	appLogger.Debugf("Decoded request: %#v", request)

	appLogger.Debug("Validating subnet format")
	if _, _, err := net.ParseCIDR(request.Subnet); err != nil {
		appLogger.Warnf("Invalid subnet format: %v", request.Subnet)
		http.Error(w, "Invalid subnet format", http.StatusBadRequest)
		return
	}

	var inList bool
	switch request.ListType {
	case "white":
		appLogger.Debugf("Checking whitelist for subnet: %s", request.Subnet)
		inList = core.WhitelistFunc(request.Subnet)
		appLogger.Debugf("Whitelist check result: %v", inList)
	case "black":
		appLogger.Debugf("Checking blacklist for subnet: %s", request.Subnet)
		inList = core.BlacklistFunc(request.Subnet)
		appLogger.Debugf("Blacklist check result: %v", inList)
	default:
		appLogger.Warnf("Invalid list type: %v", request.ListType)
		http.Error(w, "Invalid list type", http.StatusBadRequest)
		return
	}

	response := map[string]bool{"in_list": inList}

	w.Header().Set("Content-Type", "application/json")
	appLogger.Debug("Encoding response")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		appLogger.Errorf("Error encoding response: %v", err)
	}
	appLogger.Debug("Response sent successfully")
}
