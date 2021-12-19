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
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/finance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGin(serviceMock *finance.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(serviceMock, v1)

	return r
}

func TestCalculateCompoundInterests(t *testing.T) {
	output := finance.CompoundInterestsOutput{
		Total: finance.CompoundInterestsDetailOutput{
			FinalAmount:        8457.76,
			TotalContributions: 7400,
			Interests:          1057.77,
		},
		History: []finance.CompoundInterestsHistoryEntryOutput{
			{
				Period: "1",
				Totals: finance.CompoundInterestsDetailOutput{
					FinalAmount:        6660,
					Interests:          460,
					TotalContributions: 6200,
				},
			},
			{
				Period: "2",
				Totals: finance.CompoundInterestsDetailOutput{
					FinalAmount:        8457.76,
					Interests:          1057.77,
					TotalContributions: 7400,
				},
			},
		},
	}

	// arrange
	serviceMock := finance.MockInterface{}
	mockCall := serviceMock.On(
		"CalculateCompoundInterests",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)
	mockCall.Return(output, nil)

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()

	input := calculateCompoundInterestsInput{
		InterestRate:               8,
		CompoundPeriods:            12,
		InvestAmount:               5000,
		Time:                       2,
		RegularContributions:       100,
		RegularContributionsPeriod: 12,
	}
	body, _ := json.Marshal(&input)
	req, _ := http.NewRequest("POST", "/v1/finance/calculate-compound-interests", bytes.NewReader(body))

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)

	result := finance.CompoundInterestsOutput{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(t, err, "invalid body returned")
	assert.Equal(t, output, result)
}

func TestCalculateCompoundBodyMissing(t *testing.T) {
	// arrange
	error := errors.New("request body is invalid: unexpected end of JSON input")
	serviceMock := finance.MockInterface{}

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/finance/calculate-compound-interests", strings.NewReader(""))

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	expectedError := apierrors.ApiError{Message: error.Error()}
	apiError, _ := apierrors.FromResponseRecorder(w)
	assert.Equal(t, expectedError, apiError)
}

func TestCalculateCompoundBodyMissingRequired(t *testing.T) {
	// arrange
	serviceMock := finance.MockInterface{}
	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/finance/calculate-compound-interests", strings.NewReader("{}"))

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	apiError, _ := apierrors.FromResponseRecorder(w)
	assert.Contains(t, apiError.Message, "failed on the 'required'")
}

func TestCalculateCompoundInterestsOnCoreError(t *testing.T) {
	output := finance.CompoundInterestsOutput{}

	// arrange
	serviceMock := finance.MockInterface{}
	mockCall := serviceMock.On(
		"CalculateCompoundInterests",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)
	mockCall.Return(output, errors.New("unexpected error calculating interest"))

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()

	input := calculateCompoundInterestsInput{
		InterestRate:               8,
		CompoundPeriods:            12,
		InvestAmount:               5000,
		Time:                       2,
		RegularContributions:       100,
		RegularContributionsPeriod: 12,
	}
	body, _ := json.Marshal(&input)
	req, _ := http.NewRequest("POST", "/v1/finance/calculate-compound-interests", bytes.NewReader(body))

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	apiError, _ := apierrors.FromResponseRecorder(w)
	assert.Contains(t, apiError.Message, "unexpected error calculating interest")
}

// func TestPostFromUnixWithCarriageReturn(t *testing.T) {
// 	output := datetime.FromUnixTimestampOutput{
// 		UnixTimestamp: 1638964800,
// 		UtcTimestamp:  "Wed Dec  8 12:00:00 UTC 2021",
// 	}

// 	// arrange
// 	serviceMock := datetime.MockInterface{}
// 	serviceMock.On("FromUnitTimestamp", mock.Anything).Return(output, nil)

// 	r := setupGin(&serviceMock)
// 	w := httptest.NewRecorder()
// 	body := strings.NewReader(strconv.FormatInt(output.UnixTimestamp, 10) + "\n")
// 	req, _ := http.NewRequest("POST", "/v1/datetime/fromunix", body)

// 	// act
// 	r.ServeHTTP(w, req)

// 	// assert
// 	assert.Equal(t, http.StatusOK, w.Code)
// }

// func TestPostFromUnixWithInvalidTimestamp(t *testing.T) {
// 	// arrange
// 	error := errors.New("unix timestamp must be an integer number")
// 	serviceMock := datetime.MockInterface{}

// 	r := setupGin(&serviceMock)
// 	w := httptest.NewRecorder()
// 	body := strings.NewReader("invalid_timestamp")
// 	req, _ := http.NewRequest("POST", "/v1/datetime/fromunix", body)

// 	// act
// 	r.ServeHTTP(w, req)

// 	// assert
// 	assert.Equal(t, http.StatusBadRequest, w.Code)

// 	expectedError := apierrors.ApiError{Message: error.Error()}
// 	apiError := apierrors.FromResponseRecorder(w)
// 	assert.Equal(t, expectedError, apiError)
// }
