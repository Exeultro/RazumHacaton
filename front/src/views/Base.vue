<template>
  <div class="landing-page">
    <!-- Шапка -->
    <header class="landing-header">
      <div class="container">
        <div class="logo">Парламент<span>Резерв</span></div>
        <nav class="nav-links">
          <router-link to="/login" class="btn-outline">Войти</router-link>
          <router-link to="/register" class="btn-primary">Регистрация</router-link>
        </nav>
      </div>
    </header>

    <!-- Герой -->
    <section class="hero">
      <div class="container">
        <div class="hero-content">
          <h1>Развивай свои навыки и становись лидером</h1>
          <p>Платформа для участников и организаторов мероприятий. Участвуй, получай баллы, поднимайся в рейтинге и попади в кадровый резерв.</p>
          <div class="hero-buttons">
            <router-link to="/register" class="btn-primary large">Начать путь</router-link>
            <router-link to="/login" class="btn-outline large">Уже с нами</router-link>
          </div>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">500+</span>
            <span class="stat-label">мероприятий</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">1000+</span>
            <span class="stat-label">участников</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">50+</span>
            <span class="stat-label">организаторов</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Преимущества -->
    <section class="features">
      <div class="container">
        <h2 class="section-title">Для кого платформа</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">👥</div>
            <h3>Участникам</h3>
            <p>Участвуйте в мероприятиях, зарабатывайте баллы, поднимайтесь в рейтинге и получайте приглашения в кадровый резерв.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">🎯</div>
            <h3>Организаторам</h3>
            <p>Создавайте мероприятия, подтверждайте участие, начисляйте баллы и стройте рейтинг доверия.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">📊</div>
            <h3>Наблюдателям</h3>
            <p>Отслеживайте прогресс участников, анализируйте активность и формируйте кадровый резерв.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Как это работает -->
    <section class="how-it-works">
      <div class="container">
        <h2 class="section-title">Как это работает</h2>
        <div class="steps">
          <div class="step">
            <div class="step-number">1</div>
            <h3>Регистрация</h3>
            <p>Заполните профиль, выберите роль (участник/организатор).</p>
          </div>
          <div class="step">
            <div class="step-number">2</div>
            <h3>Участие или организация</h3>
            <p>Находите мероприятия или создавайте свои, привлекайте участников.</p>
          </div>
          <div class="step">
            <div class="step-number">3</div>
            <h3>Подтверждение</h3>
            <p>Организатор сканирует QR-код участника, начисляются баллы.</p>
          </div>
          <div class="step">
            <div class="step-number">4</div>
            <h3>Рейтинг и резерв</h3>
            <p>Поднимайтесь в рейтинге, попадайте в кадровый резерв.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Призыв к действию -->
    <section class="cta">
      <div class="container">
        <div class="cta-card">
          <h2>Готов стать частью сообщества?</h2>
          <p>Присоединяйся к платформе и начни свой путь к успеху уже сегодня.</p>
          <router-link to="/register" class="btn-primary large">Зарегистрироваться</router-link>
        </div>
      </div>
    </section>

    <!-- Футер -->
    <footer class="landing-footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-logo">Парламент<span>Резерв</span></div>
          <div class="footer-links">
            <router-link to="/terms">Пользовательское соглашение</router-link>
            <a href="mailto:support@parmrezerv.ru">support@parmrezerv.ru</a>
          </div>
        </div>
        <div class="footer-copyright">
          © 2025 ПарламентРезерв. Все права защищены.
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
// Landing page – не требует авторизации
</script>

<style scoped>
.landing-page {
  background: var(--bg-app);
  font-family: 'Inter', sans-serif;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Шапка */
.landing-header {
  background: white;
  border-bottom: 1px solid var(--border-color);
  padding: 16px 0;
  position: sticky;
  top: 0;
  z-index: 100;
}
.landing-header .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.logo {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--primary);
}
.logo span {
  color: #1f2937;
}
.nav-links {
  display: flex;
  gap: 16px;
}
.btn-outline, .btn-primary {
  padding: 8px 20px;
  border-radius: 8px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.2s;
}
.btn-outline {
  border: 1px solid var(--border-color);
  color: var(--text-main);
  background: transparent;
}
.btn-outline:hover {
  background: var(--bg-app);
  border-color: var(--primary);
  color: var(--primary);
}
.btn-primary {
  background: var(--primary);
  color: white;
}
.btn-primary:hover {
  background: var(--primary-hover);
}
.large {
  padding: 12px 28px;
  font-size: 1rem;
}

