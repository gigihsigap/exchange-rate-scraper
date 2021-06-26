package helpers

import (
	"backend/app/model"
)

func Contains(arr []model.Exc_rates, str string) (int, bool) {
	for n, a := range arr {
		if a.Symbol == str {
			return n, false
		}
	}
	return 0, true
}
