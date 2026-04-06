package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"
)

//go:embed templates/*.html
var content embed.FS

var indexTmpl *template.Template

func init() {
	var err error
	funcMap := template.FuncMap{
		"formatLatency": func(d time.Duration) string {
			if d == 0 {
				return "n/a"
			}
			return fmt.Sprintf("%dms", d.Milliseconds())
		},
	}

	indexTmpl, err = template.New("index.html").Funcs(funcMap).ParseFS(content, "templates/*.html")
	if err != nil {
		panic(err)
	}
}

type PageData struct {
	Version                    string
	Host                       string
	Port                       string
	CheckInterval              int
	IPCheckUrl                 string
	SimulateLatency            bool
	CheckMethod                string
	StatusCheckUrl             string
	DownloadUrl                string
	Timeout                    int
	SubscriptionUpdate         bool
	SubscriptionUpdateInterval int
	StartPort                  int
	Instance                   string
	PushUrl                    string
	Endpoints                  []EndpointInfo
	ShowServerDetails          bool
	IsPublic                   bool
	SubscriptionName           string
	ProxiesJSON                template.JS
}

type dashboardProxy struct {
	Name       string `json:"name"`
	StableID   string `json:"stableId"`
	ServerInfo string `json:"serverInfo,omitempty"`
	ProxyPort  int    `json:"proxyPort,omitempty"`
	URL        string `json:"url,omitempty"`
	Index      int    `json:"index"`
	Status     bool   `json:"status"`
	Latency    string `json:"latency"`
	LatencyMs  int64  `json:"latencyMs"`
}

func BuildProxiesJSON(endpoints []EndpointInfo, showServerDetails bool, isPublic bool) (template.JS, error) {
	proxies := make([]dashboardProxy, 0, len(endpoints))
	for _, ep := range endpoints {
		proxy := dashboardProxy{
			Name:      ep.Name,
			StableID:  ep.StableID,
			Index:     ep.Index,
			Status:    ep.Status,
			Latency:   formatLatency(ep.Latency),
			LatencyMs: ep.Latency.Milliseconds(),
		}
		if showServerDetails {
			proxy.ServerInfo = ep.ServerInfo
			proxy.ProxyPort = ep.ProxyPort
		}
		if !isPublic {
			proxy.URL = ep.URL
		}
		proxies = append(proxies, proxy)
	}

	data, err := json.Marshal(proxies)
	if err != nil {
		return "[]", err
	}

	return template.JS(data), nil
}

func formatLatency(d time.Duration) string {
	if d == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}

func RenderIndex(w io.Writer, data PageData) error {
	loader := GetAssetLoader()

	var tmpl *template.Template
	if loader != nil && loader.HasCustomTemplate() {
		tmpl = loader.GetCustomTemplate()
	} else {
		tmpl = indexTmpl
	}

	if loader == nil || !loader.HasCustomCSS() {
		return tmpl.Execute(w, data)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	html := buf.String()
	customCSSLink := `<link rel="stylesheet" href="/static/custom.css">`
	html = strings.Replace(html, "</head>", customCSSLink+"\n  </head>", 1)

	_, err := io.WriteString(w, html)
	return err
}
