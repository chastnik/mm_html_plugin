package main

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
)

type Plugin struct {
	plugin.MattermostPlugin

	configuration *configuration
}

type configuration struct {
	MaxFileSize     int  `json:"MaxFileSize"`
	EnableJavaScript bool `json:"EnableJavaScript"`
	EnableCSS       bool `json:"EnableCSS"`
}

type HTMLViewerData struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	FileId   string `json:"file_id"`
}

// OnActivate запускается при активации плагина
func (p *Plugin) OnActivate() error {
	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	return nil
}

// OnConfigurationChange обновляет конфигурацию плагина
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	p.configuration = configuration
	return nil
}

// ServeHTTP обрабатывает HTTP запросы
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/api/v1/preview":
		p.handlePreview(w, r)
	case "/api/v1/content":
		p.handleContent(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handlePreview обрабатывает запрос на получение превью HTML файла
func (p *Plugin) handlePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileId := r.URL.Query().Get("file_id")
	if fileId == "" {
		http.Error(w, "file_id parameter is required", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	fileInfo, appErr := p.API.GetFileInfo(fileId)
	if appErr != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Проверяем расширение файла
	if !p.isHTMLFile(fileInfo.Name) {
		http.Error(w, "File is not an HTML file", http.StatusBadRequest)
		return
	}

	// Проверяем размер файла
	maxSize := int64(p.configuration.MaxFileSize * 1024 * 1024) // MB to bytes
	if fileInfo.Size > maxSize {
		http.Error(w, "File is too large", http.StatusBadRequest)
		return
	}

	// Читаем содержимое файла
	fileData, appErr := p.API.GetFile(fileId)
	if appErr != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Санитизируем HTML контент
	sanitizedContent := p.sanitizeHTML(string(fileData))

	response := HTMLViewerData{
		Content:  sanitizedContent,
		Filename: fileInfo.Name,
		FileId:   fileId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleContent возвращает сырой HTML контент для iframe
func (p *Plugin) handleContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileId := r.URL.Query().Get("file_id")
	if fileId == "" {
		http.Error(w, "file_id parameter is required", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	fileInfo, appErr := p.API.GetFileInfo(fileId)
	if appErr != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Проверяем расширение файла
	if !p.isHTMLFile(fileInfo.Name) {
		http.Error(w, "File is not an HTML file", http.StatusBadRequest)
		return
	}

	// Читаем содержимое файла
	fileData, appErr := p.API.GetFile(fileId)
	if appErr != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Санитизируем HTML контент
	sanitizedContent := p.sanitizeHTML(string(fileData))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline'; script-src 'none'; object-src 'none';")
	
	io.WriteString(w, sanitizedContent)
}

// isHTMLFile проверяет, является ли файл веб-файлом
func (p *Plugin) isHTMLFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	webExtensions := []string{".html", ".htm", ".xhtml", ".xml", ".svg", ".css", ".js"}
	
	for _, webExt := range webExtensions {
		if ext == webExt {
			return true
		}
	}
	return false
}

// sanitizeHTML очищает HTML контент от потенциально опасных элементов
func (p *Plugin) sanitizeHTML(htmlContent string) string {
	policy := bluemonday.UGCPolicy()

	if p.configuration.EnableCSS {
		policy.AllowStyling()
	}

	// JavaScript и опасные атрибуты блокируются через политику безопасности

	// Разрешаем базовые HTML теги
	policy.AllowElements("html", "head", "body", "title", "meta", "link")
	policy.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	policy.AllowElements("p", "br", "hr", "div", "span")
	policy.AllowElements("ul", "ol", "li", "dl", "dt", "dd")
	policy.AllowElements("table", "thead", "tbody", "tr", "td", "th")
	policy.AllowElements("a", "img", "strong", "em", "u", "strike")
	policy.AllowElements("pre", "code", "blockquote")

	// Разрешаем атрибуты
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowAttrs("src", "alt", "width", "height").OnElements("img")
	policy.AllowAttrs("class", "id").Globally()

	if p.configuration.EnableCSS {
		policy.AllowAttrs("style").Globally()
	}

	return policy.Sanitize(htmlContent)
} 