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
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/apierrors"
	"github.com/renato0307/canivete-api/pkg/logging"
	"github.com/renato0307/canivete-core/interface/programming"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger = logging.GetLogger()

func SetRouterGroup(p programming.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/programming")
	{
		programmingGroup.GET("/uuid", getUuid(p))
		programmingGroup.POST("/jwt-debugger", postJwtDebugger(p))
	}

	return programmingGroup
}

// getUuid handles the uuid request.
// It returns 200 on success.
func getUuid(p programming.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debugw("getting a new UUID")
		output := p.NewUuid()
		logger.Debugw("new UUID created", "uuid", output.UUID)
		c.JSON(http.StatusOK, output)
	}
}

// postJwtDebugger handles the jwt-debugger request.
// It returns 500 if anything fails and a 200 otherwise.
func postJwtDebugger(p programming.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, apierrors.ApiError{Message: "request body is invalid"})
			return
		}

		tokenString, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierrors.ApiError{Message: "request body is invalid"})
			return
		}

		output, err := p.DebugJwt(string(tokenString))
		if err != nil {
			logger.Debugw("error debugging a jwt", "error", err.Error())
			c.JSON(http.StatusBadRequest, apierrors.ApiError{Message: err.Error()})
			return
		}

		c.JSON(200, output)
	}
}
