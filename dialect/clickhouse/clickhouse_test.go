package clickhouse_test

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/clickhouse"
)

func ExampleDialect() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.From("test").Where(goqu.Ex{"id": 10})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM `test` WHERE (`id` = 10) []
}

func ExampleDialect_insert() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.Insert("users").Rows(
		goqu.Record{"name": "Alice", "age": 30},
		goqu.Record{"name": "Bob", "age": 25},
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// INSERT INTO `users` (`age`, `name`) VALUES (30, 'Alice'), (25, 'Bob') []
}

func ExampleDialect_update() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.Update("users").
		Set(goqu.Record{"name": "Alice Smith"}).
		Where(goqu.C("id").Eq(1))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// UPDATE `users` SET `name`='Alice Smith' WHERE (`id` = 1) []
}

func ExampleDialect_delete() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.Delete("users").Where(goqu.C("id").Eq(1))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// DELETE FROM `users` WHERE (`id` = 1) []
}

func ExampleDialect_prepared() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.From("test").Prepared(true).Where(goqu.Ex{"id": 10})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM `test` WHERE (`id` = ?) [10]
}

func ExampleDialect_complexQuery() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.From("users").
		Select("name", goqu.COUNT("*").As("count")).
		Where(goqu.Ex{"active": true}).
		GroupBy("name").
		Having(goqu.COUNT("*").Gt(1)).
		Order(goqu.C("count").Desc())

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT `name`, COUNT(*) AS `count` FROM `users` WHERE (`active` IS TRUE) GROUP BY `name` HAVING (COUNT(*) > 1) ORDER BY `count` DESC []
}

func ExampleDialect_windowFunction() {
	dialect := goqu.Dialect("clickhouse")

	ds := dialect.From("users").
		Select(
			"name",
			"department",
			goqu.ROW_NUMBER().Over(goqu.W().PartitionBy("department").OrderBy(goqu.C("salary").Desc())).As("rank"),
		)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT `name`, `department`, ROW_NUMBER() OVER (PARTITION BY `department` ORDER BY `salary` DESC) AS `rank` FROM `users` []
}

func ExampleDialect_cte() {
	dialect := goqu.Dialect("clickhouse")

	cte := dialect.From("users").Select("id", "name").Where(goqu.C("active").IsTrue())
	ds := dialect.From("cte").With("cte", cte).Select("name")

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	// Output:
	// WITH cte AS (SELECT `id`, `name` FROM `users` WHERE (`active` IS TRUE)) SELECT `name` FROM `cte` []
}
