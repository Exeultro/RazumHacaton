import { defineStore } from 'pinia';
import api from '../api/axios';
import { ElMessage } from 'element-plus';

export const useEventStore = defineStore('event', {
    state: () => ({
        events: [],
        currentEvent: null,
        filters: {
            status: 'published',
            direction: null,
            format: null,
            date_from: null,
            date_to: null,
            page: 1,
            limit: 20
        },
        pagination: null,
        loading: false
    }),
    actions: {
        async fetchEvents() {
            this.loading = true;
            try {
                const response = await api.get('/events', { params: this.filters });
                this.events = response.data.events;
                this.pagination = response.data.pagination;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Не удалось загрузить мероприятия');
            } finally {
                this.loading = false;
            }
        },

        async fetchEventById(id) {
            this.loading = true;
            try {
                const response = await api.get(`/events/${id}`);
                this.currentEvent = response.data;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Не удалось загрузить мероприятие');
            } finally {
                this.loading = false;
            }
        },

        async createEvent(eventData) {
            try {
                const response = await api.post('/events', eventData);
                return response.data;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка создания');
                console.error('Ошибка создания:', err);
                throw err;
            }
        },

        async registerForEvent(eventId) {
            try {
                const response = await api.post(`/events/${eventId}/register`);
                ElMessage.success(response.data.message || 'Вы успешно зарегистрированы');
                return response.data;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка регистрации');
                throw err;
            }
        },

        async confirmParticipation(eventId, qrCodeToken) {
            try {
                const response = await api.post(`/events/${eventId}/confirm`, { qr_code_token: qrCodeToken });
                ElMessage.success(`Участие подтверждено, начислено ${response.data.points_earned} баллов`);
                return response.data;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка подтверждения участия');
                throw err;
            }
        },

        async cancelRegistration(eventId) {
            try {
                await api.delete(`/events/${eventId}/cancel`);
                ElMessage.success('Регистрация отменена');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка отмены регистрации');
                throw err;
            }
        },

        updateFilters(newFilters) {
            this.filters = { ...this.filters, ...newFilters, page: 1 };
            this.fetchEvents();
        }
    }
});