package server

import (
    "net/http"
    "database/sql"

    "tarea2sd/internal/database"
    "github.com/gin-gonic/gin"
)

var db *sql.DB

func initDB() {
    db = database.New()
}

func registerMember(c *gin.Context) {
    nombre := c.PostForm("nombre")
    apellido := c.PostForm("apellido")
    rut := c.PostForm("rut")
    email := c.PostForm("email")
    patente := c.PostForm("patente")
    premium := c.PostForm("premium")

    db.Prepare("INSERT INTO Miembro (nombre, apellido, rut, email, patente, premium) VALUES (?, ?, ?, ?, ?, ?)")

    _, err := db.Exec(nombre, apellido, rut, email, patente, premium)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerSale(c *gin.Context) {
    cliente := c.PostForm("cliente")
    cantidad := c.PostForm("cantidad")
    hora := c.PostForm("hora")
    stock := c.PostForm("stock")
    ubicacion := c.PostForm("ubicacion")

    db.Prepare("INSERT INTO Venta (cliente, cantidad, hora, stock, ubicacion) VALUES (?, ?, ?, ?, ?)")

    _, err := db.Exec(cliente, cantidad, hora, stock, ubicacion)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerStrange(c *gin.Context) {

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ping(c *gin.Context) {
    c.String(http.StatusOK, "pong")
}
