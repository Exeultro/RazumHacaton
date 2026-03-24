import { defineStore } from 'pinia';
import api from '../api/axios';

export const useDashboardStore = defineStore('dashboard', {
    state: () => ({
        recentEvents: [],
        ratingHistory: [],
        trendingTags: [],
        loading: false
    }),
    actions: {
        async fetchActivity() {
            this.loading = true;
            try {
                const response = await api.get('/dashboard');
                this.recentEvents = response.data.recent_events;
                this.ratingHistory = response.data.rating_history;
                this.trendingTags = response.data.trending_tags;
            } catch (err) {
                console.error(err);
            } finally {
                this.loading = false;
            }
        }
    }
});