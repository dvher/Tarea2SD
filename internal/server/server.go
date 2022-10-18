package server

import (
    "database/sql"
	"github.com/gin-gonic/gin"
)

func New() (*gin.Engine, *sql.DB) {
    initDB()
    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    r.GET("/ping", ping)
    r.POST("/member", registerMember)
    r.POST("/sale", registerSale)
    r.POST("/strange", registerStrange)

    return r, db
}
