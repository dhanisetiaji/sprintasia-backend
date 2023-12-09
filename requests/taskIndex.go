package requests

import (
	"gotham/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

type TaskIndexRequest struct {
	validation.Validatable `json:"-" form:"-" query:"-"`

	/**
	 * PathParams
	 */
	PathParams struct{}

	/**
	 * QueryParams
	 */
	QueryParams struct {
		utils.Order
		utils.Pagination
		utils.Find
	}

	/**
	 * Body
	 */
	Body struct{}
}

func (r TaskIndexRequest) Validate() error {
	return nil
}
