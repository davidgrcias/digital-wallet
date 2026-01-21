-- Database Schema

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    balance DECIMAL(18,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('WITHDRAW', 'DEPOSIT')),
    amount DECIMAL(18,2) NOT NULL CHECK (amount > 0),
    balance_before DECIMAL(18,2) NOT NULL,
    balance_after DECIMAL(18,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_users_email ON users(email);

-- Seed data for testing
INSERT INTO users (id, name, email, balance) VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'John Doe', 'john@example.com', 1000000.00),
    ('550e8400-e29b-41d4-a716-446655440002', 'Jane Smith', 'jane@example.com', 500000.00),
    ('550e8400-e29b-41d4-a716-446655440003', 'Bob Wilson', 'bob@example.com', 250000.00)
ON CONFLICT (id) DO NOTHING;
