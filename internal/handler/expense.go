package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"pim-api-go/internal/middleware"
	"pim-api-go/internal/model"
	"pim-api-go/internal/service"
	"pim-api-go/pkg/cloudinary"
)

type ExpenseHandler struct{ svc *service.ExpenseService }

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler { return &ExpenseHandler{svc} }

func (h *ExpenseHandler) ListPublic(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	result, err := h.svc.ListPublic(page, c.Query("category"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data pengeluaran"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ExpenseHandler) GetReceipt(c *gin.Context) {
	url, err := h.svc.GetReceipt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	file, header, err := c.Request.FormFile("fotoNota")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto nota wajib diupload"})
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !cloudinary.AllowedMIMEs[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Gunakan JPG, PNG, atau PDF."})
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file"})
		return
	}

	tanggalStr := c.PostForm("tanggal")
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	nominal, _ := strconv.ParseInt(c.PostForm("nominal"), 10, 64)
	claims := c.MustGet("admin").(*middleware.Claims)

	input := service.CreateExpenseInput{
		Tanggal:    tanggal,
		Keterangan: strings.TrimSpace(c.PostForm("keterangan")),
		Kategori:   model.ExpenseCategory(c.PostForm("kategori")),
		Nominal:    nominal,
		FileData:   data,
		MimeType:   mimeType,
		CreatedBy:  claims.Username,
	}

	uploadFn := func(d []byte, folder string) (string, error) {
		return cloudinary.Upload(context.Background(), d, folder)
	}

	expense, err := h.svc.Create(input, uploadFn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengeluaran berhasil disimpan", "expense": expense})
}

func (h *ExpenseHandler) ListAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	result, err := h.svc.ListAdmin(page, c.Query("category"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data pengeluaran"})
		return
	}
	c.JSON(http.StatusOK, result)
}
