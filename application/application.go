package application

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_complex-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_complex-api/http"
	"github.com/superbkibbles/realestate_complex-api/repository/db"
	complexservice "github.com/superbkibbles/realestate_complex-api/services/complexService"
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
