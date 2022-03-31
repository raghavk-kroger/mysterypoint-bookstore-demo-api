package book

import (
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
	"testing"
)

var book1 = &model.Book{
	ID:        uuid.New(),
	Title:     "Abc",
	Author:    "X",
	Summary:   "Lorem Ipsum",
	Genre:     "Action",
	Year:      2019,
	RegNum:    "ISB8726W821",
	Publisher: "saiudhiau",
	Language:  "Hebrew",
	ImageURI:  "jncj.ajcbiauadnc.com",
}

var bookRes1 = &model.BookRes{
	ID:        uuid.New(),
	Title:     "Abc",
	Author:    "X",
	Summary:   "Lorem Ipsum",
	Genre:     "Action",
	Year:      2019,
	Publisher: "saiudhiau",
	ImageURI:  "jncj.ajcbiauadnc.com",
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

func TestStore_Delete(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	id1 := uuid.New()
	tests := []struct {
		desc string
		id   uuid.UUID
		err  error
		exec *sqlmock.ExpectedExec
	}{
		{
			"Success",
			id1,
			nil,
			mock.ExpectExec(`DELETE FROM book WHERE id=$1`).WithArgs(id1).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			"DB error",
			id1,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectExec(`DELETE FROM book WHERE id=$1`).WillReturnError(errors.Error("DB Error")),
		},
	}

	for _, tc := range tests {
		err := s.Delete(ctx, tc.id)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestStore_Update(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	tests := []struct {
		desc string
		book *model.Book
		resp *model.BookRes
		err  error
		exec *sqlmock.ExpectedExec
	}{
		{
			"Success",
			book1,
			bookRes1,
			nil,
			mock.ExpectExec(getUpdateQuery(book1)).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			"DB error",
			book1,
			nil,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectExec(getUpdateQuery(book1)).WillReturnError(errors.Error("DB Error")),
		},
		{
			desc: "error",
			book: nil,
			resp: nil,
			err:  errors.Error("No object to update"),
		},
	}

	for _, tc := range tests {
		bookRes1, err := s.Update(ctx, tc.book)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, bookRes1, tc.desc)
	}
}

func TestStore_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	tests := []struct {
		desc string
		book *model.Book
		resp *model.BookRes
		err  error
		exec *sqlmock.ExpectedExec
	}{
		{
			"Success",
			book1,
			bookRes1,
			nil,
			mock.ExpectExec(createBook).WithArgs(book1.ID.String(), book1.Title, book1.Author, book1.Summary, book1.Genre, book1.Year, book1.RegNum, book1.Publisher, book1.Language, book1.ImageURI).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			"DB error",
			book1,
			nil,
			errors.Error("No object to update"),
			mock.ExpectExec(getUpdateQuery(book1)).WillReturnError(errors.Error("DB Error")),
		},
		{
			desc: "error",
			book: nil,
			resp: nil,
			err:  errors.Error("No object to update"),
		},
	}

	for _, tc := range tests {
		bookRes1, err := s.Update(ctx, tc.book)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, bookRes1, tc.desc)
	}
}
