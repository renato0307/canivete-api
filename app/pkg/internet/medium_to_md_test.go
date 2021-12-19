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
package internet

import (
	"errors"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/internet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGin(serviceMock *internet.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(serviceMock, v1)

	return r
}

func TestPostFromUnix(t *testing.T) {
	output := internet.ConvertMediumToMdOutput{
		PostId:   "1638964800",
		Markdown: "# A pretty nice markdown",
	}

	// arrange
	serviceMock := internet.MockInterface{}
	serviceMock.On("ConvertMediumToMd", mock.Anything).Return(output, nil)

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	body := strings.NewReader(output.PostId)
	req, _ := http.NewRequest("POST", "/v1/internet/medium-to-md", body)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostFromUnixWithErrorFromCore(t *testing.T) {
	output := internet.ConvertMediumToMdOutput{}
	error := errors.New("fake error")

	// arrange
	serviceMock := internet.MockInterface{}
	serviceMock.On("ConvertMediumToMd", mock.Anything).Return(output, error)

	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	body := strings.NewReader(output.PostId)
	req, _ := http.NewRequest("POST", "/v1/internet/medium-to-md", body)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedError := apierrors.ApiError{Message: error.Error()}
	apiError, _ := apierrors.FromResponseRecorder(w)
	assert.Equal(t, expectedError, apiError)

}

func TestPostFromUnixWithErrorDueToNoBody(t *testing.T) {

	// arrange
	serviceMock := internet.MockInterface{}
	r := setupGin(&serviceMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/internet/medium-to-md", nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	apiError, _ := apierrors.FromResponseRecorder(w)
	assert.Equal(t, "request body is invalid", apiError.Message)
}
