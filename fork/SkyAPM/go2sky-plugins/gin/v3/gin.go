// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//Package v3 (Gin) is a HTTP web framework written in
//Go (Golang) plugin which can be used for integration with Gin http server.
package v3

import (
	"fmt"
	"strconv"
	"time"

	"github.com/767829413/normal-frame/fork/SkyAPM/go2sky/propagation"

	"github.com/767829413/normal-frame/fork/SkyAPM/go2sky"
	v3 "github.com/767829413/normal-frame/fork/SkyAPM/go2sky/reporter/grpc/language-agent"

	"github.com/gin-gonic/gin"
)

const componentIDGINHttpServer = 5006

//Middleware gin middleware return HandlerFunc  with tracing.
func Middleware(engine *gin.Engine, tracer *go2sky.Tracer) gin.HandlerFunc {
	if engine == nil || tracer == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		span, ctx, traceID, err := tracer.CreateEntrySpan(c.Request.Context(), getOperationName(c), func(key string) (string, error) {
			return c.Request.Header.Get(key), nil
		})
		if len(traceID) > 0 {
			c.Set(propagation.Header, traceID)
		} else {
			c.Set(propagation.Header, c.Request.Header.Get(propagation.Header))
		}
		if err != nil {
			c.Next()
			return
		}
		span.SetComponent(componentIDGINHttpServer)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.Tag(go2sky.TagHTTPUserAgent, c.Request.UserAgent())
		span.SetSpanLayer(v3.SpanLayer_Http)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			span.Error(time.Now(), c.Errors.String())
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
		span.End()
	}
}

func getOperationName(c *gin.Context) string {
	return fmt.Sprintf("/%s%s", c.Request.Method, c.FullPath())
}
