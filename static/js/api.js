const API_BASE_URL = '/api';

async function fetchAPI(endpoint, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };

    // Lấy Token từ LocalStorage (do file auth.js lưu)
    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Token ${token}`;
    }

    const options = {
        method,
        headers,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
        const data = await response.json();

        if (response.status === 401 || response.status === 403) {
            alert("Phiên đăng nhập hết hạn hoặc bạn không có quyền!");
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            if (window.location.pathname.includes('/admin/')) {
                window.location.href = '/admin/login.html';
            } else {
                window.location.href = '/login.html';
            }
            return null;
        }

        if (!response.ok) {
            throw new Error(data.message || 'Có lỗi xảy ra từ máy chủ');
        }

        return data; 
    } catch (error) {
        console.error('API Error:', error);
        throw error; 
    }
}

async function fetchUploadAPI(endpoint, method = 'POST', formData) {
    const headers = {};
    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Token ${token}`;
    }

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method,
            headers,
            body: formData,
        });
        const data = await response.json();
        if (!response.ok) throw new Error(data.message);
        return data;
    } catch (error) {
        console.error('Upload Error:', error);
        throw error;
    }
}