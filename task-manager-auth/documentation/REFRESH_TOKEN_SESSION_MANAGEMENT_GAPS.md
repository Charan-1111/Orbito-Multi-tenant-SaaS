# Refresh Token and Session Management Gaps

This document explains in detail why the following items are important in a production-grade authentication system and why they are still considered missing in this project:

- No storage of refresh tokens in a session table.
- No invalidation on logout.
- No reuse detection or token revocation.
- No rotation history tracking.
- No session lifecycle management.

This is not just a code detail. It is a core part of secure authentication design.

---

## 1. Why refresh tokens exist

Access tokens are short-lived and used to authorize requests. They are usually valid for a short time such as 5 to 30 minutes.

Refresh tokens are longer-lived and used to obtain a new access token without forcing the user to log in again.

This is useful because:
- users do not have to log in repeatedly
- short-lived access tokens reduce exposure if a token is stolen
- the system can rotate credentials without forcing the user out

But refresh tokens are also sensitive. If a refresh token is stolen, an attacker may keep generating access tokens and continue acting as the user.

That is why refresh tokens must be treated as a session credential, not as a simple JWT value that can be rotated casually without tracking state.

---

## 2. No storage of refresh tokens in a session table

A refresh token should usually be linked to a user session.

In a typical design, the system stores something like:
- user ID
- session ID
- refresh token hash
- issued at time
- expires at time
- revoked flag
- replaced by token ID or parent token ID
- device or client metadata
- last used timestamp

This lets the system answer questions like:
- Is this refresh token still valid?
- Was it revoked?
- Was it replaced by a newer one?
- Is this token from a device that should no longer be trusted?

### Why this matters

If refresh tokens are not stored, then the system cannot determine:
- whether a token was logged out
- whether a token was stolen and reused
- whether multiple tokens are part of the same session
- whether a token is still active or expired

In a JWT-only implementation, the token itself may appear valid because it is cryptographically signed. But that does not mean it is still authorized by the current session state.

### Example problem

Imagine a user logs in on a browser and gets a refresh token.

Later that token is copied by a malicious script or stored in an insecure client.

If the system does not store and track refresh tokens, then the server cannot mark that token as invalid once the user logs out or changes security posture.

The token may still appear valid for the lifetime of the JWT expiry, even though the session is no longer trusted.

---

## 3. No invalidation on logout

Logout should do more than delete the client token from the browser.

A secure system should invalidate the server-side session state.

### What logout should do

On logout, the system should:
- invalidate the active refresh token
- revoke the corresponding session
- optionally revoke all active sessions for the user
- prevent future token refreshes from that session

### Why this matters

If logout does not invalidate the refresh token, then the client might still continue holding a valid refresh token, and the server may still accept it.

This creates a serious problem:
- the user believes they are logged out
- the application may still accept refresh requests from the old session
- stolen or leaked tokens remain usable

### Example

A user logs out from the app, but refresh token is still active in the database or still accepted by the API.

A malicious actor that has the token can continue to refresh the session and receive access tokens.

That means logout is ineffective.

---

## 4. No reuse detection or token revocation

A refresh token should not be reusable forever.

Many systems implement one of the following:
- rotation: every refresh generates a new refresh token and invalidates the old one
- revocation: a token can be manually or automatically marked invalid
- reuse detection: if the same refresh token is used more than once, the system treats it as suspicious and ends the session

### Why reuse detection matters

Token replay is a common attack pattern.

If the same refresh token is used twice, it often indicates:
- token theft
- interception or replay by an attacker
- client bug or duplicate requests
- compromised device or app

If the server sees the same refresh token used again, it should detect that and revoke the session.

### What happens without it

Without reuse detection:
- a stolen token can be used repeatedly
- the attacker can keep rotating tokens indefinitely
- the legitimate user may not know their session is compromised
- there is no way to identify suspicious patterns

### Best practice

When a refresh token is used:
1. look up the session record by token hash
2. confirm the session is still active
3. check whether the token has already been rotated or revoked
4. if it has, reject and terminate the session
5. generate new access and refresh tokens
6. mark the old token as rotated and record the new one

This is common in secure session management.

---

## 5. No rotation history tracking

