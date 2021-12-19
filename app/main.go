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
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/canivete-api/pkg/datetime"
	"github.com/renato0307/canivete-api/pkg/finance"
	"github.com/renato0307/canivete-api/pkg/internet"
	"github.com/renato0307/canivete-api/pkg/programming"
	datetimecore "github.com/renato0307/canivete-core/pkg/datetime"
	financecore "github.com/renato0307/canivete-core/pkg/finance"
	internetcore "github.com/renato0307/canivete-core/pkg/internet"
	programmingcore "github.com/renato0307/canivete-core/pkg/programming"
)

func main() {

	r := gin.Default()
	err := r.SetTrustedProxies(nil) // https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	if err != nil {
		log.Fatalf("error setting trusted proxies to nil: %s\n", err.Error())
	}

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to canivete-api!")
	})

	v1 := r.Group("/v1")

	programmingService := programmingcore.Service{}
	programming.SetRouterGroup(&programmingService, v1)

	datetimeService := datetimecore.Service{}
	datetime.SetRouterGroup(&datetimeService, v1)

	financeService := financecore.Service{}
	finance.SetRouterGroup(&financeService, v1)

	internetService := internetcore.Service{}
	internet.SetRouterGroup(&internetService, v1)

	err = r.Run()
	if err != nil {
		log.Fatalf("error running gin: %s\n", err.Error())
	}
}
