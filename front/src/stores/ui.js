import { defineStore } from 'pinia';

let notificationId = 0;

export const useUiStore = defineStore('ui', {
    state: () => ({
        notifications: [],
        globalLoading: false,
    }),
    actions: {
        addNotification({ type, message, timeout = 3000 }) {
            const id = ++notificationId; // гарантирует уникальность
            this.notifications.push({ id, type, message });

            setTimeout(() => {
                this.removeNotification(id);
            }, timeout);
        },
        removeNotification(id) {
            this.notifications = this.notifications.filter((n) => n.id !== id);
        },
        setGlobalLoading(status) {
            this.globalLoading = status;
        },
    },
});