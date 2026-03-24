import { defineStore } from 'pinia';
import api from '../api/axios';

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        token: localStorage.getItem('access_token') || null,
    }),
    getters: {
        isAuthenticated: (state) => !!state.token,
        userRole: (state) => state.user?.role || null,
    },
    actions: {
        async initAuth() {
            if (this.token && !this.user) {
                await this.me();
            }
        },

        async me() {
            if (!this.token) return;
            try {
                const profileRes = await api.get('/profile');

                const responseData = profileRes.data.data || profileRes.data;

                let userData = responseData.user || responseData;

                if (responseData.organizer_stats) {
                    userData.organizer_stats = responseData.organizer_stats;
                }

                const ratingRes = await api.get('/rating/me');
                if (ratingRes.data) {
                    userData.rating = ratingRes.data.rating || ratingRes.data;
                }

                this.user = userData;

            } catch (err) {
                console.error('Ошибка загрузки профиля/рейтинга:', err);
                if (err.response?.status === 401) {
                    this.logout();
                }
            }
        },

        async login(credentials) {
            const response = await api.post('/auth/login', credentials);
            this.token = response.data.token;
            this.user = response.data.user;
            localStorage.setItem('access_token', this.token);
            localStorage.setItem('user', JSON.stringify(this.user));
            await this.me();
        },

        async register(userData) {
            const response = await api.post('/auth/register', userData);
            this.token = response.data.token;
            this.user = response.data.user;
            localStorage.setItem('access_token', this.token);
            localStorage.setItem('user', JSON.stringify(this.user));
            await this.me();
        },

        logout() {
            this.token = null;
            this.user = null;
            localStorage.removeItem('access_token');
            localStorage.removeItem('user');
        },
    }
});