package api

import (
	"fmt"
	"net/http"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func success(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": obj})
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	roles := []string{"merchant", "admin"}

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker), roleBasedMiddleware(roles))

	authRoutes.POST("/v1/product", server.createProduct)
	router.GET("/v1/product/:id", server.getProduct)
	router.POST("/v1/images/:pid", server.uploadImage)
	router.GET("/v1/products", server.listProduct)
	router.POST("/v1/sign-up", server.signUp)
	router.POST("/v1/sign-in", server.login)

	server.router = router
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		store:      store,
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
