<template>
  <div class="auth-page">
    <div class="auth-container card">
      <div class="auth-header">
        <h1>Парламент<span>Резерв</span></h1>
        <p>Создание аккаунта</p>
      </div>

      <form @submit.prevent="onSubmit">
        <!-- Красивый переключатель ролей -->
        <div class="role-selector mb-4">
          <label class="form-label">Я регистрируюсь как:</label>
          <div class="role-buttons">
            <button
                type="button"
                :class="['role-btn', { active: form.role === 'participant' }]"
                @click="form.role = 'participant'"
            >
              Участник (Кандидат)
            </button>
            <button
                type="button"
                :class="['role-btn', { active: form.role === 'organizer' }]"
                @click="form.role = 'organizer'"
            >
              Организатор
            </button>
          </div>
        </div>

        <BaseInput v-model="form.full_name" label="ФИО" placeholder="Иванов Иван Иванович" required class="mb-4" />

        <div class="form-row mb-4">
          <BaseInput v-model="form.email" label="Email" type="email" placeholder="mail@example.com" required />
          <BaseInput v-model="form.password" label="Пароль" type="password" placeholder="••••••••" required />
        </div>

        <BaseInput v-model="form.city" label="Город" placeholder="Например: Москва" required class="mb-4" />

        <!-- Блок профиля с легким фоном -->
        <div v-if="form.role === 'participant'" class="extra-box mb-4">
          <p class="extra-title">Данные кандидата</p>
          <div class="form-row">
            <BaseInput v-model.number="form.age" label="Возраст" type="number" placeholder="Например: 18" required />

            <div class="custom-select-group">
              <label class="form-label">Направление</label>
              <!-- Стилизованный Select -->
              <div class="select-wrapper">
                <select v-model="form.direction" class="custom-select" required>
                  <option value="IT">IT и Инновации</option>
                  <option value="social">Социальное проектирование</option>
                  <option value="media">Медиа и Журналистика</option>
                </select>
              </div>
            </div>
          </div>
        </div>

        <!-- Блок согласия с условиями -->
        <div class="terms-checkbox mb-4">
          <label class="checkbox-label">
            <input
                type="checkbox"
                v-model="agreedToTerms"
                class="checkbox-input"
            />
            <span class="checkbox-text">
              Я соглашаюсь с <router-link to="/terms" class="terms-link">условиями использования</router-link>
            </span>
          </label>
        </div>

        <!-- Кнопка: variant="primary" решает проблему цвета текста -->
        <BaseButton
            type="submit"
            :loading="isLoading"
            :disabled="!agreedToTerms"
            block
            variant="primary"
            style="margin-top: 24px;"
        >
          Зарегистрироваться
        </BaseButton>
      </form>

      <div class="auth-footer">
        <span class="text-muted">Уже есть аккаунт? </span>
        <router-link to="/login" class="link">Войти в систему</router-link>
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

const form = ref({
  full_name: '', email: '', password: '', city: '', role: 'participant', age: '', direction: 'IT'
});
const agreedToTerms = ref(false);
const isLoading = ref(false);

const onSubmit = async () => {
  if (!agreedToTerms.value) {
    uiStore.addNotification({ type: 'error', message: 'Необходимо принять условия использования' });
    return;
  }

  isLoading.value = true;
  try {
    const payload = { ...form.value };
    if (payload.role !== 'participant') {
      delete payload.age;
      delete payload.direction;
    }
    await authStore.register(payload);
    router.push('/dashboard');
  } catch (error) {
    uiStore.addNotification({ type: 'error', message: 'Ошибка регистрации' });
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
  padding: 40px 20px;
  box-sizing: border-box;
}

.auth-container {
  width: 100%;
  max-width: 540px;
  padding: 40px;
}

.auth-header { text-align: center; margin-bottom: 32px; }
.auth-header h1 { font-size: 2rem; margin: 0 0 8px 0; font-weight: 700; color: var(--text-main); }
.auth-header h1 span { color: var(--primary); }
.auth-header p { color: var(--text-muted); margin: 0; }

.mb-4 { margin-bottom: 16px; }

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-main);
  margin-bottom: 8px;
}

/* Переключатель ролей (стиль iOS) */
.role-buttons {
  display: flex;
  background: #f1f5f9; /* Светло-серый фон */
  border-radius: var(--radius-sm);
  padding: 4px;
}
.role-btn {
  flex: 1;
  padding: 10px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}
.role-btn.active {
  background: #ffffff;
  color: var(--text-main);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

/* Сетка инпутов */
.form-row { display: flex; gap: 16px; }
.form-row > * { flex: 1; min-width: 0; }

/* Дополнительный блок для кандидата */
.extra-box {
  background: #f8fafc; /* Очень легкий серый фон */
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 20px;
}
.extra-title {
  font-size: 0.8rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 16px 0;
  font-weight: 700;
}

/* Красивый Select */
.custom-select-group {
  display: flex;
  flex-direction: column;
  flex: 1;
}
.select-wrapper {
  position: relative;
  width: 100%;
}
.custom-select {
  appearance: none; /* Убираем стандартную стрелку */
  width: 100%;
  height: 44px; /* Синхронизировано с BaseInput */
  padding: 0 36px 0 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  background-color: var(--bg-card);
  color: var(--text-main);
  font-family: inherit;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-sm);
  /* Кастомная SVG стрелочка */
  background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2216%22%20height%3D%2216%22%20viewBox%3D%220%200%2024%2024%22%20fill%3D%22none%22%20stroke%3D%22%2364748b%22%20stroke-width%3D%222%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%3E%3Cpolyline%20points%3D%226%209%2012%2015%2018%209%22%3E%3C%2Fpolyline%3E%3C%2Fsvg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  background-size: 16px;
}
.custom-select:hover { border-color: #cbd5e1; }
.custom-select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
}

/* Подвал */
.auth-footer { margin-top: 32px; text-align: center; font-size: 0.9rem; }
.link { color: var(--primary); font-weight: 600; text-decoration: none; transition: color 0.2s; }
.link:hover { color: var(--primary-hover); text-decoration: underline;}

/* Адаптивность */
@media (max-width: 520px) {
  .form-row { flex-direction: column; gap: 16px; }
  .auth-container { padding: 30px 20px; }
}

.terms-checkbox {
  margin: 16px 0;
}

.checkbox-label {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.checkbox-input {
  width: 18px;
  height: 18px;
  margin-right: 10px;
  cursor: pointer;
  accent-color: var(--primary);
}

.checkbox-text {
  font-size: 0.875rem;
  color: var(--text-main);
  line-height: 1.4;
}

.terms-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.terms-link:hover {
  text-decoration: underline;
}
</style>