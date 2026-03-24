<template>
  <div class="not-found">
    <div class="not-found__content">
      <div class="not-found__code">404</div>
      <h1 class="not-found__title">Страница не найдена</h1>
      <p class="not-found__message">
        К сожалению, запрашиваемая страница не существует или была перемещена.
      </p>
      <div class="not-found__actions">
        <router-link to="/" class="btn btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 9l9-6 9 6v11a2 2 0 0 1-2 2h-5v-7H9v7H4a2 2 0 0 1-2-2z"></path>
            <polyline points="9 22 9 12 15 12 15 22"></polyline>
          </svg>
          На главную
        </router-link>
        <button class="btn btn-outline" @click="goBack">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 12H5M12 19l-7-7 7-7"></path>
          </svg>
          Назад
        </button>
      </div>
    </div>
    <div class="not-found__illustration">
      <div class="not-found__glow"></div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goBack = () => {
  router.go(-1);
};
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f0f9ff 0%, #e6f0fa 100%);
  position: relative;
  overflow: hidden;
}

.not-found__content {
  text-align: center;
  z-index: 2;
  max-width: 500px;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  margin: 1rem;
}

.not-found__code {
  font-size: 8rem;
  font-weight: 800;
  line-height: 1;
  background: linear-gradient(135deg, var(--primary), #60a5fa);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  text-shadow: 2px 2px 20px rgba(0,0,0,0.1);
  margin-bottom: 1rem;
}

.not-found__title {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--text-main);
  margin-bottom: 1rem;
}

.not-found__message {
  color: var(--text-muted);
  margin-bottom: 2rem;
  font-size: 1rem;
}

.not-found__actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.not-found__actions .btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.not-found__illustration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.not-found__glow {
  position: absolute;
  top: 20%;
  left: 10%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(29,78,216,0.15) 0%, rgba(29,78,216,0) 70%);
  border-radius: 50%;
  filter: blur(40px);
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(20px, -20px) scale(1.05); }
  100% { transform: translate(0, 0) scale(1); }
}

@media (max-width: 768px) {
  .not-found__code {
    font-size: 5rem;
  }
  .not-found__title {
    font-size: 1.5rem;
  }
  .not-found__actions {
    flex-direction: column;
  }
  .not-found__actions .btn {
    width: 100%;
  }
}
</style>