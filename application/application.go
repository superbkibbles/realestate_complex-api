package application

import (
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_complex-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_complex-api/constants"
	"github.com/superbkibbles/realestate_complex-api/http"
	cloudstorage "github.com/superbkibbles/realestate_complex-api/repository/cloudStorage"
	"github.com/superbkibbles/realestate_complex-api/repository/db"
	complexservice "github.com/superbkibbles/realestate_complex-api/services/complexService"
)

var (
	router  = gin.Default()
	handler http.ComplexHandler
)

func StartApplication() {
	cld, err := cloudinary.NewFromParams(os.Getenv(constants.CLOUD_STORAGE_NAME), os.Getenv(constants.CLOUD_STORAGE_API_KEY), os.Getenv(constants.CLOUD_STORAGE_API_SECRET))
	if err != nil {
		panic(err)
	}
	elasticsearch.Client.Init()
	handler = http.NewComplexHandler(complexservice.NewComplexService(db.NewDbRepository(), cloudstorage.NewRepository(cld)))
	router.Use(cors.Default())
	mapUrls()
	router.Static("assets", "clients/visuals")
	router.Run(os.Getenv(constants.PORT))
}
