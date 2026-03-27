package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"huynhmanh.com/gin/internal/dto"
	"huynhmanh.com/gin/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler{
	return &ProductHandler{
		productService:productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context){
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	if fieldErrors := ValidateStruct(req);fieldErrors != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"code": "Validation_error",
			"message":"Invalid data",
			"errors": fieldErrors,
		})
	}

	product,err := h.productService.CreateProduct(c.Request.Context(),req)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
	}

	resp := dto.ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
		Stock: product.Stock,
		Category: product.Category,
		ImageURL: product.ImageURL,
		IsActive: product.IsActive,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
	c.JSON(http.StatusCreated,gin.H{"data":resp})
}

func (h *ProductHandler) UpdateProduct (c *gin.Context){
	var req dto.UpdateProductRequest
	strId := c.Param("id")

	id,err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	}

	if err:= c.ShouldBindJSON(&req); err != nil{ 
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}

	if fieldErrors := ValidateStruct(req); fieldErrors != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"code":"validation_error",
			"message":"invalid data",
			"errors":fieldErrors,
		})
		return
	}

	product,err := h.productService.UpdateProduct(c.Request.Context(),uint(id),req)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	}

	resp := dto.ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
		Stock: product.Stock,
		Category: product.Category,
		ImageURL: product.ImageURL,
		IsActive: product.IsActive,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	c.JSON(http.StatusOK,gin.H{"data":resp})
}

func (h *ProductHandler) GetAllProduct(c *gin.Context){
	 strLimit := c.Query("limit")
	 strOffset := c.Query("offset")
	 
	 limit,err := strconv.Atoi(strLimit)
	 if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	 }

	 offset,err := strconv.Atoi(strOffset)
	 if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	 }

	data,total,err := h.productService.GetAllProduct(c.Request.Context(),uint(limit),uint(offset))
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	var resp []dto.ProductResponse
	for _,product := range data{
		resp = append(resp, dto.ProductResponse{
				ID: product.ID,
				Name: product.Name,
				Description: product.Description,
				Price: product.Price,
				Stock: product.Stock,
				Category: product.Category,
				ImageURL: product.ImageURL,
				IsActive: product.IsActive,
				CreatedAt: product.CreatedAt,
				UpdatedAt: product.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK,gin.H{
		"data":resp,
		"total":total,
	})
}

func(h *ProductHandler) GetByID(c *gin.Context){
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}


	data,err := h.productService.GetByID(c.Request.Context(),uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	resp := dto.ProductResponse{
		ID: data.ID,
		Name: data.Name,
		Description: data.Description,
		Price: data.Price,
		Stock: data.Stock,
		Category: data.Category,
		ImageURL: data.ImageURL,
		IsActive: data.IsActive,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	c.JSON(http.StatusOK,gin.H{"data":resp})
}

func (h *ProductHandler) GetLowStockProducts(c *gin.Context){
	var data []dto.ProductResponse
	products,err := h.productService.GetLowStockProducts(c.Request.Context())
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _,item := range products{
		data = append(data,dto.ProductResponse{
			ID: item.ID,
			Name: item.Name,
			Description: item.Description,
			Price: item.Price,
			Stock: item.Stock,
			Category: item.Category,
			ImageURL: item.ImageURL,
			IsActive: item.IsActive,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK,gin.H{"data":data})
}


func (h *ProductHandler) SearchByName(c *gin.Context){
	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := req["name"]

	data,err := h.productService.SearchByName(c.Request.Context(),name)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ProductResponse{
		ID: data.ID,
		Name: data.Name,
		Description: data.Description,
		Price: data.Price,
		Stock: data.Stock,
		Category: data.Category,
		ImageURL: data.ImageURL,
		IsActive: data.IsActive,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	c.JSON(http.StatusOK,gin.H{"data":resp})
}