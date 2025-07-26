package controller

import (
	"ai-course/internal/base/controller"
	basemiddleware "ai-course/internal/base/middleware"
	"ai-course/internal/logger"
	"ai-course/internal/middleware"
	"ai-course/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	engine            *gin.Engine
	userService       service.UserService
	classService      service.ClassService
	assignmentService service.AssignmentService
	questionService   service.QuestionService
	submissionService service.SubmissionService
	gradingService    service.GradingService
	attachmentService service.AttachmentService
	baseCtrl          *controller.BaseController
}

// NewRouter 创建路由管理器
func NewRouter(engine *gin.Engine, userService service.UserService, classService service.ClassService, assignmentService service.AssignmentService, questionService service.QuestionService, submissionService service.SubmissionService, gradingService service.GradingService, attachmentService service.AttachmentService) *Router {
	return &Router{
		engine:            engine,
		userService:       userService,
		classService:      classService,
		assignmentService: assignmentService,
		questionService:   questionService,
		submissionService: submissionService,
		gradingService:    gradingService,
		attachmentService: attachmentService,
		baseCtrl:          &controller.BaseController{},
	}
}

// RegisterRoutes 注册所有路由
func (r *Router) RegisterRoutes() {
	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	// 注册全局中间件
	r.engine.Use(cors.New(corsConfig))
	r.engine.Use(logger.GinZapLogger(), logger.GinZapRecovery(true))
	r.engine.Use(basemiddleware.APILogger()) // 添加API日志中间件

	// 创建中间件实例
	roleMiddleware := middleware.NewRoleMiddleware(r.userService)
	
	// HealthCheck godoc
	// @Summary 健康检查
	// @Description 检查服务是否正常运行
	// @Tags 系统状态
	// @Produce json
	// @Success 200 {object} response.Response "服务正常运行"
	// @Router /health [get]
	r.engine.GET("/health", func(c *gin.Context) {
		r.baseCtrl.InitHandler(c)
		r.baseCtrl.Success(gin.H{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// ErrorTest godoc
	// @Summary 错误测试
	// @Description 测试错误处理
	// @Tags 系统状态
	// @Produce json
	// @Success 500 {object} response.Response "测试错误"
	// @Router /error [get]
	r.engine.GET("/error", func(c *gin.Context) {
		r.baseCtrl.InitHandler(c)
		r.baseCtrl.ServerError("This is a test error")
	})

	// 用户路由组（无需认证）
	userController := NewUserController(r.userService)
	userGroup := r.engine.Group("/api/user")
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
	}

	// 需要认证的API路由组
	apiGroup := r.engine.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware()) // 所有API都需要认证
	{
		// 班级路由组（需要管理员或教师权限）
		classController := NewClassController(r.classService)
		classGroup := apiGroup.Group("/class")
		classGroup.Use(roleMiddleware.RequireRole("admin", "teacher"))
		{
			classGroup.POST("/add", classController.Add)
			classGroup.PUT("/:id", classController.Edit)
			classGroup.DELETE("/:id", classController.Delete)
			classGroup.GET("/list", classController.List)
		}

		// 作业路由组
		assignmentController := NewAssignmentController(r.assignmentService)
		assignmentGroup := apiGroup.Group("/assignment")
		{
			// 教师专用路由
			teacherAssignmentGroup := assignmentGroup.Group("")
			teacherAssignmentGroup.Use(roleMiddleware.RequireTeacher())
			{
				teacherAssignmentGroup.POST("", assignmentController.Create)                        // 创建作业
				teacherAssignmentGroup.PUT("/:id", assignmentController.Update)                     // 更新作业
				teacherAssignmentGroup.DELETE("/:id", assignmentController.Delete)                  // 删除作业
				teacherAssignmentGroup.POST("/:id/publish", assignmentController.Publish)          // 发布作业
				teacherAssignmentGroup.POST("/:id/unpublish", assignmentController.Unpublish)      // 取消发布作业
				teacherAssignmentGroup.GET("/:id/statistics", assignmentController.Statistics)     // 获取作业统计
			}

			// 教师和学生都可访问的路由
			assignmentGroup.GET("/:id", assignmentController.Detail)     // 获取作业详情
			assignmentGroup.GET("/list", assignmentController.List)      // 获取作业列表
		}

		// 题目路由组（独立路由组，避免冲突）
		questionController := NewQuestionController(r.questionService, r.assignmentService)
		questionGroup := apiGroup.Group("/question")
		questionGroup.Use(roleMiddleware.RequireTeacher()) // 只有教师可以管理题目
		{
			questionGroup.POST("/assignment/:assignment_id", questionController.Create)                             // 创建题目
			questionGroup.PUT("/assignment/:assignment_id/:question_id", questionController.Update)                // 更新题目
			questionGroup.DELETE("/assignment/:assignment_id/:question_id", questionController.Delete)             // 删除题目
			questionGroup.GET("/assignment/:assignment_id", questionController.List)                               // 获取题目列表
		}

		// 提交路由组（学生专用）
		submissionController := NewSubmissionController(r.submissionService)
		submissionGroup := apiGroup.Group("/submission")
		submissionGroup.Use(roleMiddleware.RequireStudent()) // 只有学生可以提交作业
		{
			submissionGroup.POST("/draft", submissionController.SaveDraft)                                    // 保存草稿
			submissionGroup.POST("/submit", submissionController.Submit)                                      // 提交作业
			submissionGroup.GET("/student/assignments", submissionController.GetStudentAssignments)          // 获取学生作业列表
			submissionGroup.GET("/assignment/:assignment_id", submissionController.GetAssignmentForStudent)  // 获取学生特定作业详情
			submissionGroup.GET("/:id", submissionController.GetSubmissionDetail)                            // 获取提交详情
		}

		// 批改路由组（教师专用）
		gradingController := NewGradingController(r.submissionService, r.gradingService)
		gradingGroup := apiGroup.Group("/grading")
		gradingGroup.Use(roleMiddleware.RequireTeacher()) // 只有教师可以批改
		{
			gradingGroup.GET("/assignment/:assignment_id/submissions", gradingController.GetSubmissionsForGrading)    // 获取待批改提交列表
			gradingGroup.GET("/submission/:submission_id", gradingController.GetGradingDetail)                         // 获取批改详情
			gradingGroup.POST("/submission/:submission_id", gradingController.GradeSubmission)                        // 批改提交
			gradingGroup.POST("/batch", gradingController.BatchGrade)                                                  // 批量批改
			gradingGroup.POST("/assignment/:assignment_id/publish", gradingController.PublishGrades)                  // 发布成绩
			gradingGroup.GET("/assignment/:assignment_id/progress", gradingController.GetGradingProgress)             // 获取批改进度
		}

		// 附件路由组
		attachmentController := NewAttachmentController(r.attachmentService)
		attachmentGroup := apiGroup.Group("/attachment")
		{
			// 上传附件（教师专用）
			attachmentGroup.POST("/assignment/:assignment_id", roleMiddleware.RequireTeacher(), attachmentController.Upload)
			// 删除附件（教师专用）
			attachmentGroup.DELETE("/:id", roleMiddleware.RequireTeacher(), attachmentController.Delete)
			
			// 查看和下载附件（教师和学生都可以）
			attachmentGroup.GET("/assignment/:assignment_id", attachmentController.GetByAssignment)  // 获取作业附件列表
			attachmentGroup.GET("/:id/download", attachmentController.Download)                      // 下载附件
		}
	}
}
