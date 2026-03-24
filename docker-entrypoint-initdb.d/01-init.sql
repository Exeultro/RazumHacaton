-- Создание таблиц
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('organizer', 'participant', 'observer', 'admin')),
    city VARCHAR(100),
    age INTEGER CHECK (age >= 14 AND age <= 120),
    direction VARCHAR(50) CHECK (direction IN ('IT', 'social', 'media')),
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organizers (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    rating DECIMAL(3,2) DEFAULT 0,
    events_count INTEGER DEFAULT 0,
    common_prizes TEXT[] DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMP NOT NULL,
    registration_deadline TIMESTAMP NOT NULL,
    location VARCHAR(255),
    format VARCHAR(50) CHECK (format IN ('offline', 'online', 'hybrid')),
    direction VARCHAR(50) CHECK (direction IN ('IT', 'social', 'media')),
    difficulty_coefficient DECIMAL(3,2) DEFAULT 1.0,
    points_for_participation INTEGER DEFAULT 10,
    prizes JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_participations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'registered' CHECK (status IN ('registered', 'confirmed', 'attended', 'cancelled')),
    qr_code_token UUID UNIQUE DEFAULT gen_random_uuid(),
    points_earned INTEGER DEFAULT 0,
    attended_at TIMESTAMP,
    confirmed_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, user_id)
);

CREATE TABLE IF NOT EXISTS organizer_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE SET NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(participant_id, organizer_id, event_id)
);

CREATE TABLE IF NOT EXISTS observer_filters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filter_name VARCHAR(100),
    filters JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rating_cache (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_points INTEGER DEFAULT 0,
    events_count INTEGER DEFAULT 0,
    top_direction VARCHAR(50),
    global_rank INTEGER,
    it_rank INTEGER,
    social_rank INTEGER,
    media_rank INTEGER,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS points_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE SET NULL,
    points_change INTEGER NOT NULL,
    reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS difficulty_settings (
    id SERIAL PRIMARY KEY,
    direction VARCHAR(50) UNIQUE NOT NULL,
    coefficient DECIMAL(3,2) NOT NULL DEFAULT 1.0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO difficulty_settings (direction, coefficient) VALUES 
    ('IT', 2.0),
    ('social', 1.5),
    ('media', 1.2)
ON CONFLICT (direction) DO NOTHING;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_events_updated_at ON events;
CREATE TRIGGER trigger_events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE FUNCTION update_rating_cache()
RETURNS VOID AS $$
BEGIN
    TRUNCATE rating_cache;
    
    INSERT INTO rating_cache (user_id, total_points, events_count, top_direction)
    SELECT 
        u.id,
        COALESCE(SUM(ep.points_earned), 0) as total_points,
        COUNT(CASE WHEN ep.status = 'attended' THEN 1 END) as events_count,
        (
            SELECT direction
            FROM (
                SELECT 
                    e.direction,
                    SUM(ep2.points_earned) as points_by_direction
                FROM event_participations ep2
                JOIN events e ON ep2.event_id = e.id
                WHERE ep2.user_id = u.id AND ep2.status = 'attended'
                GROUP BY e.direction
                ORDER BY points_by_direction DESC
                LIMIT 1
            ) sub
        ) as top_direction
    FROM users u
    LEFT JOIN event_participations ep ON u.id = ep.user_id AND ep.status = 'attended'
    GROUP BY u.id;
    
    UPDATE rating_cache rc
    SET global_rank = sub.rank
    FROM (
        SELECT user_id, ROW_NUMBER() OVER (ORDER BY total_points DESC) as rank
        FROM rating_cache
    ) sub
    WHERE rc.user_id = sub.user_id;
    
END;
$$ LANGUAGE plpgsql;

-- Вставка тестовых пользователей 
INSERT INTO users (id, email, password, full_name, role, city, age, direction, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'admin@razum.ru', 'admin123', 'Системный Администратор', 'admin', 'Москва', NULL, NULL, NOW(), NOW()),
    (gen_random_uuid(), 'organizer@razum.ru', 'password123', 'Тестовый Организатор', 'organizer', 'Москва', NULL, NULL, NOW(), NOW()),
    (gen_random_uuid(), 'participant@razum.ru', 'password123', 'Тестовый Участник', 'participant', 'Москва', 25, 'IT', NOW(), NOW()),
    (gen_random_uuid(), 'observer@razum.ru', 'password123', 'Тестовый Наблюдатель', 'observer', 'Москва', NULL, NULL, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Создаем запись в organizers для организатора
INSERT INTO organizers (user_id, rating, events_count, common_prizes)
SELECT id, 0, 0, '{}' FROM users WHERE email = 'organizer@razum.ru'
ON CONFLICT (user_id) DO NOTHING;

-- Создаем тестовое мероприятие
DO $$
DECLARE
    org_id UUID;
BEGIN
    SELECT id INTO org_id FROM users WHERE email = 'organizer@razum.ru';
    
    INSERT INTO events (id, organizer_id, title, description, event_date, registration_deadline, 
                        location, format, direction, difficulty_coefficient, points_for_participation, 
                        prizes, status, created_at, updated_at)
    VALUES (
        gen_random_uuid(),
        org_id,
        'Хакатон по Go',
        'Крутой хакатон для разработчиков',
        '2026-04-15 10:00:00',
        '2026-04-10 23:59:59',
        'Москва, ул. Тверская 1',
        'offline',
        'IT',
        2.0,
        75,
        '[{"type": "merch", "name": "Футболка"}, {"type": "internship", "name": "Стажировка"}]'::jsonb,
        'published',
        NOW(),
        NOW()
    );
END $$;

-- Обновляем кэш рейтинга
SELECT update_rating_cache();

-- Добавляем участие участника в мероприятии с начислением баллов
DO $$
DECLARE
    participant_id UUID;
    event_id UUID;
    organizer_id UUID;
BEGIN
    -- Получаем ID участника
    SELECT id INTO participant_id FROM users WHERE email = 'participant@razum.ru';
    
    -- Получаем ID мероприятия
    SELECT id INTO event_id FROM events WHERE title = 'Хакатон по Go' LIMIT 1;
    
    -- Получаем ID организатора
    SELECT id INTO organizer_id FROM users WHERE email = 'organizer@razum.ru';
    
    -- Добавляем участие со статусом attended и начисленными баллами (75 * 2.0 = 150)
    INSERT INTO event_participations (id, event_id, user_id, status, points_earned, attended_at, confirmed_by, created_at)
    VALUES (
        gen_random_uuid(),
        event_id,
        participant_id,
        'attended',
        150,
        NOW(),
        organizer_id,
        NOW()
    );
    
    -- Добавляем запись в аудит
    INSERT INTO points_audit (id, user_id, event_id, points_change, reason, created_at)
    VALUES (
        gen_random_uuid(),
        participant_id,
        event_id,
        150,
        'Участие в мероприятии: Хакатон по Go',
        NOW()
    );
END $$;

-- Обновляем кэш рейтинга
SELECT update_rating_cache();