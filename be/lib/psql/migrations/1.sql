-- Start a transaction
BEGIN;

-- Create necessary extensions
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create USER_GENDER enum type if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_gender') THEN
        CREATE TYPE user_gender AS ENUM ('MALE', 'FEMALE', 'UNKNOWN');
    END IF;
END $$;

-- Create USER_TYPE enum type if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_type') THEN
        CREATE TYPE user_type AS ENUM ('USER', 'AGENT', 'ADMIN');
    END IF;
END $$;

-- Create ACTION_TYPE enum type if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'action_type') THEN
        CREATE TYPE action_type AS ENUM (
            'user.successful_login', 
            'user.failed_login', 
            'user.signup', 
            'user.pw_reset',
            'user.profile_update',
            'agent.signup',
            'agent.successful_login',
            'agent.failed_login'
            );
    END IF;
END $$;

-- Create the users table
CREATE TABLE IF NOT EXISTS u (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    full_name CITEXT NOT NULL,
    gender user_gender NOT NULL DEFAULT 'UNKNOWN',
    dob TIMESTAMPTZ,
    phone CITEXT,

    email CITEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role user_type NOT NULL DEFAULT 'USER',

    active BOOLEAN NOT NULL DEFAULT TRUE,
    banned BOOLEAN NOT NULL DEFAULT FALSE,

    num_user_agents INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create the contact table
CREATE TABLE IF NOT EXISTS contact (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    u_id UUID NOT NULL,
    u_agent_id UUID,
    acc_id UUID,
    contact_number VARCHAR(15) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create the AuditLog table
CREATE TABLE IF NOT EXISTS auditlog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    u_id UUID NOT NULL,
    action action_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Commit the transaction
COMMIT;
