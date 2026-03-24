import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

import './assets/main.css'

import ElementPlus, {ElMessage} from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)


ElMessage.warning = (msg) => {
    const uiStore = useUiStore();
    uiStore.addNotification({ type: 'warning', message: msg });
};

ElMessage.error = (msg) => {
    const uiStore = useUiStore();
    uiStore.addNotification({ type: 'error', message: msg });
};

ElMessage.success = (msg) => {
    const uiStore = useUiStore();
    uiStore.addNotification({ type: 'success', message: msg });
};

ElMessage.info = (msg) => {
    const uiStore = useUiStore();
    uiStore.addNotification({ type: 'info', message: msg });
};


app.use(router)
app.use(ElementPlus)

const authStore = useAuthStore(pinia)
authStore.initAuth()

app.mount('#app')

