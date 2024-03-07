package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./db/query", // 出力パス
		ModelPkgPath:      "./model",
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	sqliteConn := mysql.Open("user:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(sqliteConn, &gorm.Config{})
	if err != nil {
		return
	}

	g.UseDB(db)
	all := g.GenerateAllTable() // database to table model.

	g.ApplyBasic(all...)

	// Generate the code
	g.Execute()
}
