package mongodb

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getPaginationOptions(pagination *model.PaginationQuery) *options.FindOptions {
	var opts *options.FindOptions

	if pagination != nil {
		opts = &options.FindOptions{
			Skip:  pagination.GetSkip(),
			Limit: pagination.GetLimit(),
		}
	}
	return opts
}
