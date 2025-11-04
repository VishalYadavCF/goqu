package clickhouse_test

import (
	"regexp"
	"testing"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/clickhouse"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/suite"
)

type (
	clickhouseDialectSuite struct {
		suite.Suite
	}
	sqlTestCase struct {
		ds         exp.SQLExpression
		sql        string
		err        string
		isPrepared bool
		args       []interface{}
	}
)

func (cds *clickhouseDialectSuite) GetDs(table string) *goqu.SelectDataset {
	return goqu.Dialect("clickhouse").From(table)
}

func (cds *clickhouseDialectSuite) assertSQL(cases ...sqlTestCase) {
	for i, c := range cases {
		actualSQL, actualArgs, err := c.ds.ToSQL()
		if c.err == "" {
			cds.NoError(err, "test case %d failed", i)
		} else {
			cds.EqualError(err, c.err, "test case %d failed", i)
		}
		cds.Equal(c.sql, actualSQL, "test case %d failed", i)
		if c.isPrepared && c.args != nil || len(c.args) > 0 {
			cds.Equal(c.args, actualArgs, "test case %d failed", i)
		} else {
			cds.Empty(actualArgs, "test case %d failed", i)
		}
	}
}

func (cds *clickhouseDialectSuite) TestIdentifiers() {
	ds := cds.GetDs("test")
	cds.assertSQL(
		sqlTestCase{ds: ds.Select(
			"a",
			goqu.I("a.b.c"),
			goqu.I("c.d"),
			goqu.C("test").As("test"),
		), sql: "SELECT `a`, `a`.`b`.`c`, `c`.`d`, `test` AS `test` FROM `test`"},
	)
}

