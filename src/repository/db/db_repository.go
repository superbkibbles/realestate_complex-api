package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_employee-api/src/clients/elasticsearch"
	"github.com/superbkibbles/realestate_employee-api/src/domain/complex"
	"github.com/superbkibbles/realestate_employee-api/src/domain/query"
	"github.com/superbkibbles/realestate_employee-api/src/domain/update"
)

var (
	indexAgency = "complex"
	agencyType  = "_doc"
)

type DbRepository interface {
	Get() (complex.Complexes, rest_errors.RestErr)
	Save(complex *complex.Complex) rest_errors.RestErr
	GetByID(complexID string) (*complex.Complex, rest_errors.RestErr)
	UploadIcon(agency *complex.Complex, id string) rest_errors.RestErr
	Update(id string, updateRequest update.EsUpdate) (*complex.Complex, rest_errors.RestErr)
	Search(query query.EsQuery) (complex.Complexes, rest_errors.RestErr)
}
type dbRepository struct{}

func NewDbRepository() DbRepository {
	return &dbRepository{}
}

func (db *dbRepository) Get() (complex.Complexes, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetAllDoc(indexAgency)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to Get All Agencies Property", errors.New("databse error"))
	}

	complexes := make(complex.Complexes, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var c complex.Complex
		if err := json.Unmarshal(bytes, &c); err != nil {
			return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		}
		c.ID = hit.Id
		complexes[i] = c
	}

	return complexes, nil
}

func (db *dbRepository) GetByID(complexID string) (*complex.Complex, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetByID(indexAgency, agencyType, complexID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("no Property was found with id %s", complexID))
		}
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to id %s", complexID), errors.New("database error"))
	}

	var c complex.Complex

	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to parse database response", errors.New("database error"))
	}

	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
	}

	c.ID = result.Id
	return &c, nil
}

func (db *dbRepository) Save(c *complex.Complex) rest_errors.RestErr {
	result, err := elasticsearch.Client.Save(indexAgency, agencyType, c)
	if err != nil {
		return rest_errors.NewInternalServerErr("error when trying to save Property", errors.New("databse error"))
	}

	c.ID = result.Id
	return nil
}

func (db *dbRepository) UploadIcon(agency *complex.Complex, id string) rest_errors.RestErr {
	var es update.EsUpdate
	update := update.UpdatePropertyRequest{
		Field: "photo",
		Value: agency.Photo,
	}
	es.Fields = append(es.Fields, update)
	_, err := db.Update(id, es)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbRepository) Update(id string, updateRequest update.EsUpdate) (*complex.Complex, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Update(indexAgency, agencyType, id, updateRequest)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("no Property was found with id %s", id))
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Update Property", errors.New("databse error"))
	}

	var c complex.Complex

	bytes, err := result.GetResult.Source.MarshalJSON()
	if err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}
	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}

	c.ID = result.Id
	return &c, nil
}

func (db *dbRepository) Search(query query.EsQuery) (complex.Complexes, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexAgency, query.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to search documents", errors.New("database error"))
	}

	complexes := make(complex.Complexes, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var c complex.Complex
		if err := json.Unmarshal(bytes, &c); err != nil {
			return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		}
		c.ID = hit.Id
		complexes[i] = c
	}

	if len(complexes) == 0 {
		return nil, rest_errors.NewNotFoundErr("no items found matching given critirial")
	}

	return complexes, nil
}
