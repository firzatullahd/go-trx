package service

import (
	"context"
	"database/sql"
	"go-trx/config"
	accountModel "go-trx/domain/account/model"
	aRepository "go-trx/domain/account/repository"
	tError "go-trx/domain/transaction/error"
	"go-trx/domain/transaction/model"
	"go-trx/domain/transaction/repository"
	"testing"

	mockAccountRepo "go-trx/domain/account/repository/mock"
	mockRepo "go-trx/domain/transaction/repository/mock"

	"github.com/magiconair/properties/assert"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"
)

func Test_service_InsertTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redisrepo := mockRepo.NewMockRedisRepository(ctrl)
	repo := mockRepo.NewMockRepository(ctrl)
	accountRepo := mockAccountRepo.NewMockRepository(ctrl)

	type fields struct {
		conf        config.Config
		repo        repository.Repository
		accountRepo aRepository.Repository
		redisRepo   repository.RedisRepository
	}
	type args struct {
		ctx     context.Context
		payload model.NewTransaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		domock  func()
	}{
		{
			name: "success credit",
			fields: fields{
				conf: config.Config{
					Constant: config.Constant{
						TrxTTL: 24,
					},
				},
				repo:        repo,
				accountRepo: accountRepo,
				redisRepo:   redisrepo,
			},
			args: args{
				ctx: context.Background(),
				payload: model.NewTransaction{
					UserID:          4096,
					Amount:          decimal.New(-10, 4),
					Remark:          "cashout",
					TransactionType: model.TransactionCredit,
					ReferenceKey:    "1",
				},
			},
			wantErr: nil,
			domock: func() {
				accountRepo.EXPECT().AccountBalance(gomock.Any(), gomock.Any()).Return(accountModel.Account{
					ID:      127,
					UserID:  4096,
					Balance: decimal.New(10, 4),
				}, nil).Times(1)
				redisrepo.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
				repo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "success debit",
			fields: fields{
				conf: config.Config{
					Constant: config.Constant{
						TrxTTL: 24,
					},
				},
				repo:        repo,
				accountRepo: accountRepo,
				redisRepo:   redisrepo,
			},
			args: args{
				ctx: context.Background(),
				payload: model.NewTransaction{
					UserID:          4096,
					Amount:          decimal.New(10, 4),
					Remark:          "topup",
					TransactionType: model.TransactionDebit,
					ReferenceKey:    "1",
				},
			},
			wantErr: nil,
			domock: func() {
				accountRepo.EXPECT().AccountBalance(gomock.Any(), gomock.Any()).Return(accountModel.Account{
					ID:      127,
					UserID:  4096,
					Balance: decimal.New(10, 4),
				}, nil).Times(1)
				redisrepo.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
				repo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "success new user",
			fields: fields{
				conf: config.Config{
					Constant: config.Constant{
						TrxTTL: 24,
					},
				},
				repo:        repo,
				accountRepo: accountRepo,
				redisRepo:   redisrepo,
			},
			args: args{
				ctx: context.Background(),
				payload: model.NewTransaction{
					UserID:          4096,
					Amount:          decimal.New(10, 4),
					Remark:          "topup",
					TransactionType: model.TransactionDebit,
					ReferenceKey:    "1",
				},
			},
			wantErr: nil,
			domock: func() {
				accountRepo.EXPECT().AccountBalance(gomock.Any(), gomock.Any()).Return(accountModel.Account{}, sql.ErrNoRows).Times(1)
				redisrepo.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
				repo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "fail duplicate transaction",
			fields: fields{
				conf: config.Config{
					Constant: config.Constant{
						TrxTTL: 24,
					},
				},
				repo:        repo,
				accountRepo: accountRepo,
				redisRepo:   redisrepo,
			},
			args: args{
				ctx: context.Background(),
				payload: model.NewTransaction{
					UserID:          4096,
					Amount:          decimal.New(10, 4),
					Remark:          "topup",
					TransactionType: model.TransactionDebit,
					ReferenceKey:    "1",
				},
			},
			wantErr: tError.ErrDuplicateTrx,
			domock: func() {
				accountRepo.EXPECT().AccountBalance(gomock.Any(), gomock.Any()).Return(accountModel.Account{
					ID:      127,
					UserID:  4096,
					Balance: decimal.New(10, 4),
				}, nil).Times(1)
				redisrepo.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
			},
		},
		{
			name: "fail balance insufficient",
			fields: fields{
				conf: config.Config{
					Constant: config.Constant{
						TrxTTL: 24,
					},
				},
				repo:        repo,
				accountRepo: accountRepo,
				redisRepo:   redisrepo,
			},
			args: args{
				ctx: context.Background(),
				payload: model.NewTransaction{
					UserID:          4096,
					Amount:          decimal.New(-11, 4),
					Remark:          "topup",
					TransactionType: model.TransactionCredit,
					ReferenceKey:    "1",
				},
			},
			wantErr: tError.ErrBalanceInsufficient,
			domock: func() {
				accountRepo.EXPECT().AccountBalance(gomock.Any(), gomock.Any()).Return(accountModel.Account{
					ID:      127,
					UserID:  4096,
					Balance: decimal.New(10, 4),
				}, nil).Times(1)
				redisrepo.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				conf:        tt.fields.conf,
				repo:        tt.fields.repo,
				accountRepo: tt.fields.accountRepo,
				redisRepo:   tt.fields.redisRepo,
			}
			tt.domock()
			err := s.InsertTransaction(tt.args.ctx, tt.args.payload)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}
