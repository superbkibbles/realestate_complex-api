package complexService

import (
	"mime/multipart"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_employee-api/domain/complex"
	"github.com/superbkibbles/realestate_employee-api/domain/query"
	"github.com/superbkibbles/realestate_employee-api/domain/update"
	"github.com/superbkibbles/realestate_employee-api/repository/db"
	"github.com/superbkibbles/realestate_employee-api/utils/date_utils"
	"github.com/superbkibbles/realestate_employee-api/utils/file_utils"
)

type ComplexService interface {
	Get() (complex.Complexes, rest_errors.RestErr)
	Save(complex *complex.Complex) rest_errors.RestErr
	GetByID(complexID string) (*complex.Complex, rest_errors.RestErr)
	UploadIcon(id string, fileHeader *multipart.FileHeader) (*complex.Complex, rest_errors.RestErr)
	Update(id string, updateRequest update.EsUpdate) (*complex.Complex, rest_errors.RestErr)
	Search(updateRequest query.EsQuery) (complex.Complexes, rest_errors.RestErr)
	DeleteIcon(agencyID string) rest_errors.RestErr
}

type complexService struct {
	dbRepo db.DbRepository
}

func NewComplexService(dbRepo db.DbRepository) ComplexService {
	return &complexService{
		dbRepo: dbRepo,
	}
}

func (srv *complexService) Get() (complex.Complexes, rest_errors.RestErr) {
	return srv.dbRepo.Get()
}

func (srv *complexService) Save(c *complex.Complex) rest_errors.RestErr {
	c.Status = complex.STATUS_ACTIVE
	c.DateCreated = date_utils.GetNowDBFromat()
	return srv.dbRepo.Save(c)
}

func (cs *complexService) GetByID(complexID string) (*complex.Complex, rest_errors.RestErr) {
	return cs.dbRepo.GetByID(complexID)
}

func (srv *complexService) UploadIcon(id string, fileHeader *multipart.FileHeader) (*complex.Complex, rest_errors.RestErr) {
	complex, err := srv.GetByID(id)
	if err != nil {
		return nil, err
	}
	file, fErr := fileHeader.Open()
	if fErr != nil {
		return nil, rest_errors.NewInternalServerErr("Error while trying to open the file", nil)
	}
	filePath, err := file_utils.SaveFile(fileHeader, file)
	if err != nil {
		return nil, err
	}
	complex.Photo = "http://localhost:3040/assets/" + filePath

	srv.dbRepo.UploadIcon(complex, id)
	return complex, nil
}

func (srv *complexService) Update(id string, updateRequest update.EsUpdate) (*complex.Complex, rest_errors.RestErr) {
	return srv.dbRepo.Update(id, updateRequest)
}

func (srv *complexService) Search(updateRequest query.EsQuery) (complex.Complexes, rest_errors.RestErr) {
	return srv.dbRepo.Search(updateRequest)
}

func (srv *complexService) DeleteIcon(agencyID string) rest_errors.RestErr {
	agency, err := srv.GetByID(agencyID)
	if err != nil {
		return err
	}

	splittedPath := strings.Split(agency.Photo, "/")
	fileName := splittedPath[len(splittedPath)-1]

	file_utils.DeleteFile(fileName)

	agency.Photo = ""
	srv.dbRepo.UploadIcon(agency, agencyID)
	return nil
}
