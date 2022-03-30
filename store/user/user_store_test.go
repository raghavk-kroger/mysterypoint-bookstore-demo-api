package user

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmizerany/assert"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/krogertechnology/krogo/pkg/datastore"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/krogertechnology/krogo/pkg/krogo/config"
	"github.com/krogertechnology/krogo/pkg/log"
	"github.com/nitesh-zs/bookshelf-api/model"
)

func user1() *model.User {
	return &model.User{
		ID:    uuid.Nil,
		Email: "nitesh.saxena@zopsmart.com",
		Name:  "Nitesh",
		Type:  "admin",
	}
}

func initializeTest(t *testing.T) (sqlmock.Sqlmock, *krogo.Context, store) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("error in creating mockDB: %v", err)
	}

	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../../configs")
	k := krogo.NewWithConfig(c)
	k.ORM = nil
	gormDB, _ := gorm.Open("postgres", db)

	k.SetORM(datastore.GORMClient{DB: gormDB})

	ctx := krogo.NewContext(nil, nil, k)
	s := New()

	return mock, ctx, s
}

func TestStore_GetByEmail(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	user := user1()

	row := sqlmock.NewRows([]string{"id", "email", "name", "type"}).
		AddRow(uuid.Nil, "nitesh.saxena@zopsmart.com", "Nitesh", "admin")

	tests := []struct {
		desc  string
		email string
		res   *model.User
		err   error
		query *sqlmock.ExpectedQuery
	}{
		{
			"Exists",
			"nitesh.saxena@zopsmart.com",
			user,
			nil,
			mock.ExpectQuery(getUserByEmail).WithArgs("nitesh.saxena@zopsmart.com").WillReturnRows(row),
		},
		{
			"Not Exists",
			"abc@abc.com",
			nil,
			errors.EntityNotFound{Entity: "user", ID: "abc@abc.com"},
			mock.ExpectQuery(getUserByEmail).WithArgs("abc@abc.com").WillReturnError(sql.ErrNoRows),
		},
		{
			"DB error",
			"xyz@xyz.com",
			nil,
			errors.DB{Err: errors.Error("DB error")},
			mock.ExpectQuery(getUserByEmail).WithArgs("xyz@xyz.com").WillReturnError(errors.Error("DB error")),
		},
	}

	for _, tc := range tests {
		user, err := s.GetByEmail(ctx, tc.email)
		assert.Equal(t, tc.res, user, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestStore_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	user1 := user1()

	tests := []struct {
		desc string
		user *model.User
		err  error
		exec *sqlmock.ExpectedExec
	}{
		{
			"Success",
			user1,
			nil,
			mock.ExpectExec(createUser).WithArgs(sqlmock.AnyArg(), user1.Email, user1.Name, user1.Type).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			"DB error",
			&model.User{},
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectExec(createUser).WillReturnError(errors.Error("DB Error")),
		},
	}

	for _, tc := range tests {
		err := s.Create(ctx, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}
