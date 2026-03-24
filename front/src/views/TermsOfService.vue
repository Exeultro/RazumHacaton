<template>
  <div class="terms-page">
    <div class="container">
      <div class="card terms-card">
        <h1 class="page-title">Пользовательское соглашение</h1>
        <p class="last-updated">Последнее обновление: 24 марта 2026 г.</p>

        <section class="terms-section">
          <h2 class="section-title">1. Общие положения</h2>
          <p>Настоящее Пользовательское соглашение (далее — Соглашение) регулирует отношения между администрацией сервиса (далее — Администрация) и пользователем (далее — Пользователь) при использовании сервиса «ПарламентРезерв» (далее — Сервис).</p>
          <p>Использование Сервиса означает полное и безоговорочное принятие Пользователем условий настоящего Соглашения. Если Пользователь не согласен с условиями Соглашения, он не вправе использовать Сервис.</p>
        </section>

        <section class="terms-section">
          <h2 class="section-title">2. Термины и определения</h2>
          <ul>
            <li><strong>Сервис</strong> — совокупность программных и аппаратных средств, обеспечивающих функциональность платформы.</li>
            <li><strong>Пользователь</strong> — дееспособное физическое лицо, прошедшее регистрацию и использующее Сервис.</li>
            <li><strong>Контент</strong> — любая информация, размещаемая Пользователем на Сервисе.</li>
          </ul>
        </section>

        <section class="terms-section">
          <h2 class="section-title">3. Регистрация и учетная запись</h2>
          <p>Для доступа к основным функциям Сервиса Пользователь обязан пройти регистрацию, указав достоверные данные. Администрация вправе в любое время запросить подтверждение предоставленной информации.</p>
          <p>Пользователь несет ответственность за сохранность своих учетных данных и все действия, совершенные под его учетной записью. О любом несанкционированном доступе необходимо незамедлительно уведомить Администрацию.</p>
        </section>

        <section class="terms-section">
          <h2 class="section-title">4. Права и обязанности сторон</h2>
          <h3>4.1. Пользователь имеет право:</h3>
          <ul>
            <li>использовать Сервис в соответствии с его функциональным назначением;</li>
            <li>обращаться в службу поддержки по вопросам работы Сервиса;</li>
            <li>удалять свою учетную запись в порядке, установленном Сервисом.</li>
          </ul>
          <h3>4.2. Пользователь обязуется:</h3>
          <ul>
            <li>не нарушать законодательство РФ и настоящее Соглашение;</li>
            <li>не размещать запрещенный контент (порнография, призывы к насилию, оскорбления, спам и т.д.);</li>
            <li>не пытаться получить несанкционированный доступ к Сервису.</li>
          </ul>
          <h3>4.3. Администрация имеет право:</h3>
          <ul>
            <li>изменять функционал Сервиса, условия Соглашения в одностороннем порядке с уведомлением Пользователя;</li>
            <li>блокировать или удалять учетные записи при нарушении Соглашения;</li>
            <li>приостанавливать или прекращать работу Сервиса по техническим причинам.</li>
          </ul>
        </section>

        <section class="terms-section">
          <h2 class="section-title">5. Интеллектуальная собственность</h2>
          <p>Все элементы Сервиса (дизайн, код, базы данных, логотипы) являются объектами интеллектуальной собственности Администрации. Пользователь не вправе копировать, модифицировать или использовать их без письменного разрешения.</p>
          <p>Контент, размещенный Пользователем, остается его собственностью, но Пользователь предоставляет Администрации неисключительную лицензию на его использование в рамках работы Сервиса.</p>
        </section>

        <section class="terms-section">
          <h2 class="section-title">6. Ответственность</h2>
          <p>Сервис предоставляется «как есть». Администрация не несет ответственности за:</p>
          <ul>
            <li>прямые или косвенные убытки, возникшие в результате использования или невозможности использования Сервиса;</li>
            <li>действия других Пользователей;</li>
            <li>сбои в работе, вызванные действиями третьих лиц.</li>
          </ul>
          <p>В случае нарушения Пользователем условий Соглашения, Администрация оставляет за собой право применить санкции вплоть до блокировки учетной записи.</p>
        </section>

        <section class="terms-section">
          <h2 class="section-title">7. Заключительные положения</h2>
          <p>Настоящее Соглашение регулируется законодательством Российской Федерации. Все споры разрешаются в судебном порядке по месту нахождения Администрации.</p>
          <p>Если какое-либо положение Соглашения признается недействительным, это не влечет недействительность остальных положений.</p>
          <p>По всем вопросам, связанным с Соглашением, вы можете обратиться по адресу: <a href="mailto:support@parmrezerv.ru">support@parmrezerv.ru</a></p>
        </section>

        <div class="terms-actions">
          <button class="btn btn-secondary" @click="goBack">Назад</button>
          <router-link to="/" class="btn btn-primary">Главная</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped>
.terms-page {
  min-height: 100vh;
  background-color: var(--bg-app);
  padding: 40px 20px;
}

.container {
  max-width: 900px;
  margin: 0 auto;
}

.terms-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: var(--shadow-md);
}

.last-updated {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin-bottom: 32px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 12px;
}

.terms-section {
  margin-bottom: 32px;
}

.terms-section h2 {
  margin-bottom: 16px;
  font-size: 1.5rem;
}

.terms-section h3 {
  margin: 16px 0 8px;
  font-size: 1.2rem;
}

.terms-section p {
  margin-bottom: 12px;
  line-height: 1.6;
  color: var(--text-main);
}

.terms-section ul {
  margin: 12px 0 12px 24px;
  list-style-type: disc;
}

.terms-section li {
  margin-bottom: 8px;
  line-height: 1.5;
  color: var(--text-main);
}

.terms-actions {
  margin-top: 40px;
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 16px;
  justify-content: center;
}

.btn {
  display: inline-block;
  padding: 10px 20px;
  border-radius: var(--radius-sm);
  font-weight: 500;
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
  border: none;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover {
  background: var(--primary-hover);
}

.btn-secondary {
  background: #f1f5f9;
  color: var(--text-main);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: #e2e8f0;
}

@media (max-width: 768px) {
  .terms-card {
    padding: 20px;
  }
  .terms-section h2 {
    font-size: 1.3rem;
  }
}
</style>