// src/api/profile.js
import api from './axios';

/**
 * Получить профиль текущего пользователя
 */
export const getProfile = () => api.get('/profile');

/**
 * Обновить профиль текущего пользователя
 * @param {Object} data - обновляемые поля (full_name, city, age, direction, avatar_url)
 */
export const updateProfile = (data) => api.put('/profile', data);

/**
 * Получить публичный профиль другого пользователя
 * @param {string} userId - id пользователя
 */
export const getPublicProfile = (userId) => api.get(`/profile/${userId}`);