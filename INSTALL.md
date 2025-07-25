# Инструкция по установке и сборке HTML Viewer Plugin

## Требования к системе

Для сборки плагина необходимо установить:

1. **Go 1.19+** - https://golang.org/dl/
2. **Node.js 14+** - https://nodejs.org/
3. **npm** (поставляется с Node.js)
4. **Make** (для Unix/macOS) или **nmake** (для Windows)

## Проверка установленных компонентов

```bash
# Проверка Go
go version

# Проверка Node.js  
node --version

# Проверка npm
npm --version

# Проверка Make
make --version
```

## Пошаговая сборка

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd mm_html_plugin
```

### 2. Установка зависимостей
```bash
# Установка Go зависимостей
cd server
go mod download
cd ..

# Установка npm зависимостей
cd webapp  
npm install
cd ..
```

### 3. Сборка плагина

#### Автоматическая сборка (рекомендуется)
```bash
make bundle
```

#### Ручная сборка

##### Серверная часть
```bash
cd server

# Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/plugin-linux-amd64 ./...

# macOS  
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/plugin-darwin-amd64 ./...

# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/plugin-windows-amd64.exe ./...

cd ..
```

##### Клиентская часть
```bash
cd webapp
npm run build
cd ..
```

##### Создание bundle
```bash
rm -rf dist/
mkdir -p dist/com.mattermost.html-viewer
cp plugin.json dist/com.mattermost.html-viewer/
cp -r server/dist dist/com.mattermost.html-viewer/server/
cp -r webapp/dist dist/com.mattermost.html-viewer/webapp/
cd dist && tar -czf com.mattermost.html-viewer-1.0.0.tar.gz com.mattermost.html-viewer
```

## Установка в Mattermost

1. Откройте Mattermost как администратор
2. Перейдите в **Системная консоль** > **Плагины** > **Управление плагинами**
3. Нажмите **Выбрать файл** и загрузите `dist/com.mattermost.html-viewer-1.0.0.tar.gz`
4. Нажмите **Загрузить**
5. Найдите плагин "HTML Viewer" в списке и нажмите **Включить**

## Настройка

После установки и активации:

1. Перейдите в **Системная консоль** > **Плагины** > **HTML Viewer**
2. Настройте параметры:
   - **Максимальный размер файла (MB)**: 5 (по умолчанию)
   - **Разрешить JavaScript**: false (по умолчанию, рекомендуется)  
   - **Разрешить CSS**: true (по умолчанию)
3. Нажмите **Сохранить**

## Использование

1. Загрузите HTML файл в любой канал
2. Рядом с файлом появится кнопка **"Показать предпросмотр"**
3. Нажмите на кнопку для просмотра HTML контента
4. HTML будет отображен в безопасном iframe

## Поддерживаемые форматы

- `.html` - HTML документы
- `.htm` - HTML документы (короткое расширение)
- `.xhtml` - XHTML документы  
- `.xml` - XML документы
- `.svg` - SVG векторная графика
- `.css` - CSS стили
- `.js` - JavaScript файлы

## Решение проблем

### Ошибка "File is too large"
- Увеличьте параметр "Максимальный размер файла" в настройках плагина

### Ошибка "File is not an HTML file"  
- Убедитесь, что файл имеет поддерживаемое расширение
- Проверьте, что расширение написано правильно

### HTML не отображается корректно
- Проверьте, включен ли параметр "Разрешить CSS"
- Убедитесь, что HTML код не содержит критических ошибок
- Внешние ресурсы могут быть заблокированы политикой безопасности

### Ошибки при сборке
```bash
# Очистка кеша
make clean

# Переустановка зависимостей
rm -rf webapp/node_modules
cd webapp && npm install

# Очистка Go модулей
cd server && go clean -modcache && go mod download
```

## Отладка

Для включения отладочной информации:

```bash
MM_DEBUG=1 make bundle
```

Логи плагина можно найти в логах Mattermost:
- Linux: `/opt/mattermost/logs/mattermost.log`
- Docker: `docker logs <container_name>`

## Безопасность

⚠️ **Важно**: По умолчанию JavaScript отключен для безопасности. Включение JavaScript может представлять угрозу безопасности, поскольку позволяет выполнение произвольного кода. Используйте эту опцию только в доверенной среде. 