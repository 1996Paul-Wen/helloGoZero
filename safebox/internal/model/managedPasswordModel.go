package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ManagedPasswordModel = (*customManagedPasswordModel)(nil)

type (
	// ManagedPasswordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customManagedPasswordModel.
	ManagedPasswordModel interface {
		managedPasswordModel
		withSession(session sqlx.Session) ManagedPasswordModel
	}

	customManagedPasswordModel struct {
		*defaultManagedPasswordModel
	}
)

// NewManagedPasswordModel returns a model for the database table.
func NewManagedPasswordModel(conn sqlx.SqlConn) ManagedPasswordModel {
	return &customManagedPasswordModel{
		defaultManagedPasswordModel: newManagedPasswordModel(conn),
	}
}

func (m *customManagedPasswordModel) withSession(session sqlx.Session) ManagedPasswordModel {
	return NewManagedPasswordModel(sqlx.NewSqlConnFromSession(session))
}
