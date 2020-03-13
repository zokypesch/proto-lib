package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	utils "github.com/zokypesch/proto-lib/utils"
)

type Rules string
type TagKey string

const (
	TagKeyType  TagKey = "type"
	TagKeyRules TagKey = "rules"

	EmptyString Rules = "EMPTYSTRING"
	ZeroValue   Rules = "ZEROVALUE"
	EmptyDate   Rules = "EMPTYDATE"
)

// Decorator for decorator your query
type Decorator struct {
	query *gorm.DB
	list  interface{}
}

//NewServiceDecorator new service query decorator
func NewServiceDecorator(db *gorm.DB, list interface{}) *Decorator {
	return &Decorator{db, list}
}

// AppendWhere for appending where your query
func (decor *Decorator) AppendWhere() (*gorm.DB, error) {
	if decor.list == nil {
		return decor.query, fmt.Errorf("Failed to generate append query")
	}

	var value reflect.Value
	var typ reflect.Type

	switch reflect.TypeOf(decor.list).Kind() {
	case reflect.Ptr:
		value = reflect.ValueOf(decor.list).Elem()
		typ = reflect.TypeOf(decor.list).Elem()
	case reflect.Struct:
		value = reflect.ValueOf(decor.list)
		typ = reflect.TypeOf(decor.list)
	default:
		return decor.query, fmt.Errorf("wrong list value")
	}
	newQuery := decor.query

	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i).Interface()
		fieldProp, _ := typ.FieldByName(fieldName)

		tag, ok := fieldProp.Tag.Lookup("decorator")

		// if decorator value is not represented
		if !ok {
			continue
		}

		// set of rules
		isEmptyString := false
		isZeroValue := false
		isEmptyDate := false

		// checking
		var tagValue string
		t := strings.Split(tag, ";")
		if len(t) > 1 {
			for _, v := range t {
				if parseTagKey(v) == TagKeyType {
					tagValue = getTagValue(v)
				} else if parseTagKey(v) == TagKeyRules {
					rulesValue := getRulesValue(v)
					if rulesValue == EmptyString {
						isEmptyString = true
					}
					if rulesValue == ZeroValue {
						isZeroValue = true
					}
					if rulesValue == EmptyDate {
						isEmptyDate = true
					}
				}
			}
		} else {
			tagValue = t[0]
		}

		// rules validation
		if !isEmptyString || !isZeroValue || !isEmptyDate {
			if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
				continue
			}
		}

		// add where value
		cond, args := decor.parseCondition(tagValue, fieldName, fieldValue)
		newQuery = newQuery.Where(cond, args)
	}
	return newQuery, nil
}

func (decor *Decorator) parseCondition(cond string, field string, value interface{}) (string, interface{}) {
	newField := utils.ConvertCamelToUnderscore(field)
	switch cond {
	case "EQUAL":
		return fmt.Sprintf("%s = ?", newField), fmt.Sprintf("%s", value)
	case "LIKE":
		return fmt.Sprintf("%s LIKE ?", newField), fmt.Sprintf("%%%s%%", value)
	case "NOTEQUAL":
		return fmt.Sprintf("%s <> ?", newField), fmt.Sprintf("%s", value)
	case "GREATERTHAN":
		return fmt.Sprintf("%s > ?", newField), fmt.Sprintf("%s", value)
	case "LESSTHAN":
		return fmt.Sprintf("%s < ?", newField), fmt.Sprintf("%s", value)
	case "GREATERTHANEQUAL":
		return fmt.Sprintf("%s >= ?", newField), fmt.Sprintf("%s", value)
	case "LESSTHANEQUAL":
		return fmt.Sprintf("%s <= ?", newField), fmt.Sprintf("%s", value)
	case "IN":
		values := strings.Split(value.(string), ",")
		return fmt.Sprintf("%s IN (?)", newField), values
	default:
		return "", "="
	}
}

func parseTagKey(s string) TagKey {
	if strings.Contains(s, string(TagKeyType)) {
		return TagKeyType
	} else {
		return TagKeyRules
	}
}

func getTagValue(s string) string {
	return strings.Split(s, "=")[1]
}

func getRulesValue(s string) Rules {
	v := getTagValue(s)
	return Rules(v)
}
