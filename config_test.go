package main

import (
	"testing"

	"reflect"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryConfig_AddRule(t *testing.T) {
	config := InMemoryConfig{}

	r1 := Rule{Priority: 3}
	r2 := Rule{Priority: 4}
	r3 := Rule{Priority: 5}
	r4 := Rule{Priority: 4}
	r5 := Rule{Priority: 2}
	config.AddRule(r1)
	config.AddRule(r2)
	config.AddRule(r3)
	config.AddRule(r4)
	config.AddRule(r5)

	assert.Equal(t, len(config.Rules), 5)
	assert.Equal(t, config.Rules, []Rule{r5, r1, r2, r4, r3})
}

func SetupTests() *gorm.DB {
	mocket.Catcher.Register()
	// GORM
	db, err := gorm.Open(mocket.DRIVER_NAME, "any_string")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	return db
}

func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Name
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

func TestSqliteConfig_AddRule(t *testing.T) {
	db := SetupTests()

	config := SqliteConfig{Db: db}
	config.AddRule(Rule{})
	config.AddRule(Rule{})

	ret := []map[string]interface{}{structToMap(&Rule{}), structToMap(&Rule{})}
	mocket.Catcher.NewMock().WithReply(ret)
	rules, _ := config.GetRules()

	assert.Equal(t, 2, len(rules))
}
