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
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
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
