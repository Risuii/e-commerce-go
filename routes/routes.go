package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	User "e-commerce/internal/authentication/delivery/presenter/http"
	Product "e-commerce/internal/product/delivery/presenter/http"
	Store "e-commerce/internal/store/delivery/presenter/http"
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
	store      Store.StoreHandler
	product    Product.ProductHandler
}

func New(
	engine *gin.Engine,
	library Library.Library,
	middleware Middleware.Middleware,
	user User.UserHandler,
	store Store.StoreHandler,
	product Product.ProductHandler,
) Routes {
	return &RoutesImpl{
		engine:     engine,
		library:    library,
		middleware: middleware,
		user:       user,
		store:      store,
		product:    product,
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
	o.SetStoreRoute()
	o.SetProductRoute()
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

	logout := o.engine.Group("/api/v1/user")
	logout.Use(o.middleware.GenerateTraceID(), o.middleware.ValidateToken(), o.middleware.Logging())
	logout.POST("/logout", o.user.Logout)
}

func (o *RoutesImpl) SetStoreRoute() {
	router := o.engine.Group("/api/v1/store")

	router.Use(o.middleware.GenerateTraceID(), o.middleware.ValidateToken(), o.middleware.Logging())
	router.POST("/create", o.store.CreateStore)
	router.POST("/update", o.store.UpdateStore)
	router.GET("/", o.store.GetStore)
	router.POST("/change-status", o.store.ChangeStoreStatus)
}

func (o *RoutesImpl) SetProductRoute() {
	router := o.engine.Group("/api/v1/product")

	router.Use(o.middleware.GenerateTraceID(), o.middleware.ValidateToken(), o.middleware.Logging())
	router.POST("/create", o.product.CreateProduct)
}
