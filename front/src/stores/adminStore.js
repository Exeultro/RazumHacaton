import { defineStore } from 'pinia';
import api from '../api/axios';
import { ElMessage } from 'element-plus';

export const useAdminStore = defineStore('admin', {
    state: () => ({
        loading: false,
        stats: null,
        pendingOrganizers: [],

        users: [],
        usersPagination: { page: 1, limit: 20, total: 0, pages: 1 },

        events: [],
        eventsPagination: { page: 1, limit: 20, total: 0, pages: 1 },

        reviews: [],
        reviewsPagination: { page: 1, limit: 20, total: 0, pages: 1 },

        settings: { IT: 1.0, social: 1.0, media: 1.0 }
    }),

    actions: {
        async fetchStats() {
            try {
                const response = await api.get('/admin/stats');
                // Интерцептор уже достал данные!
                this.stats = response.data;
            } catch (err) {
                console.error('Ошибка загрузки статистики:', err);
            }
        },

        async refreshRatingCache() {
            try {
                await api.post('/rating/refresh');
                ElMessage.success('Кэш рейтинга обновлен');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка обновления кэша');
            }
        },

        async fetchPendingOrganizers() {
            this.loading = true;
            try {
                const response = await api.get('/admin/organizers/pending');
                this.pendingOrganizers = response.data.pending_organizers || [];
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки заявок');
            } finally {
                this.loading = false;
            }
        },

        async approveOrganizer(userId) {
            try {
                await api.post(`/admin/organizers/${userId}/approve`);
                this.pendingOrganizers = this.pendingOrganizers.filter(o => o.user_id !== userId);
                ElMessage.success('Организатор одобрен');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка');
            }
        },

        async rejectOrganizer(userId) {
            try {
                await api.post(`/admin/organizers/${userId}/reject`);
                this.pendingOrganizers = this.pendingOrganizers.filter(o => o.user_id !== userId);
                ElMessage.success('Организатор отклонен');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка');
            }
        },

        async fetchUsers(page = 1, limit = 20) {
            this.loading = true;
            try {
                const response = await api.get(`/admin/users?page=${page}&limit=${limit}`);
                this.users = response.data.users;
                this.usersPagination = response.data.pagination;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки пользователей');
            } finally {
                this.loading = false;
            }
        },

        async changeUserRole(userId, newRole) {
            try {
                await api.put(`/admin/users/${userId}/role`, { role: newRole });
                const userIndex = this.users.findIndex(u => u.id === userId);
                if (userIndex !== -1) {
                    this.users[userIndex].role = newRole;
                }
                ElMessage.success('Роль успешно изменена');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка изменения роли');
            }
        },

        async deleteUser(userId) {
            try {
                await api.delete(`/admin/users/${userId}`);
                this.users = this.users.filter(u => u.id !== userId);
                ElMessage.success('Пользователь удален');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка удаления пользователя');
            }
        },

        async fetchEvents(page = 1, limit = 20) {
            this.loading = true;
            try {
                const response = await api.get(`/admin/events?page=${page}&limit=${limit}`);
                this.events = response.data.events;
                this.eventsPagination = response.data.pagination;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки мероприятий');
            } finally {
                this.loading = false;
            }
        },

        async deleteEvent(eventId) {
            try {
                await api.delete(`/admin/events/${eventId}`);
                this.events = this.events.filter(e => e.id !== eventId);
                ElMessage.success('Мероприятие удалено');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка удаления');
            }
        },

        async fetchReviews(page = 1, limit = 20) {
            this.loading = true;
            try {
                const response = await api.get(`/admin/reviews?page=${page}&limit=${limit}`);
                this.reviews = response.data.reviews;
                this.reviewsPagination = response.data.pagination;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки отзывов');
            } finally {
                this.loading = false;
            }
        },

        async deleteReview(reviewId) {
            try {
                await api.delete(`/admin/reviews/${reviewId}`);
                this.reviews = this.reviews.filter(r => r.id !== reviewId);
                ElMessage.success('Отзыв удален');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка удаления');
            }
        },

        async fetchDifficultySettings() {
            try {
                const response = await api.get('/admin/settings/difficulty');
                this.settings = response.data;
            } catch (err) {
                console.error('Ошибка загрузки настроек:', err);
            }
        },

        async updateDifficultySettings(newSettings) {
            try {
                await api.put('/admin/settings/difficulty', newSettings);
                this.settings = newSettings;
                ElMessage.success('Настройки успешно обновлены');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка обновления настроек');
            }
        }
    }
});