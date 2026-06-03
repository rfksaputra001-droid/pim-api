package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pim-api-go/internal/middleware"
	"pim-api-go/internal/model"
	"pim-api-go/internal/service"
)

type UserHandler struct{ svc *service.UserService }

func NewUserHandler(svc *service.UserService) *UserHandler { return &UserHandler{svc} }

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data pengguna"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var body struct {
		Username string          `json:"username"`
		Password string          `json:"password"`
		Role     model.AdminRole `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	user, err := h.svc.Create(body.Username, body.Password, body.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Akun berhasil dibuat", "user": user})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	var body struct {
		Role model.AdminRole `json:"role"`
	}
	c.ShouldBindJSON(&body)
	claims := c.MustGet("admin").(*middleware.Claims)

	if err := h.svc.UpdateRole(c.Param("id"), claims.ID, body.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role berhasil diubah"})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var body struct {
		Password string `json:"password"`
	}
	c.ShouldBindJSON(&body)

	if err := h.svc.UpdatePassword(c.Param("id"), body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	claims := c.MustGet("admin").(*middleware.Claims)

	if err := h.svc.Delete(c.Param("id"), claims.ID); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Pengguna tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Akun berhasil dihapus"})
}
