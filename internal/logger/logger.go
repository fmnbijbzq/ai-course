package logger

import (
	"ai-course/internal/config"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// CustomTimeEncoder 自定义时间编码器
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// InitLogger 初始化日志
func InitLogger(conf *config.LoggerConfig) {
	// 配置 lumberjack 进行日志切割
	hook := &lumberjack.Logger{
		Filename:   conf.Filename,   // 日志文件路径
		MaxSize:    conf.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     conf.MaxAge,     // 文件最多保存多少天
		Compress:   conf.Compress,   // 是否压缩
	}

	// 设置日志级别
	var level zapcore.Level
	err := level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大写带颜色的日志级别
		EncodeTime:     CustomTimeEncoder,                // 自定义时间格式
		EncodeDuration: zapcore.StringDurationEncoder,    // 更易读的持续时间
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 短路径编码器
	}

	// 创建两个核心 - 一个用于控制台，一个用于文件
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台输出使用彩色编码器，文件输出使用JSON编码器
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(hook), level),
	)

	// 创建 logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// GinZapLogger 返回 gin 的日志中间件
func GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		Logger.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", cost),
		)
	}
}

// GinZapRecovery 返回 gin 的 recovery 中间件
func GinZapRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, true)
				if brokenPipe {
					Logger.Error("Recovery from panic",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					Logger.Error("Recovery from panic",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					Logger.Error("Recovery from panic",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
