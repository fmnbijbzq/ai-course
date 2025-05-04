package handler

import (
	"ai-course/internal/logger"
	"ai-course/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ClassHandler struct {
	classService service.ClassService
}

func NewClassHandler() *ClassHandler {
	return &ClassHandler{
		classService: service.NewClassService(),
	}
}

// Add godoc
// @Summary 添加班级
// @Description 添加一个新班级
// @Tags 班级管理
// @Accept json
// @Produce json
// @Param request body service.ClassAddRequest true "班级信息"
// @Success 200 {object} map[string]interface{} "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /class/add [post]
func (h *ClassHandler) Add(c *gin.Context) {
	var req service.ClassAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid add class request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	resp, err := h.classService.Add(&req)
	if err != nil {
		logger.Logger.Error("Failed to add class",
			zap.Error(err),
			zap.String("class_name", req.ClassName),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Logger.Info("Class added successfully",
		zap.String("class_name", req.ClassName),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Class added successfully",
		"data":    resp,
	})
}

// Edit godoc
// @Summary 编辑班级
// @Description 编辑班级信息
// @Tags 班级管理
// @Accept json
// @Produce json
// @Param id path int true "班级ID"
// @Param request body service.ClassEditRequest true "班级信息"
// @Success 200 {object} map[string]interface{} "编辑成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "班级不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /class/{id} [put]
func (h *ClassHandler) Edit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class ID",
		})
		return
	}

	var req service.ClassEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid edit class request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	resp, err := h.classService.Edit(uint(id), &req)
	if err != nil {
		logger.Logger.Error("Failed to edit class",
			zap.Error(err),
			zap.Uint64("class_id", id),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Logger.Info("Class edited successfully",
		zap.Uint64("class_id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Class edited successfully",
		"data":    resp,
	})
}

// Delete godoc
// @Summary 删除班级
// @Description 删除指定班级
// @Tags 班级管理
// @Param id path int true "班级ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "班级不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /class/{id} [delete]
func (h *ClassHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class ID",
		})
		return
	}

	if err := h.classService.Delete(uint(id)); err != nil {
		logger.Logger.Error("Failed to delete class",
			zap.Error(err),
			zap.Uint64("class_id", id),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Logger.Info("Class deleted successfully",
		zap.Uint64("class_id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Class deleted successfully",
	})
}

// List godoc
// @Summary 获取班级列表
// @Description 获取分页的班级列表
// @Tags 班级管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "成功获取班级列表"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /class/list [get]
func (h *ClassHandler) List(c *gin.Context) {
	// 获取分页参数，默认第1页，每页10条
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	resp, err := h.classService.List(page, pageSize)
	if err != nil {
		logger.Logger.Error("Failed to get class list",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    resp,
	})
}
