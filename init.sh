#!/usr/bin/env bash

# This file serves as a collection of example commands used to interact with the API
# Due to the way postgres generate pk's it doesn't play well with that database. It works just fine with Sqlite though

set -x

# Create event
export EVENT_TOKEN=$(echo '{
  "title": "title",
  "email": "test@mail.com",
  "date": "date"
}' | http :8080/events | jq -r '.token')

# Create group
echo '{
  "datetime": "datetime",
  "maxParticipants": 25
}' | http :8080/events/groups?token=$EVENT_TOKEN

# Create another group
echo '{
  "datetime": "datetime2",
  "maxParticipants": 25
}' | http :8080/events/groups?token=$EVENT_TOKEN

# Create participant
export PARTICIPANT_TOKEN=$(echo '{
  "name": "name",
  "email": "test1@mail.com"
}' | http :8080/participants?token=$EVENT_TOKEN | jq -r '.token')

# Find groups with participant count
http ":8080/events/groups-count?id=1&token=$PARTICIPANT_TOKEN"

# Add participant to group
http post ":8080/participants/groups?id=&token=$PARTICIPANT_TOKEN"

# Create second participant
export PARTICIPANT2_TOKEN=$(echo '{
  "name": "name2",
  "email": "test2@mail.com"
}' | http :8080/participants?token=$EVENT_TOKEN | jq -r '.token')

# Add second participant to second group
http post ":8080/participants/groups?id=2&token=$PARTICIPANT2_TOKEN"

# Find groups with participant count
http ":8080/events/groups-count?id=1&token=$PARTICIPANT_TOKEN"

# Create third participant
export PARTICIPANT3_TOKEN=$(echo '{
  "name": "name3",
  "email": "test3@mail.com"
}' | http :8080/participants?token=$EVENT_TOKEN | jq -r '.token')

# Add third participant to first group
http post ":8080/participants/groups?id=1&token=$PARTICIPANT3_TOKEN"

# Find groups with participant count
http ":8080/events/groups-count?id=1&token=$PARTICIPANT_TOKEN"

# Get event with groups and participants
http ":8080/events?token=$EVENT_TOKEN"

# Create fourth participant (which isn't in a group)
export PARTICIPANT4_TOKEN=$(echo '{
  "name": "name4",
  "email": "test4@mail.com"
}' | http :8080/participants?token=$EVENT_TOKEN | jq -r '.token')

# Get participants which aren't in any group
http ":8080/participants/not-in-groups?token=$EVENT_TOKEN"
