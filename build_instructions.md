# Завершение сборки HTML Viewer Plugin

## Клиентская часть уже собрана ✅
Файлы в `webapp/dist/`:
- main.js (8.2KB) - основной JavaScript файл плагина
- main.js.map - карта исходного кода для отладки

## Для завершения сборки установите Go:

### 1. Установка Go на macOS:
```bash
# Через Homebrew (рекомендуется)
brew install go

# Или скачайте с https://golang.org/dl/
```

### 2. Проверка установки:
```bash
go version
# Должна появиться версия Go 1.19+
```

### 3. Сборка серверной части:
```bash
# Установка зависимостей Go
cd server
go mod download

# Сборка для разных платформ
# macOS
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/plugin-darwin-amd64 ./...

# Linux  
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/plugin-linux-amd64 ./...

# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/plugin-windows-amd64.exe ./...

cd ..
```

### 4. Создание финального bundle:
```bash
# Автоматически (через Makefile)
make bundle

# Или вручную
rm -rf dist/
mkdir -p dist/com.mattermost.html-viewer
cp plugin.json dist/com.mattermost.html-viewer/
cp -r server/dist dist/com.mattermost.html-viewer/server/
cp -r webapp/dist dist/com.mattermost.html-viewer/webapp/
cd dist && tar -czf com.mattermost.html-viewer-1.0.0.tar.gz com.mattermost.html-viewer
```

### 5. Установка в Mattermost:
1. Откройте Mattermost как администратор
2. **Системная консоль** > **Плагины** > **Управление плагинами**
3. Загрузите файл `dist/com.mattermost.html-viewer-1.0.0.tar.gz`
4. Включите плагин

## Текущее состояние:
```
mm_html_plugin/
├── webapp/dist/          ✅ Собрано
│   ├── main.js          (8.2KB)
│   └── main.js.map      (16.9KB)
├── server/dist/          ❌ Требует Go
│   └── (пусто)
└── dist/                 ❌ Итоговый bundle
    └── (не создан)
```

После установки Go выполните:
```bash
make bundle
```

И получите готовый файл для установки в Mattermost! 