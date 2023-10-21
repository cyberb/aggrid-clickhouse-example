package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	database *Database
	logger   *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		database: NewDatabase(logger),
		logger:   logger,
	}
}
func (s *Server) Start() error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/clickhouse", s.Data)

	err := s.database.Connect()
	if err != nil {
		return err
	}
	err = s.database.GenerateData()
	if err != nil {
		return err
	}
	s.database.StartUpdater()
	return e.Start("Localhost:10000")
}

func (s *Server) Data(c echo.Context) error {

	var request Request
	err := c.Bind(&request)

	if err != nil {
		return err
	}
	data, total, err := s.database.GetData(request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{Data: data, RowCount: total})
}
