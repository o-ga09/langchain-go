package api

import (
	"context"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/o-ga09/langchain-go/chat"
	"github.com/o-ga09/langchain-go/pkg/config"
	"github.com/o-ga09/langchain-go/pkg/logger"
	"github.com/o-ga09/langchain-go/pkg/middleware"
)

type Server struct {
	port   string
	logger *slog.Logger
	engine *gin.Engine
}

func New() *Server {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Faild to load environment variables:", err)
	}
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Server{
		port:   cfg.Port,
		logger: logger.New(),
		engine: gin.New(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	// リクエストデータを取得
	httpLogger := middleware.RequestLogger(s.logger)
	// CORS設定
	cors := middleware.CORS()
	// リクエストタイムアウト設定
	withCtx := middleware.WithTimeout()
	// リクエストID付与
	withReqId := middleware.AddID()

	// ミドルウェア設定
	s.engine.Use(withReqId)
	s.engine.Use(withCtx)
	s.engine.Use(cors)
	s.engine.Use(httpLogger)

	// ヘルスチェック
	v1 := s.engine.Group("/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			s.logger.InfoContext(ctx.Request.Context(), "health check")
			ctx.JSON(200, gin.H{"message": "ok"})
		})
	}

	// LLMに質問を投げる
	{
		v1.POST("/question", func(ctx *gin.Context) {
			reqBody := struct {
				Question string `json:"question"`
			}{}
			if err := ctx.BindJSON(&reqBody); err != nil {
				ctx.JSON(400, gin.H{"message": "invalid request"})
				return
			}

			QaData := chat.RequestQAData{
				Question: reqBody.Question,
			}
			res, err := chat.RunWithRAG(ctx.Request.Context(), []chat.RequestQAData{QaData})
			if err != nil {
				ctx.JSON(500, gin.H{"message": "failed to get response"})
				return
			}

			s.logger.InfoContext(ctx.Request.Context(), "question completed")
			ctx.JSON(200, res)
		})
	}
	// ベクトルDBにデータを追加
	{
		v1.POST("/add/document", func(ctx *gin.Context) {
			reqBody := struct {
				PageContent string `json:"page_content"`
			}{}
			if err := ctx.BindJSON(&reqBody); err != nil {
				ctx.JSON(400, gin.H{"message": "invalid request"})
				return
			}
			data := []*chat.RequestDocumentData{
				{
					PageContent: reqBody.PageContent,
				},
			}
			if err := chat.AddDocument(ctx.Request.Context(), data); err != nil {
				ctx.JSON(500, gin.H{"message": "failed to add document"})
				return
			}
			s.logger.InfoContext(ctx.Request.Context(), "add document completed")
			ctx.JSON(200, gin.H{"message": "completed"})
		})
	}

	if err := s.engine.Run(":" + s.port); err != nil {
		return err
	}

	return nil
}
