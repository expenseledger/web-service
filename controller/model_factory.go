package controller

import (
	"errors"

	"github.com/expenseledger/web-service/orm"
)

type Modeler interface {
	Model() interface{}
}

type Intent int

const (
	Create Intent = iota
	Get
)

func makeModel(
	mapper orm.Mapper,
	obj interface{},
	intent Intent,
) (interface{}, error) {
	switch intent {
	case Create:
		return mapper.Insert(obj)
	case Get:
		return mapper.One(obj)
	}

	return nil, errors.New("unknown intent")
}
