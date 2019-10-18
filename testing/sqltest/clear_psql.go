package sqltest

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/onsi/gomega"
)

// ClearPostgreSQL allow us to truncate every table of a database.
// We can provide extra table names to be ignored.
func ClearPostgreSQL(skip int, db *sql.DB, ignoreTable ...string) {
	var ignoreTablesCondition string

	if len(ignoreTable) > 0 {
		tables := make([]string, len(ignoreTable))
		for index, table := range ignoreTable {
			tables[index] = "'" + table + "'"
		}
		ignoreTablesCondition = fmt.Sprintf(`AND relname NOT IN (%s)`, strings.Join(tables, ","))
	}

	query := `
		DO
		$func$
		BEGIN
			EXECUTE
			(
				SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
				FROM   pg_class
				WHERE  relkind = 'r'  -- only tables
				AND    relnamespace = 'public'::regnamespace
				` + ignoreTablesCondition + `
			);
		END
		$func$;
	`
	_, err := db.Exec(query)
	gomega.ExpectWithOffset(skip, err).ToNot(gomega.HaveOccurred())
}
