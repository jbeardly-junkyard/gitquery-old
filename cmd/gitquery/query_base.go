package main

import (
	"io"
	"os"

	"github.com/src-d/gitquery"
	"github.com/src-d/gitquery/internal/format"
	"github.com/src-d/gitquery/internal/function"

	"gopkg.in/src-d/go-git.v4/utils/ioutil"
	sqle "gopkg.in/src-d/go-mysql-server.v0"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

type cmdQueryBase struct {
	cmd

	Path string `short:"p" long:"path" description:"Path where the git repositories are located, one per dir"`

	engine *sqle.Engine
	name   string
}

func (c *cmdQueryBase) buildDatabase() error {
	if c.engine == nil {
		c.engine = sqle.New()
	}

	c.print("opening %q repository...\n", c.Path)

	var err error

	pool := gitquery.NewRepositoryPool()
	err = pool.AddDir(c.Path)
	if err != nil {
		println("ERR", err.Error())
		return err
	}

	c.engine.AddDatabase(gitquery.NewDatabase(c.name, &pool))
	function.Register(c.engine.Catalog)
	return nil
}

func (c *cmdQueryBase) executeQuery(sql string) (sql.Schema, sql.RowIter, error) {
	c.print("executing %q at %q\n", sql, c.name)
	return c.engine.Query(sql)
}

func (c *cmdQueryBase) printQuery(schema sql.Schema, rows sql.RowIter, formatId string) (err error) {
	defer ioutil.CheckClose(rows, &err)

	f, err := format.NewFormat(formatId, os.Stdout)
	if err != nil {
		return err
	}
	defer ioutil.CheckClose(f, &err)

	columnNames := make([]string, len(schema))
	for i, column := range schema {
		columnNames[i] = column.Name
	}

	if err := f.WriteHeader(columnNames); err != nil {
		return err
	}

	for {
		row, err := rows.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if err := f.Write(row); err != nil {
			return err
		}
	}

	return nil
}
