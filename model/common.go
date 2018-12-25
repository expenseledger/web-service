package model

type operation int

const (
	insert operation = iota
	delete
	one
	list
	clear
)
