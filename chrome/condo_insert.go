package chrome

import (
	"context"

	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Condo) insert(ctx context.Context, collection *mongo.Collection) error {
	//d := bson.D{
	//	{
	//		"Name", c.Name,
	//	},
	//	{
	//		"address", c.address,
	//	},
	//	{
	//		"district", c.district,
	//	},
	//	{
	//		"tenure", c.tenure,
	//	},
	//	{
	//		"developer", c.developer,
	//	},
	//	{
	//		"url", c.url,
	//	},
	//	{
	//		"facility", c.facility,
	//	},
	//	{
	//		"facString", c.facString,
	//	},
	//}

	doc := mapper.ConvertStructToBSONMap(*c, nil)
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
