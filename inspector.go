package inspector

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Total       int           `json:"total"`
	TotalPage   int           `json:"total_page"`
	CurrentPage int           `json:"current_page"`
	PerPage     int           `json:"per_page"`
	HasNext     bool          `json:"has_next"`
	HasPrev     bool          `json:"has_prev"`
	NextPageUrl string        `json:"next_page_url"`
	PrevPageUrl string        `json:"prev_page_url"`
	Data        []RequestStat `json:"data"`
}

type RequestStat struct {
	RequestedAt   time.Time `json:"requested_at"`
	RequestUrl    string    `json:"request_url"`
	HttpMethod    string    `json:"http_method"`
	HttpStatus    int       `json:"http_status"`
	ContentType   string    `json:"content_type"`
	GetParams     any       `json:"get_params"`
	PostParams    any       `json:"post_params"`
	PostMultipart any       `json:"post_multipart"`
	Body          any       `json:"body"`
	ClientIP      string    `json:"client_ip"`
	Cookies       any       `json:"cookies"`
	Headers       any       `json:"headers"`
}

type AllRequests struct {
	Request []RequestStat `json:"requests"`
}

var allRequests = AllRequests{}
var pagination = Pagination{}

func GetPaginator() Pagination {
	return pagination
}

func InspectorStats(inspectorEndpoint string, multipartFormMaxMemory int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if urlPath == inspectorEndpoint {
			page, _ := strconv.ParseFloat(c.DefaultQuery("page", "1"), 64)
			perPage, _ := strconv.ParseFloat(c.DefaultQuery("per_page", "10"), 64)
			total := float64(len(allRequests.Request))
			totalPage := math.Ceil(total / perPage)
			offset := (page - 1) * perPage

			if offset < 0 {
				offset = 0
			}

			pagination.HasPrev = false
			pagination.HasNext = false
			pagination.CurrentPage = int(page)
			pagination.PerPage = int(perPage)
			pagination.TotalPage = int(totalPage)
			pagination.Total = int(total)
			pagination.Data = paginate(allRequests.Request, int(offset), int(perPage))

			if pagination.CurrentPage > 1 {
				pagination.HasPrev = true
				pagination.PrevPageUrl = urlPath + "?page=" + strconv.Itoa(pagination.CurrentPage-1) + "&per_page=" + strconv.Itoa(pagination.PerPage)
			}

			if pagination.CurrentPage < pagination.TotalPage {
				pagination.HasNext = true
				pagination.NextPageUrl = urlPath + "?page=" + strconv.Itoa(pagination.CurrentPage+1) + "&per_page=" + strconv.Itoa(pagination.PerPage)
			}

		} else {

			start := time.Now()
			var bodyBytes []byte

			if strings.EqualFold(c.Request.Header.Get("Content-Type"), "application/json") {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			c.Request.ParseForm()
			c.Request.ParseMultipartForm(multipartFormMaxMemory)

			c.Next()

			request := RequestStat{
				RequestedAt:   start,
				RequestUrl:    urlPath,
				HttpMethod:    c.Request.Method,
				HttpStatus:    c.Writer.Status(),
				ContentType:   c.ContentType(),
				Headers:       c.Request.Header,
				Cookies:       c.Request.Cookies(),
				GetParams:     c.Request.URL.Query(),
				PostParams:    c.Request.PostForm,
				PostMultipart: c.Request.MultipartForm,
				Body:          string(bodyBytes),
				ClientIP:      c.ClientIP(),
			}

			allRequests.Request = append([]RequestStat{request}, allRequests.Request...)
		}
	}
}

func paginate(s []RequestStat, offset, length int) []RequestStat {
	end := offset + length
	if end < len(s) {
		return s[offset:end]
	}

	return s[offset:]
}
