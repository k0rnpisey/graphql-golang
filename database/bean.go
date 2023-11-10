package database

import "gqlgen-todos/database/model"

func AllBeans() []interface{} {
	// Modify this list when adding or removing tables
	return []interface{}{
		new(model.Product),
	}
}
