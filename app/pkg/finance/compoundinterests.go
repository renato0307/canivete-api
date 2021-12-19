/*
Copyright © 2021 Renato Torres <renato.torres@pm.me>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package finance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/finance"
)

type calculateCompoundInterestsInput struct {
	InterestRate               float64 `validate:"required"`
	CompoundPeriods            float64 `validate:"required"`
	InvestAmount               float64 `validate:"required"`
	RegularContributions       float64
	RegularContributionsPeriod float64 `validate:"gt=0,required_with=RegularContributions"`
	Time                       float64 `validate:"required"`
}

func SetRouterGroup(f finance.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/finance")
	{
		programmingGroup.POST("/calculate-compound-interests", postCalculateCompoundInterests(f))
	}

	return programmingGroup
}

func postCalculateCompoundInterests(f finance.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierrors.ApiError{Message: "unexpected error reading the body"})
			return
		}

		input := calculateCompoundInterestsInput{}
		err = json.Unmarshal(body, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierrors.ApiError{Message: "request body is invalid: " + err.Error()})
			return
		}

		validate := validator.New()
		err = validate.Struct(input)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierrors.ApiError{Message: err.Error()})
			return
		}

		output, err := f.CalculateCompoundInterests(
			input.InvestAmount,
			input.CompoundPeriods,
			input.Time,
			input.RegularContributions,
			input.RegularContributionsPeriod,
			input.InterestRate,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierrors.ApiError{Message: "unexpected error calculating interests: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
