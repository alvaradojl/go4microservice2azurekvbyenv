// Package keyvaultservice_test tests all functions of keyvaultservice
package keyvaultservice_test

import (
	"testing"

	"github.com/alvaradojl/go4microservice2azurekvbyenv/pkg/keyvault"
)

// Add_1And2_Return_3 tests that 1+2 equals 3
func Add_1And2_Return_3(t *testing.T) {
	//arrange
	num1 := 1
	num2 := 2
	expected := 3
	//act
	actual, err := keyvaultservice.Add(num1, num2)
	//assert
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Fail()
	}
}
