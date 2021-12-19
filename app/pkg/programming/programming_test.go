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
package programming

import (
	"errors"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/programming"
	programmingmocks "github.com/renato0307/canivete-core/interface/programming/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUuid(t *testing.T) {
	// arrange
	uuid := "d967aaad-1df5-485d-96b4-43d4247972e7"
	output := programming.UuidOutput{UUID: uuid}

	serviceMock := programmingmocks.Interface{}
	serviceMock.On("NewUuid", mock.Anything).Return(output, nil)

	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(&serviceMock, v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/programming/uuid", nil)
	req.Header.Add("Content-Type", "application/json")

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestPostJwtDebugger(t *testing.T) {
	// arrange
	tokenString := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYzOTgyODY0NiwiZXhwIjoxNjM5ODMyMjQ2fQ.ujQ7wTsos4hYgipdnxSjLICDdfSLq9pYbpwS0WvUKc4"

	output := programming.JwtDebuggerOutput{
		Header: map[string]interface{}{
			"alg": "HS256",
			"typ": "JWT",
		},
		Payload: map[string]interface{}{
			"admin": true,
			"exp":   1639832246,
			"iat":   1639828646,
			"name":  "John Doe",
			"sub":   "1234567890",
		},
	}

	serviceMock := programmingmocks.Interface{}
	serviceMock.On("DebugJwt", mock.Anything).Return(output, nil)

	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(&serviceMock, v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/jwt-debugger", strings.NewReader(tokenString))
	req.Header.Add("Content-Type", "application/json")

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostJwtDebuggerNoBodyShouldReturn500(t *testing.T) {
	// arrange
	serviceMock := programmingmocks.Interface{}

	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(&serviceMock, v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/jwt-debugger", nil)
	req.Header.Add("Content-Type", "application/json")

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedError := apierrors.ApiError{Message: "request body is invalid"}
	apiError := apierrors.FromResponseRecorder(w)
	assert.Equal(t, expectedError, apiError)
}

func TestPostJwtDebuggerShouldReturn500IfCoreFails(t *testing.T) {
	// arrange
	tokenString := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYzOTgyODY0NiwiZXhwIjoxNjM5ODMyMjQ2fQ.ujQ7wTsos4hYgipdnxSjLICDdfSLq9pYbpwS0WvUKc4"
	output := programming.JwtDebuggerOutput{}
	error := errors.New("fake error")

	serviceMock := programmingmocks.Interface{}
	serviceMock.On("DebugJwt", mock.Anything).Return(output, error)

	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(&serviceMock, v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/jwt-debugger", strings.NewReader(tokenString))
	req.Header.Add("Content-Type", "application/json")

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedError := apierrors.ApiError{Message: error.Error()}
	apiError := apierrors.FromResponseRecorder(w)
	assert.Equal(t, expectedError, apiError)
}
