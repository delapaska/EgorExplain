package router

import (
	"database/sql"
	numberRepo "github.com/delapaska/EgorExplain/internal/database/number"
	numberHand "github.com/delapaska/EgorExplain/internal/handlers/number"
	numberServ "github.com/delapaska/EgorExplain/internal/services/number"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	engine *gin.Engine
	db     *sql.DB
}

func NewApiServer(db *sql.DB) *APIServer {
	engine := gin.New()

	numberStore := numberRepo.New(db)
	numberService := numberServ.New(numberStore)
	numberHandler := numberHand.New(numberService)
	numberHandler.RegisterRoutes(engine)
	return &APIServer{
		engine: engine,
		db:     db,
	}
}

func (s *APIServer) Run() {

	s.engine.Run(":8000")
}
