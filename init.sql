-- Chat Application Database Schema
-- PostgreSQL initialization script

-- Create tables
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    profile_picture TEXT
);

CREATE TABLE IF NOT EXISTS rooms (
    room_id TEXT PRIMARY KEY,
    last_message_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS events (
    event_id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL REFERENCES rooms(room_id),
    title TEXT NOT NULL,
    location TEXT,
    capacity INTEGER,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS room_user (
    room_id TEXT NOT NULL REFERENCES rooms(room_id),
    user_id TEXT NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (room_id, user_id)
);

CREATE TABLE IF NOT EXISTS event_participation (
    event_id TEXT NOT NULL REFERENCES events(event_id),
    user_id TEXT NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL REFERENCES rooms(room_id),
    sender_id TEXT NOT NULL REFERENCES users(user_id),
    content TEXT,
    message_type TEXT NOT NULL,
    media_url TEXT,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_read BOOLEAN DEFAULT FALSE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages(sent_at);
CREATE INDEX IF NOT EXISTS idx_room_user_user_id ON room_user(user_id);
CREATE INDEX IF NOT EXISTS idx_event_participation_user_id ON event_participation(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert sample data
INSERT INTO users (user_id, username, email, password, is_online, created_at, updated_at, profile_picture)
VALUES 
    (gen_random_uuid()::text, 'quan', 'quan@gmail.com', 'quanpassword', false, now(), now(), 'profile_pic_url'),
    (gen_random_uuid()::text, 'alice', 'alice@example.com', 'alicepassword', true, now(), now(), 'alice_profile.jpg'),
    (gen_random_uuid()::text, 'bob', 'bob@example.com', 'bobpassword', false, now(), now(), 'bob_profile.jpg');

-- Insert sample rooms (using subqueries to get user IDs)
INSERT INTO rooms (room_id, created_at, created_by)
VALUES 
    (gen_random_uuid()::text, now(), (SELECT user_id FROM users WHERE email = 'quan@gmail.com' LIMIT 1)),
    (gen_random_uuid()::text, now(), (SELECT user_id FROM users WHERE email = 'alice@example.com' LIMIT 1))
ON CONFLICT (room_id) DO NOTHING;

-- Add users to rooms
INSERT INTO room_user (room_id, user_id)
SELECT r.room_id, u.user_id
FROM rooms r, users u
WHERE r.created_by = (SELECT user_id FROM users WHERE email = 'quan@gmail.com' LIMIT 1)
  AND u.email IN ('quan@gmail.com', 'alice@example.com', 'bob@example.com')
ON CONFLICT (room_id, user_id) DO NOTHING;

-- Insert sample messages
INSERT INTO messages (message_id, room_id, sender_id, content, message_type, sent_at, is_read)
SELECT 
    gen_random_uuid()::text,
    r.room_id,
    u.user_id,
    'Hello everyone! Welcome to the chat room.',
    'TEXT',
    now() - interval '1 hour',
    false
FROM rooms r, users u
WHERE r.created_by = (SELECT user_id FROM users WHERE email = 'quan@gmail.com' LIMIT 1)
  AND u.email = 'quan@gmail.com'
LIMIT 1;

INSERT INTO messages (message_id, room_id, sender_id, content, message_type, sent_at, is_read)
SELECT 
    gen_random_uuid()::text,
    r.room_id,
    u.user_id,
    'Thanks for creating this room!',
    'TEXT',
    now() - interval '30 minutes',
    false
FROM rooms r, users u
WHERE r.created_by = (SELECT user_id FROM users WHERE email = 'quan@gmail.com' LIMIT 1)
  AND u.email = 'alice@example.com'
LIMIT 1;

-- Insert sample event
INSERT INTO events (event_id, room_id, title, location, capacity, start_time, end_time, created_at, created_by)
SELECT 
    gen_random_uuid()::text,
    r.room_id,
    'Team Meeting',
    'Conference Room A',
    10,
    now() + interval '1 day',
    now() + interval '1 day 2 hours',
    now(),
    r.created_by
FROM rooms r
WHERE r.created_by = (SELECT user_id FROM users WHERE email = 'quan@gmail.com' LIMIT 1)
LIMIT 1;

-- Add event participants
INSERT INTO event_participation (event_id, user_id)
SELECT e.event_id, u.user_id
FROM events e, users u
WHERE e.title = 'Team Meeting'
  AND u.email IN ('quan@gmail.com', 'alice@example.com')
ON CONFLICT (event_id, user_id) DO NOTHING;

-- Display summary
SELECT 'Database initialized successfully!' as status;
SELECT count(*) as user_count FROM users;
SELECT count(*) as room_count FROM rooms;
SELECT count(*) as message_count FROM messages;
SELECT count(*) as event_count FROM events;
