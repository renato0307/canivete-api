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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-core/interface/datetime"
)

func SetRouterGroup(p datetime.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/datetime")
	{
		programmingGroup.POST("/fromunix", postFromUnix(p))
	}

	return programmingGroup
}

func postFromUnix(p datetime.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		unitTimestampBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierrors.ApiError{Message: "request body is invalid"})
			return
		}

		unitTimestampString := strings.TrimRight(string(unitTimestampBody), "\n")
		unixTimestamp, err := strconv.ParseInt(unitTimestampString, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierrors.ApiError{Message: "unix timestamp must be an integer number"})
			return
		}

		output := p.FromUnitTimestamp(unixTimestamp)
		c.JSON(http.StatusOK, output)
	}
}
