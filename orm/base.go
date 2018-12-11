package orm

import (
	"log"
	"reflect"

	"github.com/expenseledger/web-service/db"
)

// BaseMapper ...
type BaseMapper struct {
	modelType  reflect.Type
	insertStmt string
	deleteStmt string
	oneStmt    string
	manyStmt   string
	clearStmt  string
}

// Insert ...
func (mapper *BaseMapper) Insert(obj interface{}) (interface{}, error) {
	return worker(obj, mapper.modelType, mapper.insertStmt, "Error inserting")
}

// Delete ...
func (mapper *BaseMapper) Delete(obj interface{}) (interface{}, error) {
	return worker(obj, mapper.modelType, mapper.deleteStmt, "Error deleting")
}

// One ...
func (mapper *BaseMapper) One(obj interface{}) (interface{}, error) {
	return worker(obj, mapper.modelType, mapper.oneStmt, "Error geting")
}

// Many ...
func (mapper *BaseMapper) Many() (interface{}, error) {
	return sliceWorker(mapper.modelType, mapper.manyStmt, "Error selecting")
}

// Clear ...
func (mapper *BaseMapper) Clear() (interface{}, error) {
	return sliceWorker(mapper.modelType, mapper.clearStmt, "Error clearing")
}

func worker(
	obj interface{},
	t reflect.Type,
	sqlStmt string,
	logMsg string,
) (interface{}, error) {
	stmt, err := db.Conn().PrepareNamed(sqlStmt)
	if err != nil {
		log.Println(logMsg, err)
		return nil, err
	}

	newObj := reflect.New(t).Interface()
	if err := stmt.Get(newObj, obj); err != nil {
		log.Println(logMsg, err)
		return nil, err
	}
	return newObj, nil
}

func sliceWorker(t reflect.Type, sqlStmt, logMsg string) (interface{}, error) {
	stmt, err := db.Conn().Preparex(sqlStmt)
	if err != nil {
		log.Println("Error selecting", err)
		return nil, err
	}

	sliceType := reflect.SliceOf(t)
	newSlice := reflect.MakeSlice(sliceType, 0, 0)
	resultSet := reflect.New(newSlice.Type()).Interface()
	if err := stmt.Select(resultSet); err != nil {
		log.Println("Error selecting", err)
		return nil, err
	}
	return resultSet, nil
}
