import { defineStore } from 'pinia';
import api from '../api/axios';
import { ElMessage } from 'element-plus';

export const useOrganizerStore = defineStore('organizer', {
    state: () => ({
        organizer: null,
        orgStats: null,
        reviews: [],
        reviewStats: null,
        pagination: null,
        loading: false,
        submitting: false
    }),
    actions: {
        async fetchProfile(id) {
            this.loading = true;
            try {
                const res = await api.get(`/profile/${id}`);
                this.organizer = res.data.user;
                this.orgStats = res.data.stats;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Не удалось загрузить профиль организатора');
            } finally {
                this.loading = false;
            }
        },
        async fetchReviews(id, page = 1) {
            try {
                const res = await api.get(`/organizers/${id}/reviews`, {
                    params: { page, limit: 10 }
                });
                this.reviews = res.data.reviews || [];
                this.reviewStats = res.data.stats || null;
                this.pagination = res.data.pagination || null;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка при получение отзыва');
                console.error(err);
            }
        },
        async submitReview(id, reviewData) {
            this.submitting = true;
            try {
                await api.post(`/organizers/${id}/reviews`, reviewData);
                ElMessage.success('Отзыв успешно добавлен!');
                await this.fetchReviews(id, 1);
                await this.fetchProfile(id);
                return true;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка при отправке отзыва');
                return false;
            } finally {
                this.submitting = false;
            }
        }
    }
});