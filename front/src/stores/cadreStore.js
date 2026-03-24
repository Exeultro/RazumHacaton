import { defineStore } from 'pinia';
import api from '../api/axios';
import { ElMessage } from 'element-plus';

export const useCadreStore = defineStore('cadre', {
    state: () => ({
        candidates: [],
        filters: {
            age_min: null,
            age_max: null,
            city: null,
            direction: null,
            min_points: null,
            min_events: null,
            sort_by: 'points',
            sort_order: 'desc',
            page: 1,
            limit: 20
        },
        total: 0,
        pages: 0,
        savedFilters: [],
        loading: false
    }),
    actions: {
        async fetchCandidates() {
            this.loading = true;
            try {
                const response = await api.get('/cadre/candidates', { params: this.filters });
                this.candidates = response.data.candidates;
                this.total = response.data.total;
                this.pages = response.data.pages;
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки кандидатов');

            } finally {
                this.loading = false;
            }
        },
        async downloadReport(userId) {
            try {
                const response = await api.get(`/cadre/candidates/${userId}/report`, { responseType: 'blob' });
                const url = window.URL.createObjectURL(new Blob([response.data]));
                const link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', `candidate_${userId}.pdf`);
                document.body.appendChild(link);
                link.click();
                link.remove();
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка скачивания отчета');
            }
        },
        async saveFilter(name, filters) {
            try {
                const response = await api.post('/cadre/filters', { name, filters });
                this.savedFilters.push(response.data);
                ElMessage.success('Фильтр сохранен');
            } catch (err) {
                ElMessage.error(err.response?.data?.error?.message || 'Ошибка сохранения фильтра');

            }
        },
        async loadFilters() {
            try {
                const response = await api.get('/cadre/filters');
                this.savedFilters = response.data;
            } catch (err) {

                console.error(err);
            }
        },
        updateFilters(newFilters) {
            this.filters = { ...this.filters, ...newFilters, page: 1 };
            this.fetchCandidates();
        }
    }
});