// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q               = new(Query)
	Milestone       *milestone
	SchemaMigration *schemaMigration
	Tag             *tag
	User            *user
	UserTag         *userTag
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Milestone = &Q.Milestone
	SchemaMigration = &Q.SchemaMigration
	Tag = &Q.Tag
	User = &Q.User
	UserTag = &Q.UserTag
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:              db,
		Milestone:       newMilestone(db, opts...),
		SchemaMigration: newSchemaMigration(db, opts...),
		Tag:             newTag(db, opts...),
		User:            newUser(db, opts...),
		UserTag:         newUserTag(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Milestone       milestone
	SchemaMigration schemaMigration
	Tag             tag
	User            user
	UserTag         userTag
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:              db,
		Milestone:       q.Milestone.clone(db),
		SchemaMigration: q.SchemaMigration.clone(db),
		Tag:             q.Tag.clone(db),
		User:            q.User.clone(db),
		UserTag:         q.UserTag.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:              db,
		Milestone:       q.Milestone.replaceDB(db),
		SchemaMigration: q.SchemaMigration.replaceDB(db),
		Tag:             q.Tag.replaceDB(db),
		User:            q.User.replaceDB(db),
		UserTag:         q.UserTag.replaceDB(db),
	}
}

type queryCtx struct {
	Milestone       IMilestoneDo
	SchemaMigration ISchemaMigrationDo
	Tag             ITagDo
	User            IUserDo
	UserTag         IUserTagDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Milestone:       q.Milestone.WithContext(ctx),
		SchemaMigration: q.SchemaMigration.WithContext(ctx),
		Tag:             q.Tag.WithContext(ctx),
		User:            q.User.WithContext(ctx),
		UserTag:         q.UserTag.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
