// src/api/reviews.js
import api from './axios';

/**
 * Оставить отзыв на организатора
 * @param {string} organizerId - id организатора
 * @param {Object} data - { event_id?, rating, comment }
 */
export const createReview = (organizerId, data) =>
  api.post(`/organizers/${organizerId}/reviews`, data);

/**
 * Получить отзывы об организаторе
 * @param {string} organizerId - id организатора
 */
export const getReviews = (organizerId) =>
  api.get(`/organizers/${organizerId}/reviews`);