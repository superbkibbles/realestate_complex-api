package application

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_employee-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_employee-api/http"
	"github.com/superbkibbles/realestate_employee-api/repository/db"
	complexservice "github.com/superbkibbles/realestate_employee-api/services/complexService"
)

var (
	router  = gin.Default()
	handler http.ComplexHandler
)

func StartApplication() {
	elasticsearch.Client.Init()
	handler = http.NewComplexHandler(complexservice.NewComplexService(db.NewDbRepository()))
	router.Use(cors.Default())
	mapUrls()
	router.Static("assets", "clients/visuals")
	router.Run(":3040")
}
