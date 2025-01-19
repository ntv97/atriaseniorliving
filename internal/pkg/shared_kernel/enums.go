
package sharedkernel

import (
	"fmt"
)

type OrderSource int8

const (
	OrderSourceCounter OrderSource = iota
	OrderSourceWeb
)

type OrderTable int32

const (
        Table1 OrderTable = iota
        Table2
	Table3
	Table4
	Table5
	Table6
)

func (e OrderTable) String() string {
	return fmt.Sprintf("%d", int(e))
}

type Status int8

const (
	StatusPlaced Status = iota
	StatusInProcess
	StatusFulfilled
)

func (e Status) String() string {
	return fmt.Sprintf("%d", int(e))
}

type Location int8

const (
	LocationAtlanta Location = iota
	LocationCharlotte
	LocationRaleigh
)

func (e Location) String() string {
	return fmt.Sprintf("%d", int(e))
}

type CommandType int8

const (
	CommandTypePlaceOrder CommandType = iota
)

func (e CommandType) String() string {
	return fmt.Sprintf("%d", int(e))
}

type ItemType int32

const (
	ItemTypeCappuccino ItemType = iota
	ItemTypeCoffeeBlack
	ItemTypeCoffeeWithRoom
	ItemTypeEspresso
	ItemTypeEspressoDouble
	ItemTypeLatte
	ItemTypeCakePop
	ItemTypeCroissant
	ItemTypeMuffin
	ItemTypeCroissantChocolate
)

func (e ItemType) String() string {
	return []string{
		"CAPPUCCINO",
		"COFFEE_BLACK",
		"COFFEE_WITH_ROOM",
		"ESPRESSO",
		"ESPRESSO_DOUBLE",
		"LATTE",
		"CAKEPOP",
		"CROISSANT",
		"MUFFIN",
		"CROISSANT_CHOCOLATE",
		"CAPPUCCINO",
	}[e]
}

type OrderName int32

const (
        Bob OrderName = iota
        Andy
   	Betty
)

func (e OrderName) String() string {
        return []string{
                "BOB",
                "ANDY",
                "BETTY",
        }[e]
}


type OrderType int32

const (
        Cook OrderType = iota
        Chef
        Waitstaff
)

func (e OrderType) String() string {
        return []string{
                "COOK",
                "CHEF",
                "WAITSTAFF",
        }[e]
}
