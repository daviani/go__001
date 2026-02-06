package calculator

import "errors"

func Add(a, b int) (int, error) {
	result := a + b
	if result < 0 {
		return 0, errors.New("retour négatif")
	}
	return result, nil
}
