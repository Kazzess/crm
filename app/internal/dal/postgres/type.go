package postgres

import "github.com/Kazzess/libraries/queryify"

var (
	CategoryTable = queryify.NewTable("public", "category", "c", "id")
)
