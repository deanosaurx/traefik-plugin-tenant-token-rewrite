# Traefik Plugin: Tenant Token Rewrite

Rewrites a token endpoint path using tenant extracted from host.

Example:
- Request host: `beta.dev7.plainid.cloud`
- Request path: `/openid-connect/token`
- Rewritten path: `/auth/realms/beta/protocol/openid-connect/token`

## Config

- `domainSuffix` (string, required): e.g. `dev7.plainid.cloud`
- `sourcePath` (string, optional): default `/openid-connect/token`
- `targetTemplate` (string, optional): default `/auth/realms/{tenant}/protocol/openid-connect/token`

`{tenant}` is replaced by the first host label before `domainSuffix`.
