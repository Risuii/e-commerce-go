package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	Constants "e-commerce/constants"

	Library "e-commerce/library"
	UtilsPackage "e-commerce/pkg/utils"
)

type Routes interface {
	Setup()
	GetEngine() *gin.Engine
}

type RoutesImpl struct {
	engine  *gin.Engine
	library Library.Library
}

func New(
	engine *gin.Engine,
	library Library.Library,
) Routes {
	return &RoutesImpl{
		engine:  engine,
		library: library,
	}
}

func (o *RoutesImpl) Setup() {
	path := Constants.Routes
	defer UtilsPackage.CatchPanic(path, o.library)
	// SETUP CORS
	o.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //http or https
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	// EMBED ROUTES
	o.SetIndexRoute()
}

func (o *RoutesImpl) GetEngine() *gin.Engine {
	return o.engine
}

func (o *RoutesImpl) SetIndexRoute() {
	o.engine.GET("/hello", func(c *gin.Context) {

		response := gin.H{
			"status":  true,
			"message": "The request is processed successfully",
			"data":    "This is a message from API",
		}

		c.JSON(http.StatusOK, response)
	})
}
