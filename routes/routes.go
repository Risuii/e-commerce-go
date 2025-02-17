package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	User "e-commerce/internal/authentication/delivery/presenter/http"
	Library "e-commerce/library"
	Middleware "e-commerce/middlewares"
	UtilsPackage "e-commerce/pkg/utils"
)

type Routes interface {
	Setup()
	GetEngine() *gin.Engine
}

type RoutesImpl struct {
	engine     *gin.Engine
	library    Library.Library
	middleware Middleware.Middleware
	user       User.UserHandler
}

func New(
	engine *gin.Engine,
	library Library.Library,
	middleware Middleware.Middleware,
	user User.UserHandler,
) Routes {
	return &RoutesImpl{
		engine:     engine,
		library:    library,
		middleware: middleware,
		user:       user,
	}
}

func (o *RoutesImpl) Setup() {
	path := "Routes:Setup"
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
	o.SetUserRoute()
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

func (o *RoutesImpl) SetUserRoute() {
	router := o.engine.Group("/api/v1/user")

	router.Use(o.middleware.GenerateTraceID(), o.middleware.Logging())
	router.POST("/register", o.user.Register)
	router.POST("/login", o.user.Login)
}
