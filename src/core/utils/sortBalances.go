package utils

import (
	"portfolio/schema"
	"sort"
)

func SortBasedOnBalance(tokenBalRes []schema.Token) {
	sort.Slice(tokenBalRes, func(i, j int) bool {
		return tokenBalRes[i].Balance.Cmp(&tokenBalRes[i].Balance) > 0
	})
}
func SortBasedOnValue(tokenBalRes []schema.Token) {
	sort.Slice(tokenBalRes, func(i, j int) bool {
		return tokenBalRes[i].Value.Cmp(&tokenBalRes[i].Value) > 0
	})
}
