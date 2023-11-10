# DB Schema

This file lists all planned tables that the DB will have. It is split into two phases for development.

## Phase 1

### Drives Table
- `id` INT
- `title` VARCHAR(100)
- `company` VARCHAR(100)
- `description` TEXT
- `date`

### Roles Table
- `id` INT
- `title` VARCHAR(100)
- `description` TEXT
- `salary` NUMERIC
- `drive_id` FOREIGN KEY(drives.id)

## Phase 2

### Roles Table
- `id` INT
- `name` VARCHAR(50)

### Users Table
- `username` VARCHAR(32)
- `name` VARCHAR(50)
- `email` VARCHAR(255)
- `phone` VARCHAR(10)
- `birthday` DATE
- `role_id` FOREIGN KEY(roles.id)
