package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "db"
	port     = "3306"
	user     = "root"
	password = "password"
	dbname   = "testdb"
)

type Pessoa struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Endereco string `json:"endereco"`
}

func main() {
	// Abrir conexão com o banco de dados
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Criar tabela pessoa se não existir
	createTable := `
		CREATE TABLE IF NOT EXISTS pessoa (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nome VARCHAR(50) NOT NULL,
			cpf VARCHAR(11) NOT NULL,
			endereco VARCHAR(100) NOT NULL
		);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	// Inicializar o router Gin
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	// Definir rotas para o CRUD de Pessoa
	r.GET("/pessoa", listPessoaHandler(db))
	r.GET("/pessoa/:id", getPessoaHandler(db))
	r.POST("/pessoa", createPessoaHandler(db))
	r.PUT("/pessoa/:id", updatePessoaHandler(db))
	r.DELETE("/pessoa/:id", deletePessoaHandler(db))

	// Iniciar o servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func listPessoaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, nome, cpf, endereco FROM pessoa")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		pessoas := make([]Pessoa, 0)
		for rows.Next() {
			var p Pessoa
			if err := rows.Scan(&p.ID, &p.Nome, &p.CPF, &p.Endereco); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			pessoas = append(pessoas, p)
		}

		c.JSON(http.StatusOK, pessoas)
	}
}

func getPessoaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var p Pessoa
		row := db.QueryRow("SELECT id, nome, cpf, endereco FROM pessoa WHERE id=?", id)
		if err := row.Scan(&p.ID, &p.Nome, &p.CPF, &p.Endereco); err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

func createPessoaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p Pessoa
		if err := c.BindJSON(&p); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result, err := db.Exec("INSERT INTO pessoa (nome, cpf, endereco) VALUES (?, ?, ?)", p.Nome, p.CPF, p.Endereco)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		id, _ := result.LastInsertId()
		p.ID = int(id)

		c.JSON(http.StatusCreated, p)
	}
}

func updatePessoaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var p Pessoa
		if err := c.BindJSON(&p); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result, err := db.Exec("UPDATE pessoa SET nome=?, cpf=?, endereco=? WHERE id=?", p.Nome, p.CPF, p.Endereco, id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Pessoa com id %s não encontrada", id))
			return
		}

		p.ID, _ = strconv.Atoi(id)
		c.JSON(http.StatusOK, p)
	}
}

func deletePessoaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := db.Exec("DELETE FROM pessoa WHERE id=?", id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Pessoa com id %s não encontrada", id))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pessoa excluída com sucesso"})
	}
}
