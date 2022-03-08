package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_complex-api/domain/complex"
	"github.com/superbkibbles/realestate_complex-api/domain/query"
	"github.com/superbkibbles/realestate_complex-api/domain/update"
	"github.com/superbkibbles/realestate_complex-api/services/complexService"
)

type ComplexHandler interface {
	Get(*gin.Context)
	Create(*gin.Context)
	GetByID(*gin.Context)
	UploadIcon(*gin.Context)
	Update(*gin.Context)
	Search(*gin.Context)
	DeleteIcon(*gin.Context)
	Translate(*gin.Context)
}

type complexHandler struct {
	srv complexService.ComplexService
}

func NewComplexHandler(srv complexService.ComplexService) ComplexHandler {
	return &complexHandler{
		srv: srv,
	}
}

func (ch *complexHandler) Get(c *gin.Context) {
	local := c.GetHeader("local")

	complexes, err := ch.srv.Get(local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, complexes)
}

func (ch *complexHandler) Translate(c *gin.Context) {
	complexID := strings.TrimSpace(c.Param("complex_id"))
	local := c.GetHeader("local")

	var complexRequest complex.TranslateRequest
	if err := c.ShouldBindJSON(&complexRequest); err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}
	agency, err := ch.srv.Translate(complexID, complexRequest, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (ch *complexHandler) Create(c *gin.Context) {
	var complex complex.Complex

	if err := c.ShouldBindJSON(&complex); err != nil {
		logger.Info(err.Error())
		restErr := rest_errors.NewBadRequestErr("Bad JSON body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := ch.srv.Save(&complex); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, complex)
}

func (ch *complexHandler) GetByID(c *gin.Context) {
	complexID := strings.TrimSpace(c.Param("complex_id"))
	local := c.GetHeader("local")

	complex, err := ch.srv.GetByID(complexID, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, complex)
}

func (ch *complexHandler) UploadIcon(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("complex_id"))

	file, err := c.FormFile("icon")
	if err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request")
		c.JSON(restErr.Status(), restErr)
		return
	}

	agency, uploadErr := ch.srv.UploadIcon(agencyID, file)
	if uploadErr != nil {
		c.JSON(uploadErr.Status(), uploadErr)
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (ch *complexHandler) Update(c *gin.Context) {
	id := strings.TrimSpace(c.Param("complex_id"))
	var updateRequest update.EsUpdate

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	complex, err := ch.srv.Update(id, updateRequest)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, complex)
}

func (ch *complexHandler) Search(c *gin.Context) {
	local := c.GetHeader("local")
	var q query.EsQuery

	if err := c.ShouldBindJSON(&q); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	properties, err := ch.srv.Search(q, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusFound, properties)
}

func (ch *complexHandler) DeleteIcon(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("complex_id"))

	err := ch.srv.DeleteIcon(agencyID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.String(200, "Icon Deleted")
}
