import MockAdapter from 'axios-mock-adapter';
import { mocks } from './mockData';

export function setupMocks(axiosInstance) {
    const mock = new MockAdapter(axiosInstance, { delayResponse: 500 });

    // --- Авторизация ---
    mock.onPost('/auth/login').reply(200, { user: mocks.userParticipant, token: 'mock-jwt-token' });
    mock.onPost('/auth/register').reply(200, { user: mocks.userParticipant, token: 'mock-jwt-token' });
    mock.onPost('/auth/logout').reply(200, { success: true });

    // --- Профиль ---
    mock.onGet('/profile').reply(200, mocks.userParticipant);
    mock.onPut('/profile').reply(200, mocks.userParticipant);
    mock.onGet(/\/profile\/[a-zA-Z0-9-]+/).reply(200, mocks.userOrganizer);

    // --- Мероприятия ---
    // Порядок важен: сначала точные совпадения и регулярки с ID, потом общие списки
    mock.onPost(/\/events\/[a-zA-Z0-9-]+\/register/).reply(200, { success: true, message: "Вы успешно зарегистрированы" });
    mock.onPost(/\/events\/[a-zA-Z0-9-]+\/confirm/).reply(200, { success: true, points_earned: 100 });
    mock.onGet(/\/events\/[a-zA-Z0-9-]+\/participants/).reply(200, { participants: mocks.eventDetail.participants });
    mock.onGet(/\/events\/[a-zA-Z0-9-]+/).reply(200, mocks.eventDetail);

    mock.onGet('/events').reply(200, mocks.eventsList);
    mock.onPost('/events').reply(200, mocks.eventDetail);
    mock.onPut(/\/events\/[a-zA-Z0-9-]+/).reply(200, mocks.eventDetail);
    mock.onDelete(/\/events\/[a-zA-Z0-9-]+/).reply(200, { success: true });

    // --- Рейтинг ---
    mock.onGet('/rating/global').reply(200, mocks.globalRating);
    mock.onGet(/\/rating\/user\/[a-zA-Z0-9-]+/).reply(200, { ...mocks.userParticipant.rating, user_id: "uuid-1", full_name: "Иван Иванов" });

    // --- Дашборд и Кадры ---
    mock.onGet('/dashboard/activity').reply(200, mocks.dashboardActivity);
    mock.onGet('/cadre/candidates').reply(200, mocks.hrCandidates);

    // --- Отзывы ---
    mock.onGet(/\/organizers\/[a-zA-Z0-9-]+\/reviews/).reply(200, { average_rating: 4.7, reviews: [] });
    mock.onPost(/\/organizers\/[a-zA-Z0-9-]+\/reviews/).reply(200, { success: true });


    mock.onGet('/cadre/filters').reply(200, [
        { id: 'filter-1', name: 'Молодые IT-специалисты', filters: { age_min: 18, age_max: 30, direction: 'IT', min_points: 500 } },
        { id: 'filter-2', name: 'Активные медийщики', filters: { direction: 'media', min_events: 5, sort_by: 'points' } }
    ]);

    mock.onPost('/cadre/filters').reply((config) => {
        const body = JSON.parse(config.data);
        const newFilter = {
            id: `filter-${Date.now()}`,
            name: body.name,
            filters: body.filters,
        };
        return [200, newFilter];
    });

// --- Админка ---
    mock.onGet('/admin/organizers/pending').reply(200, {
        pending_organizers: [
            { user_id: 'org-1', full_name: 'Петр Петров', email: 'petr@example.com', registered_at: '2026-03-01T10:00:00Z' },
            { user_id: 'org-2', full_name: 'Анна Смирнова', email: 'anna@example.com', registered_at: '2026-03-02T12:00:00Z' }
        ]
    });

    mock.onPost(/\/admin\/organizers\/[a-zA-Z0-9-]+\/approve/).reply(200, { success: true });
    mock.onPost(/\/admin\/organizers\/[a-zA-Z0-9-]+\/reject/).reply(200, { success: true });

    mock.onPut('/admin/settings/difficulty').reply(200, (config) => {
        const { coefficients } = JSON.parse(config.data);
        return [200, { coefficients }];
    });
    mock.onGet('/admin/stats').reply(200, {
        total_users: 250,
        total_organizers: 45,
        total_events: 120,
        total_participations: 1500,
        events_by_direction: { IT: 50, social: 40, media: 30 }
    });

    // --- Fallback для ненайденных роутов ---
    mock.onAny().reply(404, { error: "Mock not found", code: 404 });

}