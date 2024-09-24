# Gin Inspector

![Gin Inspector HTML Preview](https://raw.githubusercontent.com/krobus00/gin-inspector/master/assets/preview-html.png)

![Gin Inspector HTML Preview 2](https://raw.githubusercontent.com/krobus00/gin-inspector/master/assets/preview-html-2.jpg)

Gin middleware for investigating http request.

## Usage


```sh
$ go get github.com/krobus00/gin-inspector@v1.0.0
```

### Html Template

```
package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/krobus00/gin-inspector"
	"github.com/gin-gonic/gin"
)

func formatDate(t time.Time) string {
	return t.Format(time.RFC822)
}

func main() {
	r := gin.Default()
	r.Delims("{{", "}}")

	r.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
	})

	r.LoadHTMLFiles("inspector.html")
	debug := true

	if debug {
		r.Use(inspector.InspectorStats())

		r.GET("/_inspector", func(c *gin.Context) {
			c.HTML(http.StatusOK, "inspector.html", map[string]interface{}{
				"title":      "Gin Inspector",
				"pagination": inspector.GetPaginator(),
			})

		})
	}

	r.Run(":8080")
}

```
