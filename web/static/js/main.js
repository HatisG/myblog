// API基础地址
const API_BASE_URL = 'http://localhost:8080';

// 工具函数：格式化日期
function formatDate(dateStr) {
    const date = new Date(dateStr);
    return `${date.getFullYear()}-${(date.getMonth()+1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
}

// 工具函数：显示错误提示
function showError(message) {
    alert(`错误：${message}`);
}

// 工具函数：获取URL参数
function getUrlParam(name) {
    const reg = new RegExp(`(^|&)${name}=([^&]*)(&|$)`);
    const r = window.location.search.substr(1).match(reg);
    if (r != null) return decodeURIComponent(r[2]);
    return null;
}

// 原生JS封装HTTP GET请求
function httpGet(url, params = {}) {
    return new Promise((resolve, reject) => {
        // 拼接查询参数
        const queryString = Object.keys(params)
            .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
            .join('&');
        
        const fullUrl = queryString ? `${API_BASE_URL}${url}?${queryString}` : `${API_BASE_URL}${url}`;
        
        const xhr = new XMLHttpRequest();
        xhr.open('GET', fullUrl, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        
        xhr.onload = function() {
            if (xhr.status >= 200 && xhr.status < 300) {
                try {
                    const response = JSON.parse(xhr.responseText);
                    resolve(response);
                } catch (e) {
                    reject(new Error('解析响应数据失败：' + e.message));
                }
            } else {
                reject(new Error(`请求失败：${xhr.status} ${xhr.statusText}`));
            }
        };
        
        xhr.onerror = function() {
            reject(new Error('网络请求失败，请检查连接'));
        };
        
        xhr.send();
    });
}

// 原生JS封装HTTP POST请求
function httpPost(url, data = {}) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', `${API_BASE_URL}${url}`, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        
        xhr.onload = function() {
            if (xhr.status >= 200 && xhr.status < 300) {
                try {
                    const response = xhr.responseText ? JSON.parse(xhr.responseText) : {};
                    resolve(response);
                } catch (e) {
                    reject(new Error('解析响应数据失败：' + e.message));
                }
            } else {
                try {
                    const errorData = JSON.parse(xhr.responseText);
                    reject(new Error(errorData.error || `请求失败：${xhr.status} ${xhr.statusText}`));
                } catch (e) {
                    reject(new Error(`请求失败：${xhr.status} ${xhr.statusText}`));
                }
            }
        };
        
        xhr.onerror = function() {
            reject(new Error('网络请求失败，请检查连接'));
        };
        
        xhr.send(JSON.stringify(data));
    });
}

// 原生JS封装HTTP PUT请求
function httpPut(url, data = {}) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('PUT', `${API_BASE_URL}${url}`, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        
        xhr.onload = function() {
            if (xhr.status >= 200 && xhr.status < 300) {
                try {
                    const response = xhr.responseText ? JSON.parse(xhr.responseText) : {};
                    resolve(response);
                } catch (e) {
                    reject(new Error('解析响应数据失败：' + e.message));
                }
            } else {
                try {
                    const errorData = JSON.parse(xhr.responseText);
                    reject(new Error(errorData.error || `请求失败：${xhr.status} ${xhr.statusText}`));
                } catch (e) {
                    reject(new Error(`请求失败：${xhr.status} ${xhr.statusText}`));
                }
            }
        };
        
        xhr.onerror = function() {
            reject(new Error('网络请求失败，请检查连接'));
        };
        
        xhr.send(JSON.stringify(data));
    });
}

// 加载导航栏（所有页面共用）
function loadNavbar() {
    const navbarHTML = `
        <div class="navbar-container container">
            <a href="index.html" class="logo">MyBlog</a>
            <div class="nav-links">
                <a href="index.html">首页</a>
                <a href="create-post.html">发布文章</a>
                <a href="search-post.html">查询文章</a>
                <a href="about.html">关于我们</a>
            </div>
        </div>
    `;
    document.getElementById('navbar').innerHTML = navbarHTML;
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    loadNavbar();
});