const userParticipant = {
    id: "uuid-1",
    email: "user@example.com",
    full_name: "Иван1 Иванов",
    role: "participant",
    city: "Москва",
    age: 25,
    direction: "IT",
    avatar_url: "https://i.pravatar.cc/150?u=1",
    rating: {
        total_points: 1250,
        global_rank: 15,
        direction_rank: 3,
        next_level_points: 250
    },
    portfolio: {
        events_count: 12,
        total_points: 1250,
        events: [
            {id: "e-1", title: "Хакатон по AI", date: "2026-03-15T10:00:00Z", points_earned: 100, status: "attended"},
            {id: "e-1", title: "Хакатон по AI", date: "2026-03-15T10:00:00Z", points_earned: 100, status: "attended"},
            {id: "e-1", title: "Хакатон по AI", date: "2026-03-15T10:00:00Z", points_earned: 100, status: "attended"},
            {id: "e-1", title: "Хакатон по AI", date: "2026-03-15T10:00:00Z", points_earned: 100, status: "attended"}
        ]
    }
};

const userOrganizer = {
    id: "uuid-2",
    email: "organizer@example.com",
    full_name: "Петр Петров",
    role: "organizer",
    city: "Москва",
    avatar_url: "https://i.pravatar.cc/150?u=2",
    organizer_stats: {
        events_count: 12,
        trust_rating: 4.7,
        common_prizes: ["мерч", "стажировки", "билеты на форумы"],
        recent_events: [
            {id: "e-2", title: "Мастер-класс по Go", date: "2026-03-20T10:00:00Z", participants_count: 45}
        ]
    }
};

const eventsList = {
    events: [
        {
            id: "e-1",
            title: "Хакатон по AI",
            description: "Краткое описание",
            event_date: "2026-03-25T10:00:00Z",
            format: "offline",
            direction: "IT",
            points: 50,
            prizes: [{type: "merch", name: "Футболка"}],
            organizer: {id: "uuid-2", full_name: "Петр Петров", trust_rating: 4.7},
            participants_count: 45,
            status: "published"
        },
        {
            id: "e-1",
            title: "Хакатон по AI",
            description: "Краткое описание",
            event_date: "2026-03-25T10:00:00Z",
            format: "offline",
            direction: "IT",
            points: 50,
            prizes: [{type: "merch", name: "Футболка"}],
            organizer: {id: "uuid-2", full_name: "Петр Петров", trust_rating: 4.7},
            participants_count: 45,
            status: "published"
        },
        {
            id: "e-1",
            title: "Хакатон по AI",
            description: "Краткое описание",
            event_date: "2026-03-25T10:00:00Z",
            format: "offline",
            direction: "IT",
            points: 50,
            prizes: [{type: "merch", name: "Футболка"}],
            organizer: {id: "uuid-2", full_name: "Петр Петров", trust_rating: 4.7},
            participants_count: 45,
            status: "published"
        },
        {
            id: "e-1",
            title: "Хакатон по AI",
            description: "Краткое описание",
            event_date: "2026-03-25T10:00:00Z",
            format: "offline",
            direction: "IT",
            points: 50,
            prizes: [{type: "merch", name: "Футболка"}],
            organizer: {id: "uuid-2", full_name: "Петр Петров", trust_rating: 4.7},
            participants_count: 45,
            status: "published"
        },
        {
            id: "e-1",
            title: "Хакатон по AI",
            description: "Краткое описание",
            event_date: "2026-03-25T10:00:00Z",
            format: "offline",
            direction: "IT",
            points: 50,
            prizes: [{type: "merch", name: "Футболка"}],
            organizer: {id: "uuid-2", full_name: "Петр Петров", trust_rating: 4.7},
            participants_count: 45,
            status: "published"
        }
    ],
    pagination: {page: 1, limit: 20, total: 150, pages: 8}
};

const eventDetail = {
    ...eventsList.events[0],
    description: "Полное описание с программой",
    prizes: [
        {type: "merch", name: "Футболка"},
        {type: "internship", name: "Стажировка в Думе"}
    ],
    participants: [
        {id: "uuid-1", full_name: "Иван Иванов", status: "confirmed"}
    ],
    user_status: "registered"
};

const globalRating = {
    rating: [
        {avatar_url: "https://i.pravatar.cc/150?u=2",rank: 1, user_id: "uuid-1", full_name: "Иван Иванов", total_points: 3250, direction: "IT", events_count: 24},
        {avatar_url: "https://i.pravatar.cc/150?u=2",rank: 2, user_id: "uuid-1", full_name: "Иван Иванов", total_points: 3250, direction: "IT", events_count: 24},
        {
            avatar_url: "https://i.pravatar.cc/150?u=2",
            rank: 3,
            user_id: "uuid-3",
            full_name: "Анна Смирнова",
            total_points: 3100,
            direction: "media",
            events_count: 20
        },
        {
            avatar_url: "https://i.pravatar.cc/150?u=2",
            rank: 4,
            user_id: "uuid-3",
            full_name: "Анна Смирнова",
            total_points: 3100,
            direction: "media",
            events_count: 20
        },
        {
            avatar_url: "https://i.pravatar.cc/150?u=2",
            rank: 5,
            user_id: "uuid-3",
            full_name: "Анна Смирнова",
            total_points: 3100,
            direction: "media",
            events_count: 20
        },
        {
            avatar_url: "https://i.pravatar.cc/150?u=2",
            rank: 6,
            user_id: "uuid-3",
            full_name: "Анна Смирнова",
            total_points: 3100,
            direction: "media",
            events_count: 20
        },
        {
            avatar_url: "https://i.pravatar.cc/150?u=2",
            rank: 7,
            user_id: "uuid-3",
            full_name: "User Смирнова",
            total_points: 3100,
            direction: "media",
            events_count: 20
        },
        {
            avatar_url: "https://i.pravatar.cc/150?u=2",
            rank: 8,
            user_id: "uuid-3",
            full_name: "Анна Смирнова",
            total_points: 3100,
            direction: "media",
            events_count: 20
        },
    ],
    user_rank: {rank: 7, total_points: 1250}
};

const dashboardActivity = {
    recent_events: [
        {
            id: "e-1",
            title: "Хакатон по AI",
            event_date: "2026-03-25T10:00:00Z",
            organizer_name: "Петр Петров",
            participants_count: 45
        }
    ],
    rating_history: [
        {date: "2026-03-01", points: 100},
        {date: "2026-03-05", points: 250},
        {date: "2026-03-10", points: 500},
        {date: "2026-03-15", points: 1250}
    ],
    trending_tags: [
        {tag: "IT", count: 45},
        {tag: "хакатон", count: 32},
        {tag: "мастер-класс", count: 28}
    ]
};

const hrCandidates = {
    candidates: [
        {
            id: "uuid-1",
            full_name: "Иван Иванов",
            age: 25,
            city: "Москва",
            direction: "IT",
            total_points: 1250,
            events_count: 12,
            rank: 15
        },
        {
            id: "uuid-1",
            full_name: "Иван Иванов",
            age: 25,
            city: "Москва",
            direction: "IT",
            total_points: 1250,
            events_count: 12,
            rank: 15
        },
        {
            id: "uuid-1",
            full_name: "Иван Иванов",
            age: 25,
            city: "Москва",
            direction: "IT",
            total_points: 1250,
            events_count: 12,
            rank: 15
        }
    ],
    total: 45, page: 1, pages: 5
};

export const mocks = {
    userParticipant,
    userOrganizer,
    eventsList,
    eventDetail,
    globalRating,
    dashboardActivity,
    hrCandidates
};