import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes = [
    {
        path: '/login',
        component: () => import('../views/Login.vue')
    },
    {
        path: '/register',
        component: () => import('../views/Register.vue')
    },
    {
        path: '/terms',
        name: 'TermsOfService',
        component: () => import('../views/TermsOfService.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/index',
        name: 'Base',
        component: () => import('../views/Base.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/',
        component: () => import('../layouts/MainLayout.vue'), // Обертка для всего
        meta: { requiresAuth: true },
        children: [
            { path: '', redirect: 'index' },
            { path: 'profile', component: () => import('../views/Profile.vue') },
            { path: 'record', component: () => import('../views/ParticipationConfirm.vue') },
            { path: 'events', component: () => import('../views/EventsCatalog.vue') },
            { path: 'events/:id', component: () => import('../views/EventDetail.vue') },
            { path: 'events/create', component: () => import('../views/EventCreate.vue'), meta: { roles: ['admin'] } },
            { path: 'organizer/:id', component: () => import('../views/OrganizerProfile.vue') },
            { path: 'dashboard', component: () => import('../views/Dashboard.vue') },
            { path: 'admin', component: () => import('../views/AdminPanel.vue'), meta: { roles: ['admin'] } },
            { path: 'rating', component: () => import('@/views/rating/LeaderboardView.vue') },
            { path: 'hr/inspector', component: () => import('@/views/hr/InspectorView.vue'), meta: { roles: ['observer', 'admin'] } },
            { path: 'dashboard/activity', component: () => import('@/views/dashboard/ActivityView.vue') },
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('../views/NotFound.vue'),
        meta: { requiresAuth: false }
    }
];
const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to) => {
    const authStore = useAuthStore();

    // Проверка на обязательную авторизацию
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        return '/login';
    }

    // Проверка на роль организатора (пример)
    if (to.meta.requiresOrganizer && authStore.userRole !== 'organizer') {
        return '/profile';
    }

    // Проверка на наблюдателя
    if (to.meta.requiresObserver && !['observer', 'admin'].includes(authStore.userRole)) {
        return '/profile';
    }

    // Проверка на администратора
    if (to.meta.requiresAdmin && authStore.userRole !== 'admin') {
        return '/profile';
    }

    // Проверка по массиву допустимых ролей
    if (to.meta.roles && (!authStore.userRole || !to.meta.roles.includes(authStore.userRole))) {
        return '/dashboard';
    }

    // Если все проверки пройдены
    return true;
});

export default router;