package main

//go:generate safemap -k TypeKey -v TypeValue

// These Type* are used for tests.
type (
	TypeKey   string
	TypeValue string
)
