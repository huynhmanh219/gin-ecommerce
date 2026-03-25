package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"huynhmanh.com/gin/internal/dto"
	"huynhmanh.com/gin/internal/service"
)

type UserHandler struct {
	userService         service.UserService
	registrationService service.UserRegistrationService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
func NewUserRegistratrionHandler(registrationService service.UserRegistrationService) *UserHandler {
	return &UserHandler{
		registrationService: registrationService,
	}
}

func (h *UserHandler) RegisterWithProfile(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
	}

	user, err := h.registrationService.RegisterWithProfile(
		c.Request.Context(),
		dto.CreateUserRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		},
		req.Bio,
		req.Avatar,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "REGISTRAION_FAILED",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

// CreateUser creates a new user (registration)
// @Summary      Create user
// @Description  Register a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload body     dto.CreateUserRequest true "User registration data"
// @Success      201    {object}  dto.UserResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if fieldErrors := ValidateStruct(req); fieldErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "Validation_error",
			"message": "Invalid data",
			"errors":  fieldErrors,
		})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

// GetUser retrieves a user by ID
// @Summary      Get user by ID
// @Description  Retrieve user details by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      404  {object}  dto.UserResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetUsers retrieves all users with pagination
// @Summary      Get all users
// @Description  Retrieve list of users with pagination support
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit  query  int  false  "Limit results (default: 10)"
// @Param        offset query  int  false  "Offset results (default: 0)"
// @Success      200    {object}  dto.UsersListResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed > 0 {
			offset = parsed
		}
	}

	users, total, err := h.userService.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error: ": err.Error()})
		writeError(c,err);
		return
	}
	var resp []dto.UserResponse
	for _, user := range users {
		resp = append(resp, dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  resp,
		"total": total,
	})
}

// UpdateUser updates user information
// @Summary      Update user
// @Description  Update name and/or email of a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path     int                   true  "User ID"
// @Param        payload body     dto.UpdateUserRequest true  "Update data"
// @Success      200     {object}  dto.UserResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// DeleteUser deletes a user by ID
// @Summary      Delete user
// @Description  Delete a user record permanently
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id     path     int  true  "User ID"
// @Success      200    {object}  dto.SuccessMessageResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      404    {object}  dto.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// Login authenticates a user
// @Summary      User login
// @Description  Authenticate user with email and password, returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body     dto.LoginRequest true  "Login credentials"
// @Success      200     {object}  dto.LoginResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      401     {object}  dto.ErrorResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	resp := dto.LoginResponse{Token: token}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