func (cds *clickhouseDialectSuite) TestLiteralString() {
	ds := cds.GetDs("test")
	col := goqu.C("a")
	cds.assertSQL(
		sqlTestCase{ds: ds.Where(col.Eq("test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test')"},
		sqlTestCase{ds: ds.Where(col.Eq("test'test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\'test')"},
		sqlTestCase{ds: ds.Where(col.Eq(`test\test`)), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\\\test')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\ntest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\ntest')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\rtest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\rtest')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\ttest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\ttest')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\x00test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\0test')"},
	)
}

func (cds *clickhouseDialectSuite) TestLiteralBytes() {
	col := goqu.C("a")
	ds := cds.GetDs("test")
	cds.assertSQL(
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test'test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\'test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte(`test\test`))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\\\test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\ntest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\ntest')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\rtest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\rtest')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\ttest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\ttest')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\x00test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\0test')"},
	)
}

func (cds *clickhouseDialectSuite) TestBooleanOperations() {
	col := goqu.C("a")
	ds := cds.GetDs("test")
	cds.assertSQL(
		sqlTestCase{ds: ds.Where(col.Eq(true)), sql: "SELECT * FROM `test` WHERE (`a` IS TRUE)"},
		sqlTestCase{ds: ds.Where(col.Eq(false)), sql: "SELECT * FROM `test` WHERE (`a` IS FALSE)"},
		sqlTestCase{ds: ds.Where(col.Is(true)), sql: "SELECT * FROM `test` WHERE (`a` IS TRUE)"},
		sqlTestCase{ds: ds.Where(col.Is(false)), sql: "SELECT * FROM `test` WHERE (`a` IS FALSE)"},
		sqlTestCase{ds: ds.Where(col.IsTrue()), sql: "SELECT * FROM `test` WHERE (`a` IS TRUE)"},
		sqlTestCase{ds: ds.Where(col.IsFalse()), sql: "SELECT * FROM `test` WHERE (`a` IS FALSE)"},
		sqlTestCase{ds: ds.Where(col.Neq(true)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT TRUE)"},
		sqlTestCase{ds: ds.Where(col.Neq(false)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT FALSE)"},
		sqlTestCase{ds: ds.Where(col.IsNot(true)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT TRUE)"},
		sqlTestCase{ds: ds.Where(col.IsNot(false)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT FALSE)"},
		sqlTestCase{ds: ds.Where(col.IsNotTrue()), sql: "SELECT * FROM `test` WHERE (`a` IS NOT TRUE)"},
		sqlTestCase{ds: ds.Where(col.IsNotFalse()), sql: "SELECT * FROM `test` WHERE (`a` IS NOT FALSE)"},
		sqlTestCase{ds: ds.Where(col.Like("a%")), sql: "SELECT * FROM `test` WHERE (`a` LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(col.NotLike("a%")), sql: "SELECT * FROM `test` WHERE (`a` NOT LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(col.ILike("a%")), sql: "SELECT * FROM `test` WHERE (`a` ILIKE 'a%')"},
		sqlTestCase{ds: ds.Where(col.NotILike("a%")), sql: "SELECT * FROM `test` WHERE (`a` NOT ILIKE 'a%')"},
		sqlTestCase{ds: ds.Where(col.Like(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` ~ '[ab]')"},
		sqlTestCase{ds: ds.Where(col.NotLike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` !~ '[ab]')"},
		sqlTestCase{ds: ds.Where(col.ILike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` ~* '[ab]')"},
		sqlTestCase{ds: ds.Where(col.NotILike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` !~* '[ab]')"},
	)
}

func (cds *clickhouseDialectSuite) TestBitwiseOperations() {
	col := goqu.C("a")
	ds := cds.GetDs("test")
	cds.assertSQL(
		sqlTestCase{ds: ds.Where(col.BitwiseInversion()), sql: "SELECT * FROM `test` WHERE (~ `a`)"},
		sqlTestCase{ds: ds.Where(col.BitwiseAnd(1)), sql: "SELECT * FROM `test` WHERE (`a` & 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseOr(1)), sql: "SELECT * FROM `test` WHERE (`a` | 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseXor(1)), sql: "SELECT * FROM `test` WHERE (`a` # 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseLeftShift(1)), sql: "SELECT * FROM `test` WHERE (`a` << 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseRightShift(1)), sql: "SELECT * FROM `test` WHERE (`a` >> 1)"},
	)
}

func (cds *clickhouseDialectSuite) TestSelectSQL() {
	ds := cds.GetDs("test")
	cds.assertSQL(
		sqlTestCase{ds: ds.Select("id", "name"), sql: "SELECT `id`, `name` FROM `test`"},
		sqlTestCase{ds: ds.Where(goqu.C("id").Eq(10)), sql: "SELECT * FROM `test` WHERE (`id` = 10)"},
		sqlTestCase{
			ds:  ds.Prepared(true).Where(goqu.L("? = ?", goqu.C("id"), 10)),
			sql: "SELECT * FROM `test` WHERE `id` = ?", args: []interface{}{int64(10)},
		},
	)
}

func (cds *clickhouseDialectSuite) TestInsertSQL() {
	ds := goqu.Dialect("clickhouse").Insert("test")
	cds.assertSQL(
		sqlTestCase{
			ds:  ds.Rows(goqu.Record{"name": "test"}),
			sql: "INSERT INTO `test` (`name`) VALUES ('test')",
		},
		sqlTestCase{
			ds:  ds.Rows(goqu.Record{"name": "test1"}, goqu.Record{"name": "test2"}),
			sql: "INSERT INTO `test` (`name`) VALUES ('test1'), ('test2')",
		},
	)
}

func (cds *clickhouseDialectSuite) TestUpdateSQL() {
	ds := goqu.Dialect("clickhouse").Update("test")
	cds.assertSQL(
		sqlTestCase{
			ds:  ds.Set(goqu.Record{"foo": "bar"}).Where(goqu.C("id").Eq(10)),
			sql: "UPDATE `test` SET `foo`='bar' WHERE (`id` = 10)",
		},
	)
}

func (cds *clickhouseDialectSuite) TestDeleteSQL() {
	ds := goqu.Dialect("clickhouse").Delete("test")
	cds.assertSQL(
		sqlTestCase{
			ds:  ds.Where(goqu.C("id").Eq(10)),
			sql: "DELETE FROM `test` WHERE (`id` = 10)",
		},
	)
}

func (cds *clickhouseDialectSuite) TestTruncateSQL() {
	ds := goqu.Dialect("clickhouse").Truncate("test")
	cds.assertSQL(
		sqlTestCase{
			ds:  ds,
			sql: "TRUNCATE `test`",
		},
	)
}

func (cds *clickhouseDialectSuite) TestWindowFunction() {
	ds := cds.GetDs("entry").
		Select("int", goqu.ROW_NUMBER().OverName(goqu.I("w")).As("id")).
		Window(goqu.W("w").OrderBy(goqu.I("int").Desc()))

	cds.assertSQL(
		sqlTestCase{
			ds:  ds,
			sql: "SELECT `int`, ROW_NUMBER() OVER `w` AS `id` FROM `entry` WINDOW `w` AS (ORDER BY `int` DESC)",
		},
	)
}

func (cds *clickhouseDialectSuite) TestCommonTableExpression() {
	ds := cds.GetDs("test")
	cte := goqu.Dialect("clickhouse").
		From("other").
		Select("id", "name")

	cds.assertSQL(
		sqlTestCase{
			ds:  ds.With("cte", cte).Select("id").From("cte"),
			sql: "WITH cte AS (SELECT `id`, `name` FROM `other`) SELECT `id` FROM `cte`",
		},
	)
}

func TestClickHouseDialectSuite(t *testing.T) {
	suite.Run(t, new(clickhouseDialectSuite))
}
