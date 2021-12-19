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
package datetime

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/datetime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGin(serviceMock *datetime.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(serviceMock, v1)

	return r
}

func TestPostFromUnix(t *testing.T) {
	output := datetime.FromUnixTimestampOutput{
		UnixTimestamp: 1638964800,
		UtcTimestamp:  "Wed Dec  8 12:00:00 UTC 2021",
	}

	// arrange
	serviceMock := datetime.MockInterface{}
	serviceMock.On("FromUnitTimestamp", mock.Anything).Return(output, nil)

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	body := strings.NewReader(strconv.FormatInt(output.UnixTimestamp, 10))
	req, _ := http.NewRequest("POST", "/v1/datetime/fromunix", body)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostFromUnixWithCarriageReturn(t *testing.T) {
	output := datetime.FromUnixTimestampOutput{
		UnixTimestamp: 1638964800,
		UtcTimestamp:  "Wed Dec  8 12:00:00 UTC 2021",
	}

	// arrange
	serviceMock := datetime.MockInterface{}
	serviceMock.On("FromUnitTimestamp", mock.Anything).Return(output, nil)

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	body := strings.NewReader(strconv.FormatInt(output.UnixTimestamp, 10) + "\n")
	req, _ := http.NewRequest("POST", "/v1/datetime/fromunix", body)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostFromUnixWithInvalidTimestamp(t *testing.T) {
	// arrange
	error := errors.New("unix timestamp must be an integer number")
	serviceMock := datetime.MockInterface{}

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	body := strings.NewReader("invalid_timestamp")
	req, _ := http.NewRequest("POST", "/v1/datetime/fromunix", body)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	expectedError := apierrors.ApiError{Message: error.Error()}
	apiError := apierrors.FromResponseRecorder(w)
	assert.Equal(t, expectedError, apiError)
}
