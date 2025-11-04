package clickhouse

import (
	"github.com/doug-martin/goqu/v9"
)

func DialectOptions() *goqu.SQLDialectOptions {
	opts := goqu.DefaultDialectOptions()

	// ClickHouse uses backticks for quoting identifiers
	opts.QuoteRune = '`'

	// ClickHouse does not support RETURNING clause
	opts.SupportsReturn = false

	// ClickHouse does not support ORDER BY on UPDATE/DELETE (standard SQL)
	opts.SupportsOrderByOnUpdate = false
	opts.SupportsOrderByOnDelete = false

	// ClickHouse does not support LIMIT on UPDATE/DELETE in standard form
	opts.SupportsLimitOnUpdate = false
	opts.SupportsLimitOnDelete = false

	// ClickHouse does not support ON CONFLICT clause (INSERT IGNORE is different)
	opts.SupportsConflictTarget = false
	opts.SupportsConflictUpdateWhere = false
	opts.SupportsInsertIgnoreSyntax = false

	// ClickHouse supports CTEs (Common Table Expressions)
	opts.SupportsWithCTE = true
	opts.SupportsWithCTERecursive = false // ClickHouse doesn't support recursive CTEs

	// ClickHouse doesn't support multiple tables in UPDATE
	opts.SupportsMultipleUpdateTables = false

	// ClickHouse supports DISTINCT ON
	opts.SupportsDistinctOn = true

	// ClickHouse supports window functions
	opts.SupportsWindowFunction = true

	// ClickHouse doesn't support LATERAL joins
	opts.SupportsLateral = false

	// Placeholder configuration for prepared statements
	opts.PlaceHolderFragment = []byte("?")
	opts.IncludePlaceholderNum = false

	// Default values configuration
	opts.DefaultValuesFragment = []byte("")

	// Boolean literal representation in ClickHouse (used for literal values, not IS TRUE/IS FALSE operators)
	// ClickHouse uses 1 and 0 for boolean literals in INSERT/UPDATE statements
	opts.True = []byte("1")
	opts.False = []byte("0")

	// Time format for ClickHouse
	opts.TimeFormat = "2006-01-02 15:04:05"

	// String escaping - ClickHouse uses backslash escaping
	opts.EscapedRunes = map[rune][]byte{
		'\'': []byte("\\'"),
		'\\': []byte("\\\\"),
		'\n': []byte("\\n"),
		'\r': []byte("\\r"),
		'\t': []byte("\\t"),
		0:    []byte("\\0"),
	}

	return opts
}

func init() {
	goqu.RegisterDialect("clickhouse", DialectOptions())
}
