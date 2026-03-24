<template>
  <div class="app-layout">
    <header class="top-bar">
      <div class="logo">Парламент<span>Резерв</span></div>
      <div class="user-zone">
        <div class="user-info">
          <AppAvatar :url="user?.avatar_url" :name="user?.full_name" size="sm" />
          <span>{{ user?.full_name }}</span>
        </div>
        <button class="btn-logout" @click="handleLogout">Выйти</button>
      </div>
    </header>

    <div class="layout-body">
      <aside class="sidebar">
        <nav class="sidebar-nav">
          <router-link to="/dashboard" class="nav-item">Профиль</router-link>
          <router-link to="/rating" class="nav-item">Рейтинг</router-link>
          <router-link to="/dashboard/activity" class="nav-item">Дашборд</router-link>
          <div v-if="user?.role === 'observer'" class="nav-divider"></div>
          <router-link v-if="user?.role === 'observer' || user?.role === 'admin'" to="/hr/inspector" class="nav-item">Кадры</router-link>
          <div v-if="user?.role === 'admin'" class="nav-divider"></div>
          <router-link v-if="user?.role === 'admin'" to="/admin" class="nav-item">Админ Панель</router-link>
        </nav>
        <div class="sidebar-footer">
          <router-link to="/terms" class="nav-item">Пользовательское соглашение</router-link>
        </div>
      </aside>

      <main class="content-area">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { computed } from 'vue'
import AppAvatar from "@/components/ui/AppAvatar.vue";

const authStore = useAuthStore()
const router = useRouter()
const user = computed(() => authStore.user)

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-layout { display: flex; flex-direction: column; min-height: 100vh; }

.top-bar {
  height: 64px;
  background: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0; z-index: 100;
}

.logo {
  font-weight: 800;
  color: var(--primary);
  font-size: 1.2rem;
}

.logo span {
  color: #1f2937;
}

.user-zone { display: flex; align-items: center; gap: 20px; }
.user-info { display: flex; align-items: center; gap: 10px; font-weight: 500; }
.mini-avatar { width: 32px; height: 32px; border-radius: 50%; border: 2px solid var(--primary); }

.btn-logout {
  padding: 6px 16px; border-radius: 8px; border: 1px solid #ef4444;
  color: #ef4444; background: transparent; cursor: pointer; transition: 0.2s;
}
.btn-logout:hover { background: #ef4444; color: white; }

.layout-body { display: flex; flex: 1; }

.sidebar {
  width: 240px;
  background: white;
  border-right: 1px solid var(--border-color);
  padding: 24px 12px;
  position: sticky;
  top: 64px;
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
}

.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.nav-item {
  display: block;
  padding: 12px 16px;
  color: var(--text-muted);
  text-decoration: none;
  border-radius: 8px;
  margin-bottom: 4px;
  font-weight: 500;
  transition: 0.2s;
}
.nav-item:hover { background: var(--bg-app); color: var(--primary); }
.router-link-active { background: #eff6ff; color: var(--primary); }

.sidebar-footer {
  margin-top: auto;
  border-top: 1px solid var(--border-color);
  padding-top: 16px;
}

.content-area { flex: 1; padding: 32px; background: var(--bg-app); }

@media (max-width: 768px) {
  .sidebar { display: none; }
}
</style>