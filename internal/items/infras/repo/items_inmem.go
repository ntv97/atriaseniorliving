package repo

import (
	"context"

	"github.com/google/wire"
	"github.com/ntv97/atriaseniorliving/internal/items/domain"
)

var _ domain.ItemRepo = (*itemInMemRepo)(nil)

var RepositorySet = wire.NewSet(NewOrderRepo)

type itemInMemRepo struct {
	itemTypes map[string]*domain.ItemTypeDto
}

func NewOrderRepo() domain.ItemRepo {
	return &itemInMemRepo{
		itemTypes: map[string]*domain.ItemTypeDto{
			"COFFEE": {
				Name:  "COFFEE",
				Type:  0,
				Qty: 70,
				//Image: "img/COFFEE.png",
			},
			"COFFEE_DECAF": {
				Name:  "COFFEE_BLACK",
				Type:  1,
				Qty: 70,
				//Image: "img/COFFEE_DECAF.png",
			},
			"ICE_TEA": {
				Name:  "ICE_TEA",
				Type:  2,
				Qty: 70,
				//Image: "img/ICE_TEA.png",
			},
			"APPLE_JUICE": {
				Name:  "APPLE_JUICE",
				Type:  3,
				Qty: 70,
				//iImage: "img/APPLE_JUICE.png",

			},
			"COKE": {
                                Name:  "COKE",
                                Type:  4,
                                Qty: 70,
                                //iImage: "img/APPLE_JUICE.png",

                        },
			"CHICKEN_CORDON_BLEU": {
                                Name:  "CHICKEN_CORDON_BLUE",
                                Type:  5,
                                Qty: 70,
                                //iImage: "img/APPLE_JUICE.png",

                        },
			"TURKEY_SANDWICH": {
                                Name:  "TURKEY_SANDWICH",
                                Type:  6,
                                Qty: 70,
                                //iImage: "img/APPLE_JUICE.png",

                        },
			"PEPPERONI_PIZZA": {
                                Name:  "PEPPERONI_PIZZA",
                                Type:  7,
                                Qty: 70,
                                //iImage: "img/APPLE_JUICE.png",

                        },
			"": {
                                Name:  "TURKEY_SANDWICH",
                                Type:  6,
                                Qty: 70,
                                //iImage: "img/APPLE_JUICE.png",

                        },
		},
	}
}

func (p *itemInMemRepo) GetAll(ctx context.Context) ([]*domain.ItemTypeDto, error) {
	results := make([]*domain.ItemTypeDto, 0)

	for _, v := range p.itemTypes {
		results = append(results, &domain.ItemTypeDto{
			Name:  v.Name,
			Type:  v.Type,
			Qty:   v.Qty,
			Image: v.Image,
		})
	}

	return results, nil
}

func (p *itemInMemRepo) GetByTypes(ctx context.Context, itemTypes []string) ([]*domain.ItemDto, error) {
	results := make([]*domain.ItemDto, 0)

	for _, itemType := range itemTypes {
		item := p.itemTypes[itemType]
		if item.Name != "" {
			results = append(results, &domain.ItemDto{
				Type:  item.Type,
				Qty:   item.Qty,
			})
		}
	}

	return results, nil
}
