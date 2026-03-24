<template>
  <div class="auth-page">
    <div class="auth-container card">
      <div class="auth-header">
        <h1>Парламент<span>Резерв</span></h1>
        <p>Вход в систему</p>
      </div>

      <form @submit.prevent="onSubmit">
        <BaseInput
            v-model="form.email"
            label="Email"
            type="email"
            placeholder="mail@example.com"
            class="mb-4"
        />

        <BaseInput
            v-model="form.password"
            label="Пароль"
            type="password"
            placeholder="••••••••"
            class="mb-4"
        />

        <!-- Используем variant="primary" вместо class="btn-primary" -->
        <BaseButton
            type="submit"
            :loading="isLoading"
            block
            variant="primary"
            style="margin-top: 24px;"
        >
          Войти
        </BaseButton>
      </form>

      <div class="auth-footer">
        <span class="text-muted">Нет аккаунта? </span>
        <router-link to="/register" class="link">Зарегистрироваться</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useUiStore } from '../stores/ui';
import BaseInput from '../components/ui/BaseInput.vue';
import BaseButton from '../components/ui/BaseButton.vue';

const router = useRouter();
const authStore = useAuthStore();
const uiStore = useUiStore();

const form = ref({ email: '', password: '' });
const isLoading = ref(false);

const onSubmit = async () => {
  isLoading.value = true;
  try {
    await authStore.login(form.value);
    router.push('/dashboard');
  } catch (error) {
    uiStore.addNotification({ type: 'error', message: 'Неверный email или пароль' });
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.auth-container {
  width: 100%;
  max-width: 420px;
  padding: 40px 32px; /* Увеличили отступы внутри карточки */
}
.auth-header {
  text-align: center;
  margin-bottom: 32px;
}
.auth-header h1 {
  font-size: 2rem;
  margin: 0 0 8px 0;
  font-weight: 700;
  color: var(--text-main);
}
.auth-header h1 span {
  color: var(--primary);
}
.auth-header p {
  color: var(--text-muted);
  margin: 0;
}
.auth-footer {
  margin-top: 32px;
  text-align: center;
  font-size: 0.9rem;
}
.link {
  color: var(--primary);
  font-weight: 600;
  text-decoration: none;
  transition: color 0.2s;
}
.link:hover {
  color: var(--primary-hover);
  text-decoration: underline;
}
.mb-4 { margin-bottom: 16px; }
</style>