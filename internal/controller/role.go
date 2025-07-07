package controller

import (
	"ai-course/internal/base/controller"
	"ai-course/internal/logger"
	"ai-course/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleController struct {
	controller.BaseController
	roleService service.RoleService
}

func NewRoleController(roleService service.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (c *RoleController) RegisterRoutes(r *gin.Engine) {
	roleGroup := r.Group("/api/role")
	{
		roleGroup.POST("/add", c.Add)
		roleGroup.POST("/edit", c.Edit)
		roleGroup.POST("/delete", c.Delete)
	}
}

func (c *RoleController) Add(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.CreateRoleDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid add role request",
			zap.Error(err),
		)
		c.ParamError("添加角色参数无效")
		return
	}
}

func (c *RoleController) Edit(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.UpdateRoleDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid edit role request",
			zap.Error(err),
		)
		c.ParamError("编辑角色参数无效")
		return
	}
}

func (c *RoleController) Delete(ctx *gin.Context) {
	c.InitHandler(ctx)
	var roleId string
	if err := ctx.ShouldBindQuery(&roleId); err != nil {
		logger.Logger.Warn("Invalid delete role request",
			zap.Error(err),
		)
		c.ParamError("删除角色参数无效")
		return
	}

	err := c.roleService.Delete(ctx, roleId)
	if err != nil {
		logger.Logger.Error("Failed to delete role",
			zap.Error(err),
		)
		c.ServerError(err.Error())
		return
	}
	c.Success("删除角色成功")
}

func (c *RoleController) List(ctx *gin.Context) {
	c.InitHandler(ctx)
	var roleId string
	if err := ctx.ShouldBindQuery(&roleId); err != nil {
		logger.Logger.Warn("Invalid list role request",
			zap.Error(err),
		)
		c.ParamError("获取角色列表参数无效")
		return
	}

	roles, err := c.roleService.GetById(ctx, roleId)
	if err != nil {
		logger.Logger.Error("Failed to get role list",
			zap.Error(err),
		)
		c.ServerError(err.Error())
		return
	}
	c.SuccessWithMessage("获取角色列表成功", roles)
}
