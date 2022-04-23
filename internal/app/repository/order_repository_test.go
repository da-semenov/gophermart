package repository

import (
	"context"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestOrderRepository_Save(t *testing.T) {
	tests := []struct {
		name    string
		order   models.Order
		wantErr bool
	}{
		{
			name: "OrderRepository. Save. Test 1",
			order: models.Order{
				ID:        0,
				UserID:    1,
				Num:       "11",
				Status:    "STATUS",
				UploadAt:  time.Now().Truncate(time.Microsecond),
				UpdatedAt: time.Now().Truncate(time.Microsecond),
			},
			wantErr: false,
		},
	}
	target, _ := NewOrderRepository(postgresHandler, Log)
	initDatabase(context.Background(), postgresHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save object")
			err := target.Save(context.Background(), &tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("get object")
			res, err := target.GetByNum(context.Background(), tt.order.Num)
			if err != nil {
				t.Errorf("GetByNum() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("check got object")
			assert.Equal(t, tt.order.Status, res.Status, "Compare error (order.Status): expected %s,  got %s", tt.order.Status, res.Status)
			assert.Equal(t, tt.order.UserID, res.UserID, "Compare error (order.UserID): expected %s,  got %s", tt.order.UserID, res.UserID)
			assert.Equal(t, tt.order.UpdatedAt, res.UpdatedAt, "Compare error (order.UpdatedAt): expected %s,  got %s", tt.order.UpdatedAt.String(), res.UpdatedAt.String())
			assert.Equal(t, tt.order.UploadAt, res.UploadAt, "Compare error (order.UploadAt): expected %s,  got %s", tt.order.UploadAt, res.UploadAt)
		})
	}
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		order     models.Order
		statusNew string
		wantErr   bool
	}{
		{
			name: "OrderRepository. Update.Test 1",
			order: models.Order{
				ID:        2,
				UserID:    2,
				Num:       "21",
				Status:    "STATUS",
				UploadAt:  time.Now().Truncate(time.Microsecond),
				UpdatedAt: time.Now().Truncate(time.Microsecond),
			},
			statusNew: "STATUS_NEW",
			wantErr:   false,
		},
	}
	initDatabase(context.Background(), postgresHandler)
	target, _ := NewOrderRepository(postgresHandler, Log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save object")
			if err := target.Save(context.Background(), &tt.order); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("update saved object")
			if err := target.UpdateStatus(context.Background(), tt.order.UserID, tt.order.Num, tt.statusNew); (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println("get object")
			res, err := target.GetByNum(context.Background(), tt.order.Num)
			if err != nil {
				t.Errorf("GetByNum() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("check got object")
			assert.Equal(t, tt.statusNew, res.Status, "Compare error (order.status): expected %s,  got %s", tt.statusNew, res.Status)
		})
	}
}

func TestOrderRepository_FindByUser(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		objCount  int
		objStatus string
		create    bool
		wantErr   bool
	}{
		{
			name:      "OrderRepository. FindByUser. Test 1",
			userID:    31,
			objCount:  2,
			objStatus: "NEW",
			wantErr:   false,
		},
		{
			name:      "OrderRepository. FindByUser. Test 2",
			userID:    32,
			objCount:  0,
			objStatus: "NEW",
			wantErr:   false,
		},
	}
	initDatabase(context.Background(), postgresHandler)
	target, _ := NewOrderRepository(postgresHandler, Log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save objects")
			for i := 0; i < tt.objCount; i++ {
				var obj models.Order
				obj.Num = strconv.Itoa(i)
				obj.UserID = tt.userID
				obj.Status = tt.objStatus
				timeLabel := time.Now().Truncate(time.Microsecond)
				obj.UpdatedAt = timeLabel
				obj.UploadAt = timeLabel

				if err := target.Save(context.Background(), &obj); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			fmt.Println("get saved objects")
			resArr, err := target.FindByUser(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println("check got objects")
			assert.Equal(t, tt.objCount, len(resArr), "Compare error (order.status): expected %s,  got %s", tt.objCount, len(resArr))
		})
	}
}

func TestOrderRepository_FindByNum(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		objCount  int
		objStatus string
		targetNum string
		wantErr   bool
		error     error
	}{
		{
			name:      "OrderRepository. FindByNum. Test 1",
			userID:    31,
			objCount:  2,
			objStatus: "NEW",
			targetNum: "1",
			wantErr:   false,
		},
		{
			name:      "OrderRepository. FindByNum. Test 2",
			userID:    32,
			objCount:  0,
			objStatus: "NEW",
			targetNum: "5i0ew890suf90g0-",
			wantErr:   true,
			error:     &models.NoRowFound,
		},
	}
	initDatabase(context.Background(), postgresHandler)
	target, _ := NewOrderRepository(postgresHandler, Log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save objects")
			for i := 0; i < tt.objCount; i++ {
				var obj models.Order
				obj.Num = strconv.Itoa(i)
				obj.UserID = tt.userID
				obj.Status = tt.objStatus
				timeLabel := time.Now().Truncate(time.Microsecond)
				obj.UpdatedAt = timeLabel
				obj.UploadAt = timeLabel

				if err := target.Save(context.Background(), &obj); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			fmt.Println("get saved objects")
			res, err := target.GetByNum(context.Background(), tt.targetNum)
			if (err != nil) != tt.wantErr && err != tt.error {
				t.Errorf("FindByUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println("check got objects")
			fmt.Println(res)
		})
	}
}
