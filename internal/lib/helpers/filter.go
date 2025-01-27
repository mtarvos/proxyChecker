package helpers

import "proxyChecker/internal/entity"

func Cf(val interface{}, op entity.Operand) *entity.StringFilter {
	return &entity.StringFilter{
		Val: val,
		Op:  op,
	}
}
