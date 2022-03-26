package user

import (
	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/mocks"
	"github.com/nitesh-zs/bookshelf-api/model"
	"testing"
)

func user1() *model.User {
	return &model.User{
		ID:    uuid.New(),
		Email: "nitesh.saxena@zopsmart.com",
		Name:  "Nitesh",
		Type:  "admin",
	}
}

func initializeTest(t *testing.T) (*mocks.MockUserStore, *krogo.Context, svc) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mocks.NewMockUserStore(mockCtrl)

	k := krogo.New()
	ctx := krogo.NewContext(nil, nil, k)
	s := New(mock)

	return mock, ctx, s
}

func TestSvc_Exists(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	tests := []struct {
		desc      string
		email     string
		res       bool
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Exists",
			"nitesh.saxena@zopsmart.com",
			true,
			nil,
			[]*gomock.Call{
				mock.EXPECT().Exists(ctx, "nitesh.saxena@zopsmart.com").Return(true, nil),
			},
		},
		{
			"Not Exists",
			"abc@abc.com",
			false,
			errors.EntityNotFound{Entity: "user", ID: "abc@abc.com"},
			[]*gomock.Call{
				mock.EXPECT().Exists(ctx, "abc@abc.com").Return(false, errors.EntityNotFound{Entity: "user", ID: "abc@abc.com"}),
			},
		},
		{
			"Server error",
			"xyz@xyz.com",
			false,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().Exists(ctx, "xyz@xyz.com").Return(false, errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		exists, err := s.Exists(ctx, tc.email)

		assert.Equal(t, tc.res, exists, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	user1 := user1()
	user2 := &model.User{}

	tests := []struct {
		desc      string
		user      *model.User
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Success",
			user1,
			nil,
			[]*gomock.Call{
				mock.EXPECT().Create(ctx, user1).Return(nil),
			},
		},
		{
			"Server Error",
			user2,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().Create(ctx, user2).Return(errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		err := s.Create(ctx, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}
