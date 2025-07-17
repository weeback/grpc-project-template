package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DbStruct struct {
	ObjectId  bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updated_at"`
}

func (ins *DbStruct) SetCreatedTime() {
	ins.CreatedAt = time.Now()
	ins.UpdatedAt = time.Now()
}

type DbQueryParams struct {
	Keyword      string    `json:"keyword"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	SortBy       string    `json:"sortBy"`
	SortOrder    string    `json:"sortOrder"`
	Limit        int64     `json:"limit"`
	LastObjectId string    `json:"lastObjectId"`
}
