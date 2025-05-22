package controller

import (
	"ai-course/internal/base/controller"
	"ai-course/internal/logger"
	"ai-course/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ClassController 班级控制器
type ClassController struct {
	controller.BaseController
	classService service.ClassService
}

// NewClassController 创建班级控制器
func NewClassController() *ClassController {
	return &ClassController{
		classService: service.NewClassService(),
	}
}

// Add godoc
// @Summary 添加班级
// @Description 添加新班级
// @Tags 班级管理
// @Accept json
// @Produce json
// @Param request body service.AddClassRequest true "班级信息"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /class/add [post]
func (c *ClassController) Add(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.AddClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid add class request",
			zap.Error(err),
		)
		c.ParamError("添加班级参数无效")
		return
	}

	class, err := c.classService.Add(&req)
	if err != nil {
		logger.Logger.Error("Failed to add class",
			zap.Error(err),
		)
		c.ServerError(err.Error())
		return
	}

	c.SuccessWithMessage("添加班级成功", gin.H{
		"class": class,
	})
}

// Edit godoc
// @Summary 编辑班级
// @Description 编辑现有班级
// @Tags 班级管理
// @Accept json
// @Produce json
// @Param request body service.EditClassRequest true "班级信息"
// @Success 200 {object} response.Response "编辑成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /class/{id} [put]
func (c *ClassController) Edit(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.EditClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid edit class request",
			zap.Error(err),
		)
		c.ParamError("编辑班级参数无效")
		return
	}

	class, err := c.classService.Edit(&req)
	if err != nil {
		logger.Logger.Error("Failed to edit class",
			zap.Error(err),
		)
		c.ServerError(err.Error())
		return
	}

	c.SuccessWithMessage("编辑班级成功", gin.H{
		"class": class,
	})
}

// Delete godoc
// @Summary 删除班级
// @Description 删除现有班级
// @Tags 班级管理
// @Produce json
// @Param id path int true "班级ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /class/{id} [delete]
func (c *ClassController) Delete(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("班级ID格式无效")
		return
	}

	if err := c.classService.Delete(uint(id)); err != nil {
		logger.Logger.Error("Failed to delete class",
			zap.Error(err),
			zap.String("id", idStr),
		)
		c.ServerError(err.Error())
		return
	}

	c.Success(nil)
}

// List godoc
// @Summary 获取班级列表
// @Description 分页获取班级列表
// @Tags 班级管理
// @Produce json
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /class/list [get]
func (c *ClassController) List(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.ListClassRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Logger.Warn("Invalid list class request",
			zap.Error(err),
		)
		c.ParamError("获取班级列表参数无效")
		return
	}

	classes, err := c.classService.List(req.Page, req.PageSize)
	if err != nil {
		logger.Logger.Error("Failed to get class list",
			zap.Error(err),
		)
		c.ServerError(err.Error())
		return
	}

	c.Success(gin.H{
		"classes": classes,
	})
}
