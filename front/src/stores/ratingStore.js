import { defineStore } from 'pinia';
import api from '../api/axios';
import {ElMessage} from "element-plus";

export const useRatingStore = defineStore('rating', {
  state: () => ({
    globalRating: [],
    userRank: null,
    loading: false,
    currentDirection: 'all',
  }),
  actions: {
    async fetchLeaders() {
      this.loading = true;
      try {
        let response;
        if (this.currentDirection === 'all') {
          response = await api.get('/rating/global', { params: { limit: 100, offset: 0 } });
        } else {
          response = await api.get(`/rating/direction/${this.currentDirection}`, {
            params: { limit: 100, offset: 0 },
          });
        }

        const rating = response.data.data?.rating || response.data.rating || [];

        const enriched = rating.map(user => ({
          ...user,
          direction: this._determineDirection(user),
          events_count: user.events_count || 0,
        }));

        this.globalRating = enriched;
        this.userRank = response.data.data?.user_rank || response.data.user_rank || null;
      } catch (err) {
        ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки рейтинга');
      } finally {
        this.loading = false;
      }
    },

    _determineDirection(user) {
      if (user.it_rank !== null && user.it_rank !== undefined) return 'IT';
      if (user.social_rank !== null && user.social_rank !== undefined) return 'social';
      if (user.media_rank !== null && user.media_rank !== undefined) return 'media';
      return null;
    },

    setDirection(direction) {
      if (this.currentDirection === direction) return;
      this.currentDirection = direction;
      this.fetchLeaders();
    },

    async fetchUserRating(userId) {
      try {
        const response = await api.get(`/rating/user/${userId}`);
        this.userRating = response.data.data || response.data;
      } catch (err) {
        console.error(err);
      }
    },
  },
});