# Overview

This project is a simple REST API for managing a list of todo items. It provides core functionality for users to view all existing todos, create new ones, and delete those they no longer need. The system is intended to serve as a lightweight backend service that can be consumed by web, mobile, or other client applications.

The target users are client application developers who need a straightforward todo management service, as well as end users who interact indirectly through those clients. The API emphasizes simplicity, predictability, and standard REST conventions so that it is easy to integrate and reason about.

The high-level approach is to expose a small, focused set of HTTP endpoints that operate on a single "todo" resource, with clear request/response contracts, consistent error handling, and persistent storage of todo items.

# Capabilities

## Todo Resource Model
- Each todo item must have a unique identifier that is automatically generated upon creation.
- Each todo item must have a title field containing a short textual description of the task.
- Each todo item must have a completion status indicating whether the task is done or pending.
- Each todo item must have a creation timestamp recording when it was added.
- The title must be a non-empty string with a maximum length limit (e.g., 255 characters).

## List Todos
- The API must provide an endpoint to retrieve all todo items.
- The list response must return todos in a consistent, predictable order (e.g., by creation time).
- The list endpoint must return an empty array when no todos exist, not an error.
- The list response must include all fields of each todo item.
- The list endpoint must respond with HTTP 200 on success.

## Create Todo
- The API must provide an endpoint to create a new todo item.
- The create endpoint must accept a JSON request body containing at minimum a title.
- The create endpoint must reject requests with a missing or empty title and return a clear validation error.
- The create endpoint must reject requests with malformed JSON and return a clear error.
- On success, the create endpoint must return the newly created todo, including its generated identifier and timestamp.
- The create endpoint must respond with HTTP 201 on successful creation.
- Newly created todos must default to an incomplete/pending status unless otherwise specified.

## Delete Todo
- The API must provide an endpoint to delete a todo item by its unique identifier.
- The delete endpoint must return a clear error when the specified todo does not exist.
- The delete endpoint must respond with HTTP 204 (or 200) on successful deletion.
- Deleted todos must not appear in subsequent list responses.
- Deletion must be permanent and not recoverable through the API.

## Error Handling & Responses
- All error responses must use a consistent JSON structure containing an error message and, where relevant, a code.
- The API must use appropriate HTTP status codes (400 for validation errors, 404 for not found, 500 for server errors).
- All successful and error responses must use JSON as the content type.
- The API must handle unexpected server errors gracefully without exposing internal details.

## Data Persistence
- Todo items must persist across API restarts.
- Concurrent requests must not corrupt stored data or produce duplicate identifiers.

## Non-Functional Requirements
- The API must respond to standard requests within a reasonable time (e.g., under 500ms under normal load).
- The API must be documented with endpoint definitions, request/response schemas, and example payloads.
- The API must log requests and errors sufficiently to support debugging and monitoring.
- The API must be accessible over HTTP(S) using standard REST conventions.
- The API should support a reasonable baseline of concurrent clients without degradation.
