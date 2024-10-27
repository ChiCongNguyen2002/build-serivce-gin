package binding

import (
	validate2 "build-service-gin/common/custom/validate"
	"github.com/gin-gonic/gin"
)

var (
	validate     = validate2.NewValidate()
	customBinder = &CustomBinder{}
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(c *gin.Context, i interface{}) error {
	if err := c.ShouldBind(i); err != nil {
		return err
	}

	if err := validate.ValidateStruct(i); err != nil {
		return err
	}
	return nil
}

func GetBinding() *CustomBinder {
	return customBinder
}
