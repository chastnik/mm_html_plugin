import {id as pluginId} from './manifest';

import FilePreviewComponent from './components/file_preview_component';

export default class Plugin {
    initialize(registry, store) {
        // Регистрируем обработчик для HTML файлов
        registry.registerFilePreviewComponent(
            (fileInfo) => {
                if (!fileInfo || !fileInfo.name) {
                    return false;
                }
                const ext = fileInfo.name.split('.').pop().toLowerCase();
                return ['html', 'htm', 'xhtml', 'xml', 'svg', 'css', 'js'].includes(ext);
            },
            FilePreviewComponent
        );
    }
} 