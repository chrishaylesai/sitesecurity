import type { AuthProvider, AuthUser, AuthConfig } from "./provider";
import {
  generateCodeVerifier,
  generateCodeChallenge,
  generateState,
} from "./pkce";

interface TokenData {
  accessToken: string;
  refreshToken: string;
  idToken: string;
  expiresAt: number;
}

interface TokenResponse {
  access_token: string;
  refresh_token: string;
  id_token: string;
  expires_in: number;
}

let storedTokens: TokenData | null = null;

export function getAccessToken(): string | null {
  if (!storedTokens) return null;
  if (Date.now() >= storedTokens.expiresAt) {
    storedTokens = null;
    return null;
  }
  return storedTokens.accessToken;
}

function clearTokens(): void {
  storedTokens = null;
}

function parseJWT(token: string): Record<string, unknown> {
  const base64Url = token.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
      .join("")
  );
  return JSON.parse(jsonPayload);
}

function userFromToken(token: string): AuthUser {
  const claims = parseJWT(token);
  const realmAccess = claims.realm_access as
    | { roles?: string[] }
    | undefined;
  return {
    sub: (claims.sub as string) || "",
    email: (claims.email as string) || "",
    name:
      (claims.name as string) ||
      (claims.preferred_username as string) ||
      "",
    roles: realmAccess?.roles || [],
  };
}

export class KeycloakAuthProvider implements AuthProvider {
  private config: AuthConfig;
  private currentUser: AuthUser | null = null;

  constructor(config: AuthConfig) {
    this.config = config;
    const token = getAccessToken();
    if (token) {
      try {
        this.currentUser = userFromToken(token);
      } catch {
        clearTokens();
      }
    }
  }

  private get issuerUrl(): string {
    return `${this.config.issuerUrl}/realms/${this.config.realm}`;
  }

  async login(): Promise<void> {
    const codeVerifier = generateCodeVerifier();
    const codeChallenge = await generateCodeChallenge(codeVerifier);
    const state = generateState();

    sessionStorage.setItem("pkce_code_verifier", codeVerifier);
    sessionStorage.setItem("oauth_state", state);

    const url = new URL(
      `${this.issuerUrl}/protocol/openid-connect/auth`
    );
    url.searchParams.set("client_id", this.config.clientId);
    url.searchParams.set(
      "redirect_uri",
      `${window.location.origin}/auth/callback`
    );
    url.searchParams.set("response_type", "code");
    url.searchParams.set("scope", "openid profile email");
    url.searchParams.set("state", state);
    url.searchParams.set("code_challenge", codeChallenge);
    url.searchParams.set("code_challenge_method", "S256");

    window.location.href = url.toString();
  }

  async logout(): Promise<void> {
    this.currentUser = null;
    clearTokens();

    const url = new URL(
      `${this.issuerUrl}/protocol/openid-connect/logout`
    );
    url.searchParams.set("client_id", this.config.clientId);
    url.searchParams.set(
      "post_logout_redirect_uri",
      window.location.origin
    );

    window.location.href = url.toString();
  }

  async handleCallback(code: string, state: string): Promise<void> {
    const storedState = sessionStorage.getItem("oauth_state");
    if (state !== storedState) {
      throw new Error("Invalid state parameter");
    }

    const codeVerifier = sessionStorage.getItem("pkce_code_verifier");
    if (!codeVerifier) {
      throw new Error("Missing code verifier");
    }

    const tokenUrl = `${this.issuerUrl}/protocol/openid-connect/token`;
    const body = new URLSearchParams({
      grant_type: "authorization_code",
      code,
      redirect_uri: `${window.location.origin}/auth/callback`,
      client_id: this.config.clientId,
      code_verifier: codeVerifier,
    });

    const response = await fetch(tokenUrl, {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: body.toString(),
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(`Token exchange failed: ${text}`);
    }

    const tokens: TokenResponse = await response.json();

    storedTokens = {
      accessToken: tokens.access_token,
      refreshToken: tokens.refresh_token,
      idToken: tokens.id_token,
      expiresAt: Date.now() + tokens.expires_in * 1000,
    };

    this.currentUser = userFromToken(tokens.id_token);

    sessionStorage.removeItem("pkce_code_verifier");
    sessionStorage.removeItem("oauth_state");
  }

  getUser(): AuthUser | null {
    return this.currentUser;
  }

  async getToken(): Promise<string | null> {
    return getAccessToken();
  }

  isAuthenticated(): boolean {
    return this.currentUser !== null && getAccessToken() !== null;
  }
}
