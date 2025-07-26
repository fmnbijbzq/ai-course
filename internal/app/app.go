package app

import (
	"ai-course/internal/config"
	"ai-course/internal/controller"
	"ai-course/internal/repository"
	"ai-course/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Application 应用程序结构
type Application struct {
	Engine            *gin.Engine
	Config            *config.Config
	DB                repository.DB
	UserService       service.UserService
	ClassService      service.ClassService
	AssignmentService service.AssignmentService
	QuestionService   service.QuestionService
	SubmissionService service.SubmissionService
	GradingService    service.GradingService
	AttachmentService service.AttachmentService
}

// NewApplication 创建应用程序实例
func NewApplication(
	engine *gin.Engine,
	cfg *config.Config,
	db repository.DB,
	userService service.UserService,
	classService service.ClassService,
	assignmentService service.AssignmentService,
	questionService service.QuestionService,
	submissionService service.SubmissionService,
	gradingService service.GradingService,
	attachmentService service.AttachmentService,
) *Application {
	return &Application{
		Engine:            engine,
		Config:            cfg,
		DB:                db,
		UserService:       userService,
		ClassService:      classService,
		AssignmentService: assignmentService,
		QuestionService:   questionService,
		SubmissionService: submissionService,
		GradingService:    gradingService,
		AttachmentService: attachmentService,
	}
}

// NewGinEngine 创建 Gin 引擎
func NewGinEngine(cfg *config.Config) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 gin 引擎
	engine := gin.New()
	return engine
}

// RegisterRoutes 注册路由
func (app *Application) RegisterRoutes() {
	router := controller.NewRouter(app.Engine, app.UserService, app.ClassService, app.AssignmentService, app.QuestionService, app.SubmissionService, app.GradingService, app.AttachmentService)
	router.RegisterRoutes()
}

// Run 运行应用程序
func (app *Application) Run() error {
	// 运行数据库迁移
	if err := repository.AutoMigrate(app.DB); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	// 注册路由
	app.RegisterRoutes()

	// 启动服务器
	addr := fmt.Sprintf(":%d", app.Config.Server.Port)
	return app.Engine.Run(addr)
}
