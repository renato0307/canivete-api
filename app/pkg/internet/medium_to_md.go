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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/internet"
)

func SetRouterGroup(i internet.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/internet")
	{
		programmingGroup.POST("/medium-to-md", postConvertMediumToMd(i))
	}

	return programmingGroup
}

// postConvertMediumToMd handles the medium-to-md request.
// It returns:
//
// 200 (OK) if the request succeeded;
// 400 (BadRequest) if the post id is invalid;
// 500 (InternalServerError) if anything fails and a 200 otherwise.
func postConvertMediumToMd(i internet.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, apierrors.ApiError{Message: "request body is invalid"})
			return
		}

		postIdBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierrors.ApiError{Message: "error reading the body"})
			return
		}

		postId := strings.TrimRight(string(postIdBytes), "\n")
		output, err := i.ConvertMediumToMd(postId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierrors.ApiError{Message: err.Error()})
			return
		}

		c.JSON(200, output)
	}
}