Refresh token rotation means a session does not keep using the same refresh token forever. Instead, it issues a new one each time.

But rotation needs history.

A system should track:
- original refresh token
- replaced refresh token
- parent token ID
- new token ID
- issued at time
- reason for rotation
- current active token

### Why this matters

Without rotation history, the system cannot tell:
- whether a token is the latest valid token
- whether a token was reused after being replaced
- whether a session has been hijacked
- whether a user has multiple active refresh tokens from one login

### Example

User logs in and gets refresh token A.

The app calls refresh and receives token B.

An attacker reuses token A.

If token history is not tracked, the server may accept A again because it does not know it was already replaced.

This is a classic replay vulnerability.

### secure flow

Instead of accepting any valid-looking token, a smart auth system can record this chain:

- login -> token A (active)
- refresh -> token B (new active)
- old token A marked as rotated
- if A is used again -> revoke session

That gives the system control and traceability.

---

## 6. No session lifecycle management

A session is not just a token. It is the full lifecycle of a user’s authenticated state.

Session lifecycle includes:
- session creation at login
- session expiration
- refresh flow
- rotation events
- logout or revocation
- expiration cleanup
- device tracking or client identification
- inactivity timeout

### Why lifecycle management matters

Without session lifecycle management, tokens become orphaned:
- valid tokens may continue to exist even after a user has been inactive for a long time
- sessions may continue even after the user logs out
- stale tokens may remain usable for unknown amounts of time
- there is no clear source of truth for active user sessions

### A proper session model usually includes

- session ID
- user ID
- issued at
- last activity
- expires at
- revoked at
- replaced by
- client or device metadata
- status: active, expired, revoked, rotated

This gives the system a real session state instead of relying only on token claims.

---

## 7. Why this is a security risk

The most serious issue is that the app has token generation and validation, but no real way to say:
- which refresh token is active
- which session has been revoked
- whether a refresh token was reused
- whether the user is still logged in
- whether a refresh token was stolen and is being abused

This means the app is vulnerable to problems such as:
- session persistence after logout
- replay attacks
- token hijacking
- indefinite active tokens
- inability to revoke sessions at scale

A cryptographically valid JWT is not the same as a valid active session.

---

## 8. Why the current flow is incomplete

The auth service has basic JWT logic, which is useful, but it is not session-aware.

In other words, it can validate a token structure, but it cannot answer the question:

"Is this refresh token still the current valid token for this user session?"

That missing state is exactly what session tables and rotation tracking are for.

Without that state, the app only knows that a token is signed correctly. It does not know whether that token is still authorized, active, replaced, revoked, or part of a current user session.

---

## 9. What a stronger design would look like

A proper refresh-token/session model often includes:

- `sessions` table
  - session_id
  - user_id
  - issued_at
  - last_seen_at
  - expires_at
  - revoked_at
  - status
  - device_id or client metadata

- `refresh_tokens` table
  - token_id
  - session_id
  - token_hash
  - issued_at
  - expires_at
  - revoked_at
  - rotated_from_token_id
  - is_active

The login flow would then:
1. create a session
2. create a refresh token for that session
3. hash the refresh token before storing it
4. issue access token + refresh token

On refresh:
1. read the token hash from the request
2. find matching session record
3. verify token is active and not revoked
4. detect reuse if token has already been rotated
5. issue new token pair and rotate session state

On logout:
1. mark session revoked
2. revoke active refresh token
3. prevent future refreshes

---

## 10. Summary

The missing items listed in the project are not minor cleanup tasks. They are the foundation of secure session management.

In simple terms:

- No refresh token storage means there is no source of truth for active sessions.
- No logout invalidation means users can remain logged in after they believe they logged out.
- No reuse detection means stolen tokens can be replayed.
- No rotation history means the system cannot tell whether a token was replaced or reused.
- No session lifecycle management means the app cannot safely govern login state over time.

This is why these are called out as major gaps in the roadmap and why they matter in real authentication systems.

---

## Final takeaway

A JWT-based auth system is only as strong as the session state behind it.

A system that generates access and refresh tokens but does not track their lifecycle is not a fully secure session model. It is only a basic token issuer.

That is the core concept behind these gaps.
