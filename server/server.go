package server

import (
	"go-cars/server/controllers"
	"go-cars/storage"

	docs "go-cars/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.uber.org/zap"
)

type Server struct {
	DB     *storage.Storage
	Engine *gin.Engine
	Logger *zap.Logger
}

func New(db *storage.Storage, logger *zap.Logger) *Server {
	return &Server{DB: db, Engine: gin.Default(), Logger: logger}
}

func (s *Server) setupRoutes() {
	s.Engine.GET("/get", controllers.GetCars(s.DB, s.Logger))
	s.Engine.DELETE("/delete", controllers.DeleteCarByRegNum(s.DB, s.Logger))
	s.Engine.PATCH("/update", controllers.UpdateCarByRegNum(s.DB, s.Logger))
	s.Engine.POST("/add", controllers.AddCars(s.DB, s.Logger))

	docs.SwaggerInfo.BasePath = ""
	s.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Run(addr ...string) {
	s.setupRoutes()
	s.Engine.Run(addr...)
}
