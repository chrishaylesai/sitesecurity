export interface AuthConfig {
  issuerUrl: string;
  clientId: string;
  realm: string;
}

export interface AuthUser {
  sub: string;
  email: string;
  name: string;
  roles: string[];
}

export interface AuthProvider {
  login: () => Promise<void>;
  logout: () => Promise<void>;
  getUser: () => AuthUser | null;
  getToken: () => Promise<string | null>;
  isAuthenticated: () => boolean;
}