/* Герой */
.hero {
  padding: 80px 0;
  background: linear-gradient(135deg, #f8fafc 0%, #ffffff 100%);
}
.hero .container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 48px;
  align-items: center;
}
.hero-content h1 {
  font-size: 3rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 24px;
  line-height: 1.2;
}
.hero-content p {
  font-size: 1.125rem;
  color: var(--text-muted);
  margin-bottom: 32px;
}
.hero-buttons {
  display: flex;
  gap: 16px;
}
.hero-stats {
  display: flex;
  justify-content: space-between;
  gap: 24px;
}
.stat-item {
  text-align: center;
  background: white;
  padding: 24px;
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  flex: 1;
}
.stat-number {
  font-size: 2rem;
  font-weight: 800;
  color: var(--primary);
  display: block;
}
.stat-label {
  font-size: 0.875rem;
  color: var(--text-muted);
}

/* Преимущества */
.features {
  padding: 80px 0;
}
.section-title {
  text-align: center;
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 48px;
}
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 32px;
}
.feature-card {
  background: white;
  padding: 32px;
  border-radius: 16px;
  text-align: center;
  border: 1px solid var(--border-color);
  transition: transform 0.2s, box-shadow 0.2s;
}
.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}
.feature-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}
.feature-card h3 {
  font-size: 1.5rem;
  margin-bottom: 12px;
}
.feature-card p {
  color: var(--text-muted);
  line-height: 1.5;
}

/* Как это работает */
.how-it-works {
  background: white;
  padding: 80px 0;
}
.steps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 32px;
}
.step {
  text-align: center;
}
.step-number {
  width: 48px;
  height: 48px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 20px;
}
.step h3 {
  font-size: 1.25rem;
  margin-bottom: 12px;
}
.step p {
  color: var(--text-muted);
}

/* Призыв к действию */
.cta {
  padding: 80px 0;
}
.cta-card {
  background: linear-gradient(135deg, var(--primary) 0%, #7c3aed 100%);
  border-radius: 24px;
  padding: 48px;
  text-align: center;
  color: white;
}
.cta-card h2 {
  font-size: 2rem;
  margin-bottom: 16px;
}
.cta-card p {
  font-size: 1.125rem;
  margin-bottom: 32px;
  opacity: 0.9;
}
.cta-card .btn-primary {
  background: white;
  color: var(--primary);
}
.cta-card .btn-primary:hover {
  background: #f8fafc;
}

/* Футер */
.landing-footer {
  background: #0f172a;
  color: #94a3b8;
  padding: 48px 0 24px;
}
.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 32px;
}
.footer-logo {
  font-size: 1.5rem;
  font-weight: 800;
  color: white;
}
.footer-logo span {
  color: var(--primary);
}
.footer-links {
  display: flex;
  gap: 24px;
}
.footer-links a {
  color: #cbd5e1;
  text-decoration: none;
  transition: color 0.2s;
}
.footer-links a:hover {
  color: white;
}
.footer-copyright {
  text-align: center;
  font-size: 0.875rem;
  padding-top: 24px;
  border-top: 1px solid #1e293b;
}

/* Адаптивность */
@media (max-width: 768px) {
  .hero .container {
    grid-template-columns: 1fr;
    text-align: center;
  }
  .hero-stats {
    flex-direction: column;
    align-items: center;
  }
  .stat-item {
    width: 100%;
  }
  .hero-content h1 {
    font-size: 2rem;
  }
  .hero-buttons {
    justify-content: center;
  }
  .footer-content {
    flex-direction: column;
    text-align: center;
  }
  .footer-links {
    flex-direction: column;
    align-items: center;
  }
}
</style>