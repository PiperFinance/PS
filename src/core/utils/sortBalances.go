package utils

import (
	"portfolio/core/schema"
	"sort"
)

func SortBasedOnBalance(tokenBalRes []schema.TokenBalance) {
	sort.Slice(tokenBalRes, func(i, j int) bool {
		return tokenBalRes[i].Balance.Cmp(&tokenBalRes[i].Balance) > 0
	})
}
