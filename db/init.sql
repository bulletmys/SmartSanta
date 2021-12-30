create table if not exists events
(
    event_id   uuid                     not null primary key,
    name       varchar                  not null,
    status     smallint                 not null default 0,
    created_at timestamp with time zone not null default now()
);

create table if not exists users
(
    user_id     uuid                     not null primary key,
    count_id    bigserial                not null unique,
    event_id    uuid                     not null,
    name        varchar                  not null,
    wish        text,
    is_admin    boolean                           default false,
    is_voted    boolean                           default false,
    preferences bigint[],
    created_at  timestamp with time zone not null default now(),
    FOREIGN KEY (event_id) REFERENCES events (event_id)
);

create table if not exists pairs
(
    sender_id   uuid not null,
    receiver_id uuid not null,
    event_id    uuid not null,
    FOREIGN KEY (event_id) REFERENCES events (event_id),
    FOREIGN KEY (sender_id) REFERENCES users (user_id),
    FOREIGN KEY (receiver_id) REFERENCES users (user_id),
    PRIMARY KEY (sender_id, receiver_id, event_id)
);

create table if not exists preferences
(
    sender_id   uuid not null,
    receiver_id uuid not null,
    event_id    uuid not null,
    FOREIGN KEY (event_id) REFERENCES events (event_id),
    FOREIGN KEY (sender_id) REFERENCES users (user_id),
    FOREIGN KEY (receiver_id) REFERENCES users (user_id),
    PRIMARY KEY (sender_id, receiver_id, event_id)
);

CREATE INDEX IF NOT EXISTS idx_users_events ON users (event_id);
CREATE INDEX IF NOT EXISTS idx_voted_users ON users (event_id, is_voted);
