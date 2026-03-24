import {defineStore} from 'pinia';
import {ref} from 'vue';
import api from '@/api/axios';
import {ElMessage} from 'element-plus';

export const useHrStore = defineStore('hr', () => {
    // Состояние
    const candidates = ref([]);
    const total = ref(0);
    const pages = ref(0);
    const savedFilters = ref([]);
    const isLoading = ref(false);

    // Фильтры по умолчанию
    const defaultFilters = {
        age_min: null,
        age_max: null,
        city: null,
        direction: null,
        min_points: null,
        min_events: null,
        min_avg_points: null,
        sort_by: 'points',
        sort_order: 'desc',
        page: 1,
        limit: 20,
    };
    const filters = ref({ ...defaultFilters });

    // Список городов (можно вынести в константу)
    const cities = [
        'Москва',
        'Санкт-Петербург',
        'Казань',
        'Новосибирск',
        'Екатеринбург',
        'Нижний Новгород',
    ];

    // Получение кандидатов
    const fetchCandidates = async () => {
        isLoading.value = true;
        try {
            // Очищаем null значения, чтобы они не передавались в запрос
            const params = Object.fromEntries(
                Object.entries(filters.value).filter(([_, v]) => v !== null && v !== '')
            );
            const response = await api.get('/cadre/candidates', { params });
            candidates.value = response.data.candidates || [];
            total.value = response.data.pagination.total || 0;
            pages.value = response.data.pagination.pages || 0;
        } catch (error) {
            console.error('Ошибка загрузки кандидатов:', error);
            ElMessage.error(err.response?.data?.error?.message || 'Не удалось загрузить кандидатов');
        } finally {
            isLoading.value = false;
        }
    };

    // Загрузка сохранённых фильтров
    const fetchSavedFilters = async () => {
        try {
            const response = await api.get('/filters');
            const filtersData = response.data?.filters || [];
            savedFilters.value = filtersData.map(filter => ({
                id: filter.id,
                name: filter.filter_name,
                filters: filter.filters,
                createdAt: filter.created_at,
            }));
            console.log('Сохранённые фильтры загружены:', savedFilters.value);
        } catch (error) {
            console.error('Ошибка загрузки сохранённых фильтров:', error);
            ElMessage.error(error.response?.data?.error?.message || 'Ошибка загрузки сохранённых фильтров');
        }
    };

    // Сохранение текущего фильтра
    const saveFilter = async (name) => {
        if (!name || name.trim() === '') {
            ElMessage.warning('Введите название фильтра');
            return;
        }
        try {
            const payload = {
                filter_name: name.trim(),
                filters: { ...filters.value },
            };
            const response = await api.post('/filters', payload);
            const newFilter = {
                id: response.data.id,
                name: response.data.filter_name,
                filters: response.data.filters,
                createdAt: response.data.created_at,
            };
            savedFilters.value.push(newFilter);
            ElMessage.success('Фильтр сохранён');
        } catch (error) {
            console.error('Ошибка сохранения фильтра:', error);
            ElMessage.error(error.response?.data?.error?.message || 'Не удалось сохранить фильтр');
        }
    };

    // Применение сохранённого фильтра
    const applySavedFilter = async (savedFilter) => {
        filters.value = {
            ...defaultFilters,
            ...savedFilter.filters,
            page: 1,
        };
        await fetchCandidates();
    };

    // Обновление фильтров из компонента
    const updateFilters = async (newFilters) => {
        filters.value = {
            ...filters.value,
            ...newFilters,
            page: 1,
        };
        await fetchCandidates();
    };

    // Сброс фильтров
    const resetFilters = async () => {
        filters.value = { ...defaultFilters };
        await fetchCandidates();
    };

    const deleteFilter = async (filterId) => {
        try {
            await api.delete(`/filters/${filterId}`);
            savedFilters.value = savedFilters.value.filter(f => f.id !== filterId);
            ElMessage.success('Фильтр удалён');
        } catch (error) {
            console.error('Ошибка удаления фильтра:', error);
            ElMessage.error(error.response?.data?.error?.message || 'Не удалось удалить фильтр');
        }
    };


    // Скачивание отчёта PDF
    const downloadReport = async (userId) => {
        try {
            const response = await api.get(`/cadre/candidates/${userId}/report`, {
                responseType: 'blob',
            });
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `candidate_${userId}.pdf`);
            document.body.appendChild(link);
            link.click();
            link.remove();
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Ошибка скачивания отчёта:', error);
            ElMessage.error(error.response?.data?.error?.message || 'Не удалось скачать отчёт');
        }
    };


    // Возвращаем всё, что нужно компонентам
    return {
        // состояние
        candidates,
        total,
        pages,
        savedFilters,
        isLoading,
        filters,
        cities,
        // методы
        fetchCandidates,
        fetchSavedFilters,
        saveFilter,
        applySavedFilter,
        updateFilters,
        deleteFilter,
        resetFilters,
        downloadReport,
    };
});