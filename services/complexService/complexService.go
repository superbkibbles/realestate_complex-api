package complexService

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_complex-api/domain/complex"
	"github.com/superbkibbles/realestate_complex-api/domain/query"
	"github.com/superbkibbles/realestate_complex-api/domain/update"
	cloudstorage "github.com/superbkibbles/realestate_complex-api/repository/cloudStorage"
	"github.com/superbkibbles/realestate_complex-api/repository/db"
	"github.com/superbkibbles/realestate_complex-api/utils/crypto_utils"
	"github.com/superbkibbles/realestate_complex-api/utils/date_utils"
)

type ComplexService interface {
	Get(local string) (complex.Complexes, rest_errors.RestErr)
	Save(complex *complex.Complex) rest_errors.RestErr
	GetByID(complexID string, local string) (*complex.Complex, rest_errors.RestErr)
	UploadIcon(id string, fileHeader *multipart.FileHeader) (*complex.Complex, rest_errors.RestErr)
	Update(id string, updateRequest update.EsUpdate) (*complex.Complex, rest_errors.RestErr)
	Search(updateRequest query.EsQuery, local string) (complex.Complexes, rest_errors.RestErr)
	DeleteIcon(agencyID string) rest_errors.RestErr
	Translate(complexID string, complexRequest complex.TranslateRequest, local string) (*complex.Complex, rest_errors.RestErr)
}

type complexService struct {
	dbRepo    db.DbRepository
	cloudRepo cloudstorage.CloudStorage
}

func NewComplexService(dbRepo db.DbRepository, cloudRepo cloudstorage.CloudStorage) ComplexService {
	return &complexService{
		dbRepo:    dbRepo,
		cloudRepo: cloudRepo,
	}
}

func (srv *complexService) Get(local string) (complex.Complexes, rest_errors.RestErr) {
	complexes, err := srv.dbRepo.Get()
	if err != nil {
		return nil, err
	}

	if local == "ar" {
		for i, complex := range complexes {
			if complexes[i].Ar.Address != "" {
				complexes[i].Address = complex.Ar.Address
			}

			if complexes[i].Ar.Description != "" {
				complexes[i].Description = complex.Ar.Description
			}

			if complexes[i].Ar.Name != "" {
				complexes[i].Name = complex.Ar.Name
			}
		}
	}

	if local == "kur" {
		for i, complex := range complexes {
			if complexes[i].Kur.Address != "" {
				complexes[i].Address = complex.Kur.Address
			}

			if complexes[i].Kur.Description != "" {
				complexes[i].Description = complex.Kur.Description
			}

			if complexes[i].Ar.Name != "" {
				complexes[i].Name = complex.Kur.Name
			}
		}
	}

	return complexes, nil
}

func (srv *complexService) Save(c *complex.Complex) rest_errors.RestErr {
	c.Status = complex.STATUS_ACTIVE
	c.DateCreated = date_utils.GetNowDBFromat()
	return srv.dbRepo.Save(c)
}

func (srv *complexService) Translate(complexID string, complexRequest complex.TranslateRequest, local string) (*complex.Complex, rest_errors.RestErr) {
	complex, err := srv.dbRepo.GetByID(complexID)
	if err != nil {
		return nil, err
	}

	if local == "ar" || local == "kur" {
		complex, err = srv.dbRepo.Translate(complexID, complexRequest, local)
		if err != nil {
			return nil, err
		}
	}

	return complex, nil
}

func (cs *complexService) GetByID(complexID string, local string) (*complex.Complex, rest_errors.RestErr) {
	complex, err := cs.dbRepo.GetByID(complexID)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	if local == "ar" {
		if complex.Ar.Address != "" {
			complex.Address = complex.Ar.Address
		}

		if complex.Ar.Description != "" {
			complex.Description = complex.Ar.Description
		}

		if complex.Ar.Name != "" {
			complex.Name = complex.Ar.Name
		}
	}
	if local == "kur" {
		if complex.Kur.Address != "" {
			complex.Address = complex.Kur.Address
		}

		if complex.Kur.Description != "" {
			complex.Description = complex.Kur.Description
		}

		if complex.Kur.Name != "" {
			complex.Name = complex.Kur.Name
		}
	}
	return complex, nil
}

func (srv *complexService) UploadIcon(id string, fileHeader *multipart.FileHeader) (*complex.Complex, rest_errors.RestErr) {
	complex, err := srv.GetByID(id, "")
	if err != nil {
		return nil, err
	}
	file, fErr := fileHeader.Open()
	if fErr != nil {
		return nil, rest_errors.NewInternalServerErr("Error while trying to open the file", nil)
	}
	// filePath, err := file_utils.SaveFile(fileHeader, file)
	// if err != nil {
	// 	return nil, err
	// }

	res, cloudErr := srv.cloudRepo.Save(file, id+crypto_utils.GetMd5(uuid.New().String()), id)
	if cloudErr != nil {
		return nil, err
	}

	complex.Photo = res.Url
	complex.PublicID = res.PublicID

	srv.dbRepo.UploadIcon(complex, id)
	return complex, nil
}

func (srv *complexService) Update(id string, updateRequest update.EsUpdate) (*complex.Complex, rest_errors.RestErr) {
	return srv.dbRepo.Update(id, updateRequest)
}

func (srv *complexService) Search(updateRequest query.EsQuery, local string) (complex.Complexes, rest_errors.RestErr) {
	complexes, err := srv.dbRepo.Search(updateRequest)
	if err != nil {
		return nil, err
	}
	if local == "ar" {
		for i, complex := range complexes {
			if complexes[i].Ar.Address != "" {
				complexes[i].Address = complex.Ar.Address
			}

			if complexes[i].Ar.Description != "" {
				complexes[i].Description = complex.Ar.Description
			}

			if complexes[i].Ar.Name != "" {
				complexes[i].Name = complex.Ar.Name
			}
		}
	}

	if local == "kur" {
		for i, complex := range complexes {
			if complexes[i].Kur.Address != "" {
				complexes[i].Address = complex.Kur.Address
			}

			if complexes[i].Kur.Description != "" {
				complexes[i].Description = complex.Kur.Description
			}

			if complexes[i].Ar.Name != "" {
				complexes[i].Name = complex.Kur.Name
			}
		}
	}

	return complexes, nil
}

func (srv *complexService) DeleteIcon(agencyID string) rest_errors.RestErr {
	agency, err := srv.GetByID(agencyID, "")
	if err != nil {
		return err
	}

	if err := srv.cloudRepo.Delete(agency.PublicID); err != nil {
		return err
	}

	agency.Photo = ""
	agency.PublicID = ""
	srv.dbRepo.UploadIcon(agency, agencyID)
	return nil
}
