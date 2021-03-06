package function

import (
	"reflect"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
	"gopkg.in/src-d/go-mysql-server.v0/sql/expression"
)

// IsTag checks the given string is a tag name.
type IsTag struct {
	expression.UnaryExpression
}

// NewIsTag creates a new IsTag function.
func NewIsTag(e sql.Expression) sql.Expression {
	return &IsTag{expression.UnaryExpression{Child: e}}
}

var _ sql.Expression = (*IsTag)(nil)

// Eval implements the expression interface.
func (f *IsTag) Eval(row sql.Row) (interface{}, error) {
	val, err := f.Child.Eval(row)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return false, nil
	}

	name, ok := val.(string)
	if !ok {
		return nil, sql.ErrInvalidType.New(reflect.TypeOf(val).String())
	}

	return plumbing.ReferenceName(name).IsTag(), nil
}

// Name implements the Expression interface.
func (IsTag) Name() string {
	return "is_tag"
}

// TransformUp implements the Expression interface.
func (f IsTag) TransformUp(fn func(sql.Expression) sql.Expression) sql.Expression {
	return NewIsTag(fn(f.Child))
}

// Type implements the Expression interface.
func (IsTag) Type() sql.Type {
	return sql.Boolean
}
