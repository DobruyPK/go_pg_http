package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	appDTO "go_pg_http/internal/application/user/dto"
	userService "go_pg_http/internal/application/user/service"
	domainUser "go_pg_http/internal/domain/user"
	httpDTO "go_pg_http/internal/infrastructure/http/dto"
)

type UserHandler struct {
	userService *userService.UserService
}

func NewUserHandler(userService *userService.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, httpDTO.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	defer r.Body.Close()

	var req httpDTO.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, httpDTO.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	output, err := h.userService.Create(r.Context(), appDTO.CreateUserInput{
		Name: req.Name,
	})
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, httpDTO.UserResponse{
		ID:   output.ID,
		Name: output.Name,
	})
}

func (h *UserHandler) GetUserByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, httpDTO.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	req := httpDTO.GetUserByNameRequest{
		Name: r.URL.Query().Get("name"),
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, httpDTO.ErrorResponse{
			Error: "query parameter 'name' is required",
		})
		return
	}

	output, err := h.userService.GetByName(r.Context(), appDTO.GetUserByNameInput{
		Name: req.Name,
	})
	if err != nil {
		log.Printf("GetUserByName  error: %v", err)
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, httpDTO.UserResponse{
		ID:   output.ID,
		Name: output.Name,
	})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, httpDTO.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	output, err := h.userService.List(r.Context(), appDTO.ListUsersInput{})
	if err != nil {
		log.Printf("ListUsers  error: %v", err)
		handleError(w, err)
		return
	}

	users := make([]httpDTO.UserResponse, 0, len(output.Users))
	for _, usr := range output.Users {
		users = append(users, httpDTO.UserResponse{
			ID:   usr.ID,
			Name: usr.Name,
		})
	}

	writeJSON(w, http.StatusOK, httpDTO.UsersResponse{
		Users: users,
		Total: output.Total,
	})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainUser.ErrInvalidUserName):
		writeJSON(w, http.StatusBadRequest, httpDTO.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, domainUser.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, httpDTO.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, domainUser.ErrUserExists):
		writeJSON(w, http.StatusConflict, httpDTO.ErrorResponse{
			Error: err.Error(),
		})
	default:
		writeJSON(w, http.StatusInternalServerError, httpDTO.ErrorResponse{
			Error: "internal server error",
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
