package data

import _ "embed"

//go:embed schema.sql
var Schema string

//go:embed data.sql
var Rows string
