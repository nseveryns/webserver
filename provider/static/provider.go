package static

import (
	"io/ioutil"

	"time"

	"github.com/nseveryns/webserver/configuration"
	"github.com/nseveryns/webserver/provider/page"
)

type StaticProvider struct {
	configuration.Configuration
	Pages map[string]StaticPage
}

func (staticProvide *StaticProvider) ProvidePage(pageName string) *page.Page {
	p, ok := staticProvide.Pages[pageName]
	if ok {
		return &p.Page
	}
	p, err := staticProvide.loadPage(pageName)
	if err != nil {
		if pageName == staticProvide.ErrorPage {
			return &page.Page{[]byte{}} //Return empty page in case of loop
		}
		return staticProvide.ProvidePage(staticProvide.ErrorPage)
	}
	staticProvide.Pages[pageName] = p
	return &p.Page
}

func (provider StaticProvider) loadPage(page string) (StaticPage, error) {
	var err error
	static := StaticPage{}
	static.file = page
	static.Content, err = ioutil.ReadFile(provider.Directory + page)
	static.accessTime = time.Now().Unix()
	return static, err
}
