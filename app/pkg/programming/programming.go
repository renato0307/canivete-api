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
	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-core/interface/programming"
)

func SetRouterGroup(p programming.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/programming")
	{
		programmingGroup.GET("/uuid", getUuid(p))
	}

	return programmingGroup
}

func getUuid(p programming.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		output := p.NewUuid()
		c.JSON(200, output)
	}
}
