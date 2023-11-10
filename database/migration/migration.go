package main

import (
	"database/sql"
	"fmt"
	. "gqlgen-todos/database"
	"reflect"
)

func main() {
	_, err := Connect()
	if err != nil {
		fmt.Printf("Failed to connect to the DB instance: %+v\n", err)
		return
	}
	sqlDB, err := GetConnection().DB()
	if err != nil {
		fmt.Printf("Failed to get sql.DB because of [%T] %+v\n", err, err)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			fmt.Printf("Failed to close sql.DB because of [%T] %+v\n", err, err)
		}
	}(sqlDB)

	// GORM auto migration
	err = migrate(AllBeans())
	if err != nil {
		fmt.Printf("Migration error: %+v\n", err)
	}
}

func migrate(beans []interface{}) (err error) {
	if err = checkBeans(beans); err != nil {
		return
	}

	err = GetConnection().AutoMigrate(beans...)
	if err != nil {
		return err
	}

	return
}

func checkBeans(beans []interface{}) (err error) {
	m := make(map[string]bool)
	for _, b := range beans {
		elm := reflect.ValueOf(b).Type().Elem()
		name := fmt.Sprintf("%s/%s", elm.PkgPath(), elm.Name())
		if m[name] {
			err = fmt.Errorf("%s is duplicated in `func AllBeans()`", name)
			return
		} else {
			m[name] = true
		}
	}
	return
}
