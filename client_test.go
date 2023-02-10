package pingen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	//good case
	createData := &CreateData{}
	createData.Data.Attributes.FileOriginalName = "DeleteTest.pdf"
	createData.Data.Attributes.AddressPosition = "left"
	createData.Data.Attributes.AutoSend = false

	err := validateCreateData(createData)
	assert.NoError(t, err)
	assert.Equal(t, createData.Data.Attributes.DeliveryProduct, "cheap")
	assert.Equal(t, createData.Data.Attributes.PrintSpectrum, "grayscale")
	assert.Equal(t, createData.Data.Attributes.PrintMode, "simplex")

	//bad case
	createData = &CreateData{}
	createData.Data.Attributes.FileOriginalName = "DeleteTest.pdf"
	createData.Data.Attributes.AddressPosition = "wrong"
	createData.Data.Attributes.AutoSend = false

	err = validateCreateData(createData)
	assert.Error(t, err)
}
