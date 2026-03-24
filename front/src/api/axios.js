import axios from 'axios';
import { setupMocks } from './mockSetup';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8000/api',
    timeout: 10000,
});

// if (import.meta.env.VITE_USE_MOCK === 'true') {
//     setupMocks(api);
//     console.warn('Запущено с MOCK API. Реальные запросы не отправляются.');
// }

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    (response) => {
        if (response.data && typeof response.data.data !== 'undefined') {
            response.data = response.data.data;
        }
        return response;
    },
    (error) => {
        if (error.response && error.response.status === 401) {
            localStorage.removeItem('access_token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default api;