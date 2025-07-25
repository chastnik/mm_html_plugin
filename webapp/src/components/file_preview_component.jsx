import React, {useState, useEffect} from 'react';
import PropTypes from 'prop-types';

const FilePreviewComponent = ({fileInfo}) => {
    const [htmlContent, setHtmlContent] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showPreview, setShowPreview] = useState(false);

    useEffect(() => {
        if (fileInfo && fileInfo.id) {
            loadHtmlContent();
        }
    }, [fileInfo]);

    const loadHtmlContent = async () => {
        try {
            setLoading(true);
            setError(null);

            const url = `/plugins/com.mattermost.html-viewer/api/v1/preview?file_id=${fileInfo.id}`;
            const response = await fetch(url);
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const data = await response.json();
            setHtmlContent(data.content);
        } catch (err) {
            setError(err.message);
            console.error('Failed to load HTML content:', err);
        } finally {
            setLoading(false);
        }
    };

    const getContentUrl = () => {
        return `/plugins/com.mattermost.html-viewer/api/v1/content?file_id=${fileInfo.id}`;
    };

    const handleTogglePreview = () => {
        setShowPreview(!showPreview);
    };

    const formatFileSize = (bytes) => {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    if (loading) {
        return (
            <div className="html-viewer-preview">
                <div className="html-viewer-loading">
                    <i className="fa fa-spinner fa-spin"></i>
                    <span style={{marginLeft: '8px'}}>Загрузка HTML файла...</span>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="html-viewer-preview">
                <div className="html-viewer-error">
                    <i className="fa fa-exclamation-triangle" style={{color: '#d24b47'}}></i>
                    <span style={{marginLeft: '8px', color: '#d24b47'}}>
                        Ошибка загрузки: {error}
                    </span>
                </div>
            </div>
        );
    }

    return (
        <div className="html-viewer-preview">
            <div className="html-viewer-header">
                <div className="html-viewer-info">
                    <i className="fa fa-file-code-o" style={{marginRight: '8px', color: '#3f4350'}}></i>
                    <span className="html-viewer-filename">{fileInfo.name}</span>
                    <span className="html-viewer-filesize">
                        ({formatFileSize(fileInfo.size)})
                    </span>
                </div>
                <button 
                    className="btn btn-sm btn-primary html-viewer-toggle"
                    onClick={handleTogglePreview}
                    style={{marginLeft: '12px'}}
                >
                    {showPreview ? 'Скрыть предпросмотр' : 'Показать предпросмотр'}
                </button>
            </div>

            {showPreview && (
                <div className="html-viewer-content">
                    <iframe
                        src={getContentUrl()}
                        style={{
                            width: '100%',
                            height: '400px',
                            border: '1px solid #e6e6e6',
                            borderRadius: '3px',
                            backgroundColor: '#fff'
                        }}
                        sandbox="allow-same-origin"
                        title={`HTML Preview: ${fileInfo.name}`}
                    />
                    <div className="html-viewer-warning" style={{
                        marginTop: '8px',
                        padding: '8px',
                        backgroundColor: '#fef7e0',
                        border: '1px solid #f0d000',
                        borderRadius: '3px',
                        fontSize: '12px',
                        color: '#7d6608'
                    }}>
                        <i className="fa fa-info-circle" style={{marginRight: '4px'}}></i>
                        HTML контент очищен от потенциально опасных элементов для безопасности
                    </div>
                </div>
            )}
        </div>
    );
};

FilePreviewComponent.propTypes = {
    fileInfo: PropTypes.object.isRequired,
};

export default FilePreviewComponent; 