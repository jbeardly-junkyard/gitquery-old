package gitquery

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
	"gopkg.in/src-d/go-mysql-server.v0/sql/expression"
	"gopkg.in/src-d/go-mysql-server.v0/sql/plan"

	"gopkg.in/src-d/go-git-fixtures.v3"
)

func TestReferencesTable_Name(t *testing.T) {
	require := require.New(t)

	f := fixtures.Basic().One()
	table := getTable(require, f, referencesTableName)
	require.Equal(referencesTableName, table.Name())
}

func TestReferencesTable_Children(t *testing.T) {
	require := require.New(t)

	f := fixtures.Basic().One()
	table := getTable(require, f, referencesTableName)
	require.Equal(0, len(table.Children()))
}

func TestReferencesTable_RowIter(t *testing.T) {
	require := require.New(t)

	f := fixtures.Basic().One()
	table := getTable(require, f, referencesTableName)

	rows, err := sql.NodeToRows(plan.NewSort(
		[]plan.SortField{{Column: expression.NewGetField(0, sql.Text, "name", false), Order: plan.Ascending}},
		table))
	require.Nil(err)

	expected := []sql.Row{
		sql.NewRow("HEAD", "symbolic-reference", nil, "refs/heads/master", false, false, false, false),
		sql.NewRow("refs/heads/branch", "hash-reference", "e8d3ffab552895c19b9fcf7aa264d277cde33881", nil, true, false, false, false),
		sql.NewRow("refs/heads/master", "hash-reference", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5", nil, true, false, false, false),
		sql.NewRow("refs/remotes/origin/HEAD", "symbolic-reference", nil, "refs/remotes/origin/master", false, false, true, false),
		sql.NewRow("refs/remotes/origin/branch", "hash-reference", "e8d3ffab552895c19b9fcf7aa264d277cde33881", nil, false, false, true, false),
		sql.NewRow("refs/remotes/origin/master", "hash-reference", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5", nil, false, false, true, false),
		sql.NewRow("refs/tags/v1.0.0", "hash-reference", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5", nil, false, false, false, true),
	}
	require.Equal(expected, rows)

	schema := table.Schema()
	for idx, row := range rows {
		err := schema.CheckRow(row)
		require.Nil(err, "row %d doesn't conform to schema", idx)
	}
}
