package knife4g

import (
	"embed"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var (
	//go:embed front
	front   embed.FS
	docJson []byte
	s       service
)

type Config struct {
	RelativePath string
}

type service struct {
	Name           string `json:"name"`
	Url            string `json:"url"`
	SwaggerVersion string `json:"swaggerVersion"`
	Location       string `json:"location"`
}

func init() {
	var err error
	docJson, err = os.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Println("no swagger.json found in ./docs")
	}
}

func Handler(config Config) gin.HandlerFunc {
	docJsonPath := config.RelativePath + "/docJson"
	servicesPath := config.RelativePath + "/front/service"
	docPath := config.RelativePath + "/index"
	appjsPath := config.RelativePath + "/front/webjars/js/app.42aa019b.js"

	s.Url = "/docJson"
	s.Location = "/docJson"
	s.Name = "API Documentation"
	s.SwaggerVersion = "2.0"

	appjsTemplate, err := template.New("app.42aa019b.js").
		Delims("{[(", ")]}").
		ParseFS(front, "front/webjars/js/app.42aa019b.js")
	if err != nil {
		log.Println(err)
	}
	docTemplate, err := template.New("doc.html").
		Delims("{[(", ")]}").
		ParseFS(front, "front/doc.html")
	if err != nil {
		log.Println(err)
	}

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}
		switch ctx.Request.RequestURI {
		case appjsPath:
			err := appjsTemplate.Execute(ctx.Writer, config)
			if err != nil {
				log.Println(err)
			}
		case servicesPath:
			ctx.JSON(http.StatusOK, []service{s})
		case docPath:
			err := docTemplate.Execute(ctx.Writer, config)
			if err != nil {
				log.Println(err)
			}
		case docJsonPath:
			ctx.Data(http.StatusOK, "application/json; charset=utf-8", docJson)
		default:
			ctx.FileFromFS(strings.TrimPrefix(ctx.Request.RequestURI, config.RelativePath), http.FS(front))
		}

	}
}
