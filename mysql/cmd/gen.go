package main

import (
	"a-project-backend/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
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
	g.Execute()
	all := g.GenerateAllTable() // database to table model.

	all = append(all, g.GenerateModel(
		model.TableNameUser,
		gen.FieldRelateModel(
			field.Many2Many,
			"Tags",
			model.UserTag{},
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag: field.GormTag{}.
					Set("many2many", "user_tags"),
			}),
	))
	g.ApplyBasic(all...)

	// Generate the code
	g.Execute()
}
