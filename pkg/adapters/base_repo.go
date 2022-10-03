package adapters

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

var QueryMap = map[string]string{"Equal": "=", "Like": "LIKE", "Inc": "IN", "Lte": "<=", "Gte": ">="}

type BaseRepo struct{}

func (s *BaseRepo) QueryCondition(c Condition, tx *gorm.DB) *gorm.DB {
	if c.Fields != nil {
		tx = tx.Select(c.Fields)
	}
	valueOfCondition := reflect.ValueOf(c)

	for key, value := range QueryMap {
		field := valueOfCondition.FieldByName(key)
		if !field.IsValid() || field.IsNil() {
			continue
		}
		result := field.Interface()
		for k, v := range result.(map[string]interface{}) {
			tx = tx.Where(fmt.Sprintf("%s %s ?", k, value), v)
		}
	}

	if c.Not != nil {
		tx = tx.Not(c.Not)
	}
	if c.Or != nil {
		tx = tx.Or(c.Or)
	}
	if c.Limit != 0 {
		tx = tx.Offset(c.Offset).Limit(c.Limit)
	}
	return tx
}
