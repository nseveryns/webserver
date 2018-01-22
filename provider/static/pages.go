package static

import "github.com/nseveryns/webserver/provider/page"

type StaticPage struct {
	page.Page
	file       string
	accessTime int64
}
