function getBackendUrl() {
    if (window.location.protocol === 'file:') {
        // dev version assumes CORS is enabled
        return 'http://localhost:9000';
    }
    return '';
}

export default getBackendUrl