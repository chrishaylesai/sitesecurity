"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import type { AuthUser } from "./provider";
import { KeycloakAuthProvider } from "./keycloak";

interface AuthContextValue {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: () => Promise<void>;
  logout: () => Promise<void>;
  getToken: () => Promise<string | null>;
  handleCallback: (code: string, state: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function createProvider(): KeycloakAuthProvider {
  return new KeycloakAuthProvider({
    issuerUrl:
      process.env.NEXT_PUBLIC_AUTH_URL || "http://localhost:8180",
    clientId:
      process.env.NEXT_PUBLIC_AUTH_CLIENT_ID || "sitesecurity-frontend",
    realm: process.env.NEXT_PUBLIC_AUTH_REALM || "sitesecurity",
  });
}

export function AuthContextProvider({ children }: { children: ReactNode }) {
  const [provider] = useState<KeycloakAuthProvider>(createProvider);
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setUser(provider.getUser());
    setIsLoading(false);
  }, [provider]);

  const handleCallback = useCallback(
    async (code: string, state: string) => {
      await provider.handleCallback(code, state);
      setUser(provider.getUser());
    },
    [provider]
  );

  const value: AuthContextValue = {
    user,
    isAuthenticated: user !== null && provider.isAuthenticated(),
    isLoading,
    login: provider.login.bind(provider),
    logout: provider.logout.bind(provider),
    getToken: provider.getToken.bind(provider),
    handleCallback,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthContextProvider");
  }
  return context;
}
